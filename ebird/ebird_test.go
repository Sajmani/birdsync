package ebird

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
			w.Header().Set("Content-Type", "audio/mpeg")
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
// machine: mime.ExtensionsByType returns a sorted set, so taking its first
// element gave .jpe on macOS, .jfif on Linux, and .m2a for MP3 audio
// everywhere — a filename that depended on whose laptop ran the sync (T-004).
//
// Verifies: P-045, T-004.
func TestFileExtension(t *testing.T) {
	tests := []struct {
		contentType string
		want        string
		wantAny     bool // any non-empty extension is acceptable
		wantErr     bool
	}{
		{contentType: "image/jpeg", want: ".jpg"},
		{contentType: "audio/mpeg", want: ".mp3"},
		{contentType: "image/png", want: ".png"},
		// Parameters are legal on the header and must not defeat the lookup.
		{contentType: "image/jpeg; charset=binary", want: ".jpg"},
		{contentType: "IMAGE/JPEG", want: ".jpg"},
		// Unmapped but known to the system: falls back rather than failing.
		// No specific extension is asserted, because the answer genuinely
		// depends on the machine's mime database — text/plain resolves to
		// .conf on macOS. Asserting one would rebuild the fragility this
		// mapping exists to remove.
		{contentType: "text/plain", wantAny: true},
		{contentType: "not a media type", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			got, err := fileExtension(tt.contentType)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("fileExtension(%q) = %q, want an error", tt.contentType, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("fileExtension(%q) error = %v", tt.contentType, err)
			}
			if tt.wantAny {
				if got == "" {
					t.Errorf("fileExtension(%q) returned no extension", tt.contentType)
				}
				return
			}
			if got != tt.want {
				t.Errorf("fileExtension(%q) = %q, want %q", tt.contentType, got, tt.want)
			}
		})
	}
}
