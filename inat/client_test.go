package inat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Verifies: T-016.
func TestClient_CreateObservation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/observations" {
			t.Errorf("Expected path /observations, got %s", r.URL.Path)
		}
		// Every request identifies the client and carries the token
		// (T-016). iNaturalist may throttle a client that omits either.
		if got := r.Header.Get("User-Agent"); got != "test-user-agent" {
			t.Errorf("User-Agent = %q, want %q (T-016)", got, "test-user-agent")
		}
		if got := r.Header.Get("Authorization"); got != "test-token" {
			t.Errorf("Authorization = %q, want %q", got, "test-token")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", "test-user-agent")

	obs := Observation{UUID: uuid.New()}
	if err := client.CreateObservation(obs); err != nil {
		t.Errorf("CreateObservation() error = %v", err)
	}
}

// TestClient_UpdateObservation checks that an update always sets
// ignore_photos. Without it, iNaturalist treats the absent photo list as an
// instruction to detach the observation's photos — so birdsync appending an
// asset URL to a description would silently destroy the media it just
// uploaded. Nothing asserted this before.
//
// Verifies: T-009, P-046.
func TestClient_UpdateObservation(t *testing.T) {
	obsUUID := uuid.New()
	var body UpdateObservation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/observations/%s", obsUUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Decoding request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", "test-user-agent")

	obs := Observation{UUID: obsUUID, Description: "updated"}
	if err := client.UpdateObservation(obs); err != nil {
		t.Errorf("UpdateObservation() error = %v", err)
	}
	if !body.IgnorePhotos {
		t.Error("ignore_photos = false, want true: an update must not detach attached media (T-009)")
	}
	if body.Observation.Description != "updated" {
		t.Errorf("Description = %q, want %q", body.Observation.Description, "updated")
	}
}

func TestClient_DeleteObservation(t *testing.T) {
	obsUUID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/observations/%s", obsUUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", "test-user-agent")

	if err := client.DeleteObservation(obsUUID); err != nil {
		t.Errorf("DeleteObservation() error = %v", err)
	}
}

// TestStatusErrorIncludesBody checks that a refusal carries iNaturalist's
// explanation. Without it every failure read "bad HTTP status: 422
// Unprocessable Entity", whether the file was too large, the format
// unsupported, or the asset withdrawn — and the user had no way to tell which.
//
// Verifies: T-034.
func TestStatusErrorIncludesBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "File is too large (max 50 MB)", http.StatusUnprocessableEntity)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", "test-user-agent")
	err := client.CreateObservation(Observation{})
	if err == nil {
		t.Fatal("CreateObservation() against a refusing server returned no error")
	}
	if !strings.Contains(err.Error(), "File is too large") {
		t.Errorf("Error %q doesn't include the server's explanation (T-034)", err)
	}

	var statusErr *StatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("Error %v is not a *StatusError, so callers can't classify it (P-063)", err)
	}
	if !statusErr.Permanent() {
		t.Error("422 should be permanent: the service rejected the file itself (P-063)")
	}
}

// TestStatusErrorPermanence pins the classification down case by case. Getting
// it wrong in one direction retries forever; in the other it discards a user's
// photo on a transient error.
//
// Verifies: P-063.
func TestStatusErrorPermanence(t *testing.T) {
	for _, tt := range []struct {
		code int
		want bool
	}{
		{http.StatusUnprocessableEntity, true},   // 422: file rejected
		{http.StatusRequestEntityTooLarge, true}, // 413: too big
		{http.StatusUnsupportedMediaType, true},  // 415: wrong format
		{http.StatusNotFound, true},              // 404: gone
		{http.StatusUnauthorized, false},         // 401: refresh the token
		{http.StatusRequestTimeout, false},       // 408: try again
		{http.StatusTooManyRequests, false},      // 429: explicitly try later
		{http.StatusInternalServerError, false},  // 5xx: their problem, not the file's
		{http.StatusServiceUnavailable, false},   //
		{http.StatusGatewayTimeout, false},       //
	} {
		e := &StatusError{StatusCode: tt.code}
		if got := e.Permanent(); got != tt.want {
			t.Errorf("StatusError{%d}.Permanent() = %v, want %v", tt.code, got, tt.want)
		}
	}
}
