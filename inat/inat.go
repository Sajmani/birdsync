// Package inat provides types and helper functions for the iNaturalist API.
package inat

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DownloadObservations downloads and returns all observations for inatUserID.
// The dates d1 and d2 specify the start and end of the observation date range if nonzero.
// The fields list specifies which fields are populated in the results.
func (c *Client) DownloadObservations(inatUserID string, d1, d2 time.Time, fields ...string) []Result {
	const dateFormat = "2006-01-02"
	var d1str, d2str string
	if !d1.IsZero() {
		d1str = " after " + d1.Format(dateFormat)
	}
	if !d2.IsZero() {
		d2str = " before " + d2.Format(dateFormat)
	}
	log.Printf("Downloading observations for %s%s%s", inatUserID, d1str, d2str)

	// From https://www.inaturalist.org/pages/api+recommended+practices:
	// If using the API to fetch a lot of results, please use the highest supported per_page value.
	// For example you can get up to 200 observations in a single request,
	// which would be faster and more efficient than fetching the default 30 results at a time.
	const perPage = 200

	var results []Result
	var totalResults int
	for page := 1; ; page++ {
		u, err := url.Parse(c.baseURL + "/observations")
		if err != nil {
			log.Fatal(err)
		}
		q := u.Query()
		q.Set("user_id", inatUserID)
		q.Set("page", strconv.Itoa(page))
		q.Set("per_page", strconv.Itoa(perPage))
		if !d1.IsZero() {
			q.Set("d1", d1.Format(dateFormat))
		}
		if !d2.IsZero() {
			q.Set("d2", d2.Format(dateFormat))
		}
		if len(fields) > 0 {
			q.Set("fields", strings.Join(fields, ","))
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			log.Fatal(err)
		}
		body, err := c.roundTrip(req)
		if err != nil {
			log.Fatal(err)
		}

		var observations Observations
		err = json.Unmarshal([]byte(body), &observations)
		if err != nil {
			log.Fatal(err)
		}

		if observations.TotalResults == 0 {
			break
		}
		if totalResults == 0 { // first loop
			totalResults = observations.TotalResults
		}
		results = append(results, observations.Results...)
		log.Printf("Downloaded %d of %d observations", len(results), totalResults)
		if len(results) >= totalResults {
			break
		}
	}
	return results
}

func TestObservation() Observation {
	return Observation{
		UUID:         uuid.New(),
		CaptiveFlag:  true, // casual observation for testing
		Description:  "Testing github.com/Sajmani/birdsync tools",
		SpeciesGuess: "Homo Sapiens",
	}
}
