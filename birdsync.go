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

func main() {
	if len(os.Args) != 2 {
		log.Println("usage: birdsync MyEBirdData.csv")
		os.Exit(1)
	}
	eBirdCSVFile := os.Args[1]
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

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
		elapsed := time.Since(last)
		last = time.Now()
		log.Println("Line", line, "of", len(recs), "-- estimate", elapsed*time.Duration(len(recs)-i), "remaining")
		key := ebirdSpecies{
			ebirdChecklist:      rec[field["Submission ID"]],
			ebirdScientificName: rec[field["Scientific Name"]],
		}
		if r, ok := previouslySynced[key]; ok {
			log.Printf("Already synced %s(%s) to iNaturalist: http://inaturalist.org/observations/%s\n",
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
		log.Printf("Syncing eBird observation %s(%s) to iNaturalist (%d photos)\n",
			key.ebirdChecklist, key.ebirdScientificName, len(photoIDs))
		err = client.CreateObservation(obs)
		if err != nil {
			log.Fatalf("CreateObservation: %v", err)
		}
		for _, id := range photoIDs {
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
