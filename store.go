package main

import (
	"sync"
	"time"
)

type TTSRequestStore struct {
	Request TTSRequestInput
	Expires time.Time
}

type TTSRequestsStore struct {
	mu      sync.RWMutex
	entries map[string]TTSRequestStore
}

func (s *TTSRequestsStore) set(id string, v TTSRequestStore) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[id] = v
}

func (s *TTSRequestsStore) get(id string) (TTSRequestStore, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.entries[id]
	return v, ok
}

func (s *TTSRequestsStore) delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, id)
}

func (s *TTSRequestsStore) expireOld() {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range s.entries {
		if value.Expires.Before(now) {
			delete(s.entries, key)
		}
	}
}

func initTTSRequestsStore() *TTSRequestsStore {
	s := &TTSRequestsStore{entries: make(map[string]TTSRequestStore)}
	go func() {
		for {
			s.expireOld()
			time.Sleep(15 * time.Minute)
		}
	}()
	return s
}
