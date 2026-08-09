// Birdsync syncs eBird observations, photos, and sounds to iNaturalist.
//
// See README.md for detailed documentation.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
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
)

func init() {
	flag.BoolVar(&debug, "debug", false,
		"Log verbosely")
	flag.BoolVar(&dryRun, "dryrun", false,
		"Don't actually sync any observations, just log what birdsync would do")
	flag.BoolVar(&verifiable, "verifiable", true,
		"Sync only observations that include Macaulay Catalog Numbers (photos or sound)")
	flag.BoolVar(&fuzzy, "fuzzy", false,
		"Don't create a birdsync observation if a non-birdsync observation already exists for the same bird on the same date."+
			"This fuzzy matching is useful when you've entered the same observation manually into both eBird and iNaturalist, "+
			"but it may skip legitimate uploads if you saw the same bird twice on the same day.")
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
	// should not end a sync that has already created observations (P-062).
	invalidSkips                                           int
	totalRecords, createdObservations, updatedObservations int
	uploadedPhotos, uploadedSounds                         int
	// pendingMedia counts the media assets a --dryrun would have uploaded.
	// A Macaulay Library asset ID doesn't say whether it's a photo or a sound,
	// and a dry run doesn't download it to find out, so these can't be split
	// into uploadedPhotos and uploadedSounds.
	pendingMedia int
	errors       int
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
	// effect, so those lines are conditional (P-055).
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
		// (P-060, T-007).
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
	return lines
}

func birdsync(eBirdCSVFilename string, ebirdClient ebirdClient, inatUserID string, inatClient inatClient) stats {
	results := inatClient.DownloadObservations(inatUserID, after.Time(), before.Time(),
		"description", "observed_on", "photos.all", "sounds.all", "taxon.all", "ofvs.all")

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
		}
	}
	debugf("Previously synced %d observations\n", len(previouslySynced))

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
			// leave the sync half done (P-062).
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
			// Upload the media
			for _, id := range assetIDs.ids {
				obs.Description += "Macaulay Library Asset: " + mlAssetURL(id) + "\n"
				if dryRun {
					log.Printf("DRYRUN: Download ML Asset %s and upload to iNaturalist", id)
					s.pendingMedia++
				} else {
					filename, isPhoto, err := ebirdClient.DownloadMLAsset(id)
					if err != nil {
						log.Printf("Couldn't download ML asset %s from eBird: %v", id, err)
						s.errors++
						continue
					}
					err = inatClient.UploadMedia(filename, isPhoto, id, obs.UUID.String())
					// The download is a temp file that belongs to us now, so
					// remove it whether or not the upload worked. Syncing an
					// account with thousands of assets used to leave one file
					// per asset behind (T-023). Removing here rather than in a
					// defer keeps at most one asset on disk at a time.
					if rmErr := os.Remove(filename); rmErr != nil {
						debugf("Couldn't remove temp file %s: %v", filename, rmErr)
					}
					if err != nil {
						log.Printf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
						s.errors++
						continue
					}
					if isPhoto {
						s.uploadedPhotos++
					} else {
						s.uploadedSounds++
					}
				}
			}
			// Update the description
			if dryRun {
				log.Printf("DRYRUN: Updating observation %s with %d added media assets\n",
					obs.URLWithSpecies(), assetIDs.Len())
				prettyPrintln(obs)
			} else {
				err = inatClient.UpdateObservation(obs)
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
			if checkFuzzy(rec.CommonName) || checkFuzzy(rec.ScientificName) {
				continue
			}
		}

		assetIDs := eBirdMLAssets(rec.MLCatalogNumbers)
		// Skip records without media assets if --verifiable is set. This is
		// checked before the record is parsed any further, so that a record is
		// counted against the rule that actually skipped it (P-026).
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
			debugf("Syncing eBird observation %s to iNaturalist (%d media assets)\n",
				key, assetIDs.Len())
			err = inatClient.CreateObservation(obs)
			if err != nil {
				log.Fatalf("CreateObservation: %v", err)
			}
		}
		s.createdObservations++
		addMedia(obs.UUID, obs.Description, assetIDs)
	}
	return s
}
