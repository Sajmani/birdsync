// Package ebird provides helper functions for working with eBird data.
package ebird

import (
	"encoding/csv"
	"fmt"
	"io"
	"iter"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// PositionalAccuracy is the default positional accuracy in meters
// that we use for eBird observations. This is intended to serve as
// an approximation of the radius of a typical eBird hotspot.
const PositionalAccuracy = 1000 // meters

// Record contains the fields in MyEBirdData.csv records.
type Record struct {
	Line               int // line in the CSV file
	SubmissionID       string
	CommonName         string
	ScientificName     string
	TaxonomicOrder     string
	Count              string // "X" or integer
	StateProvince      string
	County             string
	LocationID         string
	Location           string
	Latitude           string
	Longitude          string
	Date               string // YYYY-MM-DD
	Time               string // 07:00 AM
	Protocol           string
	DurationMin        string
	AllObsReported     string // "1" means yes
	DistanceTraveledKm string
	AreaCoveredHa      string
	NumberOfObservers  string
	BreedingCode       string
	ObservationDetails string
	ChecklistComments  string
	MLCatalogNumbers   string
}

func (r Record) URL() string {
	return "https://ebird.org/checklist/" + r.SubmissionID
}

func (r Record) URLWithSpecies() string {
	return fmt.Sprintf("%s [%s] (%s)", r.URL(), r.ScientificName, r.CommonName)
}

// Observed returns the observation time for this record.
// The record always includes the date but might not include the time.
// The date and time formats vary between users for reasons I don't understand.
func (r Record) Observed() (time.Time, error) {
	if r.Time == "" {
		if strings.Contains(r.Date, "/") {
			return time.Parse("1/2/2006", r.Date)
		} else {
			return time.Parse("2006-01-02", r.Date)
		}
	}
	if strings.Contains(r.Date, "/") {
		return time.Parse("1/2/2006 3:04 PM", r.Date+" "+r.Time)
	} else {
		return time.Parse("2006-01-02 03:04 PM", r.Date+" "+r.Time)
	}
}

func (r Record) ObservationID() ObservationID {
	return ObservationID{r.SubmissionID, r.ScientificName}
}

func Records(filename string) (iter.Seq[Record], error) {
	// Check the path before opening it. eBird's download arrives as a zip that
	// extracts to a folder, and users pass the folder or the zip by mistake
	// (issue #1). Letting that reach the CSV reader produces the operating
	// system's message — "is a directory" on Unix, "Incorrect function." on
	// Windows — which tells nobody what to do (P-066).
	info, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("Records(%s): %w", filename, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("Records(%s): that is a folder, not a file; pass the MyEBirdData.csv inside it", filename)
	}
	if strings.EqualFold(filepath.Ext(filename), ".zip") {
		return nil, fmt.Errorf("Records(%s): that is a zip archive. Extract it and pass the MyEBirdData.csv inside", filename)
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Records(%s): %w", filename, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	// eBird's CSV export returns a variable number of fields per record,
	// so disable this check. This means we need to explicitly check len(rec)
	// before accessing fields that might not be there.
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Records(%s): reading CSV: %w", filename, err)
	}
	if len(recs) < 1 {
		return nil, fmt.Errorf("Records(%s): file is empty", filename)
	}
	field := make(map[string]int)
	for i, f := range recs[0] {
		field[f] = i
	}
	// Without this the file parses happily and produces nonsense: a lookup of
	// an absent column returns index 0, so every record would take its
	// submission ID from whatever the first column holds (P-066).
	if _, ok := field["Submission ID"]; !ok {
		return nil, fmt.Errorf("Records(%s): no %q column; this does not look like an eBird MyEBirdData.csv export",
			filename, "Submission ID")
	}
	recs = recs[1:]
	log.Printf("Read %d eBird observations", len(recs))
	return func(yield func(Record) bool) {
		for i, rec := range recs {
			stringField := func(key string) string {
				if f := field[key]; f < len(rec) {
					return rec[f]
				}
				return ""
			}
			if !yield(Record{
				Line:               i + 2, // header was line 1
				SubmissionID:       stringField("Submission ID"),
				CommonName:         stringField("Common Name"),
				ScientificName:     stringField("Scientific Name"),
				TaxonomicOrder:     stringField("Taxonomic Order"),
				Count:              stringField("Count"),
				StateProvince:      stringField("State/Province"),
				County:             stringField("County"),
				LocationID:         stringField("Location ID"),
				Location:           stringField("Location"),
				Latitude:           stringField("Latitude"),
				Longitude:          stringField("Longitude"),
				Date:               stringField("Date"),
				Time:               stringField("Time"),
				Protocol:           stringField("Protocol"),
				DurationMin:        stringField("Duration (Min)"),
				AllObsReported:     stringField("All Obs Reported"),
				DistanceTraveledKm: stringField("Distance Traveled (km)"),
				AreaCoveredHa:      stringField("Area Covered (ha)"),
				NumberOfObservers:  stringField("Number of Observers"),
				BreedingCode:       stringField("Breeding Code"),
				ObservationDetails: stringField("Observation Details"),
				ChecklistComments:  stringField("Checklist Comments"),
				MLCatalogNumbers:   stringField("ML Catalog Numbers"),
			}) {
				return
			}
		}
	}, nil
}

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

// macaulayBaseURL is the base URL for the Macaulay Library CDN.
const macaulayBaseURL = "https://cdn.download.ams.birds.cornell.edu/api/v2"

// DownloadMLAsset downloads the photo or sound with the provided ML asset ID
// (numbers only) and returns the local filename and whether it's a photo.
// This file is temporary and may be deleted at any time.
//
// Since the ML asset ID doesn't indicate whether this is a photo or sound file,
// we try downloading the photo file first, and if it's not there,
// we try downloading the sound file.
func DownloadMLAsset(mlAssetID string) (string, bool, error) {
	return downloadMLAsset(macaulayBaseURL, mlAssetID)
}

// downloadMLAsset is the implementation of DownloadMLAsset, accepting a base
// URL so it can be tested against a local HTTP server.
func downloadMLAsset(baseURL, mlAssetID string) (string, bool, error) {
	// Try fetching this ML asset as a photo
	url := fmt.Sprintf("%s/asset/%s/2400", baseURL, mlAssetID)
	resp, err := http.Get(url)
	if err != nil {
		return "", false, fmt.Errorf("DownloadMLAsset(%s): %s: %w", mlAssetID, url, err)
	}
	defer resp.Body.Close()
	isPhoto := resp.StatusCode == http.StatusOK
	if resp.StatusCode == http.StatusNotFound {
		// Photo not found; try fetching it as a sound
		url = fmt.Sprintf("%s/asset/%s/mp3", baseURL, mlAssetID)
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
	tmpName := tmpFile.Name()
	// Clean up on every path but success. A failed download used to abandon the
	// file, and on Windows its open handle with it — the reported symptom being
	// "The process cannot access the file because it is being used by another
	// process" on a later attempt (issue #1, T-023).
	renamed := false
	defer func() {
		tmpFile.Close() // no-op once already closed below
		if !renamed {
			os.Remove(tmpName)
		}
	}()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to copy asset data to file: %w", mlAssetID, err)
	}

	// Detect the file extension from the Content-Type response header.
	ext := fileExtension(resp.Header.Get("Content-Type"), isPhoto)
	// Close before renaming: Windows refuses to rename a file that is still
	// open, which is how issue #1 first surfaced.
	if err := tmpFile.Close(); err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): closing temp file: %w", mlAssetID, err)
	}

	newPath := tmpName + ext
	if err := os.Rename(tmpName, newPath); err != nil {
		return "", isPhoto, fmt.Errorf("DownloadMLAsset(%s): failed to rename file: %w", mlAssetID, err)
	}
	renamed = true
	return newPath, isPhoto, nil
}

// canonicalExtensions gives one extension per content type the Macaulay
// Library serves. The values are what the CDN actually sends, checked against
// it, not what the standards suggest it ought to send: sounds arrive as
// audio/mpeg3, which is not a registered type at all.
//
// The system mime database is deliberately not consulted. mime.ExtensionsByType
// returns every extension registered for a type, sorted, and which ones those
// are depends on the machine — it gave .jpe for a JPEG on macOS, .jfif on
// Linux, and nothing whatsoever for audio/mpeg3. Both of this file's extension
// bugs came from trusting it (T-004, P-045).
var canonicalExtensions = map[string]string{
	"image/jpeg":  ".jpg",
	"image/png":   ".png",
	"audio/mpeg":  ".mp3",
	"audio/mpeg3": ".mp3",
	"audio/wav":   ".wav",
}

// fileExtension returns the filename extension for an asset served with the
// given Content-Type. isPhoto says which endpoint answered.
//
// It cannot fail. An unrecognized content type falls back to what the endpoint
// implies, because the photo URL serves images and the sound URL serves audio;
// getting the extension slightly wrong is a much smaller harm than refusing to
// download a user's media over it, which is what the previous version did to
// every sound file.
func fileExtension(contentType string, isPhoto bool) string {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err == nil {
		if ext, ok := canonicalExtensions[mediaType]; ok {
			return ext
		}
	}
	// Worth knowing about: it means the CDN has started serving something new,
	// and the map should be extended rather than left to the fallback.
	log.Printf("Unrecognized Content-Type %q for a %s; falling back by endpoint",
		contentType, map[bool]string{true: "photo", false: "sound"}[isPhoto])
	if isPhoto {
		return ".jpg"
	}
	return ".mp3"
}
