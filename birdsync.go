// Birdsync syncs eBird observations and photos to iNaturalist.
//
// As a prerequisite, you must download your data from eBird using
// https://ebird.org/downloadMyData — save the zip file and unzip it
// to get the MyEBirdData.csv file.
//
// Birdsync works as follows:
//  1. Download all iNaturalist observations for INAT_USER_ID into memory
//  2. Index these observations by (eBird checklist ID, species name)
//  3. Read eBird observations from the CSV file provided as a command line argument
//
// For each eBird observation:
//  4. Skip any eBird observations that have already been uploaded
//  5. Create the iNaturalist observation
//
// For each Macaulay Library ID for this eBird observation:
//  6. Download the image from the Macaulay Library
//  7. Upload the image to iNaturalist, associated with the new observation
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
		"Sync only observations created before the provided DateTime (2006-01-02 15:04:05)")
	flag.Var(&after, "after",
		"Sync only observations created after the provided DateTime (2006-01-02 15:04:05)")
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
	eBirdCSVFile := flag.Arg(0)
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

	// TODO: optimize the observation download using before and after filters
	log.Println("Downloading observations for", inatUserID)
	results := inat.DownloadObservations(inatUserID, "description", "taxon.name", "ofvs.all")
	log.Println("Downloaded", len(results), "observations")
	type ebirdSpecies struct {
		ebirdChecklist      string
		ebirdScientificName string
	}
	previouslySynced := map[ebirdSpecies]inat.Result{}
	for _, r := range results {
		key := ebirdSpecies{
			ebirdChecklist:      r.ObservationFieldValue(inat.EBirdField),
			ebirdScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if key.ebirdChecklist == "" || key.ebirdScientificName == "" {
			// not a synced observation, skip this one
			continue
		}
		previouslySynced[key] = r
	}
	log.Printf("Previously synced %d observations\n", len(previouslySynced))
	log.Println("Reading eBird observations from", eBirdCSVFile)
	f, err := os.Open(eBirdCSVFile)
	if err != nil {
		log.Fatalf("Error opening %s: %v", eBirdCSVFile, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// iNaturalist's CSV export returns a variable number of fields per record,
	// so disable this check. This means we need to explicitly check len(rec)
	// before accessing fields that might not be there.
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records from %s: %v", eBirdCSVFile, err)
	}
	if len(recs) < 1 {
		log.Fatalf("No records found in %s", eBirdCSVFile)
	}
	field := make(map[string]int)
	for i, f := range recs[0] {
		field[f] = i
	}
	recs = recs[1:]
	log.Println("Read", len(recs), "eBird observations")
	last := time.Now()
	for i, rec := range recs {
		line := i + 2 // header was line 1

		// Skip records that were not created between --after and --before.
		d, t := rec[field["Date"]], rec[field["Time"]]
		created, err := parseEBirdDateTime(d, t)
		if err != nil {
			log.Fatalf("line %d: could not parse Date %q Time %q: %v", line, d, t, err)
		}
		if !after.Time().IsZero() && created.Before(after.Time()) {
			log.Printf("line %d: SKIPPING record created on %s (before --after=%s)",
				line, created, after.Time())
			continue
		}
		if !before.Time().IsZero() && created.After(before.Time()) {
			log.Printf("line %d: SKIPPING record created on %s (after --before=%s)",
				line, created, before.Time())
			continue
		}

		elapsed := time.Since(last)
		last = time.Now()
		log.Println("Line", line, "of", len(recs), "-- estimate", elapsed*time.Duration(len(recs)-i), "remaining")
		key := ebirdSpecies{
			ebirdChecklist:      rec[field["Submission ID"]],
			ebirdScientificName: rec[field["Scientific Name"]],
		}
		if r, ok := previouslySynced[key]; ok {
			log.Printf("Already synced %s(%s) to iNaturalist as http://inaturalist.org/observations/%s\n",
				key.ebirdChecklist, key.ebirdScientificName, r.UUID)
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
			log.Printf("line %d: SKIPPING record that has no photos or sounds (--verifiable=true)", line)
			continue
		}
		if dryRun {
			log.Printf("DRYRUN: Syncing eBird observation %s(%s) to iNaturalist (%d photos)\n",
				key.ebirdChecklist, key.ebirdScientificName, len(photoIDs))
		} else {
			log.Printf("Syncing eBird observation %s(%s) to iNaturalist (%d photos)\n",
				key.ebirdChecklist, key.ebirdScientificName, len(photoIDs))
			err = client.CreateObservation(obs)
			if err != nil {
				log.Fatalf("CreateObservation: %v", err)
			}
		}
		for _, id := range photoIDs {
			if dryRun {
				log.Printf("DRYRUN: Download ML Asset %s and upload to iNaturalist", id)
				continue
			}
			filename, err := ebird.DownloadMLAsset(id)
			if err != nil {
				log.Fatalf("Couldn't download ML asset %s from eBird: %v", id, err)
			}
			err = client.UploadImage(filename, id, obs.UUID.String())
			if err != nil {
				log.Fatalf("Couldn't upload ML asset %s to iNaturalist: %v", id, err)
			}
		}
	}
}
