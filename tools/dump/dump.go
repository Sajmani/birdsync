// Dump is a tool for testing downloads of iNaturalist observations.
package main

import (
	"time"

	"github.com/Sajmani/birdsync/inat"
	"github.com/kr/pretty"
)

func main() {
	inatUserID := inat.GetUserID()
	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"description", "taxon.name", "ofvs.all")

	for _, r := range results {
		pretty.Println(r)
	}
}
