package inat

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestClient_CreateObservation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/observations" {
			t.Errorf("Expected path /observations, got %s", r.URL.Path)
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

func TestClient_UpdateObservation(t *testing.T) {
	obsUUID := uuid.New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/observations/%s", obsUUID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token", "test-user-agent")

	obs := Observation{UUID: obsUUID}
	if err := client.UpdateObservation(obs); err != nil {
		t.Errorf("UpdateObservation() error = %v", err)
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
