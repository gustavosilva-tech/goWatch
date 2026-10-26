package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunSchedulerRunsRepeatedlyAndStops(t *testing.T) {
	var hits atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			w.WriteHeader(http.StatusOK)
		},
	))
	defer server.Close()

	store := NewStore()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		defer close(done)
		runScheduler(ctx, server.Client(), []string{server.URL}, 1, 10*time.Millisecond, store)
	}()

	deadline := time.After(2 * time.Second)
	for hits.Load() < 3 {
		select {
		case <-deadline:
			t.Fatalf("scheduler ran only %d times", hits.Load())
		case <-time.After(5 * time.Millisecond):
		}
	}

	cancel()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("scheduler did not stop after cancel")
	}

	if latest := store.Latest(); len(latest) != 1 || !latest[0].Online {
		t.Errorf("unexpected store contents: %+v", latest)
	}
}
