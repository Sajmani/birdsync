// Package inat provides types for messages defined in the iNaturalist API.
package inat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DownloadObservations downloads and returns all observations for inatUserID.
// The fields list specifies which fields are populated in the results.
func DownloadObservations(inatUserID string, fields ...string) []Result {
	// From https://www.inaturalist.org/pages/api+recommended+practices:
	// If using the API to fetch a lot of results, please use the highest supported per_page value.
	// For example you can get up to 200 observations in a single request,
	// which would be faster and more efficient than fetching the default 30 results at a time.
	const perPage = 200

	var results []Result
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
		var observations Observations
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
