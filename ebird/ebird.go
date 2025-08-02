// Package ebird provides helper functions for working with eBird data.
package ebird

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

// PositionalAccuracy is the default positional accuracy in meters
// that we use for eBird observations. This is intended to serve as
// an approximation of the radius of a typical eBird hotspot.
const PositionalAccuracy = 500 // meters

// ObservationID identifies a unique eBird observation
// as a submission ID and eBird's scientific name. EBird's
// scientific names may differ from iNaturalist's taxa
// in various ways, notably for "slashes" and "spuhs".
type ObservationID struct {
	// Submission ID is the eBird checklist ID, including leading "S".
	// Example: "S193523301"
	SubmissionID string

	// ScientificName examples:
	// - "Struthio camelus"
	// - "Cairina moschata (Domestic type)"
	// - "Anas platyrhynchos x rubripes"
	// - "Aythya marila/affinis"
	// - "Melanitta sp."
	ScientificName string
}

// Valid returns whether this observation ID has all fields set.
func (o ObservationID) Valid() bool {
	return o.SubmissionID != "" && o.ScientificName != ""
}

func (o ObservationID) String() string {
	return fmt.Sprintf("%s[%s]", o.SubmissionID, o.ScientificName)
}

// DownloadMLAsset downloads the image with the provided ML asset ID
// (numbers only) and returns the local filename.
// This file is temporary and may be deleted at any time.
//
// TODO: support downloading sound assets, too.
func DownloadMLAsset(mlAssetID string) (string, error) {
	url := fmt.Sprintf("https://cdn.download.ams.birds.cornell.edu/api/v2/asset/%s/2400", mlAssetID)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close() // Close the response body when the function exits
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "birdsync")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy image data to file: %w", err)
	}

	// Re-open the file to detect content type
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to seek to beginning of temp file: %w", err)
	}

	buf := make([]byte, 512) // 512 bytes is the required size for DetectContentType
	n, err := tmpFile.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read from temp file for content type detection: %w", err)
	}
	buf = buf[:n]

	mimeType := http.DetectContentType(buf)
	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {
		return "", fmt.Errorf("failed to find file extension for mime type %s: %w", mimeType, err)
	}
	tmpFile.Close() // Close the file before renaming it.

	newPath := tmpFile.Name() + extensions[0]
	err = os.Rename(tmpFile.Name(), newPath)
	if err != nil {
		return "", fmt.Errorf("failed to rename file: %w", err)
	}
	return newPath, nil
}
