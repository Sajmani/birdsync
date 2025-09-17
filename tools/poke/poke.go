// Poke is a tool for testing creating iNaturalist observations and photos.
package main

import (
	"log"
	"os"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-poke/0.1"

func usage() {
	log.Print(`poke [create|image]
poke create
poke image <Macaulay Library Asset ID> <iNat Observation UUID>
`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	c := inat.NewClient(inat.BaseURL, inat.GetAPIToken(), UserAgent)
	switch os.Args[1] {
	case "create":
		c.CreateObservation(inat.TestObservation())
	case "image":
		if len(os.Args) < 4 {
			usage()
		}
		mlAssetID := os.Args[2]
		obsUUID := os.Args[3]
		filename, isPhoto, err := ebird.DownloadMLAsset(mlAssetID)
		if err != nil {
			log.Fatal(err)
		}
		err = c.UploadMedia(filename, isPhoto, mlAssetID, obsUUID)
		if err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}
