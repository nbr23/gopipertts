package main

import (
	"testing"
	"time"
)

func TestExpireTTSRequests_RemovesExpiredEntry(t *testing.T) {
	r := map[string]TTSRequestStore{
		"expired": {Expires: time.Now().Add(-1 * time.Second)},
	}
	expireTTSRequests(r)
	if len(r) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(r))
	}
}

func TestExpireTTSRequests_KeepsFutureEntry(t *testing.T) {
	r := map[string]TTSRequestStore{
		"future": {Expires: time.Now().Add(10 * time.Minute)},
	}
	expireTTSRequests(r)
	if len(r) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r))
	}
}

func TestExpireTTSRequests_MixedEntries(t *testing.T) {
	r := map[string]TTSRequestStore{
		"expired": {Expires: time.Now().Add(-1 * time.Second)},
		"future":  {Expires: time.Now().Add(10 * time.Minute)},
	}
	expireTTSRequests(r)
	if len(r) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r))
	}
	if _, ok := r["future"]; !ok {
		t.Fatal("expected 'future' entry to remain")
	}
}

func TestExpireTTSRequests_EmptyMap(t *testing.T) {
	r := map[string]TTSRequestStore{}
	expireTTSRequests(r)
	if len(r) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(r))
	}
}
