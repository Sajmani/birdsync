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

// DownloadMLAsset downloads the photo or sound with the provided ML asset ID
// (numbers only) and returns the local filename and whether it's a photo.
// This file is temporary and may be deleted at any time.
//
// Since the ML asset ID doesn't indicate whether this is a photo or sound file,
// we try downloading the photo file first, and if it's not there,
// we try downloading the sound file.
func DownloadMLAsset(mlAssetID string) (string, bool, error) {
	// Try fetching this ML asset as a photo
	url := fmt.Sprintf("https://cdn.download.ams.birds.cornell.edu/api/v2/asset/%s/2400", mlAssetID)
	resp, err := http.Get(url)
	if err != nil {
		return "", false, fmt.Errorf("DownloadMLAsset(%s): %s: %w", mlAssetID, url, err)
	}
	defer resp.Body.Close()
	isPhoto := resp.StatusCode == http.StatusOK
	if resp.StatusCode == http.StatusNotFound {
		// Photo not found; try fetching it as a sound
		url = fmt.Sprintf("https://cdn.download.ams.birds.cornell.edu/api/v2/asset/%s/mp3", mlAssetID)
		resp, err = http.Get(url)
		if err != nil {
			return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): %s: %w", mlAssetID, url, err)
		}
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): %s: %s", mlAssetID, url, resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "birdsync")
	if err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): CreateTemp: %w", mlAssetID, err)
	}
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to copy asset data to file: %w", mlAssetID, err)
	}

	ext := ".mp3"
	if isPhoto {
		// For photos only: re-open the file to detect content type
		_, err = tmpFile.Seek(0, 0)
		if err != nil {
			return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to seek to beginning of temp file: %w", mlAssetID, err)
		}

		buf := make([]byte, 512) // 512 bytes is the required size for DetectContentType
		n, err := tmpFile.Read(buf)
		if err != nil && err != io.EOF {
			return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to read from temp file for content type detection: %w", mlAssetID, err)
		}
		buf = buf[:n]

		mimeType := http.DetectContentType(buf)
		extensions, err := mime.ExtensionsByType(mimeType)
		if err != nil || len(extensions) == 0 {
			return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to find file extension for mime type %s: %w", mlAssetID, mimeType, err)
		}
		tmpFile.Close() // Close the file before renaming it.
		ext = extensions[0]
	}

	newPath := tmpFile.Name() + ext
	err = os.Rename(tmpFile.Name(), newPath)
	if err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to rename file: %w", mlAssetID, err)
	}
	return newPath, isPhoto, nil
}
