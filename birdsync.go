// Birdsync syncs eBird observations and photos to iNaturalist.
//
// See README.md for detailed documentation.
package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
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
	dryRun     bool
	verifiable bool
	before     dateTimeFlag
	after      dateTimeFlag
)

func init() {
	flag.BoolVar(&dryRun, "dryrun", false,
		"Don't actually sync any observations, just log what birdsync would do")
	flag.BoolVar(&verifiable, "verifiable", false,
		"Sync only observations that include Macaulay Catalog Numbers (photos or sound)")
	flag.Var(&before, "before",
		"Sync only observations observed before the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).")
	flag.Var(&after, "after",
		"Sync only observations observed after the provided DateTime (2006-01-02 15:04:05). The time can be omitted (2006-01-02).")
}

func parseEBirdDateTime(d, t string) (time.Time, error) {
	if t == "" {
		return time.Parse("2006-01-02", d)
	}
	return time.Parse("2006-01-02 03:04 PM", d+" "+t)
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
	eBirdCSVFile, err := os.Open(eBirdCSVFilename)
	if err != nil {
		log.Fatalf("Error opening %s: %v", eBirdCSVFilename, err)
	}
	defer eBirdCSVFile.Close()

	inatUserID := inat.GetUserID()
	inatAPIToken := inat.GetAPIToken()
	client := inat.NewClient(inatAPIToken, UserAgent)

	results := inat.DownloadObservations(inatUserID, after.Time(), before.Time(),
		"description", "taxon.name", "ofvs.all")

	previouslySynced := map[ebird.ObservationID]inat.Result{}
	for _, r := range results {
		key := ebird.ObservationID{
			SubmissionID:   r.ObservationFieldValue(inat.EBirdField),
			ScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if !key.Valid() {
			continue // not a synced observation, skip this one
		}
		previouslySynced[key] = r
	}
	debugf("Previously synced %d observations\n", len(previouslySynced))

	log.Printf("Reading eBird observations from %s", eBirdCSVFilename)
	r := csv.NewReader(eBirdCSVFile)
	// iNaturalist's CSV export returns a variable number of fields per record,
	// so disable this check. This means we need to explicitly check len(rec)
	// before accessing fields that might not be there.
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records from %s: %v", eBirdCSVFilename, err)
	}
	if len(recs) < 1 {
		log.Fatalf("No records found in %s", eBirdCSVFilename)
	}
	field := make(map[string]int)
	for i, f := range recs[0] {
		field[f] = i
	}
	recs = recs[1:]
	debugf("Read %d eBird observations", len(recs))
	var stats struct {
		afterSkips, beforeSkips, verifiableSkips, previouslySkips int
		createdObservations, uploadedPhotos                       int
	}
	for i, rec := range recs {
		line := i + 2 // header was line 1

		// Skip records that were not observed between --after and --before.
		d, t := rec[field["Date"]], rec[field["Time"]]
		observed, err := parseEBirdDateTime(d, t)
		if err != nil {
			log.Fatalf("line %d: could not parse Date %q Time %q: %v", line, d, t, err)
		}
		if !after.Time().IsZero() && observed.Before(after.Time()) {
			debugf("line %d: SKIPPING record observed on %s (before --after=%s)",
				line, observed, after.Time())
			stats.afterSkips++
			continue
		}
		if !before.Time().IsZero() && observed.After(before.Time()) {
			debugf("line %d: SKIPPING record observed on %s (after --before=%s)",
				line, observed, before.Time())
			stats.beforeSkips++
			continue
		}

		key := ebird.ObservationID{
			SubmissionID:   rec[field["Submission ID"]],
			ScientificName: rec[field["Scientific Name"]],
		}
		if r, ok := previouslySynced[key]; ok {
			debugf("line %d: Already synced %s to iNaturalist as http://inaturalist.org/observations/%s\n",
				line, key, r.UUID)
			stats.previouslySkips++
			continue
		}
		parseFloat64 := func(key string) float64 {
			s := rec[field[key]]
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("line %d: invalid float64 for %s: %q: %v", line, key, s, err)
			}
			return f
		}
		keyField := func(id int, key string) inat.ObservationFieldValue {
			return inat.ObservationFieldValue{
				ObservationFieldID: id,
				Value:              rec[field[key]],
			}
		}
		obs := inat.Observation{
			UUID:             uuid.New(),
			CaptiveFlag:      false, // eBird checklists should only include wild birds
			Latitude:         parseFloat64("Latitude"),
			Longitude:        parseFloat64("Longitude"),
			LocationIsExact:  false,
			SpeciesGuess:     rec[field["Scientific Name"]],
			ObservedOnString: rec[field["Date"]] + " " + rec[field["Time"]],
			ObservationFieldValuesAttributes: []inat.ObservationFieldValue{
				keyField(inat.CountField, "Count"),
				keyField(inat.CommonNameField, "Common Name"),
				keyField(inat.LocationField, "Location"),
				keyField(inat.CountyField, "County"),
				keyField(inat.StateOrProvinceField, "State/Province"),
				keyField(inat.NumObserversField, "Number of Observers"),
				// EBirdField and EBirdScientificNameField are used to match iNaturalist observations
				// to the corresponding eBird checklist and species entry. We cannot rely on the taxon
				// in the iNaturalist observation because it may be changed after upload.
				keyField(inat.EBirdField, "Submission ID"),
				keyField(inat.EBirdScientificNameField, "Scientific Name"),
			},
		}
		obs.Description = "Observation created using github.com/Sajmani/birdsync \n"
		if field["Observation Details"] < len(rec) && len(rec[field["Observation Details"]]) > 0 {
			obs.Description += "eBird observation details:\n" +
				rec[field["Observation Details"]] + "\n"
		}
		obs.Description += "Checklist: https://ebird.org/checklist/" + rec[field["Submission ID"]] + "\n"
		obs.Description += "Protocol: " + rec[field["Protocol"]] + "\n"
		if field["Checklist Comments"] < len(rec) && len(rec[field["Checklist Comments"]]) > 0 {
			obs.Description += "eBird checklist comments:\n" +
				rec[field["Checklist Comments"]] + "\n"
		}
		var photoIDs []string
		if field["ML Catalog Numbers"] < len(rec) && len(rec[field["ML Catalog Numbers"]]) > 0 {
			photoIDs = strings.Split(rec[field["ML Catalog Numbers"]], " ")
			for _, id := range photoIDs {
				obs.Description += "Macaulay Library Asset: https://macaulaylibrary.org/asset/" + id + "\n"
			}
		}
		if verifiable && len(photoIDs) == 0 {
			debugf("line %d: SKIPPING record that has no photos or sounds (--verifiable=true)", line)
			stats.verifiableSkips++
			continue
		}
		if dryRun {
			log.Printf("DRYRUN: Syncing eBird observation %s to iNaturalist (%d photos)\n",
				key, len(photoIDs))
		} else {
			debugf("Syncing eBird observation %s to iNaturalist (%d photos)\n",
				key, len(photoIDs))
			err = client.CreateObservation(obs)
			if err != nil {
				log.Fatalf("CreateObservation: %v", err)
			}
		}
		stats.createdObservations++

		for _, id := range photoIDs {
			if dryRun {
				log.Printf("DRYRUN: Download ML Asset %s and upload to iNaturalist", id)
			} else {
				filename, err := ebird.DownloadMLAsset(id)
				if err != nil {
					log.Fatalf("Couldn't download ML asset %s from eBird: %v", id, err)
				}
				err = client.UploadImage(filename, id, obs.UUID.String())
				if err != nil {
					log.Fatalf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
				}
			}
			stats.uploadedPhotos++
		}
	}
	log.Printf("Finished processing %d eBird observations", len(recs))
	log.Printf("Skipped %d previously uploaded by birdsync", stats.previouslySkips)
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
}
