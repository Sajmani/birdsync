package inat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestDownloadObservations checks the basics of a download: the request
// identifies the user and asks for the largest page size, and every result is
// returned. Paging itself is covered by TestDownloadObservationsPagesByID.
//
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
		json.NewEncoder(w).Encode(Observations{
			TotalResults: 2,
			Results:      []Result{{ID: 1, Description: "obs 1"}, {ID: 2, Description: "obs 2"}},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	results, err := client.DownloadObservations("testuser", time.Time{}, time.Time{})
	if err != nil {
		t.Fatalf("DownloadObservations() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	if results[0].Description != "obs 1" || results[1].Description != "obs 2" {
		t.Errorf("Got %q, %q", results[0].Description, results[1].Description)
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
	if _, err := client.DownloadObservations("testuser", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("DownloadObservations() error = %v", err)
	}

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
	if _, err := client.DownloadObservations("testuser", d1, d2); err != nil {
		t.Fatalf("DownloadObservations() error = %v", err)
	}

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
	if _, err := client.DownloadObservations("testuser", time.Time{}, time.Time{}); err != nil {
		t.Fatalf("DownloadObservations() error = %v", err)
	}

	if hasD1 || hasD2 {
		t.Errorf("Unset date window sent d1=%v d2=%v, want neither (P-025)", hasD1, hasD2)
	}
}

// TestDownloadObservationsError checks that a failed download returns an error
// rather than ending the process. It used to call log.Fatal, which meant a
// transient server error killed a sync that had already created observations,
// and killed the test binary in any test that provoked it.
//
// Verifies: T-027.
func TestDownloadObservationsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	results, err := client.DownloadObservations("testuser", time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("DownloadObservations() with a failing server returned no error")
	}
	if results != nil {
		t.Errorf("DownloadObservations() returned %d results alongside an error", len(results))
	}
}

// TestDownloadObservationsPagesByID checks that the download walks the result
// set with id_above rather than page numbers.
//
// iNaturalist refuses to page past 10,000 results — "An error will be thrown if
// results beyond 10k are requested" — so page-number paging fails outright for
// a large account, which is what issue #5 reports. The recommended alternative
// is to sort by id ascending and pass the last id seen as id_above, which has
// no ceiling.
//
// The fake honours per_page, as the real API does, and serves 450 results so
// the cursor has to advance twice and then stop on a short page. An earlier
// version of this test imposed its own page size and so made the client look
// broken when it was the fake breaking the contract.
//
// Verifies: P-065, T-036.
func TestDownloadObservationsPagesByID(t *testing.T) {
	const total = 450
	var all []Result
	for i := 1; i <= total; i++ {
		all = append(all, Result{ID: i * 1000, Description: fmt.Sprintf("obs %d", i)})
	}

	var sawPageParam bool
	var idAboves []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "" {
			sawPageParam = true
		}
		if q.Get("order_by") != "id" || q.Get("order") != "asc" {
			t.Errorf("order_by=%q order=%q, want id/asc: id_above is meaningless without it",
				q.Get("order_by"), q.Get("order"))
		}
		if !strings.Contains(q.Get("fields"), "id") {
			t.Errorf("fields=%q must request id, or the cursor can't be read", q.Get("fields"))
		}
		perPage, err := strconv.Atoi(q.Get("per_page"))
		if err != nil {
			t.Fatalf("per_page=%q is not a number", q.Get("per_page"))
		}
		above := q.Get("id_above")
		idAboves = append(idAboves, above)

		var start int
		if above != "" {
			n, err := strconv.Atoi(above)
			if err != nil {
				t.Errorf("id_above=%q is not a number", above)
			}
			for start < len(all) && all[start].ID <= n {
				start++
			}
		}
		end := min(start+perPage, len(all))
		// total_results shrinks as the cursor advances, as the real API's does.
		json.NewEncoder(w).Encode(Observations{
			TotalResults: len(all) - start,
			Results:      all[start:end],
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	got, err := client.DownloadObservations("testuser", time.Time{}, time.Time{}, "description")
	if err != nil {
		t.Fatalf("DownloadObservations() error = %v", err)
	}

	if sawPageParam {
		t.Error("Download sent a page parameter; page-number paging fails past 10k results (T-036)")
	}
	if len(got) != total {
		t.Fatalf("Got %d observations, want %d; cursors were %v", len(got), total, idAboves)
	}
	for i := range all {
		if got[i].ID != all[i].ID {
			t.Fatalf("result[%d].ID = %d, want %d (order or cursor wrong)", i, got[i].ID, all[i].ID)
		}
	}
	// First request has no cursor; each later one resumes from the last id seen.
	want := []string{"", "200000", "400000"}
	if !slices.Equal(idAboves, want) {
		t.Errorf("id_above sequence = %v, want %v", idAboves, want)
	}
}

// TestDownloadObservationsCursorMustAdvance checks that a server which ignores
// id_above produces an error rather than an infinite loop.
//
// This was found by mutation-testing the paging change: removing the cursor
// didn't fail the test, it hung, because a full page always looks like "there
// is more". In production that would re-download the same page forever and
// hammer a service whose recommended practices ask for one request a second.
//
// The fake caps its responses so a regression here fails fast instead of
// wedging the suite.
//
// Verifies: T-036.
func TestDownloadObservationsCursorMustAdvance(t *testing.T) {
	const cap = 10
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests > cap {
			t.Errorf("Client made %d requests without advancing; it is looping (T-036)", requests)
			json.NewEncoder(w).Encode(Observations{TotalResults: 0})
			return
		}
		// A full page, always the same: id_above ignored.
		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		results := make([]Result, perPage)
		for i := range results {
			results[i] = Result{ID: i + 1}
		}
		json.NewEncoder(w).Encode(Observations{TotalResults: 100000, Results: results})
	}))
	defer server.Close()

	client := NewClient(server.URL, "", "")
	_, err := client.DownloadObservations("testuser", time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("Expected an error when the cursor doesn't advance, got nil")
	}
	if !strings.Contains(err.Error(), "cursor did not advance") {
		t.Errorf("Error %q doesn't explain that paging stalled", err)
	}
	if requests > 2 {
		t.Errorf("Made %d requests before giving up, want 2", requests)
	}
}
