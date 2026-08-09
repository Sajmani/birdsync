// Taxonfilter answers the question behind CR-003 in spec/decisions.md: does
// iNaturalist's iconic_taxa[]=Aves filter return observations that have no
// taxon at all?
//
// It matters because birdsync creates untaxoned observations by design. When
// iNaturalist cannot resolve an eBird name — a "slash" like
// Aythya marila/affinis, a "spuh" like Melanitta sp., a domestic form — the
// resulting observation has no taxon and therefore no iconic taxon. If the
// filtered query skips those, birdsync never sees them on the next run,
// re-creates them, and the user accumulates duplicates.
//
// The tool fetches the same account twice, once unfiltered and once with the
// Aves filter, and compares the two by UUID.
//
// This tool is READ-ONLY. It issues GET requests only: it cannot create,
// update, or delete anything, which is why it has no --dryrun flag. It is the
// one tool in this directory that is safe to point at an account you care
// about.
//
// Usage:
//
//	go run ./tools/taxonfilter
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

const UserAgent = "birdsync-taxonfilter/0.1"

// fields limits the response to what the comparison needs. taxon.all is
// required because the whole question is which observations have a taxon.
var fields = []string{"observed_on", "taxon.all", "ofvs.all"}

// fetchAll downloads every observation for userID, paging until the server
// stops returning results. When iconicTaxa is non-empty it is sent as the
// iconic_taxa[] parameter; when empty the parameter is omitted entirely, which
// is the difference the tool exists to measure.
//
// It takes a base URL so it can be tested against an httptest server, the same
// seam inat.NewClient uses (T-014).
func fetchAll(baseURL, apiToken, userID, iconicTaxa string) ([]inat.Result, error) {
	const perPage = 200 // T-017: the largest page size the API supports

	var results []inat.Result
	var totalResults int
	for page := 1; ; page++ {
		u, err := url.Parse(baseURL + "/observations")
		if err != nil {
			return nil, fmt.Errorf("fetchAll: %w", err)
		}
		q := u.Query()
		q.Set("user_id", userID)
		q.Set("page", strconv.Itoa(page))
		q.Set("per_page", strconv.Itoa(perPage))
		q.Set("fields", strings.Join(fields, ","))
		if iconicTaxa != "" {
			q.Set("iconic_taxa[]", iconicTaxa)
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("fetchAll: %w", err)
		}
		req.Header.Set("User-Agent", UserAgent) // T-016
		req.Header.Set("Authorization", apiToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetchAll: %w", err)
		}
		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return nil, fmt.Errorf("fetchAll: %s: refresh your INAT_API_TOKEN from https://www.inaturalist.org/users/api_token", resp.Status)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("fetchAll: bad HTTP status: %s", resp.Status)
		}

		var observations inat.Observations
		err = json.NewDecoder(resp.Body).Decode(&observations)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("fetchAll: decoding response: %w", err)
		}

		if observations.TotalResults == 0 || len(observations.Results) == 0 {
			break
		}
		if totalResults == 0 {
			totalResults = observations.TotalResults
		}
		results = append(results, observations.Results...)
		log.Printf("  downloaded %d of %d", len(results), totalResults)
		if len(results) >= totalResults {
			break
		}
	}
	return results, nil
}

// report is what the comparison found. Counts are of observations in the
// unfiltered download.
type report struct {
	total        int // all observations
	aves         int // returned by the iconic_taxa[]=Aves query
	unidentified int // no iconic taxon: the population CR-003 is about
	nonAves      int // an iconic taxon other than Aves
	birdsync     int // carries a complete birdsync sync key

	// missing are observations present unfiltered but absent from the filtered
	// download. These are invisible to birdsync when the filter is in use.
	missing []inat.Result

	// missingUnidentified is the decisive subset: observations with no iconic
	// taxon that the filter hid. An observation with a non-Aves iconic taxon is
	// *supposed* to be filtered out, so only these say anything about CR-003.
	missingUnidentified []inat.Result

	// missingBirdsync is the subset of missing that birdsync created. These are
	// the ones that get duplicated on the next run: birdsync looks for its own
	// sync key, doesn't find it, and creates the observation again.
	missingBirdsync []inat.Result
}

// hasSyncKey reports whether r carries both halves of the birdsync sync key
// (P-019). Matches ebird.ObservationID.Valid, without importing the CSV reader.
func hasSyncKey(r inat.Result) bool {
	return r.ObservationFieldValue(inat.EBirdField) != "" &&
		r.ObservationFieldValue(inat.EBirdScientificNameField) != ""
}

// unidentified reports whether r belongs to no iconic taxon.
//
// Classifying on the iconic taxon rather than on the taxon itself is
// deliberate. It is unclear whether iNaturalist represents an observation
// nobody has identified as a null taxon or as a placeholder such as "Life" or
// "Unknown", and the two would need different tests. It doesn't matter: under
// either representation the observation belongs to no iconic taxon, which is
// exactly the property iconic_taxa[]=Aves filters on. This predicate is
// therefore correct without settling the question.
func unidentified(r inat.Result) bool {
	switch r.Taxon.IconicTaxonName {
	case "", "Unknown", "Life":
		return true
	}
	return false
}

// analyze compares an unfiltered download against a filtered one. It is a pure
// function so the comparison can be tested without a network, which matters
// because AGENTS.md forbids running anything in tools/ (T-011): this logic is
// verified by taxonfilter_test.go and never by a live run.
func analyze(all, filtered []inat.Result) report {
	inFiltered := make(map[uuid.UUID]bool, len(filtered))
	for _, r := range filtered {
		inFiltered[r.UUID] = true
	}

	rep := report{total: len(all), aves: len(filtered)}
	for _, r := range all {
		switch {
		case unidentified(r):
			rep.unidentified++
		case r.Taxon.IconicTaxonName != "Aves":
			rep.nonAves++
		}
		if hasSyncKey(r) {
			rep.birdsync++
		}
		if !inFiltered[r.UUID] {
			rep.missing = append(rep.missing, r)
			if unidentified(r) {
				rep.missingUnidentified = append(rep.missingUnidentified, r)
			}
			if hasSyncKey(r) {
				rep.missingBirdsync = append(rep.missingBirdsync, r)
			}
		}
	}
	return rep
}

func (rep report) print() {
	fmt.Println()
	fmt.Println("Observations in the account (unfiltered):", rep.total)
	fmt.Println("Returned by iconic_taxa[]=Aves:          ", rep.aves)
	fmt.Println("  with no iconic taxon (unidentified):   ", rep.unidentified)
	fmt.Println("  with a non-Aves iconic taxon:          ", rep.nonAves)
	fmt.Println("  created by birdsync (has sync key):    ", rep.birdsync)
	fmt.Println()
	fmt.Println("Present unfiltered but MISSING from the filtered query:", len(rep.missing))
	fmt.Println("  of those, unidentified:      ", len(rep.missingUnidentified))
	fmt.Println("  of those, created by birdsync:", len(rep.missingBirdsync))
	fmt.Println()

	for _, r := range rep.missingBirdsync {
		fmt.Printf("  %s  taxon=%q iconic=%q  eBird=%s [%s]\n",
			r.URL(), r.Taxon.Name, r.Taxon.IconicTaxonName,
			r.ObservationFieldValue(inat.EBirdField),
			r.ObservationFieldValue(inat.EBirdScientificNameField))
	}

	fmt.Println()
	switch {
	case rep.total == 0:
		fmt.Println("VERDICT: INCONCLUSIVE. No observations were downloaded at all.")

	case len(rep.missingUnidentified) > 0:
		fmt.Printf("VERDICT: CR-003 CONFIRMED. The filter hides %d unidentified observations,\n"+
			"%d of them created by birdsync. Any birdsync observation in that state is\n"+
			"re-created on the next run.\n",
			len(rep.missingUnidentified), len(rep.missingBirdsync))

	case rep.unidentified == 0:
		// The decisive case cannot be observed: there is nothing in the account
		// for the filter to hide. Saying "not confirmed" here would be a false
		// all-clear, so say what was actually not tested.
		fmt.Printf("VERDICT: INCONCLUSIVE on the question that matters. This account has no\n"+
			"unidentified observations, so the filter was never given one to drop. The\n"+
			"%d observations it did drop all have a non-Aves iconic taxon, which is the\n"+
			"filter working as intended. CR-003 is neither confirmed nor refuted.\n",
			len(rep.missing))

	default:
		fmt.Printf("VERDICT: CR-003 REFUTED. The account has %d unidentified observations and\n"+
			"the filtered query returned every one of them, so iconic_taxa[]=Aves does not\n"+
			"hide observations that lack an iconic taxon.\n", rep.unidentified)
	}

	if rep.nonAves > 0 {
		// Expressed in requests, not observations: the optimization that
		// removing the filter gives up is worth whole pages, not records.
		const perPage = 200
		pagesAll := (rep.total + perPage - 1) / perPage
		pagesAves := (rep.aves + perPage - 1) / perPage
		fmt.Printf("\nCost of removing the filter (P-061): %d of %d observations are non-Aves.\n"+
			"At %d per page that is %d requests instead of %d — %d extra.\n",
			rep.nonAves, rep.total, perPage, pagesAll, pagesAves, pagesAll-pagesAves)
	}
}

func main() {
	userID := inat.GetUserID()
	apiToken := inat.GetAPIToken()

	log.Println("Downloading all observations (no taxon filter)")
	all, err := fetchAll(inat.BaseURL, apiToken, userID, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Downloading observations with iconic_taxa[]=Aves")
	aves, err := fetchAll(inat.BaseURL, apiToken, userID, "Aves")
	if err != nil {
		log.Fatal(err)
	}

	analyze(all, aves).print()
}
