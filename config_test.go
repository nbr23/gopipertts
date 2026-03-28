package main

import (
	"testing"
)

func TestGetEnv_ReturnsEnvValue(t *testing.T) {
	t.Setenv("TEST_KEY_GOSTREAM", "myvalue")
	result := getEnv("TEST_KEY_GOSTREAM", "default")
	if result != "myvalue" {
		t.Fatalf("expected 'myvalue', got %q", result)
	}
}

func TestGetEnv_ReturnsDefaultWhenUnset(t *testing.T) {
	result := getEnv("TEST_KEY_GOSTREAM_UNSET_XYZ", "default")
	if result != "default" {
		t.Fatalf("expected 'default', got %q", result)
	}
}

func TestGetIntEnv_ReturnsIntValue(t *testing.T) {
	t.Setenv("TEST_INT_GOSTREAM", "30")
	result := getIntEnv("TEST_INT_GOSTREAM", "15")
	if result != 30 {
		t.Fatalf("expected 30, got %d", result)
	}
}

func TestGetIntEnv_ReturnsDefaultWhenUnset(t *testing.T) {
	result := getIntEnv("TEST_INT_GOSTREAM_UNSET_XYZ", "15")
	if result != 15 {
		t.Fatalf("expected 15, got %d", result)
	}
}
