package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

const UserAgent = "birdsync-testing/0.1"

func main() {
	apiToken := getAPIToken()
	createTestObservation(apiToken)
}

func createTestObservation(apiToken string) {
	create := inat.CreateObservation{
		Observation: inat.Observation{
			UUID:         uuid.New(),
			CaptiveFlag:  true, // casual observation for testing
			Description:  "Testing github.com/Sajmani/birdsync tools",
			SpeciesGuess: "Homo Sapiens",
		},
	}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(create)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://api.inaturalist.org/v2/observations", buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Authorization", apiToken)
	fmt.Printf("\nREQUEST: %+v\n", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nRESPONSE: %+v\n", resp)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nBODY: " + string(b))
}

func prettyPrint(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func getAPIToken() string {
	apiToken := os.Getenv("INAT_API_TOKEN")
	if apiToken != "" {
		return apiToken
	}
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
		return apiToken
	}
}
