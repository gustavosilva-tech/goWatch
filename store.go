package main

import (
	"sync"
	"time"
)

type Summary struct {
	Total   int `json:"total"`
	Online  int `json:"online"`
	Offline int `json:"offline"`
}

type Snapshot struct {
	UpdatedAt *time.Time    `json:"updated_at"`
	Summary   Summary       `json:"summary"`
	Results   []CheckResult `json:"results"`
}

type Store struct {
	mu        sync.RWMutex
	latest    []CheckResult
	updatedAt time.Time
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(results []CheckResult) {
	copied := make([]CheckResult, len(results))
	copy(copied, results)

	s.mu.Lock()
	s.latest = copied
	s.updatedAt = time.Now().UTC()
	s.mu.Unlock()
}

func (s *Store) Latest() []CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copied := make([]CheckResult, len(s.latest))
	copy(copied, s.latest)

	return copied
}

// Snapshot returns the latest results with a summary. UpdatedAt is nil
// until the first round of checks has completed.
func (s *Store) Snapshot() Snapshot {
	s.mu.RLock()
	updatedAt := s.updatedAt
	s.mu.RUnlock()

	results := s.Latest()
	snap := Snapshot{
		Summary: summarize(results),
		Results: results,
	}

	if !updatedAt.IsZero() {
		snap.UpdatedAt = &updatedAt
	}

	return snap
}

func summarize(results []CheckResult) Summary {
	summary := Summary{Total: len(results)}

	for _, result := range results {
		if result.Online {
			summary.Online++
		} else {
			summary.Offline++
		}
	}

	return summary
}
