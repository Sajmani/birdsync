// Repair is a tool for fixing up observations
// that are missing the eBird scientific name field.
//
// Procedure:
//   - Create a list of scientific names for each eBird checklist.
//   - Download iNaturalist observations.
//   - Identify observations that have eBird checklist IDs but not eBird scientific names.
//   - If the observation's taxon matches a scientific name in its checklist,
//     set that observation's eBird scientific name to the taxon name.
//   - Otherwise use the taxonToEBirdName mapping in this file to set that observation's
//     eBird scientific name based on the mapped taxon name.
package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-repair/0.1"

func main() {
	if len(os.Args) != 2 {
		log.Println("usage: repair MyEBirdData.csv")
		os.Exit(1)
	}
	eBirdCSVFile := os.Args[1]
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

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

	checklistScientificNames := map[string]map[string]bool{}
	for _, rec := range recs {
		checklist := rec[field["Submission ID"]]
		scientificName := rec[field["Scientific Name"]]
		if checklistScientificNames[checklist] == nil {
			checklistScientificNames[checklist] = map[string]bool{}
		}
		checklistScientificNames[checklist][scientificName] = true
	}

	taxonToEBirdName := map[string]string{
		"Columba livia domestica": "Columba livia (Feral Pigeon)",
		"Larinae":                 "Larinae sp.",
	}

	log.Println("Downloading observations for", inatUserID)
	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"taxon.name", "ofvs.all")
	log.Println("Downloaded", len(results), "observations")
	for _, r := range results {
		ebirdChecklist := r.ObservationFieldValue(inat.EBirdField)
		ebirdScientificName := r.ObservationFieldValue(inat.EBirdScientificNameField)
		if ebirdChecklist == "" {
			// observation not created by birdsync
			continue
		}
		if ebirdScientificName != "" {
			// observation has its eBird scientific name set, all good
			continue
		}
		// This observation was created before we added EBirdScientificNameField.
		// Check whether the observation taxon name matches any in the checklist.
		if checklistScientificNames[ebirdChecklist][r.Taxon.Name] {
			log.Printf("Set %s eBird sci name to obs taxon name %s", r.UUID, r.Taxon.Name)
			err := client.UpdateObservation(inat.Observation{
				UUID: r.UUID,
				ObservationFieldValuesAttributes: []inat.ObservationFieldValue{{
					ObservationFieldID: inat.EBirdScientificNameField,
					Value:              r.Taxon.Name,
				}},
			})
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		mappedName := taxonToEBirdName[r.Taxon.Name]
		if checklistScientificNames[ebirdChecklist][mappedName] {
			log.Printf("Set %s eBird sci name to mapped name %s", r.UUID, mappedName)

			err := client.UpdateObservation(inat.Observation{
				UUID: r.UUID,
				ObservationFieldValuesAttributes: []inat.ObservationFieldValue{{
					ObservationFieldID: inat.EBirdScientificNameField,
					Value:              mappedName,
				}},
			})
			if err != nil {
				log.Fatal(err)
			}
			continue
		}
		log.Printf("DON'T KNOW HOW to set %s eBird sci name: taxon name %s, checklist %s", r.UUID, r.Taxon.Name, ebirdChecklist)
	}
}
