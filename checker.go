package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

// maxDrainBytes limits how much of a response body is read so the
// connection can be reused without downloading huge pages.
const maxDrainBytes = 64 << 10

type CheckResult struct {
	URL          string    `json:"url"`
	Online       bool      `json:"online"`
	StatusCode   int       `json:"status_code"`
	ResponseTime int64     `json:"response_time_ms"`
	CheckedAt    time.Time `json:"checked_at"`
	Error        string    `json:"error,omitempty"`
}

func checkURL(ctx context.Context, client *http.Client, target string) CheckResult {
	result := CheckResult{
		URL:       target,
		CheckedAt: time.Now().UTC(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	req.Header.Set("User-Agent", "gowatch/1.0")

	start := time.Now()
	resp, err := client.Do(req)
	result.ResponseTime = time.Since(start).Milliseconds()

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, maxDrainBytes))

	result.StatusCode = resp.StatusCode
	result.Online = isOnline(resp.StatusCode)

	return result
}

// isOnline decides whether a status code counts as "up". Redirects are
// followed by the client, so only the final status code reaches here.
func isOnline(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
