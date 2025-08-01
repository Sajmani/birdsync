// Dedupe is a tool for deleting duplicate observations.
package main

import (
	"cmp"
	"log"
	"slices"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-dedupe/0.1"

const debug = true

func main() {
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(apiToken, UserAgent)

	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"created_at", "identifications_count", "ofvs.all")

	log.Println("downloaded", len(results), "results")
	m := map[ebird.ObservationID][]inat.Result{}
	for _, r := range results {
		key := ebird.ObservationID{
			SubmissionID:   r.ObservationFieldValue(inat.EBirdField),
			ScientificName: r.ObservationFieldValue(inat.EBirdScientificNameField),
		}
		if !key.Valid() {
			continue // not a synced observation, skip this one
		}
		m[key] = append(m[key], r)
	}
	for key, rs := range m {
		if len(rs) == 1 {
			continue
		}
		slices.SortFunc(rs, func(a, b inat.Result) int {
			countCmp := cmp.Compare(a.IdentificationsCount, b.IdentificationsCount)
			if countCmp != 0 {
				return -countCmp // prefer more identifications
			}
			return cmp.Compare(a.CreatedAt, b.CreatedAt) // prefer earlier
		})
		log.Println(key, "keeping", rs[0].UUID)
		for _, r := range rs[1:] {
			log.Println(key, "deleting duplicate", r.UUID)
			if !debug {
				err := client.DeleteObservation(r.UUID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
