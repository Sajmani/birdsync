// Birdsync syncs eBird observations, photos, and sounds to iNaturalist.
//
// See README.md for detailed documentation.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

const UserAgent = "birdsync/0.1"

var (
	debug              bool
	dryRun             bool
	verifiable         bool
	fuzzy              bool
	before             dateTimeFlag
	after              dateTimeFlag
	positionalAccuracy int

	// Seams for testing the 24-hour iNaturalist API request budget
	// (tech.dailyRequestBudget) deterministically without real clocks or
	// nine thousand mock records.
	dailyRequestBudget = inat.DefaultMaxRequests
	nowFunc            = time.Now
)

func init() {
	flag.BoolVar(&debug, "debug", false,
		"Log verbosely")
	flag.BoolVar(&dryRun, "dryrun", false,
		"Don't actually sync any observations, just log what birdsync would do")
	flag.BoolVar(&verifiable, "verifiable", true,
		"Sync only observations that include Macaulay Catalog Numbers (photos or sound)")
	flag.BoolVar(&fuzzy, "fuzzy", false,
		"Don't create a birdsync observation if an observation without a complete birdsync sync key already exists for the same species on the same date. "+
			"This fuzzy matching is useful when you've entered the same observation manually into both eBird and iNaturalist, "+
			"but it may skip legitimate uploads if you saw the same species (or two subspecies of it) twice on the same day.")
	flag.Var(&before, "before",
		"Sync only observations observed before the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).")
	flag.Var(&after, "after",
		"Sync only observations observed after the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).")
	flag.IntVar(&positionalAccuracy, "positional_accuracy_meters", ebird.PositionalAccuracy,
		"Positional accuracy in meters of the iNaturalist observations created by birdsync.")
}

func debugf(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

func prettyPrintln(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

type stats struct {
	afterSkips, beforeSkips, verifiableSkips, previouslySkips, fuzzySkips int
	// invalidSkips counts records whose date, time, or coordinates could not be
	// parsed. They are skipped rather than fatal: one bad row in a large export
	// should not end a sync that has already created observations (product.badRecordsSkipped).
	invalidSkips                                           int
	totalRecords, createdObservations, updatedObservations int
	uploadedPhotos, uploadedSounds                         int
	// pendingMedia counts the media assets a --dryrun would have uploaded.
	// A Macaulay Library asset ID doesn't say whether it's a photo or a sound,
	// and a dry run doesn't download it to find out, so these can't be split
	// into uploadedPhotos and uploadedSounds.
	pendingMedia         int
	errors               int
	stoppedAtDailyBudget bool
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("usage: birdsync MyEBirdData.csv")
		flag.Usage()
		os.Exit(1)
	}
	if !after.Time().IsZero() && !before.Time().IsZero() && after.Time().After(before.Time()) {
		log.Fatalf("--after (%s) is after --before (%s), won't match any records",
			after.Time(), before.Time())
	}

	eBirdCSVFilename := flag.Arg(0)
	if f, err := os.Open(eBirdCSVFilename); err != nil {
		log.Fatalf("Can't open %s: %v", eBirdCSVFilename, err)
	} else {
		f.Close()
	}

	inatAPIClient := inatClientImpl{
		client: inat.NewClient(inat.BaseURL, inat.GetAPIToken(), UserAgent),
	}
	ebirdAPIClient := ebirdClientImpl{}

	stats := birdsync(eBirdCSVFilename, ebirdAPIClient, inat.GetUserID(), inatAPIClient)

	for _, line := range stats.summary() {
		log.Print(line)
	}
}

// summary returns the end-of-run report, one line per entry.
//
// It lives outside main so it can be tested. A summary is the only account of
// the run most users read, and it is the one place where saying something
// untrue is both easy and invisible: the counters and their labels are written
// in different places, so nothing stops them drifting apart.
func (s stats) summary() []string {
	var lines []string
	add := func(format string, args ...any) {
		lines = append(lines, fmt.Sprintf(format, args...))
	}

	add("Finished processing %d eBird observations", s.totalRecords)
	add("Skipped %d previously uploaded by birdsync", s.previouslySkips)
	// A skip counter is only meaningful when the rule producing it was in
	// effect, so those lines are conditional (product.conditionalSkipCounters).
	if fuzzy {
		add("Skipped %d eBird observations with --fuzzy matching", s.fuzzySkips)
	}
	if !after.Time().IsZero() {
		add("Skipped %d eBird observations before --after", s.afterSkips)
	}
	if !before.Time().IsZero() {
		add("Skipped %d eBird observations after --before", s.beforeSkips)
	}
	if verifiable {
		add("Skipped %d unverifiable eBird observations", s.verifiableSkips)
	}
	if s.invalidSkips > 0 {
		add("Skipped %d eBird observations with unparseable fields", s.invalidSkips)
	}

	if dryRun {
		// The counters are incremented outside the --dryrun gates, so they
		// count what a real run would have done. That number is worth
		// reporting, but reporting it as "Created" would be a lie, and
		// --dryrun is the one feature whose whole value is being trusted
		// (product.dryRunCounterLabels, tech.honestDryRunCounters).
		add("Would create %d new iNaturalist observations", s.createdObservations)
		add("Would update %d iNaturalist observations", s.updatedObservations)
		// A dry run doesn't download the assets, so it can't tell photos from sounds.
		add("Would upload %d media assets to iNaturalist", s.pendingMedia)
	} else {
		add("Created %d new iNaturalist observations", s.createdObservations)
		add("Updated %d iNaturalist observations", s.updatedObservations)
		add("Uploaded %d photos to iNaturalist", s.uploadedPhotos)
		add("Uploaded %d sounds to iNaturalist", s.uploadedSounds)
	}
	if s.errors > 0 {
		add("Failed to upload %d media assets", s.errors)
	}
	if s.stoppedAtDailyBudget {
		add("Stopped at daily iNaturalist API request budget (%d requests/24h); re-run in 24 hours to continue", dailyRequestBudget)
	}
	return lines
}

// downloadPageRequests returns the number of GET /observations requests
// DownloadObservations made to return resultCount observations. It always stops
// on the first page shorter than DownloadPageSize (including an empty page
// after an exact multiple), which is (resultCount / DownloadPageSize) + 1.
func downloadPageRequests(resultCount int) int {
	return (resultCount / inat.DownloadPageSize) + 1
}

// parseCreatedAt parses the created_at timestamp returned by the iNaturalist
// API. v2 returns ISO-8601/RFC3339 timestamps, with fallbacks for alternate
// layouts seen across API responses.
func parseCreatedAt(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05-0700",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05 MST",
		time.DateOnly,
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// estimateRecentWriteRequests estimates how many write requests have already
// been made on the user's iNaturalist account in the 24 hours ending at now,
// using only the observations birdsync already downloaded
// (decisions.dailyRequestBudget). Each observation created in that window took
// one create request, plus one upload per attached or permanently failed asset,
// plus one description update when assets were handled.
func estimateRecentWriteRequests(results []inat.Result, now time.Time) int {
	cutoff := now.Add(-24 * time.Hour)
	var total int
	for _, r := range results {
		createdAt, ok := parseCreatedAt(r.CreatedAt)
		if !ok || createdAt.Before(cutoff) {
			continue
		}
		total++ // POST /observations
		uploaded, failed := iNatMLAssets(r)
		mediaAttempts := max(len(r.Photos)+len(r.Sounds), uploaded.Len()) + failed.Len()
		total += mediaAttempts
		if mediaAttempts > 0 {
			total++ // PUT /observations/{uuid}
		}
	}
	return total
}

// speciesBinomial returns the two-word "Genus species" stem of sciName when it
// names a single species or subspecies (for example "Garrulus glandarius" from
// "Garrulus glandarius leucotis" or "Junco hyemalis [oreganus Group]"). Names
// that don't identify one species — genera or families ("Corvidae"),
// unidentified sp. entries ("Empidonax sp."), species slashes
// ("Larus glaucoides/thayeri"), and hybrids ("x" or "×") — return "", false
// (product.fuzzyMatching, decisions.fuzzySpeciesBinomial).
func speciesBinomial(sciName string) (string, bool) {
	tokens := strings.Fields(sciName)
	if len(tokens) < 2 {
		return "", false
	}
	for _, tok := range tokens {
		if tok == "x" || tok == "×" {
			return "", false
		}
	}
	genus, epithet := tokens[0], tokens[1]
	if !isCapitalizedWord(genus) || !isSpeciesEpithet(epithet) || epithet == "sp" {
		return "", false
	}
	return genus + " " + epithet, true
}

func isCapitalizedWord(s string) bool {
	if len(s) < 2 || s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if s[i] < 'a' || s[i] > 'z' {
			return false
		}
	}
	return true
}

func isSpeciesEpithet(s string) bool {
	if len(s) < 2 || s[0] < 'a' || s[0] > 'z' || s[len(s)-1] < 'a' || s[len(s)-1] > 'z' {
		return false
	}
	for i := 0; i < len(s); i++ {
		if (s[i] < 'a' || s[i] > 'z') && s[i] != '-' {
			return false
		}
	}
	return true
}

func birdsync(eBirdCSVFilename string, ebirdClient ebirdClient, inatUserID string, inatClient inatClient) stats {
	results, err := inatClient.DownloadObservations(inatUserID, after.Time(), before.Time(),
		"created_at", "description", "observed_on", "photos.all", "sounds.all", "taxon.all", "ofvs.all")
	if err != nil {
		// Nothing useful can happen without the existing observations: syncing
		// blind would duplicate everything the user already has.
		log.Fatalf("Downloading iNaturalist observations: %v", err)
	}

	previouslySynced := map[ebird.ObservationID]inat.Result{}
	type fuzzyKey struct {
		observedDate string // 2006-01-02
		name         string
	}
	fuzzyMatch := map[fuzzyKey][]string{}
	for _, r := range results {
		key := ebird.ObservationID{
			SubmissionID:   r.ObservationFieldValue(inat.EBirdField),
			ScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if key.Valid() {
			previouslySynced[key] = r
		} else {
			// This iNaturalist observation was not created by birdsync.
			// Record its date and common name for fuzzy matching.
			addFuzzy := func(name string) {
				if name == "" {
					return // an empty name would match every unnamed eBird record
				}
				key := fuzzyKey{
					observedDate: r.ObservedOn, // iNaturalist always uses format 2006-01-02
					name:         name,
				}
				fuzzyMatch[key] = append(fuzzyMatch[key], r.UUID.String())
				slices.Sort(fuzzyMatch[key])
				debugf("fuzzy match: add %s to %+v", r.UUID, key)
			}
			addFuzzy(r.Taxon.PreferredCommonName)
			addFuzzy(r.Taxon.Name)
			if stem, ok := speciesBinomial(r.Taxon.Name); ok && stem != r.Taxon.Name {
				addFuzzy(stem)
			}
		}
	}
	debugf("Previously synced %d observations\n", len(previouslySynced))

	usedRequests := downloadPageRequests(len(results)) + estimateRecentWriteRequests(results, nowFunc())

	log.Printf("Reading eBird observations from %s", eBirdCSVFilename)
	records, err := ebirdClient.Records(eBirdCSVFilename)
	if err != nil {
		log.Fatal(err)
	}
	var s stats
	for rec := range records {
		s.totalRecords++
		observed, err := rec.Observed()
		if err != nil {
			// A malformed row costs that row, not the run. eBird's export
			// varies between users, so one unparseable date in a twelve
			// thousand row file is realistic, and aborting partway would
			// leave the sync half done (product.badRecordsSkipped).
			log.Printf("line %d: SKIPPING record with bad date/time: %v", rec.Line, err)
			s.invalidSkips++
			continue
		}
		// Skip records that were not observed between --after and --before.
		if !after.Time().IsZero() && observed.Before(after.Time()) {
			debugf("line %d: SKIPPING record observed on %s (before --after=%s)",
				rec.Line, observed, after.Time())
			s.afterSkips++
			continue
		}
		if !before.Time().IsZero() && observed.After(before.Time()) {
			debugf("line %d: SKIPPING record observed on %s (after --before=%s)",
				rec.Line, observed, before.Time())
			s.beforeSkips++
			continue
		}

		// addMedia uploads the Maculay Library assets in assetIDs to iNaturalist
		// then appends the asset URLs to the description of observation u.
		addMedia := func(u uuid.UUID, desc string, assetIDs mlAssetSet) {
			if assetIDs.Len() == 0 {
				return
			}
			debugf("Adding %d media assets to %s\n",
				assetIDs.Len(), inat.ObservationURL(u))

			obs := inat.Observation{
				UUID:        u,
				Description: desc,
			}
			// uploaded collects the assets that actually made it, so the
			// description can be built from those alone. The description is
			// how birdsync remembers what it has uploaded — iNatMLAssets reads
			// the URLs back out of it on the next run — so listing an asset
			// that failed makes the failure permanent as well as untrue
			// (product.descriptionContents, decisions.descriptionListsOnlyUploaded).
			var uploaded, permanentlyFailed mlAssetSet
			// Upload the media
			for _, id := range assetIDs.ids {
				if dryRun {
					log.Printf("DRYRUN: Download ML Asset %s and upload to iNaturalist", id)
					s.pendingMedia++
					// A dry run reports what a successful run would do, so the
					// printed description shows the asset as listed.
					uploaded.Add(id)
				} else {
					filename, isPhoto, err := ebirdClient.DownloadMLAsset(id)
					if err != nil {
						log.Printf("Couldn't download ML asset %s from eBird: %v", id, err)
						s.errors++
						continue
					}
					err = inatClient.UploadMedia(filename, isPhoto, id, obs.UUID.String())
					usedRequests++
					// The download is a temp file that belongs to us now, so
					// remove it whether or not the upload worked. Syncing an
					// account with thousands of assets used to leave one file
					// per asset behind (tech.tempFilesDeleted). Removing here rather than in a
					// defer keeps at most one asset on disk at a time.
					if rmErr := os.Remove(filename); rmErr != nil {
						debugf("Couldn't remove temp file %s: %v", filename, rmErr)
					}
					if err != nil {
						log.Printf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
						s.errors++
						// A refusal of the file itself won't come good on a
						// later run, so record it rather than re-downloading
						// and re-uploading it forever (product.permanentUploadFailureRecorded).
						var statusErr *inat.StatusError
						if errors.As(err, &statusErr) && statusErr.Permanent() {
							permanentlyFailed.Add(id)
						}
						continue
					}
					if isPhoto {
						s.uploadedPhotos++
					} else {
						s.uploadedSounds++
					}
					uploaded.Add(id)
				}
			}
			if uploaded.Len() == 0 && permanentlyFailed.Len() == 0 {
				// Everything failed, and might yet succeed. Don't write an
				// unchanged description back, and don't count an update that
				// didn't happen (tech.honestDryRunCounters). The next run tries again.
				return
			}
			for _, id := range uploaded.ids {
				obs.Description += assetLine(id, true)
			}
			for _, id := range permanentlyFailed.ids {
				obs.Description += assetLine(id, false)
			}
			// Update the description
			if dryRun {
				log.Printf("DRYRUN: Updating observation %s with %d added media assets\n",
					obs.URLWithSpecies(), uploaded.Len())
				prettyPrintln(obs)
			} else {
				err = inatClient.UpdateObservation(obs)
				usedRequests++
				if err != nil {
					log.Fatalf("UpdateObservation %s: %v", obs.URLWithSpecies(), err)
				}
			}
			s.updatedObservations++
		}

		// Skip records that have previously been uploaded by birdsync.
		key := rec.ObservationID()
		if r, ok := previouslySynced[key]; ok {
			debugf("line %d: Already synced %s to iNaturalist as %s\n",
				rec.Line, key, r.URLWithSpecies())
			addedMediaIDs, summary := mediaChange(rec, r)
			if summary != "" {
				log.Printf("Media assets differ between eBird %s and iNaturalist %s: %s",
					rec.URLWithSpecies(), r.URLWithSpecies(), summary)
			}
			if addedMediaIDs.Len() == 0 {
				s.previouslySkips++
				continue
			}
			if !dryRun {
				needed := addedMediaIDs.Len() + 1
				if usedRequests+needed > dailyRequestBudget {
					log.Printf("Stopping before line %d (%s): %d estimated iNaturalist API requests in the last 24 hours; adding %d media asset(s) needs %d more, exceeding the %d daily budget. Re-run in 24 hours to continue.",
						rec.Line, key, usedRequests, addedMediaIDs.Len(), needed, dailyRequestBudget)
					s.totalRecords--
					s.stoppedAtDailyBudget = true
					break
				}
			}
			addMedia(r.UUID, r.Description, addedMediaIDs)
			continue
		}

		if fuzzy {
			// Skip records for the same bird and date as an existing non-birdsync observation.
			// eBird writes dates in several formats, so compare against the parsed
			// observation date rather than the raw CSV field, which may be "1/2/2006".
			checkFuzzy := func(name string) bool {
				if name == "" {
					return false // an empty name would match every unnamed taxon
				}
				key := fuzzyKey{
					name:         name,
					observedDate: observed.Format(time.DateOnly),
				}
				debugf("line %d: fuzzy match: check %+v", rec.Line, key)
				if _, ok := fuzzyMatch[key]; ok {
					log.Printf("line %d: SKIPPING fuzzy match: observation for same bird and date: %+v", rec.Line, key)
					s.fuzzySkips++
					return true
				}
				return false
			}
			matched := checkFuzzy(rec.CommonName) || checkFuzzy(rec.ScientificName)
			if !matched {
				if stem, ok := speciesBinomial(rec.ScientificName); ok && stem != rec.ScientificName {
					matched = checkFuzzy(stem)
				}
			}
			if matched {
				continue
			}
		}

		assetIDs := eBirdMLAssets(rec.MLCatalogNumbers)
		// Skip records without media assets if --verifiable is set. This is
		// checked before the record is parsed any further, so that a record is
		// counted against the rule that actually skipped it (product.skipOrder).
		if verifiable && assetIDs.Len() == 0 {
			debugf("line %d: SKIPPING record that has no photos or sounds (--verifiable=true)", rec.Line)
			s.verifiableSkips++
			continue
		}

		// Create the iNaturalist observation from the eBird record.
		coordinate := func(s string) (float64, error) {
			if s == "" {
				return 0, nil // eBird omits coordinates for some locations
			}
			return strconv.ParseFloat(s, 64)
		}
		latitude, err := coordinate(rec.Latitude)
		if err != nil {
			log.Printf("line %d: SKIPPING record with bad latitude %q: %v", rec.Line, rec.Latitude, err)
			s.invalidSkips++
			continue
		}
		longitude, err := coordinate(rec.Longitude)
		if err != nil {
			log.Printf("line %d: SKIPPING record with bad longitude %q: %v", rec.Line, rec.Longitude, err)
			s.invalidSkips++
			continue
		}
		keyField := func(id int, s string) inat.ObservationFieldValue {
			return inat.ObservationFieldValue{
				ObservationFieldID: id,
				Value:              s,
			}
		}
		obs := inat.Observation{
			UUID:               uuid.New(),
			CaptiveFlag:        false, // eBird checklists should only include wild birds
			Latitude:           latitude,
			Longitude:          longitude,
			LocationIsExact:    false,
			PositionalAccuracy: float64(positionalAccuracy),
			SpeciesGuess:       rec.ScientificName,
			ObservedOnString:   rec.Date + " " + rec.Time,
			ObservationFieldValuesAttributes: []inat.ObservationFieldValue{
				keyField(inat.CountField, rec.Count),
				keyField(inat.CommonNameField, rec.CommonName),
				keyField(inat.LocationField, rec.Location),
				keyField(inat.CountyField, rec.County),
				keyField(inat.StateOrProvinceField, rec.StateProvince),
				keyField(inat.NumObserversField, rec.NumberOfObservers),
				// EBirdField and EBirdScientificNameField are used to match iNaturalist observations
				// to the corresponding eBird checklist and species entry. We cannot rely on the taxon
				// in the iNaturalist observation because it may be changed after upload.
				keyField(inat.EBirdField, rec.SubmissionID),
				keyField(inat.EBirdScientificNameField, rec.ScientificName),
			},
		}
		obs.Description = "Observation created using github.com/Sajmani/birdsync \n"
		if len(rec.ObservationDetails) > 0 {
			obs.Description += "eBird observation details:\n" +
				rec.ObservationDetails + "\n"
		}
		obs.Description += "Checklist: " + rec.URL() + "\n"
		obs.Description += "Protocol: " + rec.Protocol + "\n"
		if len(rec.ChecklistComments) > 0 {
			obs.Description += "eBird checklist comments:\n" +
				rec.ChecklistComments + "\n"
		}
		if dryRun {
			log.Printf("DRYRUN: Syncing eBird observation %s to iNaturalist (%d media assets)\n",
				key, assetIDs.Len())
			prettyPrintln(obs)
		} else {
			needed := 1
			if assetIDs.Len() > 0 {
				needed += assetIDs.Len() + 1
			}
			if usedRequests+needed > dailyRequestBudget {
				log.Printf("Stopping before line %d (%s): %d estimated iNaturalist API requests in the last 24 hours; syncing this observation needs up to %d more, exceeding the %d daily budget. Re-run in 24 hours to continue.",
					rec.Line, key, usedRequests, needed, dailyRequestBudget)
				s.totalRecords--
				s.stoppedAtDailyBudget = true
				break
			}
			debugf("Syncing eBird observation %s to iNaturalist (%d media assets)\n",
				key, assetIDs.Len())
			err = inatClient.CreateObservation(obs)
			usedRequests++
			if err != nil {
				log.Fatalf("CreateObservation: %v", err)
			}
		}
		s.createdObservations++
		addMedia(obs.UUID, obs.Description, assetIDs)
	}
	return s
}
