package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Sajmani/birdsync/inat"
)

// If using the API to fetch a lot of results, please use the highest supported per_page value.
// For example you can get up to 200 observations in a single request,
// which would be faster and more efficient than fetching the default 30 results at a time.
// From https://www.inaturalist.org/pages/api+recommended+practices
const perPage = 200

func main() {
	inatUserID := os.Getenv("INAT_USER_ID")
	if inatUserID == "" {
		log.Fatal("INAT_USER_ID environment variable not set")
	}
	results := downloadObservations(inatUserID, "ofvs.id", "ofvs.name", "ofvs.value")

	fmt.Println("downloaded", len(results), "results")
}

func downloadObservations(inatUserID string, fields ...string) []inat.Result {
	var results []inat.Result
	var totalResults int
	for page := 1; ; page++ {
		u, err := url.Parse("https://api.inaturalist.org/v2/observations")
		if err != nil {
			log.Fatal(err)
		}
		q := u.Query()
		q.Set("user_id", inatUserID)
		q.Set("page", strconv.Itoa(page))
		q.Set("per_page", strconv.Itoa(perPage))
		if len(fields) > 0 {
			q.Set("fields", strings.Join(fields, ","))
		}
		u.RawQuery = q.Encode()
		fmt.Println("Fetching page", page)
		resp, err := http.Get(u.String())
		if err != nil {
			log.Fatal(err)
		}
		var observations inat.Observations
		err = json.NewDecoder(resp.Body).Decode(&observations)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		if observations.TotalResults == 0 {
			break
		}
		if totalResults == 0 { // first loop
			totalResults = observations.TotalResults
		}
		results = append(results, observations.Results...)
		if len(results) >= totalResults {
			break
		}
	}
	return results
}

func prettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
