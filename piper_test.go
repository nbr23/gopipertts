package main

import (
	"strings"
	"testing"
)

func TestSpeedToLengthScale(t *testing.T) {
	tests := []struct {
		speed float64
		want  float64
	}{
		{1.0, 1.0},
		{2.0, 0.5},
		{0.5, 2.0},
		{0, 1.0},
		{-1.0, 1.0},
	}
	for _, tt := range tests {
		if got := speedToLengthScale(tt.speed); got != tt.want {
			t.Errorf("speedToLengthScale(%v) = %v, want %v", tt.speed, got, tt.want)
		}
	}
}

func hasFlagWithValue(args []string, flag, value string) bool {
	for i, a := range args {
		if a == flag {
			return i+1 < len(args) && args[i+1] == value
		}
	}
	return false
}

func TestBuildPiperCmd_LengthScale(t *testing.T) {
	cmd := buildPiperCmd("en_US-amy-low", 0, 0.5)
	if !hasFlagWithValue(cmd.Args, "--length-scale", "0.5") {
		t.Fatalf("expected --length-scale 0.5 in args, got %v", cmd.Args)
	}
}

func TestBuildPiperCmd_NormalSpeedOmitsLengthScale(t *testing.T) {
	cmd := buildPiperCmd("en_US-amy-low", 0, 1.0)
	for _, a := range cmd.Args {
		if a == "--length-scale" {
			t.Fatalf("expected no --length-scale at speed 1.0, got %v", cmd.Args)
		}
	}
}

func TestBuildPiperCmd_NonPositiveLengthScaleOmitsFlag(t *testing.T) {
	cmd := buildPiperCmd("en_US-amy-low", 0, 0)
	for _, a := range cmd.Args {
		if a == "--length-scale" {
			t.Fatalf("expected no --length-scale for non-positive value, got %v", cmd.Args)
		}
	}
}

func TestBuildPiperCmd_SpeakerId(t *testing.T) {
	cmd := buildPiperCmd("en_US-amy-low", 3, 1.0)
	if !hasFlagWithValue(cmd.Args, "--speaker-id", "3") {
		t.Fatalf("expected --speaker-id 3 in args, got %v", cmd.Args)
	}
	if !strings.Contains(strings.Join(cmd.Args, " "), "en_US-amy-low.onnx") {
		t.Fatalf("expected model path in args, got %v", cmd.Args)
	}
}
