package inat

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// ParseAPIToken extracts the bare iNaturalist API token from raw.
// It accepts either a bare token or the {"api_token":"TOKEN"} JSON object
// served by https://www.inaturalist.org/users/api_token.
//
// Error messages deliberately omit raw so credentials cannot leak into logs.
func ParseAPIToken(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", errors.New("empty API token")
	}
	if strings.HasPrefix(s, "{") {
		var obj struct {
			APIToken string `json:"api_token"`
		}
		if err := json.Unmarshal([]byte(s), &obj); err != nil {
			return "", errors.New(`invalid API token JSON; expected {"api_token":"..."}`)
		}
		tok := strings.TrimSpace(obj.APIToken)
		if tok == "" {
			return "", errors.New("empty api_token field in JSON")
		}
		if strings.ContainsAny(tok, " \t\r\n") {
			return "", errors.New("invalid API token: unexpected whitespace")
		}
		return tok, nil
	}
	if strings.ContainsAny(s, " \t\r\n") {
		return "", errors.New("invalid API token: unexpected whitespace")
	}
	return s, nil
}

// UserIDFromToken extracts the user_id claim from an iNaturalist JWT without
// contacting iNaturalist. It returns ("", false) if token is not a JWT or does
// not contain a valid user_id claim.
func UserIDFromToken(token string) (string, bool) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", false
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(parts[1], "="))
	if err != nil {
		return "", false
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.UseNumber()
	var claims struct {
		UserID any `json:"user_id"`
	}
	if err := dec.Decode(&claims); err != nil {
		return "", false
	}
	switch v := claims.UserID.(type) {
	case json.Number:
		id, err := v.Int64()
		if err != nil || id <= 0 {
			return "", false
		}
		return strconv.FormatInt(id, 10), true
	case string:
		s := strings.TrimSpace(v)
		if s == "" || strings.ContainsAny(s, " \t\r\n") {
			return "", false
		}
		return s, true
	default:
		return "", false
	}
}

// GetUserID returns the user's iNaturalist user ID.
//
// First, it checks the INAT_USER_ID environment variable.
// Next, it reads the user_id claim directly from apiToken if it is a JWT.
// Otherwise, it prompts the user to enter their user ID from
// https://www.inaturalist.org/home.
func GetUserID(apiToken string) (string, error) {
	return getUserID(os.Getenv("INAT_USER_ID"), apiToken, os.Stdin)
}

func getUserID(envUserID, apiToken string, in io.Reader) (string, error) {
	if userID := strings.TrimSpace(envUserID); userID != "" {
		return userID, nil
	}
	if userID, ok := UserIDFromToken(apiToken); ok {
		return userID, nil
	}
	scanner := bufio.NewScanner(in)
	for {
		log.Print(`Birdsync needs your iNaturalist user ID.
Your iNaturalist user ID allows this tool to act on your behalf.
Copy your iNaturalist user ID from the top of https://www.inaturalist.org/home
(next to your profile picture) and paste it below.
To skip this step in the future, set the INAT_USER_ID environment variable to
your user ID.
`)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("reading user ID: %w", err)
			}
			return "", errors.New("no user ID provided")
		}
		userID := strings.TrimSpace(scanner.Text())
		if userID == "" {
			log.Println("Empty user ID")
			continue
		}
		if strings.ContainsAny(userID, " \t\r\n") {
			log.Println("User ID must not contain whitespace")
			continue
		}
		return userID, nil
	}
}

// GetAPIToken returns the user's iNaturalist API token.
//
// First, it checks the INAT_API_TOKEN environment variable.
// Otherwise, it prompts the user to paste their token from
// https://www.inaturalist.org/users/api_token.
// Both bare tokens and {"api_token":"TOKEN"} JSON objects are accepted
// in either place.
func GetAPIToken() (string, error) {
	return getAPIToken(os.Getenv("INAT_API_TOKEN"), os.Stdin)
}

func getAPIToken(envToken string, in io.Reader) (string, error) {
	if envToken != "" {
		tok, err := ParseAPIToken(envToken)
		if err != nil {
			return "", fmt.Errorf("INAT_API_TOKEN: %w", err)
		}
		return tok, nil
	}
	scanner := bufio.NewScanner(in)
	for {
		log.Print(`Birdsync needs your iNaturalist API token.
Your iNaturalist API token allows this tool to act on your behalf.
The API token needs to be refreshed every 24 hours.
Copy your current iNaturalist API token from
https://www.inaturalist.org/users/api_token and paste it below
(either the full {"api_token":"..."} text or just the token).
To skip this step in the future, set the INAT_API_TOKEN environment variable,
but remember to refresh your token every 24 hours!
`)
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("reading API token: %w", err)
			}
			return "", errors.New("no API token provided")
		}
		apiToken, err := ParseAPIToken(scanner.Text())
		if err != nil {
			log.Printf("Bad API token: %v", err)
			continue
		}
		return apiToken, nil
	}
}
