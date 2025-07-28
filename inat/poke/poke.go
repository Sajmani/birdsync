package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

const UserAgent = "birdsync-testing/0.1"

func main() {
	switch os.Args[1] {
	case "fetch":
		downloadImage(os.Args[2])
	case "create":
		createTestObservation(getAPIToken())
	}
}

// downloadImage downloads an image from a given URL and saves it to a specified file path.
func downloadImage(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close() // Close the response body when the function exits
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "birdsync")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy image data to file: %w", err)
	}

	// Re-open the file to detect content type
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to beginning of temp file: %w", err)
	}

	buf := make([]byte, 512) // 512 bytes is the required size for DetectContentType
	n, err := tmpFile.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read from temp file for content type detection: %w", err)
	}
	buf = buf[:n]

	mimeType := http.DetectContentType(buf)
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {
		return fmt.Errorf("failed to find file extension for mime type %s: %w", mimeType, err)
	}
	tmpFile.Close() // Close the file before renaming it.

	newPath := tmpFile.Name() + extensions[0]
	err = os.Rename(tmpFile.Name(), newPath)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}
	fmt.Printf("Image downloaded and saved to %s\n", newPath)
	return nil
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

func getAPIToken() string {
	// If you use INAT_API_TOKEN, set it to the TOKEN without the surrounding {"api_token":"TOKEN"}
	apiToken := os.Getenv("INAT_API_TOKEN")
	if apiToken != "" {
		return apiToken
	}
	for {
		fmt.Println("Your iNaturalist API token allows this tool to act on your behalf. The token needs to be refreshed every 24 hours.")
		fmt.Println(`The token is a long string of characters starting and ending with curly braces, like this: {"api_token":"..."}`)
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
