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
	apiToken, err := inat.GetAPIToken()
	if err != nil {
		log.Fatal(err)
	}
	inatUserID, err := inat.GetUserID(apiToken)
	if err != nil {
		log.Fatal(err)
	}
	client := inat.NewClient(inat.BaseURL, apiToken, UserAgent)

	results, err := client.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"description", "photos.all", "sounds.all", "taxon.name", "ofvs.all")
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range results {
		prettyPrintln(r)
	}
}
