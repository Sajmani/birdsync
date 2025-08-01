// Dedupe is a tool for deleting duplicate observations.
package main

import (
	"log"
	"os"

	"github.com/Sajmani/birdsync/inat"
	"github.com/kr/pretty"
)

func main() {
	inatUserID := os.Getenv("INAT_USER_ID")
	if inatUserID == "" {
		log.Fatal("INAT_USER_ID environment variable not set")
	}
	results := inat.DownloadObservations(inatUserID, "description", "ofvs.all")

	log.Println("downloaded", len(results), "results")
	for _, r := range results {
		pretty.Println(r)
	}
}
