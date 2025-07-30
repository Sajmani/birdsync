package main

import (
	"log"
	"os"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
)

const UserAgent = "birdsync-testing/0.1"

func usage() {
	log.Print(`poke [create|image]
poke image <Macaulay Library Asset ID> <iNat Observation UUID>
`)
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	c := inat.NewClient(inat.GetAPIToken(), UserAgent)
	switch os.Args[1] {
	case "create":
		c.CreateObservation(inat.TestObservation())
	case "image":
		if len(os.Args) < 4 {
			usage()
		}
		mlAssetID := os.Args[2]
		obsUUID := os.Args[3]
		filename, err := ebird.DownloadMLAsset(mlAssetID)
		if err != nil {
			log.Fatal(err)
		}
		err = c.UploadImage(filename, mlAssetID, obsUUID)
		if err != nil {
			log.Fatal(err)
		}
	default:
		usage()
	}
}
