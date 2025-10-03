// Dump is a tool for testing downloads of iNaturalist observations.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Sajmani/birdsync/inat"
)

func prettyPrintln(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

const UserAgent = "birdsync-dump/0.1"

func main() {
	inatUserID := inat.GetUserID()
	apiToken := inat.GetAPIToken()
	client := inat.NewClient(inat.BaseURL, apiToken, UserAgent)

	results := client.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"description", "photos.all", "sounds.all", "taxon.name", "ofvs.all")

	for _, r := range results {
		prettyPrintln(r)
	}
}
