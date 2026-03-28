package main

import (
	"log"
	"os"
	"strconv"
)

var VOICES_PATH = getEnv("VOICES_PATH", "/voices")
var VOICES_JSON_PATH = getEnv("VOICES_JSON_PATH", "/app/voices.json")
var PIPER_BINARY = getEnv("PIPER_BINARY", "/usr/share/piper/piper")
var STREAM_EXPIRATION_MINUTES = getIntEnv("STREAM_EXPIRATION_MINUTES", "15")
var logInput = os.Getenv("LOG_INPUT") != ""

const VOICES_REPO_BASE_URL = "https://huggingface.co/rhasspy/piper-voices/resolve/main"

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getIntEnv(key, defaultValue string) int {
	value := getEnv(key, defaultValue)
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", key, value)
	}
	return intValue
}
