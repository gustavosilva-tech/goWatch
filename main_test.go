package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux(NewStore()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d; want 200", rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
}

func TestStatusEndpoint(t *testing.T) {
	store := NewStore()
	store.Set([]CheckResult{{URL: "http://a.example", Online: true, StatusCode: 200}})

	rec := httptest.NewRecorder()
	newMux(store).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/status", nil))

	var snap Snapshot
	if err := json.NewDecoder(rec.Body).Decode(&snap); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if len(snap.Results) != 1 || snap.Summary.Online != 1 {
		t.Errorf("unexpected snapshot: %+v", snap)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	rec := httptest.NewRecorder()
	newMux(NewStore()).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/status", nil))

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d; want 405", rec.Code)
	}

	if allow := rec.Header().Get("Allow"); allow == "" {
		t.Error("Allow header should be set")
	}
}
