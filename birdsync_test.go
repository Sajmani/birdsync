package main

import (
	"errors"
	"iter"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Sajmani/birdsync/ebird"
	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

type mockEBirdClient struct {
	records []ebird.Record

	// tempDir, when set, makes DownloadMLAsset write a real file there and
	// record its path in downloaded, so a test can check it gets cleaned up.
	tempDir    string
	downloaded []string
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

// DownloadMLAsset returns an empty filename by default. Set tempDir to make it
// behave like the real one and write an actual file, which is what a test of
// temp-file cleanup needs.
func (m *mockEBirdClient) DownloadMLAsset(id string) (string, bool, error) {
	if m.tempDir == "" {
		return "", false, nil // isPhoto is false (media are sounds)
	}
	f, err := os.CreateTemp(m.tempDir, "ML"+id+"-*.mp3")
	if err != nil {
		return "", false, err
	}
	defer f.Close()
	if _, err := f.WriteString("fake asset data"); err != nil {
		return "", false, err
	}
	m.downloaded = append(m.downloaded, f.Name())
	return f.Name(), false, nil
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
	downloadErr    error
	createObsErr   error
	updateObsErr   error
	uploadMediaErr error

	// failUploads makes UploadMedia fail for particular Macaulay Library asset
	// IDs with a chosen error, so a test can exercise a partial failure rather
	// than an all-or-nothing one, and can distinguish a permanent refusal from
	// a transient one. Until this existed no test set any of the error fields.
	failUploads map[string]error

	// Every mutating call is recorded, so a test can assert both what birdsync
	// sent and — for --dryrun — that it sent nothing at all. Without this the
	// dry-run guarantee (tech.dryRunGateAtCallSite, product.dryRunNoWrites) can only be checked indirectly through
	// the counters, which is exactly the mistake decisions.dryRunCountersRelabeled records.
	created        []inat.Observation
	updated        []inat.Observation
	uploaded       []uploadedMedia
	downloadFields []string
}

func (m *mockINatClient) GetUserID() string {
	return m.userID
}

func (m *mockINatClient) GetAPIToken() string {
	return m.apitoken
}

func (m *mockINatClient) DownloadObservations(userID string, after, before time.Time, fields ...string) ([]inat.Result, error) {
	m.downloadFields = slices.Clone(fields)
	return m.observations, m.downloadErr
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
	if err, ok := m.failUploads[assetID]; ok {
		return err
	}
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
	dailyRequestBudget = inat.DefaultMaxRequests
	nowFunc = time.Now
}

// TestBirdsync exercises the full skip order against one set of records:
// every rule fires exactly once.
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
	// the window TestBirdsync left behind (tech.resetFlagsInTests).
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
// separate claim, and decisions.dryRunCountersRelabeled records what happens when the two are conflated.
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

// TestInvalidRecordsAreSkipped checks that a record birdsync can't parse costs
// that record, not the run. A single malformed date used to call log.Fatalf,
// ending the sync after an arbitrary number of observations had already been
// created — and eBird exports vary enough between users that one bad row in
// twelve thousand is a realistic thing to hit.
func TestInvalidRecordsAreSkipped(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{
		{
			SubmissionID:     "S600",
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			Latitude:         "37.4",
			Longitude:        "-122.0",
			MLCatalogNumbers: "60001",
		},
		{
			SubmissionID:     "S601",
			ScientificName:   "Turdus migratorius",
			CommonName:       "American Robin",
			Date:             "not-a-date", // unparseable
			Time:             "03:00 PM",
			MLCatalogNumbers: "60002",
		},
		{
			SubmissionID:     "S602",
			ScientificName:   "Zenaida macroura",
			CommonName:       "Mourning Dove",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			Latitude:         "not-a-number", // unparseable
			Longitude:        "-122.0",
			MLCatalogNumbers: "60003",
		},
		{
			// Must still be processed: the run continues past both failures.
			SubmissionID:     "S603",
			ScientificName:   "Sturnus vulgaris",
			CommonName:       "European Starling",
			Date:             "2023-01-03",
			Time:             "04:00 PM",
			Latitude:         "37.5",
			Longitude:        "-122.1",
			MLCatalogNumbers: "60004",
		},
	}}
	mockInat := &mockINatClient{}

	resetFlags()

	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.totalRecords != 4 {
		t.Errorf("totalRecords = %d, want 4", stats.totalRecords)
	}
	if stats.invalidSkips != 2 {
		t.Errorf("invalidSkips = %d, want 2 (product.badRecordsSkipped)", stats.invalidSkips)
	}
	if stats.createdObservations != 2 {
		t.Errorf("createdObservations = %d, want 2: the run must continue past a bad record", stats.createdObservations)
	}
	if len(mockInat.created) != 2 {
		t.Fatalf("Created %d observations, want 2", len(mockInat.created))
	}
	// The record after both failures is the one that proves the run continued.
	if got := mockInat.created[1].ObservationFieldValuesAttributes; len(got) == 0 {
		t.Error("Second created observation has no fields")
	} else if last := mockInat.created[1]; last.SpeciesGuess != "Sturnus vulgaris" {
		t.Errorf("Second created observation = %q, want the record following the bad ones", last.SpeciesGuess)
	}
}

// TestUntaxonedObservationIsRecognized checks that a previously synced
// observation is matched by its sync key whatever its taxon — including no
// taxon at all, which is what iNaturalist leaves behind when it can't resolve
// an eBird name like "Aythya marila/affinis".
//
// This states the requirement rather than the mechanism, so it passes against
// the code that had the decisions.noTaxonFilterOnDownload bug too: that bug was in the download query, not
// in the matching. It is here to stop the next optimization of the download
// from quietly reintroducing the same failure.
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
		t.Errorf("Re-created %d observations that already existed, want 0 (product.idempotentRerun)", len(mockInat.created))
	}
}

// TestCreatedObservationContent pins down what birdsync actually sends to
// iNaturalist. Until the mock recorded its arguments, every requirement in
// "What birdsync writes" was unverified: the tests could only see counters.
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
		t.Error("CaptiveFlag = true, want false: eBird checklists record wild birds (product.createdObservationWild)")
	}
	if obs.LocationIsExact {
		t.Error("LocationIsExact = true, want false: the checklist location is not the bird's (product.checklistCoordinates)")
	}
	if obs.Latitude != 37.4321 || obs.Longitude != -122.0789 {
		t.Errorf("Coordinates = (%v, %v), want (37.4321, -122.0789) (product.checklistCoordinates)", obs.Latitude, obs.Longitude)
	}
	if obs.PositionalAccuracy != float64(ebird.PositionalAccuracy) {
		t.Errorf("PositionalAccuracy = %v, want %v (product.checklistCoordinates)", obs.PositionalAccuracy, ebird.PositionalAccuracy)
	}
	if obs.SpeciesGuess != rec.ScientificName {
		t.Errorf("SpeciesGuess = %q, want %q (product.speciesGuessScientificName)", obs.SpeciesGuess, rec.ScientificName)
	}
	if want := "2023-01-03 03:00 PM"; obs.ObservedOnString != want {
		t.Errorf("ObservedOnString = %q, want %q (product.observedDateTime)", obs.ObservedOnString, want)
	}

	// product.observationFieldMapping: the eBird columns copied into observation fields. The two eBird
	// fields are also the sync key (product.syncKey), so their absence would silently
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
			t.Errorf("Observation field %d = %q, want %q (product.observationFieldMapping)", id, got, want)
		}
	}

	// product.descriptionContents: the description. Asset lines are deliberately absent here — they
	// are appended by the later update, not at creation (product.createBeforeMedia, product.descriptionUpdateKeepsMedia).
	for _, want := range []string{
		"github.com/Sajmani/birdsync",
		"perched on a snag",
		"https://ebird.org/checklist/S400",
		"Protocol: Stationary",
		"windy morning",
	} {
		if !strings.Contains(obs.Description, want) {
			t.Errorf("Description missing %q (product.descriptionContents):\n%s", want, obs.Description)
		}
	}
	if strings.Contains(obs.Description, "Macaulay Library Asset:") {
		t.Errorf("Created description should not list assets yet (product.createBeforeMedia, product.descriptionUpdateKeepsMedia):\n%s", obs.Description)
	}

	// The update that follows carries them (product.descriptionUpdateKeepsMedia), and the upload is keyed by
	// the Macaulay Library asset ID (product.mediaFilename).
	if len(mockInat.updated) != 1 {
		t.Fatalf("Expected 1 update after media upload, got %d", len(mockInat.updated))
	}
	if !strings.Contains(mockInat.updated[0].Description, mlAssetURL("33333")) {
		t.Errorf("Updated description missing the asset URL (product.descriptionUpdateKeepsMedia):\n%s", mockInat.updated[0].Description)
	}
	if len(mockInat.uploaded) != 1 {
		t.Fatalf("Expected 1 media upload, got %d", len(mockInat.uploaded))
	}
	if got := mockInat.uploaded[0].assetID; got != "33333" {
		t.Errorf("Uploaded asset ID = %q, want %q (product.mediaFilename)", got, "33333")
	}
}

// TestTempFilesAreCleanedUp checks that downloaded media doesn't accumulate in
// the system temp directory. Each asset is downloaded to a temp file and
// uploaded; nothing deleted them afterwards, so a full sync of an account with
// thousands of photos left thousands of files behind.
func TestTempFilesAreCleanedUp(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	dir := t.TempDir()
	mockEbird := &mockEBirdClient{
		tempDir: dir,
		records: []ebird.Record{{
			SubmissionID:     "S700",
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2023-01-03",
			Time:             "03:00 PM",
			MLCatalogNumbers: "70001 70002",
		}},
	}
	mockInat := &mockINatClient{}

	resetFlags()

	birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if len(mockEbird.downloaded) != 2 {
		t.Fatalf("Downloaded %d assets, want 2", len(mockEbird.downloaded))
	}
	for _, name := range mockEbird.downloaded {
		if _, err := os.Stat(name); !os.IsNotExist(err) {
			t.Errorf("Temp file %s still exists after the run (tech.tempFilesDeleted)", name)
		}
	}

	// The uploads must still have happened: deleting the file too early would
	// pass this test and break the feature.
	if len(mockInat.uploaded) != 2 {
		t.Errorf("Uploaded %d assets, want 2", len(mockInat.uploaded))
	}
}

// TestSummaryDryRunLabels checks that a dry run doesn't claim credit for work
// it didn't do. The counters were incremented outside the --dryrun branch, so
// a run that touched nothing still ended with "Created 57 new iNaturalist
// observations". The count is worth reporting; the claim is not.
func TestSummaryDryRunLabels(t *testing.T) {
	resetFlags()
	dryRun = true
	defer func() { dryRun = false }()

	s := stats{totalRecords: 9, createdObservations: 3, updatedObservations: 2, pendingMedia: 7}
	got := strings.Join(s.summary(), "\n")

	for _, want := range []string{
		"Would create 3 new iNaturalist observations",
		"Would update 2 iNaturalist observations",
		"Would upload 7 media assets",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Dry-run summary missing %q (product.dryRunCounterLabels):\n%s", want, got)
		}
	}
	for _, forbidden := range []string{"Created 3", "Updated 2", "Uploaded "} {
		if strings.Contains(got, forbidden) {
			t.Errorf("Dry-run summary claims %q, but a dry run writes nothing (tech.honestDryRunCounters):\n%s", forbidden, got)
		}
	}
}

// TestSummaryRealRunLabels is the other half: a real run must still report
// plainly, so the fix for the dry-run wording can't just hedge everything.
func TestSummaryRealRunLabels(t *testing.T) {
	resetFlags()

	s := stats{totalRecords: 9, createdObservations: 3, updatedObservations: 2,
		uploadedPhotos: 4, uploadedSounds: 1}
	got := strings.Join(s.summary(), "\n")

	for _, want := range []string{
		"Created 3 new iNaturalist observations",
		"Updated 2 iNaturalist observations",
		"Uploaded 4 photos",
		"Uploaded 1 sounds",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Summary missing %q (product.endOfRunSummary):\n%s", want, got)
		}
	}
	if strings.Contains(got, "Would ") {
		t.Errorf("Real run summary hedges with \"Would\":\n%s", got)
	}
}

// TestSummaryConditionalLines checks that a skip counter appears only when the
// rule that produces it was in effect, and the media-failure line only when
// something failed. Otherwise every run reports a wall of zeroes and people
// stop reading it.
func TestSummaryConditionalLines(t *testing.T) {
	resetFlags()
	verifiable = false

	quiet := strings.Join(stats{}.summary(), "\n")
	for _, forbidden := range []string{"--fuzzy", "--after", "--before", "unverifiable", "unparseable", "Failed to upload"} {
		if strings.Contains(quiet, forbidden) {
			t.Errorf("Summary mentions %q when the rule wasn't in effect (product.conditionalSkipCounters, product.conditionalFailureCount):\n%s", forbidden, quiet)
		}
	}

	resetFlags()
	fuzzy = true
	if err := after.Set("2023-01-01"); err != nil {
		t.Fatal(err)
	}
	if err := before.Set("2023-12-31"); err != nil {
		t.Fatal(err)
	}
	loud := strings.Join(stats{invalidSkips: 1, errors: 2}.summary(), "\n")
	for _, want := range []string{"--fuzzy", "--after", "--before", "unverifiable", "unparseable", "Failed to upload 2"} {
		if !strings.Contains(loud, want) {
			t.Errorf("Summary missing %q when the rule was in effect (product.conditionalSkipCounters, product.conditionalFailureCount):\n%s", want, loud)
		}
	}
}

// TestFailedUploadIsRetriedNextRun is the end-to-end statement of decisions.descriptionListsOnlyUploaded: a
// media upload that fails must be retried on the next run, not recorded as
// done.
//
// The description is how birdsync remembers which assets it has uploaded, so
// writing an asset's URL there before the upload is attempted makes a failure
// permanent -- the next run parses that URL back out, concludes the asset is
// already synced, and never tries again. The observation stays unverifiable
// until someone notices by hand.
//
// The test runs the sync twice, feeding the first run's own output back in as
// the second run's starting state, because that round trip is where the defect
// lives. Asserting only that the description omits the failed asset would test
// the symptom.
func TestFailedUploadIsRetriedNextRun(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	const good, bad = "80001", "80002"
	rec := ebird.Record{
		SubmissionID:     "S800",
		ScientificName:   "Corvus brachyrhynchos",
		CommonName:       "American Crow",
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: good + " " + bad,
	}

	// Run 1: the second asset fails to upload.
	firstEbird := &mockEBirdClient{records: []ebird.Record{rec}}
	firstInat := &mockINatClient{failUploads: map[string]error{bad: errors.New("upload failed")}}

	resetFlags()
	first := birdsync("MyEBirdData.csv", firstEbird, "myUserID", firstInat)

	if first.errors != 1 {
		t.Errorf("First run: errors = %d, want 1", first.errors)
	}
	if len(firstInat.updated) != 1 {
		t.Fatalf("First run: %d updates, want 1", len(firstInat.updated))
	}
	desc := firstInat.updated[0].Description
	if !strings.Contains(desc, mlAssetURL(good)) {
		t.Errorf("First run: description omits the asset that uploaded:\n%s", desc)
	}
	if strings.Contains(desc, mlAssetURL(bad)) {
		t.Errorf("First run: description claims asset %s, which failed to upload (product.descriptionContents):\n%s", bad, desc)
	}

	// Run 2: the same record, against the observation run 1 produced, with the
	// upload now succeeding. The failed asset must be picked up.
	created := firstInat.created[0]
	secondEbird := &mockEBirdClient{records: []ebird.Record{rec}}
	secondInat := &mockINatClient{observations: []inat.Result{{
		UUID:        created.UUID,
		ObservedOn:  "2023-01-03",
		Taxon:       inat.Taxon{Name: "Corvus brachyrhynchos"},
		Description: desc,
		Photos:      []inat.Photo{{}}, // the one asset that made it
		Ofvs: []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S800"},
			{FieldID: inat.EBirdScientificNameField, Value: "Corvus brachyrhynchos"},
		},
	}}}

	resetFlags()
	second := birdsync("MyEBirdData.csv", secondEbird, "myUserID", secondInat)

	if second.errors != 0 {
		t.Errorf("Second run: errors = %d, want 0", second.errors)
	}
	if len(secondInat.created) != 0 {
		t.Errorf("Second run: created %d observations, want 0 (the observation already exists)", len(secondInat.created))
	}
	if len(secondInat.uploaded) != 1 {
		t.Fatalf("Second run: uploaded %d assets, want 1 (the one that failed before)", len(secondInat.uploaded))
	}
	if got := secondInat.uploaded[0].assetID; got != bad {
		t.Errorf("Second run: retried asset %s, want %s", got, bad)
	}
}

// TestAllUploadsFailingLeavesNoUpdate checks the degenerate case: if nothing
// uploaded, there is nothing to say, so the observation is not updated and the
// counter doesn't claim it was.
func TestAllUploadsFailingLeavesNoUpdate(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	mockEbird := &mockEBirdClient{records: []ebird.Record{{
		SubmissionID:     "S801",
		ScientificName:   "Corvus brachyrhynchos",
		CommonName:       "American Crow",
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: "80003 80004",
	}}}
	mockInat := &mockINatClient{failUploads: map[string]error{"80003": errors.New("upload failed"), "80004": errors.New("upload failed")}}

	resetFlags()
	stats := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if stats.errors != 2 {
		t.Errorf("errors = %d, want 2", stats.errors)
	}
	if len(mockInat.updated) != 0 {
		t.Errorf("Updated the observation with nothing to add: %+v", mockInat.updated)
	}
	if stats.updatedObservations != 0 {
		t.Errorf("updatedObservations = %d, want 0: no update was made (tech.honestDryRunCounters)", stats.updatedObservations)
	}
}

// TestPermanentUploadFailureIsNotRetried covers the case decisions.descriptionListsOnlyUploaded's fix made
// worse before decisions.permanentFailuresRecorded fixed it: an asset iNaturalist will never accept, such as
// a sound file over its size limit. Retrying it means re-downloading tens of
// megabytes from the Macaulay Library on every run, forever.
//
// A permanent refusal is recorded in the description so it isn't retried, and
// reported each run so the user knows. A transient failure is not recorded, so
// it is.
func TestPermanentUploadFailureIsNotRetried(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	const tooBig = "90001"
	rec := ebird.Record{
		SubmissionID:     "S900",
		ScientificName:   "Corvus brachyrhynchos",
		CommonName:       "American Crow",
		Date:             "2023-01-03",
		Time:             "03:00 PM",
		MLCatalogNumbers: tooBig,
	}

	// Run 1: iNaturalist refuses the file itself.
	firstEbird := &mockEBirdClient{records: []ebird.Record{rec}}
	firstInat := &mockINatClient{failUploads: map[string]error{
		tooBig: &inat.StatusError{
			StatusCode: 422,
			Status:     "422 Unprocessable Entity",
			Body:       "File is too large (max 50 MB)",
		},
	}}

	resetFlags()
	birdsync("MyEBirdData.csv", firstEbird, "myUserID", firstInat)

	if len(firstInat.updated) != 1 {
		t.Fatalf("Expected the observation to be updated to record the failure, got %d updates", len(firstInat.updated))
	}
	desc := firstInat.updated[0].Description
	if !strings.Contains(desc, failedMarker) {
		t.Errorf("Description doesn't record the permanent failure (product.permanentUploadFailureRecorded):\n%s", desc)
	}

	// Run 2: same record, against the observation run 1 produced. The upload
	// would now succeed, but it must not be attempted — the description says
	// the service refused it.
	secondEbird := &mockEBirdClient{records: []ebird.Record{rec}}
	secondInat := &mockINatClient{observations: []inat.Result{{
		UUID:        firstInat.created[0].UUID,
		ObservedOn:  "2023-01-03",
		Taxon:       inat.Taxon{Name: "Corvus brachyrhynchos"},
		Description: desc,
		Ofvs: []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S900"},
			{FieldID: inat.EBirdScientificNameField, Value: "Corvus brachyrhynchos"},
		},
	}}}

	resetFlags()
	second := birdsync("MyEBirdData.csv", secondEbird, "myUserID", secondInat)

	if len(secondInat.uploaded) != 0 {
		t.Errorf("Retried an asset the service permanently refused: %+v (product.permanentUploadFailureRecorded)", secondInat.uploaded)
	}
	if second.errors != 0 {
		t.Errorf("errors = %d, want 0: nothing should have been attempted", second.errors)
	}
	if second.previouslySkips != 1 {
		t.Errorf("previouslySkips = %d, want 1", second.previouslySkips)
	}
}

// TestTransientUploadFailureIsRetried is the other side of product.permanentUploadFailureRecorded: a failure
// that might come good must not be recorded as permanent, or a network blip
// would cost a photo.
func TestTransientUploadFailureIsRetried(t *testing.T) {
	origDebug := debug
	debug = true
	defer func() { debug = origDebug }()

	for _, tc := range []struct {
		name string
		err  error
	}{
		{"server error", &inat.StatusError{StatusCode: 503, Status: "503 Service Unavailable"}},
		{"rate limited", &inat.StatusError{StatusCode: 429, Status: "429 Too Many Requests"}},
		{"token expired", &inat.StatusError{StatusCode: 401, Status: "401 Unauthorized"}},
		{"network error", errors.New("connection reset")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockEbird := &mockEBirdClient{records: []ebird.Record{{
				SubmissionID:     "S901",
				ScientificName:   "Corvus brachyrhynchos",
				CommonName:       "American Crow",
				Date:             "2023-01-03",
				Time:             "03:00 PM",
				MLCatalogNumbers: "90002",
			}}}
			mockInat := &mockINatClient{failUploads: map[string]error{"90002": tc.err}}

			resetFlags()
			birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

			for _, u := range mockInat.updated {
				if strings.Contains(u.Description, failedMarker) {
					t.Errorf("Recorded %v as permanent; it may come good on the next run (product.permanentUploadFailureRecorded):\n%s",
						tc.err, u.Description)
				}
			}
		})
	}
}

// TestDailyRequestBudget checks tech.dailyRequestBudget and birdsync.requestRate
// across both recent-account history and in-run writes:
//   - birdsync asks DownloadObservations for created_at so the estimate works
//     against the real API, not only mocks;
//   - observations created in the last 24 hours count toward the budget
//     (create + media uploads + description update), while older ones do not;
//   - the sync stops cleanly before an observation whose writes would push the
//     24-hour total past the budget, leaving no observation half-synced without
//     its media, and reports the stop in the summary.
func TestDailyRequestBudget(t *testing.T) {
	resetFlags()
	defer resetFlags()

	fixedNow := time.Date(2026, 11, 8, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }
	// Budget of 9 requests:
	//   - 1 download page for 2 existing observations
	//   - 4 requests from an observation created 2h ago with 2 photos (1 create + 2 uploads + 1 update)
	//   - 0 requests from an observation created 30h ago with 3 photos (outside 24h window)
	//   => 5 requests already used; 4 remaining.
	// Record 1 needs 3 requests (1 create + 1 media upload + 1 update) -> fits (total 8 <= 9).
	// Record 2 needs 3 requests (1 create + 1 media upload + 1 update) -> would hit 11 > 9, so birdsync stops before creating it.
	dailyRequestBudget = 9

	existing := []inat.Result{
		{
			ID:          1,
			UUID:        uuid.New(),
			CreatedAt:   fixedNow.Add(-2 * time.Hour).Format(time.RFC3339),
			ObservedOn:  "2023-05-01",
			Photos:      []inat.Photo{{ID: 10}, {ID: 11}},
			Description: assetLine("1001", true) + assetLine("1002", true),
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S100"},
				{FieldID: inat.EBirdScientificNameField, Value: "Turdus migratorius"},
			},
		},
		{
			ID:          2,
			UUID:        uuid.New(),
			CreatedAt:   fixedNow.Add(-30 * time.Hour).Format(time.RFC3339),
			ObservedOn:  "2023-05-02",
			Photos:      []inat.Photo{{ID: 20}, {ID: 21}, {ID: 22}},
			Description: assetLine("2001", true) + assetLine("2002", true) + assetLine("2003", true),
			Ofvs: []inat.Ofv{
				{FieldID: inat.EBirdField, Value: "S101"},
				{FieldID: inat.EBirdScientificNameField, Value: "Cardinalis cardinalis"},
			},
		},
	}

	records := []ebird.Record{
		{
			Line:             2,
			SubmissionID:     "S200",
			ScientificName:   "Corvus brachyrhynchos",
			CommonName:       "American Crow",
			Date:             "2026-11-01",
			Time:             "08:00 AM",
			MLCatalogNumbers: "3001",
		},
		{
			Line:             3,
			SubmissionID:     "S201",
			ScientificName:   "Cyanocitta cristata",
			CommonName:       "Blue Jay",
			Date:             "2026-11-01",
			Time:             "08:15 AM",
			MLCatalogNumbers: "3002",
		},
	}

	mockEbird := &mockEBirdClient{records: records}
	mockInat := &mockINatClient{observations: existing}

	s := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)

	if !slices.Contains(mockInat.downloadFields, "created_at") {
		t.Errorf("DownloadObservations fields %v missing \"created_at\"; prior 24h usage would always estimate 0", mockInat.downloadFields)
	}
	if len(mockInat.created) != 1 || len(mockInat.uploaded) != 1 || len(mockInat.updated) != 1 {
		t.Fatalf("got %d creates, %d uploads, %d updates; want 1, 1, 1 (must finish S200 and stop cleanly before S201)",
			len(mockInat.created), len(mockInat.uploaded), len(mockInat.updated))
	}
	if !s.stoppedAtDailyBudget {
		t.Error("stoppedAtDailyBudget = false, want true")
	}
	if s.totalRecords != 1 || s.createdObservations != 1 {
		t.Errorf("totalRecords=%d createdObservations=%d, want 1 and 1", s.totalRecords, s.createdObservations)
	}
	summary := strings.Join(s.summary(), "\n")
	if !strings.Contains(summary, "Stopped at daily iNaturalist API request budget") {
		t.Errorf("Summary does not report budget stop:\n%s", summary)
	}
}

// TestFuzzyMatchSubspecies checks that --fuzzy matches subspecies and eBird
// subspecific groups to the parent species (and to each other) via their
// binomial species stem, while refusing to match across species on higher
// taxa, spuhs ("Empidonax sp."), species slashes ("Larus glaucoides/thayeri"),
// or hybrids ("Anas platyrhynchos x Anas acuta").
func TestFuzzyMatchSubspecies(t *testing.T) {
	origDebug := debug
	defer func() {
		debug = origDebug
		resetFlags()
	}()

	cases := []struct {
		name        string
		ebirdCommon string
		ebirdSci    string
		inatCommon  string
		inatSci     string
		wantSkip    bool
	}{
		{
			name:        "eBird species matches iNaturalist subspecies (issue #7)",
			ebirdCommon: "Eurasian Jay",
			ebirdSci:    "Garrulus glandarius",
			inatCommon:  "White-faced Jay",
			inatSci:     "Garrulus glandarius leucotis",
			wantSkip:    true,
		},
		{
			name:        "eBird slash-subspecies group matches iNaturalist subspecies (issue #7 workaround)",
			ebirdCommon: "Eurasian Jay (White-faced)",
			ebirdSci:    "Garrulus glandarius leucotis/oatesi",
			inatCommon:  "White-faced Jay",
			inatSci:     "Garrulus glandarius leucotis",
			wantSkip:    true,
		},
		{
			name:        "eBird bracketed group matches iNaturalist subspecies",
			ebirdCommon: "Dark-eyed Junco (Oregon)",
			ebirdSci:    "Junco hyemalis [oreganus Group]",
			inatCommon:  "Oregon Junco",
			inatSci:     "Junco hyemalis shufeldti",
			wantSkip:    true,
		},
		{
			name:        "eBird domestic form matches iNaturalist species",
			ebirdCommon: "Rock Pigeon (Feral Pigeon)",
			ebirdSci:    "Columba livia (Feral Pigeon)",
			inatCommon:  "Rock Pigeon",
			inatSci:     "Columba livia",
			wantSkip:    true,
		},
		{
			name:        "different species in same genus does not match",
			ebirdCommon: "Black-headed Jay",
			ebirdSci:    "Garrulus lanceolatus",
			inatCommon:  "White-faced Jay",
			inatSci:     "Garrulus glandarius leucotis",
			wantSkip:    false,
		},
		{
			name:        "genus-level iNaturalist identification does not match species",
			ebirdCommon: "Eurasian Jay",
			ebirdSci:    "Garrulus glandarius",
			inatCommon:  "Typical Jays",
			inatSci:     "Garrulus",
			wantSkip:    false,
		},
		{
			name:        "eBird spuh does not match species in same genus",
			ebirdCommon: "Empidonax sp.",
			ebirdSci:    "Empidonax sp.",
			inatCommon:  "Least Flycatcher",
			inatSci:     "Empidonax minimus",
			wantSkip:    false,
		},
		{
			name:        "eBird species slash does not match single species",
			ebirdCommon: "Iceland/Thayer's Gull",
			ebirdSci:    "Larus glaucoides/thayeri",
			inatCommon:  "Iceland Gull",
			inatSci:     "Larus glaucoides",
			wantSkip:    false,
		},
		{
			name:        "eBird hybrid with x does not match parent species",
			ebirdCommon: "Mallard x Northern Pintail (hybrid)",
			ebirdSci:    "Anas platyrhynchos x acuta",
			inatCommon:  "Mallard",
			inatSci:     "Anas platyrhynchos",
			wantSkip:    false,
		},
		{
			name:        "iNaturalist hybrid with multiplication sign does not match parent species",
			ebirdCommon: "Mallard",
			ebirdSci:    "Anas platyrhynchos",
			inatCommon:  "Mallard × Northern Pintail",
			inatSci:     "Anas platyrhynchos × acuta",
			wantSkip:    false,
		},
		{
			name:        "subgenus parenthetical in second token does not match another species in subgenus",
			ebirdCommon: "Herring Gull",
			ebirdSci:    "Larus (Larus) argentatus",
			inatCommon:  "Lesser Black-backed Gull",
			inatSci:     "Larus (Larus) fuscus",
			wantSkip:    false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetFlags()
			fuzzy = true

			mockEbird := &mockEBirdClient{records: []ebird.Record{{
				Line:             2,
				SubmissionID:     "S212386574",
				CommonName:       tc.ebirdCommon,
				ScientificName:   tc.ebirdSci,
				Date:             "2025-02-01",
				Time:             "09:30 AM",
				MLCatalogNumbers: "99999",
			}}}
			mockInat := &mockINatClient{observations: []inat.Result{{
				ID:         100,
				UUID:       uuid.New(),
				ObservedOn: "2025-02-01",
				Taxon: inat.Taxon{
					PreferredCommonName: tc.inatCommon,
					Name:                tc.inatSci,
				},
			}}}

			s := birdsync("MyEBirdData.csv", mockEbird, "myUserID", mockInat)
			if tc.wantSkip {
				if s.fuzzySkips != 1 || s.createdObservations != 0 {
					t.Errorf("fuzzySkips=%d createdObservations=%d, want 1 and 0", s.fuzzySkips, s.createdObservations)
				}
			} else {
				if s.fuzzySkips != 0 || s.createdObservations != 1 {
					t.Errorf("fuzzySkips=%d createdObservations=%d, want 0 and 1", s.fuzzySkips, s.createdObservations)
				}
			}
		})
	}
}
