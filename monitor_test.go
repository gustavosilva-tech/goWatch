package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckURLs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/health" {
				w.WriteHeader(http.StatusOK)
				return
			}

			w.WriteHeader(http.StatusNotFound)
		},
	))
	defer server.Close()

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	targets := []string{
		server.URL + "/health",
		server.URL + "/missing",
	}

	results := checkURLs(context.Background(), client, targets, 2)

	if len(results) != len(targets) {
		t.Fatalf(
			"got %d results; want %d",
			len(results),
			len(targets),
		)
	}

	for i, result := range results {
		if result.URL != targets[i] {
			t.Errorf(
				"result %d URL = %q; want %q",
				i,
				result.URL,
				targets[i],
			)
		}

		if result.Error != "" {
			t.Errorf("result %d unexpected error: %s", i, result.Error)
		}
	}

	if !results[0].Online || results[0].StatusCode != http.StatusOK {
		t.Errorf("unexpected health result: %+v", results[0])
	}

	if results[1].Online || results[1].StatusCode != http.StatusNotFound {
		t.Errorf("unexpected missing-page result: %+v", results[1])
	}
}
