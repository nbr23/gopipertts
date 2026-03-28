package main

import (
	"testing"
	"time"
)

func TestExpireTTSRequests_RemovesExpiredEntry(t *testing.T) {
	r := initTTSRequestsStore()
	r.set("expired", TTSRequestStore{Expires: time.Now().Add(-1 * time.Second)})
	r.expireOld()
	if _, exists := r.get("expired"); exists {
		t.Fatal("expected expired entry to be removed")
	}
}

func TestExpireTTSRequests_KeepsFutureEntry(t *testing.T) {
	r := initTTSRequestsStore()
	r.set("future", TTSRequestStore{Expires: time.Now().Add(10 * time.Minute)})
	r.expireOld()
	if _, exists := r.get("future"); !exists {
		t.Fatal("expected future entry to remain")
	}
}

func TestExpireTTSRequests_MixedEntries(t *testing.T) {
	r := initTTSRequestsStore()
	r.set("expired", TTSRequestStore{Expires: time.Now().Add(-1 * time.Second)})
	r.set("future", TTSRequestStore{Expires: time.Now().Add(10 * time.Minute)})
	r.expireOld()
	if _, exists := r.get("expired"); exists {
		t.Fatal("expected expired entry to be removed")
	}
	if _, exists := r.get("future"); !exists {
		t.Fatal("expected future entry to remain")
	}
}

func TestExpireTTSRequests_EmptyMap(t *testing.T) {
	r := initTTSRequestsStore()
	r.expireOld()
}
