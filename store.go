package main

import "time"

type TTSRequestStore struct {
	Request TTSRequestInput
	Expires time.Time
}

func expireTTSRequests(r map[string]TTSRequestStore) {
	for key, value := range r {
		if value.Expires.Before(time.Now()) {
			delete(r, key)
		}
	}
}

func initTTSRequestsStore() map[string]TTSRequestStore {
	r := make(map[string]TTSRequestStore)

	go func() {
		for {
			expireTTSRequests(r)
			time.Sleep(15 * time.Minute)
		}
	}()

	return r
}
