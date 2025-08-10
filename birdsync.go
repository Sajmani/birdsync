// Birdsync syncs eBird observations, photos, and sounds to iNaturalist.
//
// See README.md for detailed documentation.
package main

import (
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const UserAgent = "birdsync/0.1"

const debug = false

func debugf(format string, args ...any) {
	if debug {
		log.Printf(format, args...)
	}
}

type dateTimeFlag struct {
	t time.Time
}

func (f *dateTimeFlag) String() string {
	return f.t.Format(time.DateTime)
}

func (f *dateTimeFlag) Set(s string) error {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		t, err = time.Parse(time.DateOnly, s)
		if err != nil {
			return err
		}
	}
	f.t = t
	return nil
}

func (f *dateTimeFlag) Time() time.Time {
	return f.t
}

var (
	dryRun             bool
	verifiable         bool
	fuzzy              bool
	before             dateTimeFlag
	after              dateTimeFlag
	positionalAccuracy int
)

func init() {
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
	inatUserID := inat.GetUserID()
	inatAPIToken := inat.GetAPIToken()
	client := inat.NewClient(inatAPIToken, UserAgent)

	results := inat.DownloadObservations(inatUserID, after.Time(), before.Time(),
		"description", "observed_on", "taxon.all", "ofvs.all")

	previouslySynced := map[ebird.ObservationID]inat.Result{}
	type fuzzyKey struct {
		observedDate string // 2006-01-02
		commonName   string
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
			key := fuzzyKey{
				observedDate: r.ObservedOn,
				commonName:   r.Taxon.PreferredCommonName,
			}
			fuzzyMatch[key] = append(fuzzyMatch[key], r.UUID.String())
			slices.Sort(fuzzyMatch[key])
			debugf("fuzzy match: add %s to %+v", r.UUID, key)
		}
	}
	debugf("Previously synced %d observations\n", len(previouslySynced))

	log.Printf("Reading eBird observations from %s", eBirdCSVFilename)
	records, err := ebird.Records(eBirdCSVFilename)
	if err != nil {
		log.Fatal(err)
	}
	var stats struct {
		afterSkips, beforeSkips, verifiableSkips, previouslySkips, fuzzySkips int
		totalRecords, createdObservations, uploadedPhotos, uploadedSounds     int
	}
	for rec := range records {
		stats.totalRecords++
		observed, err := rec.Observed()
		if err != nil {
			log.Fatalf("line %d: bad date/time: %v", rec.Line, err)
		}
		// Skip records that were not observed between --after and --before.
		if !after.Time().IsZero() && observed.Before(after.Time()) {
			debugf("line %d: SKIPPING record observed on %s (before --after=%s)",
				rec.Line, observed, after.Time())
			stats.afterSkips++
			continue
		}
		if !before.Time().IsZero() && observed.After(before.Time()) {
			debugf("line %d: SKIPPING record observed on %s (after --before=%s)",
				rec.Line, observed, before.Time())
			stats.beforeSkips++
			continue
		}

		// Skip records that have previously been uploaded by birdsync.
		key := rec.ObservationID()
		if r, ok := previouslySynced[key]; ok {
			debugf("line %d: Already synced %s to iNaturalist as http://inaturalist.org/observations/%s\n",
				rec.Line, key, r.UUID)
			stats.previouslySkips++
			if debug {
				// Determine whether photos or sounds have changed.
				// TODO: sync the changes
				eSet := eBirdMLAssets(rec.MLCatalogNumbers)
				iSet := iNatMLAssets(r)
				if !maps.Equal(eSet, iSet) {
					fmt.Printf("line %d: %s different ML assets: eBird %+v; iNat %+v\n", rec.Line, key, eSet, iSet)
					fmt.Printf("added: %+v\n", mlAssetDiff(eSet, iSet))
					fmt.Printf("deleted: %+v\n", mlAssetDiff(iSet, eSet))
				}
			}
			continue
		}

		if fuzzy {
			// Skip records for the same bird and date as an existing non-birdsync observation.
			key := fuzzyKey{
				commonName:   rec.CommonName,
				observedDate: rec.Date,
			}
			debugf("line %d: fuzzy match: check %+v", rec.Line, key)
			if _, ok := fuzzyMatch[key]; ok {
				log.Printf("line %d: SKIPPING fuzzy match: observation for same bird and date: %+v", rec.Line, key)
				stats.fuzzySkips++
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
		obs.Description += "Checklist: https://ebird.org/checklist/" + rec.SubmissionID + "\n"
		obs.Description += "Protocol: " + rec.Protocol + "\n"
		if len(rec.ChecklistComments) > 0 {
			obs.Description += "eBird checklist comments:\n" +
				rec.ChecklistComments + "\n"
		}
		var assetIDs []string
		if len(rec.MLCatalogNumbers) > 0 {
			assetIDs = strings.Split(rec.MLCatalogNumbers, " ")
			for _, id := range assetIDs {
				obs.Description += "Macaulay Library Asset: https://macaulaylibrary.org/asset/" + id + "\n"
			}
		}
		// Skip records without media assets if --verifiable is set.
		if verifiable && len(assetIDs) == 0 {
			debugf("line %d: SKIPPING record that has no photos or sounds (--verifiable=true)", rec.Line)
			stats.verifiableSkips++
			continue
		}
		if dryRun {
			log.Printf("DRYRUN: Syncing eBird observation %s to iNaturalist (%d media assets)\n",
				key, len(assetIDs))
			pretty.Println(obs)
		} else {
			debugf("Syncing eBird observation %s to iNaturalist (%d media assets)\n",
				key, len(assetIDs))
			err = client.CreateObservation(obs)
			if err != nil {
				log.Fatalf("CreateObservation: %v", err)
			}
		}
		stats.createdObservations++

		for _, id := range assetIDs {
			if dryRun {
				log.Printf("DRYRUN: Download ML Asset %s and upload to iNaturalist", id)
				stats.uploadedPhotos++
			} else {
				filename, isPhoto, err := ebird.DownloadMLAsset(id)
				if err != nil {
					log.Fatalf("Couldn't download ML asset %s from eBird: %v", id, err)
				}
				err = client.UploadMedia(filename, isPhoto, id, obs.UUID.String())
				if err != nil {
					log.Fatalf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
				}
				if isPhoto {
					stats.uploadedPhotos++
				} else {
					stats.uploadedSounds++
				}
			}
		}
	}
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
	log.Printf("Uploaded %d photos to iNaturalist", stats.uploadedPhotos)
	log.Printf("Uploaded %d sounds to iNaturalist", stats.uploadedSounds)
}

type mlAssetSet map[string]bool

func eBirdMLAssets(mlAssets string) mlAssetSet {
	set := mlAssetSet{}
	for _, id := range strings.Split(mlAssets, " ") {
		set[strings.TrimSpace(id)] = true
	}
	return set
}

func iNatMLAssets(r inat.Result) mlAssetSet {
	set := mlAssetSet{}
	for _, line := range strings.Split(r.Description, "\n") {
		if i := strings.Index(line, "macaulaylibrary.org/asset/"); i >= 0 {
			id := line[i+len("macaulaylibrary.org/asset/"):]
			set[strings.TrimSpace(id)] = true
		}
	}
	return set
}

func mlAssetDiff(a, b mlAssetSet) mlAssetSet {
	diff := mlAssetSet{}
	for id := range a {
		if !b[id] {
			diff[id] = true
		}
	}
	return diff
}
