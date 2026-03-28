package main

import (
	"encoding/binary"
	"testing"
)

func TestGenerateWAVHeader_Length(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	if len(h) != 44 {
		t.Fatalf("expected 44 bytes, got %d", len(h))
	}
}

func TestGenerateWAVHeader_RIFFChunkId(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	if string(h[0:4]) != "RIFF" {
		t.Fatalf("expected RIFF, got %q", h[0:4])
	}
}

func TestGenerateWAVHeader_RIFFChunkSize(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	for i := 4; i < 8; i++ {
		if h[i] != 0xFF {
			t.Fatalf("byte %d: expected 0xFF, got 0x%02X", i, h[i])
		}
	}
}

func TestGenerateWAVHeader_WAVEMarker(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	if string(h[8:12]) != "WAVE" {
		t.Fatalf("expected WAVE, got %q", h[8:12])
	}
}

func TestGenerateWAVHeader_FmtChunkId(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	if string(h[12:16]) != "fmt " {
		t.Fatalf("expected 'fmt ', got %q", h[12:16])
	}
}

func TestGenerateWAVHeader_FmtChunkSize(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	val := binary.LittleEndian.Uint32(h[16:20])
	if val != 16 {
		t.Fatalf("expected fmt chunk size 16, got %d", val)
	}
}

func TestGenerateWAVHeader_AudioFormat(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	val := binary.LittleEndian.Uint16(h[20:22])
	if val != 1 {
		t.Fatalf("expected audio format 1 (PCM), got %d", val)
	}
}

func TestGenerateWAVHeader_Channels(t *testing.T) {
	h := generateWAVHeader(22050, 2, 16)
	if h[22] != 2 {
		t.Fatalf("expected channels=2, got %d", h[22])
	}
	if h[23] != 0 {
		t.Fatalf("expected h[23]=0, got %d", h[23])
	}
}

func TestGenerateWAVHeader_SampleRate(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	val := binary.LittleEndian.Uint32(h[24:28])
	if val != 22050 {
		t.Fatalf("expected sample rate 22050, got %d", val)
	}
}

func TestGenerateWAVHeader_ByteRate(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	expected := uint32(22050 * 1 * 16 / 8)
	val := binary.LittleEndian.Uint32(h[28:32])
	if val != expected {
		t.Fatalf("expected byte rate %d, got %d", expected, val)
	}
}

func TestGenerateWAVHeader_BlockAlign(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	expected := uint16(1 * 16 / 8)
	val := binary.LittleEndian.Uint16(h[32:34])
	if val != expected {
		t.Fatalf("expected block align %d, got %d", expected, val)
	}
}

func TestGenerateWAVHeader_BitsPerSample(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	val := binary.LittleEndian.Uint16(h[34:36])
	if val != 16 {
		t.Fatalf("expected bits per sample 16, got %d", val)
	}
}

func TestGenerateWAVHeader_DataChunkId(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	if string(h[36:40]) != "data" {
		t.Fatalf("expected 'data', got %q", h[36:40])
	}
}

func TestGenerateWAVHeader_DataChunkSize(t *testing.T) {
	h := generateWAVHeader(22050, 1, 16)
	for i := 40; i < 44; i++ {
		if h[i] != 0xFF {
			t.Fatalf("byte %d: expected 0xFF, got 0x%02X", i, h[i])
		}
	}
}
