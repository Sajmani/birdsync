package inat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/google/uuid"
)

type Client struct {
	apiToken  string
	userAgent string
}

func NewClient(apiToken, userAgent string) *Client {
	return &Client{
		apiToken:  apiToken,
		userAgent: userAgent,
	}
}

func (c *Client) UploadImage(filename string, mlAssetID string, obsUUID string) error {
	destFilename := "ML" + mlAssetID + path.Ext(filename)
	fmt.Println("uploading image as", destFilename)
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	fileWriter, err := writer.CreateFormFile("file", destFilename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return err
	}
	err = writer.WriteField("observation_photo[observation_id]", obsUUID)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.inaturalist.org/v2/observation_photos", &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	// Set the Content-Type header to the multipart writer's boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	//fmt.Printf("\nREQUEST: %+v\n", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status code: %s", resp.Status)
	}
	//fmt.Printf("\nRESPONSE: %+v\n", resp)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nBODY: " + string(b))

	return nil
}

func (c *Client) CreateTestObservation() {
	create := CreateObservation{
		Observation: Observation{
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
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	//fmt.Printf("\nREQUEST: %+v\n", req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status code: %s", resp.Status)
	}
	//fmt.Printf("\nRESPONSE: %+v\n", resp)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nBODY: " + string(b))
}
