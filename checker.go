package main

import (
	"net/http"
	"time"
)

type CheckResult struct {
	URL          string    `json:"url"`
	Online       bool      `json:"online"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	CheckedAt    time.Time `json:"checked_at"`
	Error        string    `json:"error,omitempty"`
}

func checkURL(client *http.Client, target string) CheckResult {
	result := CheckResult{
		URL:       target,
		CheckedAt: time.Now().UTC(),
	}

	start := time.Now()
	resp, err := client.Get(target)
	result.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Online = resp.StatusCode >= 200 && resp.StatusCode < 300

	return result
}
