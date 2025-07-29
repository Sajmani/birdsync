package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sajmani/birdsync/inat"
)

func main() {
	inatUserID := os.Getenv("INAT_USER_ID")
	if inatUserID == "" {
		log.Fatal("INAT_USER_ID environment variable not set")
	}
	results := inat.DownloadObservations(inatUserID, "ofvs.id", "ofvs.name", "ofvs.value")

	fmt.Println("downloaded", len(results), "results")
}
