package main

import "sync"

type Store struct {
	mu     sync.RWMutex
	latest []CheckResult
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Set(results []CheckResult) {
	copied := make([]CheckResult, len(results))
	copy(copied, results)

	s.mu.Lock()
	s.latest = copied
	s.mu.Unlock()
}

func (s *Store) Latest() []CheckResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	copied := make([]CheckResult, len(s.latest))
	copy(copied, s.latest)

	return copied
}
