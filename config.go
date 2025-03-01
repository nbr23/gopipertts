package main

import "os"

var VOICES_PATH = getEnv("VOICES_PATH", "/voices")
var VOICES_JSON_PATH = getEnv("VOICES_JSON_PATH", "/app/voices.json")

const VOICES_REPO_BASE_URL = "https://huggingface.co/rhasspy/piper-voices/resolve/main"

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
