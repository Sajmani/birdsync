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

func (c *Client) CreateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(CreateObservation{
		Observation: obs,
	})
	if err != nil {
		return fmt.Errorf("json.Encode CreateObservation: %w", err)
	}
	req, err := http.NewRequest("POST", "https://api.inaturalist.org/v2/observations", buf)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	if debug {
		log.Printf("\nREQUEST: %+v\n", req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP POST: bad status: %s", resp.Status)
	}
	if debug {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading HTTP response: %w", err)
		}
		log.Println("\nBODY: " + string(b))
	}
	log.Printf("Uploaded as http://inaturalist.org/observations/%s\n", obs.UUID)
	return nil
}

func (c *Client) UpdateObservation(obs Observation) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(UpdateObservation{
		Observation: obs,
	})
	if err != nil {
		return fmt.Errorf("json.Encode UpdateObservation: %w", err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.inaturalist.org/v2/observations/%s", obs.UUID), buf)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", c.apiToken)

	if debug {
		log.Printf("\nREQUEST: %+v\n", req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP POST: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP POST: bad status: %s", resp.Status)
	}
	if debug {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading HTTP response: %w", err)
		}
		log.Println("\nBODY: " + string(b))
	}
	log.Printf("Updated http://inaturalist.org/observations/%s\n", obs.UUID)
	return nil
}

func (c *Client) UploadImage(filename string, mlAssetID string, obsUUID string) error {
	destFilename := "ML" + mlAssetID + path.Ext(filename)
	log.Println("Uploading image as", destFilename)
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

	if debug {
		log.Printf("\nREQUEST: %+v\n", req)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("bad status code: %s", resp.Status)
	}
	if debug {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("\nBODY: " + string(b))
	}
	// TODO: log the image URL
	return nil
}
