package ebird

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Verifies: P-043, P-044.
func TestDownloadMLAsset_Sound(t *testing.T) {
	const assetID = "12345"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset/" + assetID + "/2400":
			// Photo not found; this is a sound asset.
			http.NotFound(w, r)
		case "/asset/" + assetID + "/mp3":
			// audio/mpeg3 is what the Macaulay Library CDN actually sends,
			// verified against it. The previous fixture said audio/mpeg, a
			// plausible invention, and that is why nobody noticed that no
			// sound file could be downloaded at all.
			w.Header().Set("Content-Type", "audio/mpeg3")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake mp3 data"))
		default:
			t.Errorf("Unexpected request path: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	filename, isPhoto, err := downloadMLAsset(server.URL, assetID)
	if err != nil {
		t.Fatalf("downloadMLAsset() error = %v", err)
	}
	defer os.Remove(filename)

	if isPhoto {
		t.Errorf("Expected isPhoto=false for a sound asset, got true")
	}
	// The extension must be exactly .mp3, not merely one of the several the
	// system mime database associates with audio/mpeg. Accepting any of them
	// is what let the filename vary by machine (P-045).
	if ext := filepath.Ext(filename); ext != ".mp3" {
		t.Errorf("Sound extension = %q, want %q (P-045)", ext, ".mp3")
	}
}

// Verifies: P-043, P-044.
func TestDownloadMLAsset_Photo(t *testing.T) {
	const assetID = "67890"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset/" + assetID + "/2400":
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("fake jpeg data"))
		default:
			t.Errorf("Unexpected request path: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	filename, isPhoto, err := downloadMLAsset(server.URL, assetID)
	if err != nil {
		t.Fatalf("downloadMLAsset() error = %v", err)
	}
	defer os.Remove(filename)

	if !isPhoto {
		t.Errorf("Expected isPhoto=true for a photo asset, got false")
	}
	if ext := filepath.Ext(filename); ext != ".jpg" {
		t.Errorf("Photo extension = %q, want %q (P-045)", ext, ".jpg")
	}
}

// Verifies: T-019.
func TestRecord_Observed(t *testing.T) {
	testCases := []struct {
		name     string
		record   Record
		expected time.Time
		hasError bool
	}{
		{
			name: "Date only with dash",
			record: Record{
				Date: "2023-01-02",
			},
			expected: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name: "Date only with slash",
			record: Record{
				Date: "1/2/2023",
			},
			expected: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name: "Date and time with dash",
			record: Record{
				Date: "2023-01-02",
				Time: "03:04 PM",
			},
			expected: time.Date(2023, 1, 2, 15, 4, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name: "Date and time with slash",
			record: Record{
				Date: "1/2/2023",
				Time: "3:04 PM",
			},
			expected: time.Date(2023, 1, 2, 15, 4, 0, 0, time.UTC),
			hasError: false,
		},
		{
			name: "Invalid date",
			record: Record{
				Date: "invalid-date",
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			observed, err := tc.record.Observed()
			if tc.hasError {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !observed.Equal(tc.expected) {
					t.Errorf("Expected %v, but got %v", tc.expected, observed)
				}
			}
		})
	}
}

// Verifies: T-018.
func TestRecords(t *testing.T) {
	csvData := `Submission ID,Common Name,Scientific Name,Taxonomic Order,Count,State/Province,County,Location ID,Location,Latitude,Longitude,Date,Time,Protocol,Duration (Min),All Obs Reported,Distance Traveled (km),Area Covered (ha),Number of Observers,Breeding Code,Observation Details,Checklist Comments,ML Catalog Numbers
S123,American Robin,Turdus migratorius,1,1,CA,Santa Clara,L123,Some Park,37.123,-122.123,2023-01-02,03:04 PM,Stationary,60,1,0,0,1,,,
`
	tmpfile, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(csvData)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	records, err := Records(tmpfile.Name())
	if err != nil {
		t.Fatalf("Records() error: %v", err)
	}

	var recs []Record
	for rec := range records {
		recs = append(recs, rec)
	}

	if len(recs) != 1 {
		t.Fatalf("Expected 1 record, but got %d", len(recs))
	}

	expectedRecord := Record{
		Line:           2,
		SubmissionID:   "S123",
		CommonName:     "American Robin",
		ScientificName: "Turdus migratorius",
		Count:          "1",
		StateProvince:  "CA",
		County:         "Santa Clara",
		LocationID:     "L123",
		Location:       "Some Park",
		Latitude:       "37.123",
		Longitude:      "-122.123",
		Date:           "2023-01-02",
		Time:           "03:04 PM",
		Protocol:       "Stationary",
	}

	rec := recs[0]
	if rec.Line != expectedRecord.Line ||
		rec.SubmissionID != expectedRecord.SubmissionID ||
		rec.CommonName != expectedRecord.CommonName ||
		rec.ScientificName != expectedRecord.ScientificName ||
		rec.Count != expectedRecord.Count ||
		rec.StateProvince != expectedRecord.StateProvince ||
		rec.County != expectedRecord.County ||
		rec.LocationID != expectedRecord.LocationID ||
		rec.Location != expectedRecord.Location ||
		rec.Latitude != expectedRecord.Latitude ||
		rec.Longitude != expectedRecord.Longitude ||
		rec.Date != expectedRecord.Date ||
		rec.Time != expectedRecord.Time ||
		rec.Protocol != expectedRecord.Protocol {
		t.Errorf("Expected record %+v, but got %+v", expectedRecord, rec)
	}
}

// Verifies: P-019, P-022.
func TestObservationID_Valid(t *testing.T) {
	testCases := []struct {
		name string
		id   ObservationID
		want bool
	}{
		{
			name: "valid",
			id:   ObservationID{SubmissionID: "S123", ScientificName: "Turdus migratorius"},
			want: true,
		},
		{
			name: "missing submission id",
			id:   ObservationID{ScientificName: "Turdus migratorius"},
			want: false,
		},
		{
			name: "missing scientific name",
			id:   ObservationID{SubmissionID: "S123"},
			want: false,
		},
		{
			name: "both missing",
			id:   ObservationID{},
			want: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.id.Valid(); got != tc.want {
				t.Errorf("ObservationID.Valid() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestFileExtension covers the mapping directly, including the cases the two
// download tests can't reach. The extension has to be the same on every
// machine, and the function must never fail: refusing to download a user's
// sound file because its content type wasn't recognised is a far worse outcome
// than naming it slightly wrong (T-004, P-045).
//
// Verifies: P-045, T-004.
func TestFileExtension(t *testing.T) {
	tests := []struct {
		contentType  string
		isPhoto      bool
		want         string
		wantFallback bool // reached the endpoint default rather than the map
	}{
		// Content types confirmed against the Macaulay Library CDN.
		{contentType: "image/jpeg", isPhoto: true, want: ".jpg"},
		{contentType: "audio/mpeg3", want: ".mp3"},
		// Registered spelling, in case the CDN ever switches to it.
		{contentType: "audio/mpeg", want: ".mp3"},
		{contentType: "image/png", isPhoto: true, want: ".png"},
		// Parameters are legal on the header and must not defeat the lookup.
		{contentType: "image/jpeg; charset=binary", isPhoto: true, want: ".jpg"},
		{contentType: "IMAGE/JPEG", isPhoto: true, want: ".jpg"},
		// Unrecognised: fall back to what the endpoint implies rather than
		// failing, which is what lost every sound file.
		{contentType: "audio/x-something-new", want: ".mp3", wantFallback: true},
		{contentType: "image/webp", isPhoto: true, want: ".jpg", wantFallback: true},
		{contentType: "", want: ".mp3", wantFallback: true},
		{contentType: "not a media type", isPhoto: true, want: ".jpg", wantFallback: true},
	}
	for _, tt := range tests {
		name := tt.contentType
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			// The fallback returns the same answer as the map for the types
			// birdsync sees, so comparing extensions alone can't tell whether
			// a type is actually mapped. Watch for the warning instead:
			// without this, deleting audio/mpeg3 from the map changes nothing
			// that any test can see.
			var logged bytes.Buffer
			log.SetOutput(&logged)
			defer log.SetOutput(os.Stderr)

			got := fileExtension(tt.contentType, tt.isPhoto)
			if got != tt.want {
				t.Errorf("fileExtension(%q, isPhoto=%v) = %q, want %q",
					tt.contentType, tt.isPhoto, got, tt.want)
			}
			fellBack := strings.Contains(logged.String(), "Unrecognized Content-Type")
			if fellBack != tt.wantFallback {
				t.Errorf("fileExtension(%q) fell back to the endpoint default = %v, want %v: %s",
					tt.contentType, fellBack, tt.wantFallback, logged.String())
			}
		})
	}
}

// TestRecordsRejectsBadInput covers the three ways the input path goes wrong,
// each of which used to produce either an operating-system message that
// explains nothing or, worse, no error at all.
//
// The checks are on the path and the header rather than on error text: issue #1
// reported "Incorrect function.", which is what Windows returns for reading a
// directory, where Unix says "is a directory" (T-004).
//
// Verifies: P-066.
func TestRecordsRejectsBadInput(t *testing.T) {
	dir := t.TempDir()

	zipPath := filepath.Join(dir, "ebird_1757602033437.zip")
	if err := os.WriteFile(zipPath, []byte("PK\x03\x04 not really a zip"), 0o600); err != nil {
		t.Fatal(err)
	}
	// A CSV that parses fine but isn't an eBird export. Without a header check
	// this is the dangerous case: every lookup of a missing column returns
	// index 0, so each record takes its submission ID from the first column.
	foreign := filepath.Join(dir, "something-else.csv")
	if err := os.WriteFile(foreign, []byte("Name,Date\nAmerican Robin,2023-01-02\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name string
		path string
		want string
	}{
		// The expected text is birdsync's own remedy, not a word that might
		// appear in the operating system's message by luck. Matching
		// "directory" alone passes on Unix, whose error says "is a directory",
		// and fails on Windows, whose error says "Incorrect function." — which
		// is the bug being fixed.
		{"directory", dir, "MyEBirdData.csv"},
		{"zip archive", zipPath, "Extract"},
		{"not an eBird export", foreign, "Submission ID"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Records(tt.path)
			if err == nil {
				t.Fatalf("Records(%s) returned no error", tt.path)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("Error %q doesn't mention %q, so it doesn't tell the user what to fix (P-066)",
					err, tt.want)
			}
		})
	}
}
