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

func main() {
	inatUserID := inat.GetUserID()
	results := inat.DownloadObservations(inatUserID, time.Time{}, time.Time{},
		"description", "taxon.name", "ofvs.all")

	for _, r := range results {
		prettyPrintln(r)
	}
}
