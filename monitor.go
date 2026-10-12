package main

import (
	"net/http"
	"sync"
)

type checkJob struct {
	index int
	url   string
}

type indexedResult struct {
	index  int
	result CheckResult
}

func checkURLs(
	client *http.Client,
	targets []string,
	workers int,
) []CheckResult {
	results := make([]CheckResult, len(targets))

	if len(targets) == 0 {
		return results
	}

	if workers < 1 {
		workers = 1
	}

	if workers > len(targets) {
		workers = len(targets)
	}

	jobs := make(chan checkJob)
	completed := make(chan indexedResult)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for job := range jobs {
				completed <- indexedResult{
					index:  job.index,
					result: checkURL(client, job.url),
				}
			}
		}()
	}

	go func() {
		for i, target := range targets {
			jobs <- checkJob{
				index: i,
				url:   target,
			}
		}

		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(completed)
	}()

	for item := range completed {
		results[item.index] = item.result
	}

	return results
}
