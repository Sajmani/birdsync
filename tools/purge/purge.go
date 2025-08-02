// Purge is a tool for deleting casual observations that have neither photos nor sounds.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-purge/0.1"

const debug = false

func main() {
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"photos", "sounds", "quality_grade", "ofvs.all")

	for _, r := range results {
		key := ebird.ObservationID{
			SubmissionID:   r.ObservationFieldValue(inat.EBirdField),
			ScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if !key.Valid() {
			fmt.Println("SKIP", r.UUID)
			continue // not a synced observation, skip this one
		}
		fmt.Println(r.UUID, r.QualityGrade, len(r.Photos), len(r.Sounds))
		if len(r.Photos) == 0 && len(r.Sounds) == 0 {
			fmt.Println("DELETE", r.UUID)
		} else {
			fmt.Println("KEEP", r.UUID)
			continue
		}
		if !debug {
			err := client.DeleteObservation(r.UUID)
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}
