package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
	"github.com/kr/pretty"
)

const (
	userAgent = "github-com-Sajmani-birdsync/0.1"

	// iNaturalist observation fields. Look up IDs using:
	// https://www.inaturalist.org/observation_fields?order=asc&order_by=created_at
	countField           = 1   // https://www.inaturalist.org/observation_fields/1
	locationField        = 157 // https://www.inaturalist.org/observation_fields/157
	countyField          = 245
	commonNameField      = 256
	distanceField        = 396
	protocolField        = 1285
	numObserversField    = 2527
	eBirdField           = 6033
	stateOrProvinceField = 7739
)

func main() {
	eBirdCSVFile := os.Getenv("EBIRD_CSV_FILE")
	if eBirdCSVFile == "" {
		log.Fatal("EBIRD_CSV_FILE environment variable not set")
	}
	observations := eBirdExportToiNatObservations(eBirdCSVFile)
	if len(observations) > 0 {
		pretty.Println("%+v\n", observations[rand.Intn(len(observations))])
	} else {
		fmt.Println("no observations in", eBirdCSVFile)
	}
}

func eBirdExportToiNatObservations(exportFile string) (observations []inat.Observation) {
	f, err := os.Open(exportFile)
	if err != nil {
		log.Fatalf("Error opening %s: %v", exportFile, err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	// iNaturalist's CSV export returns a variable number of fields per record,
	// so disable this check. This means we need to explicitly check len(rec)
	// before accessing fields that might not be there.
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV records from %s: %v", exportFile, err)
	}
	if len(recs) < 1 {
		log.Fatalf("No records found in %s", exportFile)
	}
	field := make(map[string]int)
	for i, f := range recs[0] {
		field[f] = i
	}
	recs = recs[1:]
	for i, rec := range recs {
		line := i + 2 // header was line 1
		parseFloat64 := func(key string) float64 {
			s := rec[field[key]]
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatalf("line %d: invalid float64 for %s: %q: %v", line, key, s, err)
			}
			return f
		}
		stringField := func(id int, val string) inat.ObservationFieldValue {
			return inat.ObservationFieldValue{
				ObservationFieldID: id,
				Value:              val,
			}
		}
		keyField := func(id int, key string) inat.ObservationFieldValue {
			return stringField(id, rec[field[key]])
		}
		obs := inat.Observation{
			UUID:             uuid.New(),
			CaptiveFlag:      false, // eBird checklists should only include wild birds
			Latitude:         parseFloat64("Latitude"),
			Longitude:        parseFloat64("Longitude"),
			LocationIsExact:  false,
			SpeciesGuess:     rec[field["Scientific Name"]],
			ObservedOnString: rec[field["Date"]] + " " + rec[field["Time"]],
			ObservationFieldValuesAttributes: []inat.ObservationFieldValue{
				keyField(countField, "Count"),
				keyField(commonNameField, "Common Name"),
				keyField(locationField, "Location"),
				keyField(countyField, "County"),
				keyField(stateOrProvinceField, "State/Province"),
				keyField(protocolField, "Protocol"),
				keyField(numObserversField, "Number of Observers"),
				stringField(eBirdField,
					"https://ebird.org/checklist/"+rec[field["Submission ID"]]),
			},
		}
		if field["Observation Details"] < len(rec) && rec[field["Observation Details"]] != "" {
			obs.Description += "eBird observation details:|n" +
				rec[field["Observation Details"]] + "\n"
		}
		if field["Checklist Comments"] < len(rec) && rec[field["Checklist Comments"]] != "" {
			obs.Description += "eBird checklist comments:\n" +
				rec[field["Checklist Comments"]] + "\n"
		}
		observations = append(observations, obs)
	}
	return
}
