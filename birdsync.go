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
	totalRecords, createdObservations, updatedObservations                int
	uploadedPhotos, uploadedSounds                                        int
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
		client: inat.NewClient(inat.GetAPIToken(), UserAgent),
	}
	ebirdAPIClient := ebirdClientImpl{}

	stats := birdsync(eBirdCSVFilename, ebirdAPIClient, inat.GetUserID(), inatAPIClient)

	log.Printf("Finished processing %d eBird observations", stats.totalRecords)
	log.Printf("Skipped %d previously uploaded by birdsync", stats.previouslySkips)
	if fuzzy {
		log.Printf("Skipped %d eBird observations with --fuzzy matching", stats.fuzzySkips)
	}
	if !after.Time().IsZero() {
		log.Printf("Skipped %d eBird observations before --after", stats.afterSkips)
	}
	if !before.Time().IsZero() {
		log.Printf("Skipped %d eBird observations after --before", stats.beforeSkips)
	}
	if verifiable {
		log.Printf("Skipped %d unverifiable eBird observations", stats.verifiableSkips)
	}
	log.Printf("Created %d new iNaturalist observations", stats.createdObservations)
	log.Printf("Updated %d iNaturalist observations", stats.updatedObservations)
	log.Printf("Uploaded %d photos to iNaturalist", stats.uploadedPhotos)
	log.Printf("Uploaded %d sounds to iNaturalist", stats.uploadedSounds)
}

func birdsync(eBirdCSVFilename string, ebirdClient ebirdClient, inatUserID string, inatClient inatClient) stats {
	results := inatClient.DownloadObservations(inatUserID, after.Time(), before.Time(),
		"description", "observed_on", "photos", "sounds", "taxon.all", "ofvs.all")

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
				key := fuzzyKey{
					observedDate: r.ObservedOn,
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
			log.Fatalf("line %d: bad date/time: %v", rec.Line, err)
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
					s.uploadedPhotos++
				} else {
					filename, isPhoto, err := ebirdClient.DownloadMLAsset(id)
					if err != nil {
						log.Fatalf("Couldn't download ML asset %s from eBird: %v", id, err)
					}
					err = inatClient.UploadMedia(filename, isPhoto, id, obs.UUID.String())
					if err != nil {
						log.Fatalf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
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
			checkFuzzy := func(name string) bool {
				key := fuzzyKey{
					name:         name,
					observedDate: rec.Date,
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

		// Create the iNaturalist observation from the eBird record.
		floatField := func(line int, s string) float64 {
			if s == "" {
				return 0
			}
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("line %d: Invalid float64 %q: %v", line, s, err)
			}
			return f
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
			Latitude:           floatField(rec.Line, rec.Latitude),
			Longitude:          floatField(rec.Line, rec.Longitude),
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
		assetIDs := eBirdMLAssets(rec.MLCatalogNumbers)
		// Skip records without media assets if --verifiable is set.
		if verifiable && assetIDs.Len() == 0 {
			debugf("line %d: SKIPPING record that has no photos or sounds (--verifiable=true)", rec.Line)
			s.verifiableSkips++
			continue
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
