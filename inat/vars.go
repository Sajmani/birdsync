package inat

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// GetUserID returns the user's iNaturalist user ID.
//
// First, it checks the INAT_USER_ID environment variable.
//
// Otheriwse, it prompts the user to enter their user ID
// from the URL https://www.inaturalist.org/home.
func GetUserID() string {
	userID := os.Getenv("INAT_USER_ID")
	if userID != "" {
		return userID
	}
	for {
		fmt.Print(`
Your iNaturalist user ID allows this tool to act on your behalf.
Copy your iNaturalist user ID from the top of https://www.inaturalist.org/home
(next to your profile picture) and paste it below.
To skip this step in the future, set the INAT_USER_ID environment variable to
your user ID.
`)
		var userID string
		_, err := fmt.Scan(&userID)
		if err != nil {
			fmt.Println("Didn't get your user ID: ", err)
			continue
		}
		if userID == "" {
			fmt.Println("Empty user ID")
			continue
		}
		return userID
	}
}

// GetAPIToken returns the user's iNaturalist API token.
//
// First, it checks the INAT_API_TOKEN environment variable,
// which should contain the API token string without the surrounding JSON framing
// (TOKEN without the surrounding {"api_token":"TOKEN"}).
//
// Otherwise, it prompts the user to paste in their token (with JSON framing)
// from the URL https://www.inaturalist.org/users/api_token.
func GetAPIToken() string {
	apiToken := os.Getenv("INAT_API_TOKEN")
	if apiToken != "" {
		return apiToken
	}
	for {
		fmt.Print(`
Your iNaturalist API token allows this tool to act on your behalf.
The API token needs to be refreshed every 24 hours.
The token is a long string of characters starting and ending with curly braces,
like this: {"api_token":"..."}. Copy your current iNaturalist API token from
https://www.inaturalist.org/users/api_token and paste it below.
To skip this step in the future, set the INAT_API_TOKEN environment variable to
TOKEN, without the surrounding {"api_token":"TOKEN"}, but remember to refresh
your token every 24 hours!
`)
		var tokenJSON string
		_, err := fmt.Scan(&tokenJSON)
		if err != nil {
			fmt.Println("Didn't get your API token: ", err)
			continue
		}
		m := make(map[string]string)
		err = json.NewDecoder(strings.NewReader(tokenJSON)).Decode(&m)
		if err != nil {
			fmt.Println("Bad API token: ", err)
			continue
		}
		apiToken = m["api_token"]
		if apiToken == "" {
			fmt.Println("Empty API token")
			continue
		}
		return apiToken
	}
}
