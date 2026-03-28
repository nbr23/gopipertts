package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseVoiceDetails_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "voice.json")
	content := `{"audio":{"sample_rate":22050},"speaker_id_map":{"s0":0}}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	v, err := parseVoiceDetails(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Audio.SampleRate != 22050 {
		t.Fatalf("expected sample rate 22050, got %d", v.Audio.SampleRate)
	}
	if id, ok := v.SpeakerIdMap["s0"]; !ok || id != 0 {
		t.Fatalf("expected SpeakerIdMap[s0]=0, got %v", v.SpeakerIdMap)
	}
}

func TestParseVoiceDetails_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "voice.json")
	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := parseVoiceDetails(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseVoiceDetails_FileNotFound(t *testing.T) {
	_, err := parseVoiceDetails("/nonexistent/path/voice.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
