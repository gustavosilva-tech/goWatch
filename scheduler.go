package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

func runScheduler(
	ctx context.Context,
	client *http.Client,
	targets []string,
	workers int,
	interval time.Duration,
	store *Store,
) {
	runOnce := func() {
		results := checkURLs(client, targets, workers)
		store.Set(results)
		log.Printf("Completed %d checks", len(results))
	}

	runOnce()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped")
			return
		case <-ticker.C:
			runOnce()
		}
	}
}
