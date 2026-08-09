package main

import (
	"iter"
	"strings"
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

// uploadedMedia records one call to mockINatClient.UploadMedia.
type uploadedMedia struct {
	filename string
	isPhoto  bool
	assetID  string
	obsUUID  string
}

type mockINatClient struct {
	userID         string
	apitoken       string
	observations   []inat.Result
	createObsErr   error
	updateObsErr   error
	uploadMediaErr error

	// Every mutating call is recorded, so a test can assert both what birdsync
	// sent and — for --dryrun — that it sent nothing at all. Without this the
	// dry-run guarantee (T-005, P-051) can only be checked indirectly through
	// the counters, which is exactly the mistake CR-001 records.
	created  []inat.Observation
	updated  []inat.Observation
	uploaded []uploadedMedia
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
	m.created = append(m.created, obs)
	return m.createObsErr
}

func (m *mockINatClient) UpdateObservation(obs inat.Observation) error {
	m.updated = append(m.updated, obs)
	return m.updateObsErr
}

func (m *mockINatClient) UploadMedia(filename string, isPhoto bool, assetID, obsUUID string) error {
	m.uploaded = append(m.uploaded, uploadedMedia{filename, isPhoto, assetID, obsUUID})
	return m.uploadMediaErr
}

// resetFlags restores the package-level flag variables to their defaults, so a
// test doesn't inherit state from whichever test ran before it. The date flags
// must be zeroed directly: dateTimeFlag.Set rejects the empty string, so
// after.Set("") silently leaves the previous value in place.
func resetFlags() {
	dryRun = false
	verifiable = true
	fuzzy = false
	after = dateTimeFlag{}
	before = dateTimeFlag{}
	positionalAccuracy = ebird.PositionalAccuracy
}

// TestBirdsync exercises the full skip order against one set of records:
// every rule fires exactly once.
//
// Verifies: P-020, P-026, P-027, P-028, P-030, P-031.
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
			UUID:        uuid.New(),
			ObservedOn:  "2023-01-01",
			Taxon:       inat.Taxon{PreferredCommonName: "Ring-billed Gull"},
			Description: mlAssetURL("12345"),
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S123"},
				{FieldID: inat.EBirdScientificNameField, Value: "Larus delawarensis"},
			},
		},
		{ // previously uploaded, without media
			UUID:        uuid.New(),
			ObservedOn:  "2023-01-04",
			Taxon:       inat.Taxon{PreferredCommonName: "American Crow"},
			Description: mlAssetURL("67890"),
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
	resetFlags()
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

// Verifies: P-047.
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
			UUID:        uuid.New(),
			ObservedOn:  "2023-01-03",
			Taxon:       inat.Taxon{PreferredCommonName: "American Crow"},
			Description: mlAssetURL("67890"),
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S129"},
				{FieldID: inat.EBirdScientificNameField, Value: "Corvus brachyrhynchos"},
			},
		},
	}

	mockEbird := &mockEBirdClient{records: ebirdRecords}
	mockInat := &mockINatClient{userID: "testuser", observations: inatObservations}

	// Reset flags to default. This must go through resetFlags: after.Set("")
	// returns an error and leaves the previous test's date window in place, so
	// this test used to pass only because 2023-01-03 happened to fall inside
	// the window TestBirdsync left behind (T-015).
	resetFlags()
	verifiable = false

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

// TestFuzzyMatchDateFormats checks that fuzzy matching works regardless of which
// date format the user's eBird export uses. iNaturalist always reports
// observed_on as 2006-01-02, but eBird writes either 2006-01-02 or 1/2/2006,
// so the raw CSV field can't be compared directly.
func TestFuzzyMatchDateFormats(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	for _, tc := range []struct {
		name string
		date string
		time string
	}{
		{"ISO date", "2023-01-03", "02:00 PM"},
		{"slash date", "1/3/2023", "2:00 PM"},
		{"ISO date, no time", "2023-01-03", ""},
		{"slash date, no time", "1/3/2023", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockEbird := &mockEBirdClient{records: []ebird.Record{{
				SubmissionID:     "S200",
				ScientificName:   "Zenaida macroura",
				CommonName:       "Mourning Dove",
				Date:             tc.date,
				Time:             tc.time,
				MLCatalogNumbers: "11111",
			}}}
			// A non-birdsync observation of the same bird on the same day:
			// it has no eBird observation fields, so it's a fuzzy-match candidate.
			mockInat := &mockINatClient{observations: []inat.Result{{
				UUID:       uuid.New(),
				ObservedOn: "2023-01-03",
				Taxon:      inat.Taxon{PreferredCommonName: "Mourning Dove"},
			}}}

			resetFlags()
			fuzzy = true

			stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

			if stats.fuzzySkips != 1 {
				t.Errorf("Expected 1 fuzzy skip for date %q, got %d", tc.date, stats.fuzzySkips)
			}
			if stats.createdObservations != 0 {
				t.Errorf("Expected 0 created observations for date %q, got %d", tc.date, stats.createdObservations)
			}
		})
	}
}

// TestFuzzyMatchIgnoresEmptyNames checks that an iNaturalist observation with no
// taxon name doesn't fuzzy-match every eBird record that happens to be missing a
// name on the same date, which would silently drop legitimate observations.
func TestFuzzyMatchIgnoresEmptyNames(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{{
		SubmissionID:     "S201",
		ScientificName:   "Corvus brachyrhynchos",
		CommonName:       "", // no common name in this export
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: "22222",
	}}}
	// An unidentified observation: iNaturalist returns an empty taxon.
	mockInat := &mockINatClient{observations: []inat.Result{{
		UUID:       uuid.New(),
		ObservedOn: "2023-01-03",
		Taxon:      inat.Taxon{},
	}}}

	resetFlags()
	fuzzy = true

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.fuzzySkips != 0 {
		t.Errorf("Expected 0 fuzzy skips, got %d", stats.fuzzySkips)
	}
	if stats.createdObservations != 1 {
		t.Errorf("Expected 1 created observation, got %d", stats.createdObservations)
	}
}

// TestDryRunMediaCount checks that a dry run reports media assets as an
// unclassified count. It can't report photos and sounds separately, because it
// never downloads the asset and the Macaulay Library ID doesn't reveal its type.
func TestDryRunMediaCount(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{{
		SubmissionID:     "S202",
		ScientificName:   "Corvus brachyrhynchos",
		CommonName:       "American Crow",
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: "33333 44444",
	}}}
	mockInat := &mockINatClient{}

	resetFlags()
	dryRun = true
	defer func() { dryRun = false }()

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.pendingMedia != 2 {
		t.Errorf("Expected 2 pending media assets, got %d", stats.pendingMedia)
	}
	if stats.uploadedPhotos != 0 {
		t.Errorf("Expected 0 uploaded photos in a dry run, got %d", stats.uploadedPhotos)
	}
	if stats.uploadedSounds != 0 {
		t.Errorf("Expected 0 uploaded sounds in a dry run, got %d", stats.uploadedSounds)
	}
	if stats.createdObservations != 1 {
		t.Errorf("Expected 1 created observation, got %d", stats.createdObservations)
	}
}

// TestDryRunIssuesNoWrites is the direct check on the guarantee the whole tool
// rests on: a dry run may read, but must not create, update, or upload
// anything. The counters are deliberately not consulted here — they are a
// separate claim, and CR-001 records what happens when the two are conflated.
//
// Verifies: P-051, T-005.
func TestDryRunIssuesNoWrites(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{
		{
			// A new record: would be created, and its media uploaded.
			SubmissionID:     "S300",
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			MLCatalogNumbers: "33333 44444",
		},
		{
			// A previously synced record with an asset added since: would be
			// updated. This exercises the second write path, which a test
			// using only new records would miss.
			SubmissionID:     "S301",
			ScientificName:   "Turdus migratorius",
			CommonName:       "American Robin",
			Date:             "2023-01-03",
			Time:             "09:00 AM",
			MLCatalogNumbers: "55555 66666",
		},
	}}
	mockInat := &mockINatClient{observations: []inat.Result{{
		UUID:        uuid.New(),
		ObservedOn:  "2023-01-03",
		Taxon:       inat.Taxon{PreferredCommonName: "American Robin"},
		Description: mlAssetURL("55555"),
		Ofvs: []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S301"},
			{FieldID: inat.EBirdScientificNameField, Value: "Turdus migratorius"},
		},
	}}}

	resetFlags()
	dryRun = true
	defer func() { dryRun = false }()

	birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if len(mockInat.created) != 0 {
		t.Errorf("--dryrun created %d observations, want 0: %+v", len(mockInat.created), mockInat.created)
	}
	if len(mockInat.updated) != 0 {
		t.Errorf("--dryrun updated %d observations, want 0: %+v", len(mockInat.updated), mockInat.updated)
	}
	if len(mockInat.uploaded) != 0 {
		t.Errorf("--dryrun uploaded %d media assets, want 0: %+v", len(mockInat.uploaded), mockInat.uploaded)
	}
}

// TestUntaxonedObservationIsRecognized checks that a previously synced
// observation is matched by its sync key whatever its taxon — including no
// taxon at all, which is what iNaturalist leaves behind when it can't resolve
// an eBird name like "Aythya marila/affinis".
//
// This states the requirement rather than the mechanism, so it passes against
// the code that had the CR-003 bug too: that bug was in the download query, not
// in the matching. It is here to stop the next optimization of the download
// from quietly reintroducing the same failure.
//
// Verifies: P-020, P-061.
func TestUntaxonedObservationIsRecognized(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{{
		SubmissionID:     "S500",
		ScientificName:   "Aythya marila/affinis",
		CommonName:       "Greater/Lesser Scaup",
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: "77777",
	}}}
	mockInat := &mockINatClient{observations: []inat.Result{{
		UUID:        uuid.New(),
		ObservedOn:  "2023-01-03",
		Taxon:       inat.Taxon{}, // unresolvable name: no taxon, no iconic taxon
		Description: mlAssetURL("77777"),
		Ofvs: []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S500"},
			{FieldID: inat.EBirdScientificNameField, Value: "Aythya marila/affinis"},
		},
	}}}

	resetFlags()

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.previouslySkips != 1 {
		t.Errorf("Expected the untaxoned observation to be recognized as already synced, got %d skips", stats.previouslySkips)
	}
	if len(mockInat.created) != 0 {
		t.Errorf("Re-created %d observations that already existed, want 0 (P-020)", len(mockInat.created))
	}
}

// TestCreatedObservationContent pins down what birdsync actually sends to
// iNaturalist. Until the mock recorded its arguments, every requirement in
// "What birdsync writes" was unverified: the tests could only see counters.
//
// Verifies: P-019, P-035, P-036, P-037, P-038, P-039, P-040.
func TestCreatedObservationContent(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	rec := ebird.Record{
		SubmissionID:       "S400",
		CommonName:         "American Crow",
		ScientificName:     "Corvus brachyrhynchos",
		Count:              "3",
		StateProvince:      "US-CA",
		County:             "Santa Clara",
		Location:           "Shoreline Lake",
		Latitude:           "37.4321",
		Longitude:          "-122.0789",
		Date:               "2023-01-03",
		Time:               "03:00 PM",
		Protocol:           "Stationary",
		NumberOfObservers:  "2",
		ObservationDetails: "perched on a snag",
		ChecklistComments:  "windy morning",
		MLCatalogNumbers:   "33333",
	}
	mockEbird := &mockEBirdClient{records: []ebird.Record{rec}}
	mockInat := &mockINatClient{}

	resetFlags()

	birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if len(mockInat.created) != 1 {
		t.Fatalf("Expected 1 created observation, got %d", len(mockInat.created))
	}
	obs := mockInat.created[0]

	if obs.CaptiveFlag {
		t.Error("CaptiveFlag = true, want false: eBird checklists record wild birds (P-035)")
	}
	if obs.LocationIsExact {
		t.Error("LocationIsExact = true, want false: the checklist location is not the bird's (P-036)")
	}
	if obs.Latitude != 37.4321 || obs.Longitude != -122.0789 {
		t.Errorf("Coordinates = (%v, %v), want (37.4321, -122.0789) (P-036)", obs.Latitude, obs.Longitude)
	}
	if obs.PositionalAccuracy != float64(ebird.PositionalAccuracy) {
		t.Errorf("PositionalAccuracy = %v, want %v (P-036)", obs.PositionalAccuracy, ebird.PositionalAccuracy)
	}
	if obs.SpeciesGuess != rec.ScientificName {
		t.Errorf("SpeciesGuess = %q, want %q (P-037)", obs.SpeciesGuess, rec.ScientificName)
	}
	if want := "2023-01-03 03:00 PM"; obs.ObservedOnString != want {
		t.Errorf("ObservedOnString = %q, want %q (P-038)", obs.ObservedOnString, want)
	}

	// P-039: the eBird columns copied into observation fields. The two eBird
	// fields are also the sync key (P-019), so their absence would silently
	// break idempotence rather than producing a visible error.
	wantFields := map[int]string{
		inat.CountField:               "3",
		inat.CommonNameField:          "American Crow",
		inat.LocationField:            "Shoreline Lake",
		inat.CountyField:              "Santa Clara",
		inat.StateOrProvinceField:     "US-CA",
		inat.NumObserversField:        "2",
		inat.EBirdField:               "S400",
		inat.EBirdScientificNameField: "Corvus brachyrhynchos",
	}
	gotFields := map[int]string{}
	for _, ofv := range obs.ObservationFieldValuesAttributes {
		s, ok := ofv.Value.(string)
		if !ok {
			t.Errorf("Observation field %d has non-string value %v", ofv.ObservationFieldID, ofv.Value)
			continue
		}
		gotFields[ofv.ObservationFieldID] = s
	}
	for id, want := range wantFields {
		if got := gotFields[id]; got != want {
			t.Errorf("Observation field %d = %q, want %q (P-039)", id, got, want)
		}
	}

	// P-040: the description. Asset lines are deliberately absent here — they
	// are appended by the later update, not at creation (P-042, P-046).
	for _, want := range []string{
		"github.com/Sajmani/birdsync",
		"perched on a snag",
		"https://ebird.org/checklist/S400",
		"Protocol: Stationary",
		"windy morning",
	} {
		if !strings.Contains(obs.Description, want) {
			t.Errorf("Description missing %q (P-040):\n%s", want, obs.Description)
		}
	}
	if strings.Contains(obs.Description, "Macaulay Library Asset:") {
		t.Errorf("Created description should not list assets yet (P-042, P-046):\n%s", obs.Description)
	}

	// The update that follows carries them (P-046), and the upload is keyed by
	// the Macaulay Library asset ID (P-045).
	if len(mockInat.updated) != 1 {
		t.Fatalf("Expected 1 update after media upload, got %d", len(mockInat.updated))
	}
	if !strings.Contains(mockInat.updated[0].Description, mlAssetURL("33333")) {
		t.Errorf("Updated description missing the asset URL (P-046):\n%s", mockInat.updated[0].Description)
	}
	if len(mockInat.uploaded) != 1 {
		t.Fatalf("Expected 1 media upload, got %d", len(mockInat.uploaded))
	}
	if got := mockInat.uploaded[0].assetID; got != "33333" {
		t.Errorf("Uploaded asset ID = %q, want %q (P-045)", got, "33333")
	}
}
