// Repair is a tool for updating the positional accuracy of observations.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-position/0.1"

const debug = true

func main() {
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"ofvs.all", "positional_accuracy")

	for _, r := range results {
		key := ebird.ObservationID{
			SubmissionID:   r.ObservationFieldValue(inat.EBirdField),
			ScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if !key.Valid() {
			fmt.Println("SKIP", r.UUID)
			continue // not a synced observation, skip this one
		}
		if r.PositionalAccuracy == ebird.PositionalAccuracy {
			continue // not changing
		}
		fmt.Println(r.UUID, "update", r.PositionalAccuracy, "to", ebird.PositionalAccuracy)
		if !debug {
			err := client.UpdateObservation(inat.Observation{
				UUID:               r.UUID,
				PositionalAccuracy: ebird.PositionalAccuracy,
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
