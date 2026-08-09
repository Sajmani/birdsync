package inat

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

// newTestClient returns a client with pacing switched off. Only the pacing
// tests want a real delay between requests; everywhere else it would add
// minutes to the suite and cover nothing.
func newTestClient(baseURL, apiToken, userAgent string) *Client {
	c := NewClient(baseURL, apiToken, userAgent)
	c.SetMinRequestInterval(0)
	return c
}

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

	client := newTestClient(server.URL, "test-token", "test-user-agent")

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

	client := newTestClient(server.URL, "test-token", "test-user-agent")

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

	client := newTestClient(server.URL, "test-token", "test-user-agent")

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

	client := newTestClient(server.URL, "test-token", "test-user-agent")
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

// TestStatusErrorDropsHTMLBody checks that a proxy's error page doesn't end up
// in the log. iNaturalist's own refusals are short text; a 413 comes from nginx
// as a seven-line HTML document that says nothing the status line doesn't.
//
// Verifies: T-034.
func TestStatusErrorDropsHTMLBody(t *testing.T) {
	const nginx = `<html>
<head><title>413 Request Entity Too Large</title></head>
<body>
<center><h1>413 Request Entity Too Large</h1></center>
<hr><center>nginx</center>
</body>
</html>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		w.Write([]byte(nginx))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "test-user-agent")
	err := client.CreateObservation(Observation{})
	if err == nil {
		t.Fatal("CreateObservation() against a refusing server returned no error")
	}
	if strings.Contains(err.Error(), "<html") || strings.Contains(err.Error(), "\n") {
		t.Errorf("Error carries the proxy's HTML page (T-034): %q", err)
	}
	if !strings.Contains(err.Error(), "413") {
		t.Errorf("Error %q doesn't name the status", err)
	}
}

// TestStatusErrorCollapsesWhitespace keeps a genuine multi-line message on one
// log line.
//
// Verifies: T-034.
func TestStatusErrorCollapsesWhitespace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("{\n  \"error\": \"File is too large\"\n}"))
	}))
	defer server.Close()

	client := newTestClient(server.URL, "test-token", "test-user-agent")
	err := client.CreateObservation(Observation{})
	if err == nil {
		t.Fatal("Expected an error")
	}
	if strings.Contains(err.Error(), "\n") {
		t.Errorf("Error spans multiple lines: %q", err)
	}
	if !strings.Contains(err.Error(), "File is too large") {
		t.Errorf("Error %q lost the explanation", err)
	}
}

// TestClientPacesRequests checks that requests are spaced out. iNaturalist asks
// for about one per second and warns that it may block IPs that persistently
// exceed that; a first sync of a media-heavy account is thousands of requests,
// and nothing else in birdsync slows them down.
//
// Verifies: T-035.
func TestClientPacesRequests(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// A real second per request would make this test unbearable, so pace at a
	// scaled-down interval and check the spacing rather than the constant.
	const interval = 40 * time.Millisecond
	client := NewClient(server.URL, "", "")
	client.SetMinRequestInterval(interval)

	start := time.Now()
	for range 3 {
		if err := client.CreateObservation(Observation{}); err != nil {
			t.Fatalf("CreateObservation() error = %v", err)
		}
	}
	elapsed := time.Since(start)

	if requests != 3 {
		t.Fatalf("Server saw %d requests, want 3", requests)
	}
	// Three requests means two gaps; the first goes out immediately.
	if min := 2 * interval; elapsed < min {
		t.Errorf("3 requests took %v, want at least %v: they are not being paced (T-035)", elapsed, min)
	}
}

// TestDefaultMinRequestIntervalMatchesGuidance pins the constant to the rate
// the governing source asks for, so a future change to it has to be deliberate.
//
// Verifies: T-035.
func TestDefaultMinRequestIntervalMatchesGuidance(t *testing.T) {
	if DefaultMinRequestInterval != time.Second {
		t.Errorf("DefaultMinRequestInterval = %v, want 1s: iNaturalist asks for about one request per second",
			DefaultMinRequestInterval)
	}
	if NewClient("http://example.invalid", "", "").minRequestInterval != DefaultMinRequestInterval {
		t.Error("NewClient does not pace by default; a caller that forgets to set it would run unthrottled")
	}
}
