package ebird

import (
	"os"
	"testing"
	"time"
)

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
