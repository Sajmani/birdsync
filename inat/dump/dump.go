package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sajmani/birdsync/inat"
)

func main() {
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
		break
	}
	fmt.Println("Got token:", apiToken)
}
