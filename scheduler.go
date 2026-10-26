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
		results := checkURLs(ctx, client, targets, workers)
		if ctx.Err() != nil {
			return
		}

		store.Set(results)
		summary := summarize(results)
		log.Printf("Completed %d checks: %d online, %d offline",
			summary.Total, summary.Online, summary.Offline)
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
