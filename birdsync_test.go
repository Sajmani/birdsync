package main

import (
	"iter"
	"testing"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

type mockEBirdClient struct {
	records []ebird.Record
}

func (m *mockEBirdClient) Records(path string) (iter.Seq[ebird.Record], error) {
	return func(yield func(ebird.Record) bool) {
		for _, r := range m.records {
			if !yield(r) {
				return
			}
		}
	}, nil
}

func (m *mockEBirdClient) DownloadMLAsset(id string) (string, bool, error) {
	return "", false, nil // isPhoto is false (media are sounds)
}

type mockINatClient struct {
	userID         string
	apitoken       string
	observations   []inat.Result
	createObsErr   error
	updateObsErr   error
	uploadMediaErr error
}

func (m *mockINatClient) GetUserID() string {
	return m.userID
}

func (m *mockINatClient) GetAPIToken() string {
	return m.apitoken
}

func (m *mockINatClient) DownloadObservations(userID string, after, before time.Time, fields ...string) []inat.Result {
	return m.observations
}

func (m *mockINatClient) CreateObservation(obs inat.Observation) error {
	return m.createObsErr
}

func (m *mockINatClient) UpdateObservation(obs inat.Observation) error {
	return m.updateObsErr
}

func (m *mockINatClient) UploadMedia(filename string, isPhoto bool, assetID, obsUUID string) error {
	return m.uploadMediaErr
}

func TestBirdsync(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	// Mock eBird records
	ebirdRecords := []ebird.Record{
		{
			SubmissionID:     "S123",               // skip: previously uploaded
			ScientificName:   "Larus delawarensis", // skip: previously uploaded
			CommonName:       "Ring-billed Gull",
			Date:             "2023-01-01",
			Time:             "12:00 PM",
			MLCatalogNumbers: "12345",
		},
		{
			SubmissionID:     "S124",
			ScientificName:   "Buteo jamaicensis",
			CommonName:       "Red-tailed Hawk",
			Date:             "2023-01-02",
			Time:             "01:00 PM",
			MLCatalogNumbers: "", // skip: unverifiable
		},
		{
			SubmissionID:   "S125",
			ScientificName: "Cardinalis cardinalis",
			CommonName:     "Northern Cardinal",
			Date:           "2022-12-31", // skip: before --after
			Time:           "10:00 AM",
		},
		{
			SubmissionID:   "S126",
			ScientificName: "Turdus migratorius",
			CommonName:     "American Robin",
			Date:           "2023-01-04", // skip: after --before
			Time:           "11:00 AM",
		},
		{
			SubmissionID:     "S127",
			ScientificName:   "Zenaida macroura",
			CommonName:       "Mourning Dove", // skip: fuzzy match
			Date:             "2023-01-03",    // skip: fuzzy match
			Time:             "02:00 PM",
			MLCatalogNumbers: "54321",
		},
		{
			SubmissionID:     "S128", // successful upload
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			MLCatalogNumbers: "67890",
		},
		{
			SubmissionID:     "S129", // previously uploaded with new media
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "1/3/2023", // test alternate date format
			Time:             "3:00 PM",  // test alternate time format
			MLCatalogNumbers: "67891 67890",
		},
		{
			SubmissionID:     "S130",
			ScientificName:   "Turdus migratorius", // Fuzzy match on scientific name
			CommonName:       "Robin",
			Date:             "2023-01-03",
			Time:             "09:00 AM",
			MLCatalogNumbers: "98765",
		},
	}

	// Mock iNaturalist observations
	inatObservations := []inat.Result{
		{ // previously uploaded with media
			UUID:       uuid.New(),
			ObservedOn: "2023-01-01",
			Taxon:      inat.Taxon{PreferredCommonName: "Ring-billed Gull"},
			Sounds:     []inat.Sound{{OriginalFilename: "ML12345.wav"}},
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S123"},
				{FieldID: inat.EBirdScientificNameField, Value: "Larus delawarensis"},
			},
		},
		{ // previously uploaded, without media
			UUID:       uuid.New(),
			ObservedOn: "2023-01-04",
			Taxon:      inat.Taxon{PreferredCommonName: "American Crow"},
			Sounds:     []inat.Sound{{OriginalFilename: "ML67890.wav"}},
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S129"},
				{FieldID: inat.EBirdScientificNameField, Value: "Corvus brachyrhynchos"},
			},
		},
		{ // fuzzy match
			UUID:       uuid.New(),
			ObservedOn: "2023-01-03",
			Taxon:      inat.Taxon{PreferredCommonName: "Mourning Dove"},
		},
		{ // fuzzy match on scientific name
			UUID:       uuid.New(),
			ObservedOn: "2023-01-03",
			Taxon:      inat.Taxon{Name: "Turdus migratorius", PreferredCommonName: "American Robin"},
		},
	}

	mockEbird := &mockEBirdClient{records: ebirdRecords}
	mockInat := &mockINatClient{userID: "testuser", observations: inatObservations}

	// Set flags
	after.Set("2023-01-01")
	before.Set("2023-01-04")
	verifiable = true
	fuzzy = true

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.totalRecords != 8 {
		t.Errorf("Expected 8 total records, got %d", stats.totalRecords)
	}
	if stats.previouslySkips != 1 {
		t.Errorf("Expected 1 previously skipped, got %d", stats.previouslySkips)
	}
	if stats.verifiableSkips != 1 {
		t.Errorf("Expected 1 verifiable skipped, got %d", stats.verifiableSkips)
	}
	if stats.afterSkips != 1 {
		t.Errorf("Expected 1 after skipped, got %d", stats.afterSkips)
	}
	if stats.beforeSkips != 1 {
		t.Errorf("Expected 1 before skipped, got %d", stats.beforeSkips)
	}
	if stats.fuzzySkips != 2 {
		t.Errorf("Expected 2 fuzzy skips, got %d", stats.fuzzySkips)
	}
	if stats.createdObservations != 1 {
		t.Errorf("Expected 1 created observations, got %d", stats.createdObservations)
	}
	if stats.updatedObservations != 2 {
		t.Errorf("Expected 2 updated observations, got %d", stats.updatedObservations)
	}
	if stats.uploadedPhotos != 0 {
		t.Errorf("Expected 0 uploaded photos, got %d", stats.uploadedPhotos)
	}
	if stats.uploadedSounds != 2 {
		t.Errorf("Expected 2 updated sounds, got %d", stats.uploadedSounds)
	}
}

func TestUpdateMedia(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	// Mock eBird records
	ebirdRecords := []ebird.Record{
		{
			SubmissionID:     "S129",
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			MLCatalogNumbers: "67891 67890",
		},
	}

	// Mock iNaturalist observations
	inatObservations := []inat.Result{
		{ // previously uploaded, without one of the media assets
			UUID:       uuid.New(),
			ObservedOn: "2023-01-03",
			Taxon:      inat.Taxon{PreferredCommonName: "American Crow"},
			Sounds:     []inat.Sound{{OriginalFilename: "ML67890.wav"}},
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S129"},
				{FieldID: inat.EBirdScientificNameField, Value: "Corvus brachyrhynchos"},
			},
		},
	}

	mockEbird := &mockEBirdClient{records: ebirdRecords}
	mockInat := &mockINatClient{userID: "testuser", observations: inatObservations}

	// Reset flags to default
	after.Set("")
	before.Set("")
	verifiable = false
	fuzzy = false

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.totalRecords != 1 {
		t.Errorf("Expected 1 total records, got %d", stats.totalRecords)
	}
	if stats.previouslySkips != 0 {
		t.Errorf("Expected 0 previously skipped, got %d", stats.previouslySkips)
	}
	if stats.createdObservations != 0 {
		t.Errorf("Expected 0 created observations, got %d", stats.createdObservations)
	}
	if stats.updatedObservations != 1 {
		t.Errorf("Expected 1 updated observations, got %d", stats.updatedObservations)
	}
	if stats.uploadedPhotos != 0 {
		t.Errorf("Expected 0 uploaded photos, got %d", stats.uploadedPhotos)
	}
	if stats.uploadedSounds != 1 {
		t.Errorf("Expected 1 uploaded sound, got %d", stats.uploadedSounds)
	}
}