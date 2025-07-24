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

func main() {
	inatUserID := os.Getenv("INAT_USER_ID")
	if inatUserID == "" {
		log.Fatal("INAT_USER_ID environment variable not set")
	}
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
		q.Set("fields", strings.Join([]string{
			"ofvs.id",
			"ofvs.name",
			"ofvs.value",
		}, ","))
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
			fmt.Println("no observations")
		}
		if totalResults == 0 {
			totalResults = observations.TotalResults
		}
		results = append(results, observations.Results...)
		if len(results) >= totalResults {
			break
		}
	}
	b, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func getToken() string {
	var _ inat.Observation
	var apiToken string
	for {
		fmt.Println("Your iNaturalist API token allows this tool to act on your behalf. The token needs to be refreshed every 24 hours.")
		fmt.Println(`The token is a long string of characters starting and ending with curly braces, like this: {"api_token":"...""}`)
		fmt.Println("Copy your current iNaturalist API token from https://www.inaturalist.org/users/api_token and paste it here:")
		var tokenJSON string
		_, err := fmt.Scan(&tokenJSON)
		if err != nil {
			fmt.Println("Didn't get your token: ", err)
			continue
		}
		m := make(map[string]string)
		err = json.NewDecoder(strings.NewReader(tokenJSON)).Decode(&m)
		if err != nil {
			fmt.Println("Bad token: ", err)
			continue
		}
		apiToken = m["api_token"]
		if apiToken == "" {
			fmt.Println("Empty token")
			continue
		}
		fmt.Println("Got token:", apiToken)
		return apiToken
	}
}
