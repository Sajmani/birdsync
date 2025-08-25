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

const debug = false

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

func (c *Client) roundTrip(req *http.Request) (string, error) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	if debug {
		log.Printf("\nREQUEST: %+v\n", req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("making HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("%s: refresh your INAT_API_TOKEN", resp.Status)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad HTTP status: %s", resp.Status)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading HTTP response: %w", err)
	}
	body := string(b)
	if debug {
		log.Println("\nBODY: " + body)
	}
	return body, nil
}

func (c *Client) CreateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(CreateObservation{
		Observation: obs,
	})
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	req, err := http.NewRequest("POST", "https://api.inaturalist.org/v2/observations", buf)
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("CreateObservation: %w", err)
	}
	log.Printf("Created http://inaturalist.org/observations/%s\n", obs.UUID)
	return nil
}

func (c *Client) UpdateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(UpdateObservation{
		Observation: obs,
	})
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.inaturalist.org/v2/observations/%s", obs.UUID), buf)
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("UpdateObservation: %w", err)
	}
	log.Printf("Updated http://inaturalist.org/observations/%s\n", obs.UUID)
	return nil
}

func (c *Client) DeleteObservation(id uuid.UUID) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://api.inaturalist.org/v2/observations/%s", id), nil)
	if err != nil {
		return fmt.Errorf("DeleteObservation: %w", err)
	}
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("DeleteObservation: %w", err)
	}
	log.Printf("Deleted http://inaturalist.org/observations/%s\n", id)
	return nil
}

func (c *Client) UploadMedia(filename string, isPhoto bool, mlAssetID string, obsUUID string) error {
	destFilename := "ML" + mlAssetID + path.Ext(filename)
	var fieldName string
	var postURL string
	if isPhoto {
		fieldName = "observation_photo[observation_id]"
		postURL = "https://api.inaturalist.org/v2/observation_photos"
		log.Println("Uploading photo as", destFilename)
	} else {
		fieldName = "observation_sound[observation_id]"
		postURL = "https://api.inaturalist.org/v2/observation_sounds"
		log.Println("Uploading sound as", destFilename)
	}
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	fileWriter, err := writer.CreateFormFile("file", destFilename)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	defer f.Close()
	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	err = writer.WriteField(fieldName, obsUUID)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}

	req, err := http.NewRequest("POST", postURL, &requestBody)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	// Set the Content-Type header to the multipart writer's boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err = c.roundTrip(req)
	if err != nil {
		return fmt.Errorf("UploadMedia: %w", err)
	}
	// TODO: log the media URL from the response body
	return nil
}
