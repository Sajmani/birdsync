package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sajmani/birdsync/inat"
	"github.com/google/uuid"
)

// syncKeyed returns a Result carrying a complete birdsync sync key, as an
// observation birdsync created would.
func syncKeyed(taxonName, iconic, submissionID string) inat.Result {
	return inat.Result{
		UUID: uuid.New(),
		Taxon: inat.Taxon{
			Name:            taxonName,
			IconicTaxonName: iconic,
		},
		Ofvs: []inat.Ofv{
			{FieldID: inat.EBirdField, Value: submissionID},
			{FieldID: inat.EBirdScientificNameField, Value: taxonName + " (eBird)"},
		},
	}
}

// manual returns a Result with no birdsync observation fields.
func manual(taxonName, iconic string) inat.Result {
	return inat.Result{
		UUID:  uuid.New(),
		Taxon: inat.Taxon{Name: taxonName, IconicTaxonName: iconic},
	}
}

func TestAnalyze(t *testing.T) {
	// A birdsync observation iNaturalist couldn't resolve, so it belongs to no
	// iconic taxon. This is the population CR-003 is about. Whether iNaturalist
	// represents it as a null taxon or as a "Life"/"Unknown" placeholder is
	// unsettled, so both shapes are exercised below.
	untaxonedBirdsync := syncKeyed("", "", "S100")
	untaxonedBirdsync.Taxon = inat.Taxon{} // no taxon at all
	lifeBirdsync := syncKeyed("Life", "", "S102")
	lifeBirdsync.Taxon = inat.Taxon{Name: "Life", ID: 48460, IconicTaxonName: "Unknown"}
	resolvedBirdsync := syncKeyed("Corvus brachyrhynchos", "Aves", "S101")
	manualBird := manual("Turdus migratorius", "Aves")
	plant := manual("Quercus agrifolia", "Plantae")

	tests := []struct {
		name                string
		all                 []inat.Result
		filtered            []inat.Result
		wantTotal           int
		wantUnidentified    int
		wantNonAves         int
		wantBirdsync        int
		wantMissing         int
		wantMissingUnid     int
		wantMissingBirdsync int
	}{
		{
			name:                "filter drops an untaxoned birdsync observation",
			all:                 []inat.Result{untaxonedBirdsync, resolvedBirdsync, manualBird},
			filtered:            []inat.Result{resolvedBirdsync, manualBird},
			wantTotal:           3,
			wantUnidentified:    1,
			wantBirdsync:        2,
			wantMissing:         1,
			wantMissingUnid:     1,
			wantMissingBirdsync: 1,
		},
		{
			// Same situation, but with iNaturalist reporting a "Life"
			// placeholder instead of a null taxon. The conclusion must not
			// depend on which representation the API uses.
			name:                "filter drops a Life-placeholder birdsync observation",
			all:                 []inat.Result{lifeBirdsync, resolvedBirdsync},
			filtered:            []inat.Result{resolvedBirdsync},
			wantTotal:           2,
			wantUnidentified:    1,
			wantBirdsync:        2,
			wantMissing:         1,
			wantMissingUnid:     1,
			wantMissingBirdsync: 1,
		},
		{
			// The observed shape of the maintainer's account: the filter drops
			// only non-Aves observations and there are no unidentified ones, so
			// nothing can be concluded either way.
			name:            "inconclusive: nothing unidentified to hide",
			all:             []inat.Result{resolvedBirdsync, manualBird, plant},
			filtered:        []inat.Result{resolvedBirdsync, manualBird},
			wantTotal:       3,
			wantNonAves:     1,
			wantBirdsync:    1,
			wantMissing:     1,
			wantMissingUnid: 0,
		},
		{
			name:         "filter returns everything",
			all:          []inat.Result{resolvedBirdsync, manualBird},
			filtered:     []inat.Result{resolvedBirdsync, manualBird},
			wantTotal:    2,
			wantBirdsync: 1,
		},
		{
			name:        "filter drops only a non-birdsync observation",
			all:         []inat.Result{resolvedBirdsync, plant},
			filtered:    []inat.Result{resolvedBirdsync},
			wantTotal:   2,
			wantNonAves: 1,
			// The plant is legitimately excluded and birdsync never created
			// it, so this is the filter working as intended, not CR-003.
			wantBirdsync:        1,
			wantMissing:         1,
			wantMissingBirdsync: 0,
		},
		{
			name:     "empty account",
			all:      nil,
			filtered: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyze(tt.all, tt.filtered)
			if got.total != tt.wantTotal {
				t.Errorf("total = %d, want %d", got.total, tt.wantTotal)
			}
			if got.aves != len(tt.filtered) {
				t.Errorf("aves = %d, want %d", got.aves, len(tt.filtered))
			}
			if got.unidentified != tt.wantUnidentified {
				t.Errorf("unidentified = %d, want %d", got.unidentified, tt.wantUnidentified)
			}
			if len(got.missingUnidentified) != tt.wantMissingUnid {
				t.Errorf("missingUnidentified = %d, want %d", len(got.missingUnidentified), tt.wantMissingUnid)
			}
			if got.nonAves != tt.wantNonAves {
				t.Errorf("nonAves = %d, want %d", got.nonAves, tt.wantNonAves)
			}
			if got.birdsync != tt.wantBirdsync {
				t.Errorf("birdsync = %d, want %d", got.birdsync, tt.wantBirdsync)
			}
			if len(got.missing) != tt.wantMissing {
				t.Errorf("missing = %d, want %d", len(got.missing), tt.wantMissing)
			}
			if len(got.missingBirdsync) != tt.wantMissingBirdsync {
				t.Errorf("missingBirdsync = %d, want %d", len(got.missingBirdsync), tt.wantMissingBirdsync)
			}
		})
	}
}

func TestHasSyncKey(t *testing.T) {
	tests := []struct {
		name string
		ofvs []inat.Ofv
		want bool
	}{
		{"both fields", []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S1"},
			{FieldID: inat.EBirdScientificNameField, Value: "Corvus corax"},
		}, true},
		{"checklist only", []inat.Ofv{
			{FieldID: inat.EBirdField, Value: "S1"},
		}, false},
		{"scientific name only", []inat.Ofv{
			{FieldID: inat.EBirdScientificNameField, Value: "Corvus corax"},
		}, false},
		{"neither", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasSyncKey(inat.Result{Ofvs: tt.ofvs}); got != tt.want {
				t.Errorf("hasSyncKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFetchAllQuery checks the difference the whole tool depends on: the
// iconic_taxa[] parameter is sent when asked for and absent otherwise. If the
// unfiltered fetch quietly sent the filter too, the comparison would always
// report "not confirmed".
func TestFetchAllQuery(t *testing.T) {
	for _, tt := range []struct {
		name       string
		iconicTaxa string
		wantParam  string
		wantSent   bool
	}{
		{"unfiltered", "", "", false},
		{"aves", "Aves", "Aves", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var gotParam string
			var sent bool
			var gotUserAgent, gotAuth, gotPerPage, gotUserID, gotFields string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				gotParam = q.Get("iconic_taxa[]")
				_, sent = q["iconic_taxa[]"]
				gotPerPage = q.Get("per_page")
				gotUserID = q.Get("user_id")
				gotFields = q.Get("fields")
				gotUserAgent = r.Header.Get("User-Agent")
				gotAuth = r.Header.Get("Authorization")
				json.NewEncoder(w).Encode(inat.Observations{TotalResults: 0})
			}))
			defer server.Close()

			if _, err := fetchAll(server.URL, "test-token", "testuser", tt.iconicTaxa); err != nil {
				t.Fatalf("fetchAll() error = %v", err)
			}

			if sent != tt.wantSent {
				t.Errorf("iconic_taxa[] sent = %v, want %v", sent, tt.wantSent)
			}
			if gotParam != tt.wantParam {
				t.Errorf("iconic_taxa[] = %q, want %q", gotParam, tt.wantParam)
			}
			if gotPerPage != "200" {
				t.Errorf("per_page = %q, want %q (T-017)", gotPerPage, "200")
			}
			if gotUserID != "testuser" {
				t.Errorf("user_id = %q, want %q", gotUserID, "testuser")
			}
			if !strings.Contains(gotFields, "taxon.all") {
				t.Errorf("fields = %q, want it to include taxon.all", gotFields)
			}
			if gotUserAgent != UserAgent {
				t.Errorf("User-Agent = %q, want %q (T-016)", gotUserAgent, UserAgent)
			}
			if gotAuth != "test-token" {
				t.Errorf("Authorization = %q, want %q", gotAuth, "test-token")
			}
		})
	}
}

// TestFetchAllPaginates checks that the tool reads past the first page. An
// account large enough to matter here is larger than one page by definition.
func TestFetchAllPaginates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp inat.Observations
		switch r.URL.Query().Get("page") {
		case "1":
			resp = inat.Observations{TotalResults: 3, Results: []inat.Result{
				{Taxon: inat.Taxon{Name: "a"}}, {Taxon: inat.Taxon{Name: "b"}},
			}}
		case "2":
			resp = inat.Observations{TotalResults: 3, Results: []inat.Result{
				{Taxon: inat.Taxon{Name: "c"}},
			}}
		default:
			t.Errorf("Unexpected page %q", r.URL.Query().Get("page"))
			resp = inat.Observations{TotalResults: 3}
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	results, err := fetchAll(server.URL, "", "testuser", "")
	if err != nil {
		t.Fatalf("fetchAll() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("Got %d results, want 3", len(results))
	}
	for i, want := range []string{"a", "b", "c"} {
		if results[i].Taxon.Name != want {
			t.Errorf("results[%d].Taxon.Name = %q, want %q", i, results[i].Taxon.Name, want)
		}
	}
}

// TestFetchAllUnauthorized checks the message a user is most likely to see,
// since the API token expires every 24 hours (P-017).
func TestFetchAllUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	_, err := fetchAll(server.URL, "stale-token", "testuser", "")
	if err == nil {
		t.Fatal("fetchAll() with a 401 returned no error")
	}
	if !strings.Contains(err.Error(), "refresh your INAT_API_TOKEN") {
		t.Errorf("Error %q doesn't tell the user to refresh their token", err)
	}
}
