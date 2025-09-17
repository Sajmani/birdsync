package inat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDownloadObservations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
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

	results := DownloadObservations(server.URL, "testuser", time.Time{}, time.Time{})
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
