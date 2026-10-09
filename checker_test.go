package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckURLStatus(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantOnline bool
	}{
		{
			name:       "200 is online",
			statusCode: http.StatusOK,
			wantOnline: true,
		},
		{
			name:       "404 is offline",
			statusCode: http.StatusNotFound,
			wantOnline: false,
		},
		{
			name:       "500 is offline",
			statusCode: http.StatusInternalServerError,
			wantOnline: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.statusCode)
				},
			))
			defer server.Close()

			client := &http.Client{
				Timeout: 2 * time.Second,
			}

			result := checkURL(client, server.URL)

			if result.Online != tt.wantOnline {
				t.Errorf(
					"Online = %v; want %v",
					result.Online,
					tt.wantOnline,
				)
			}

			if result.StatusCode != tt.statusCode {
				t.Errorf(
					"StatusCode = %d; want %d",
					result.StatusCode,
					tt.statusCode,
				)
			}

			if result.Error != "" {
				t.Errorf("unexpected request error: %s", result.Error)
			}

			if result.URL != server.URL {
				t.Errorf("URL = %q; want %q", result.URL, server.URL)
			}

			if result.CheckedAt.IsZero() {
				t.Error("CheckedAt should not be empty")
			}
		})
	}
}

func TestCheckURLTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			select {
			case <-time.After(500 * time.Millisecond):
				w.WriteHeader(http.StatusOK)
			case <-r.Context().Done():
				return
			}
		},
	))
	defer server.Close()

	client := &http.Client{
		Timeout: 50 * time.Millisecond,
	}

	result := checkURL(client, server.URL)

	if result.Online {
		t.Error("expected Online to be false after timeout")
	}

	if result.Error == "" {
		t.Error("expected a request error after timeout")
	}

	if result.StatusCode != 0 {
		t.Errorf("StatusCode = %d; want 0", result.StatusCode)
	}
}
