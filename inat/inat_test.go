package inat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Verifies: T-017, P-023.
func TestDownloadObservations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		// The API's recommended practices require the largest supported page
		// size; fetching 30 at a time makes a large account many times more
		// expensive for iNaturalist to serve (T-017).
		if got := q.Get("per_page"); got != "200" {
			t.Errorf("per_page = %q, want %q (T-017)", got, "200")
		}
		if got := q.Get("user_id"); got != "testuser" {
			t.Errorf("user_id = %q, want %q", got, "testuser")
		}
		page := q.Get("page")
		switch page {
		case "1":
			resp := Observations{
				TotalResults: 2,
				Page:         1,
				PerPage:      1,
				Results:      []Result{{Description: "obs 1"}},
			}
			json.NewEncoder(w).Encode(resp)
		case "2":
			resp := Observations{
				TotalResults: 2,
				Page:         2,
				PerPage:      1,
				Results:      []Result{{Description: "obs 2"}},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			resp := Observations{
				TotalResults: 2,
				Page:         3,
				PerPage:      1,
				Results:      []Result{},
			}
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	results := client.DownloadObservations("testuser", time.Time{}, time.Time{})
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
	if results[0].Description != "obs 1" {
		t.Errorf("Expected obs 1, got %s", results[0].Description)
	}
	if results[1].Description != "obs 2" {
		t.Errorf("Expected obs 2, got %s", results[1].Description)
	}
}

// TestDownloadObservationsNoTaxonFilter checks that the download is not
// restricted by taxon. Filtering on iconic_taxa[]=Aves hid any observation
// without an iconic taxon — the state an unresolvable eBird name produces — so
// birdsync couldn't see its own work and created it again on the next run.
//
// Verifies: P-061.
func TestDownloadObservationsNoTaxonFilter(t *testing.T) {
	var hasIconic bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasIconic = r.URL.Query()["iconic_taxa[]"]
		json.NewEncoder(w).Encode(Observations{TotalResults: 0})
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	client.DownloadObservations("testuser", time.Time{}, time.Time{})

	if hasIconic {
		t.Error("Download sent iconic_taxa[]; it must not filter by taxon (P-061)")
	}
}

// TestDownloadObservationsDateWindow checks that --after and --before narrow
// the download itself, not just the local comparison. This is what makes a
// date-limited run cheap, and it is also why duplicate detection only
// considers observations inside the window.
//
// Verifies: P-025.
func TestDownloadObservationsDateWindow(t *testing.T) {
	var gotD1, gotD2 string
	var sawParams bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotD1 = r.URL.Query().Get("d1")
		gotD2 = r.URL.Query().Get("d2")
		sawParams = true
		json.NewEncoder(w).Encode(Observations{TotalResults: 0})
	}))
	defer server.Close()

	d1 := time.Date(2023, 1, 2, 15, 4, 5, 0, time.UTC)
	d2 := time.Date(2023, 3, 4, 6, 7, 8, 0, time.UTC)
	client := NewClient(server.URL, "", "")
	client.DownloadObservations("testuser", d1, d2)

	if !sawParams {
		t.Fatal("Server received no request")
	}
	// The API takes dates, not timestamps, so the time component is dropped.
	if gotD1 != "2023-01-02" {
		t.Errorf("d1 = %q, want %q (P-025)", gotD1, "2023-01-02")
	}
	if gotD2 != "2023-03-04" {
		t.Errorf("d2 = %q, want %q (P-025)", gotD2, "2023-03-04")
	}
}

// TestDownloadObservationsNoDateWindow checks the other half of P-025: with no
// window set, the parameters are absent rather than sent as zero values, which
// the API would read as a date in year 1.
//
// Verifies: P-025.
func TestDownloadObservationsNoDateWindow(t *testing.T) {
	var hasD1, hasD2 bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, hasD1 = r.URL.Query()["d1"]
		_, hasD2 = r.URL.Query()["d2"]
		json.NewEncoder(w).Encode(Observations{TotalResults: 0})
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	client.DownloadObservations("testuser", time.Time{}, time.Time{})

	if hasD1 || hasD2 {
		t.Errorf("Unset date window sent d1=%v d2=%v, want neither (P-025)", hasD1, hasD2)
	}
}
