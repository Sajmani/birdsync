// Package inat provides types and helper functions for the iNaturalist API.
package inat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DownloadObservations downloads and returns all observations for inatUserID.
// The dates d1 and d2 specify the start and end of the observation date range if nonzero.
// The fields list specifies which fields are populated in the results.
func (c *Client) DownloadObservations(inatUserID string, d1, d2 time.Time, fields ...string) ([]Result, error) {
	const dateFormat = "2006-01-02"
	var d1str, d2str string
	if !d1.IsZero() {
		d1str = " after " + d1.Format(dateFormat)
	}
	if !d2.IsZero() {
		d2str = " before " + d2.Format(dateFormat)
	}
	log.Printf("Downloading observations for %s%s%s", inatUserID, d1str, d2str)

	// From https://www.inaturalist.org/pages/api+recommended+practices:
	// If using the API to fetch a lot of results, please use the highest supported per_page value.
	// For example you can get up to 200 observations in a single request,
	// which would be faster and more efficient than fetching the default 30 results at a time.
	const perPage = 200

	// Page with id_above rather than page numbers. From the same page:
	//
	//	The page and per_page parameters can be used to fetch up to (for many
	//	endpoints) 10k results. An error will be thrown if results beyond 10k
	//	are requested.
	//
	// A birder with more than 10,000 iNaturalist observations could not sync at
	// all (issue #5). Sorting by id ascending and resuming from the last id
	// seen has no such ceiling. Note that total_results shrinks as the cursor
	// advances, so it can't decide when to stop; a short page can.
	var results []Result
	var totalResults int
	var idAbove int
	for {
		u, err := url.Parse(c.baseURL + "/observations")
		if err != nil {
			return nil, fmt.Errorf("DownloadObservations: %w", err)
		}
		q := u.Query()
		q.Set("user_id", inatUserID)
		q.Set("per_page", strconv.Itoa(perPage))
		q.Set("order_by", "id")
		q.Set("order", "asc")
		if idAbove > 0 {
			q.Set("id_above", strconv.Itoa(idAbove))
		}
		// Deliberately unfiltered by taxon. Restricting to iconic_taxa[]=Aves
		// downloaded fewer observations, but hid any observation without an
		// iconic taxon — the state iNaturalist leaves an observation in when it
		// can't resolve the eBird name birdsync gave it. Those became invisible
		// to duplicate detection and were re-created on every run. The saving
		// was two HTTP requests on a 1478-observation account; see CR-003 in
		// spec/decisions.md.
		if !d1.IsZero() {
			q.Set("d1", d1.Format(dateFormat))
		}
		if !d2.IsZero() {
			q.Set("d2", d2.Format(dateFormat))
		}
		// "id" is the paging cursor, and unlike uuid the v2 API returns it only
		// when asked. The client adds it rather than trusting every caller to
		// know that paging depends on it.
		q.Set("fields", strings.Join(append([]string{"id"}, fields...), ","))
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("DownloadObservations: %w", err)
		}
		body, err := c.roundTrip(req)
		if err != nil {
			return nil, fmt.Errorf("DownloadObservations: after id %d: %w", idAbove, err)
		}

		var observations Observations
		err = json.Unmarshal([]byte(body), &observations)
		if err != nil {
			return nil, fmt.Errorf("DownloadObservations: decoding results after id %d: %w", idAbove, err)
		}

		if len(observations.Results) == 0 {
			break
		}
		if totalResults == 0 { // only the first response sees the whole set
			totalResults = observations.TotalResults
		}
		results = append(results, observations.Results...)
		// Insist the cursor moves. If id_above were ignored — stripped by a
		// proxy, unsupported by a future API version — or if "id" came back
		// zero because the field wasn't returned, the loop would otherwise
		// fetch the same page forever, hammering the service and breaching the
		// rate limit it asks clients to respect.
		last := observations.Results[len(observations.Results)-1].ID
		if last <= idAbove {
			return nil, fmt.Errorf(
				"DownloadObservations: cursor did not advance past id %d; the API may be ignoring id_above or omitting id",
				idAbove)
		}
		idAbove = last
		log.Printf("Downloaded %d of %d observations", len(results), totalResults)
		if len(observations.Results) < perPage {
			break
		}
	}
	return results, nil
}

func TestObservation() Observation {
	return Observation{
		UUID:         uuid.New(),
		CaptiveFlag:  true, // casual observation for testing
		Description:  "Testing github.com/Sajmani/birdsync tools",
		SpeciesGuess: "Homo Sapiens",
	}
}
