package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGetTTSStrParameter_PostValueTakesPrecedence(t *testing.T) {
	c, _ := newTestContext("GET", "/?voice=fromQuery", "")
	result := getTTSStrParameter(c, "fromPost", "voice", "default")
	if result != "fromPost" {
		t.Fatalf("expected 'fromPost', got %q", result)
	}
}

func TestGetTTSStrParameter_FallsBackToQuery(t *testing.T) {
	c, _ := newTestContext("GET", "/?voice=fromQuery", "")
	result := getTTSStrParameter(c, "", "voice", "default")
	if result != "fromQuery" {
		t.Fatalf("expected 'fromQuery', got %q", result)
	}
}

func TestGetTTSStrParameter_FallsBackToDefault(t *testing.T) {
	c, _ := newTestContext("GET", "/", "")
	result := getTTSStrParameter(c, "", "voice", "default")
	if result != "default" {
		t.Fatalf("expected 'default', got %q", result)
	}
}

func TestGetTTSFloatParameter_PostValueTakesPrecedence(t *testing.T) {
	c, _ := newTestContext("GET", "/?speed=2.0", "")
	result := getTTSFloatParameter(c, 1.5, "speed", 1.0)
	if result != 1.5 {
		t.Fatalf("expected 1.5, got %f", result)
	}
}

func TestGetTTSFloatParameter_FallsBackToQuery(t *testing.T) {
	c, _ := newTestContext("GET", "/?speed=1.5", "")
	result := getTTSFloatParameter(c, 0, "speed", 1.0)
	if result != 1.5 {
		t.Fatalf("expected 1.5, got %f", result)
	}
}

func TestGetTTSFloatParameter_FallsBackToDefault(t *testing.T) {
	c, _ := newTestContext("GET", "/", "")
	result := getTTSFloatParameter(c, 0, "speed", 1.0)
	if result != 1.0 {
		t.Fatalf("expected 1.0, got %f", result)
	}
}

func TestGetTTSFloatParameter_InvalidQueryValue(t *testing.T) {
	c, _ := newTestContext("GET", "/?speed=notanumber", "")
	result := getTTSFloatParameter(c, 0, "speed", 1.0)
	if result != 1.0 {
		t.Fatalf("expected default 1.0, got %f", result)
	}
}

func TestGetTTSRequestInput_GET_AllParams(t *testing.T) {
	c, _ := newTestContext("GET", "/?text=hello&voice=myvoice&speaker=0&speed=1.5&outputFormat=mp3", "")
	input, err := getTTSRequestInput(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Text != "hello" {
		t.Fatalf("expected text='hello', got %q", input.Text)
	}
	if input.Voice != "myvoice" {
		t.Fatalf("expected voice='myvoice', got %q", input.Voice)
	}
	if input.Speaker != "0" {
		t.Fatalf("expected speaker='0', got %q", input.Speaker)
	}
	if input.Speed != 1.5 {
		t.Fatalf("expected speed=1.5, got %f", input.Speed)
	}
	if input.OutputFormat != "mp3" {
		t.Fatalf("expected outputFormat='mp3', got %q", input.OutputFormat)
	}
}

func TestGetTTSRequestInput_GET_DefaultsApplied(t *testing.T) {
	c, _ := newTestContext("GET", "/?text=hello", "")
	input, err := getTTSRequestInput(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Voice != "en_US-amy-low" {
		t.Fatalf("expected default voice, got %q", input.Voice)
	}
	if input.Speed != 1.0 {
		t.Fatalf("expected default speed 1.0, got %f", input.Speed)
	}
	if input.OutputFormat != "wav" {
		t.Fatalf("expected default outputFormat 'wav', got %q", input.OutputFormat)
	}
	if input.Speaker != "" {
		t.Fatalf("expected empty speaker, got %q", input.Speaker)
	}
}

func TestGetTTSRequestInput_POST_JSONBody(t *testing.T) {
	body := `{"text":"hello","voice":"en_US-kathleen-low","speed":1.2,"outputFormat":"mp3"}`
	c, _ := newTestContext("POST", "/", body)
	input, err := getTTSRequestInput(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.Text != "hello" {
		t.Fatalf("expected text='hello', got %q", input.Text)
	}
	if input.Voice != "en_US-kathleen-low" {
		t.Fatalf("expected voice='en_US-kathleen-low', got %q", input.Voice)
	}
	if input.Speed != 1.2 {
		t.Fatalf("expected speed=1.2, got %f", input.Speed)
	}
	if input.OutputFormat != "mp3" {
		t.Fatalf("expected outputFormat='mp3', got %q", input.OutputFormat)
	}
}

func TestVoicesHandler_ReturnsJSON(t *testing.T) {
	voices := Voices{
		"en_US-amy-low": Voice{Key: "en_US-amy-low", Name: "amy"},
	}
	c, w := newTestContext("GET", "/api/voices", "")
	voicesHandler(&voices)(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result Voices
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := result["en_US-amy-low"]; !ok {
		t.Fatal("expected 'en_US-amy-low' in response")
	}
}

func TestVoicesHandler_EmptyVoices(t *testing.T) {
	voices := Voices{}
	c, w := newTestContext("GET", "/api/voices", "")
	voicesHandler(&voices)(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "{}" {
		t.Fatalf("expected '{}', got %q", w.Body.String())
	}
}

func TestHomeHandler_ReturnsHTML(t *testing.T) {
	c, w := newTestContext("GET", "/", "")
	homeHandler(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("expected text/html content type, got %q", ct)
	}
}

func TestWriteWavStreamHttpHeaders_Headers(t *testing.T) {
	c, w := newTestContext("GET", "/", "")
	err := writeWavStreamHttpHeaders(c, 22050, 1, 16)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "audio/wav" {
		t.Fatalf("expected Content-Type audio/wav, got %q", ct)
	}
	if len(w.Body.Bytes()) != 44 {
		t.Fatalf("expected 44 byte WAV header in body, got %d", len(w.Body.Bytes()))
	}
}

func TestTTSPostStreamHandler_MissingText(t *testing.T) {
	r := map[string]TTSRequestStore{}
	c, w := newTestContext("POST", "/api/tts/stream", `{"voice":"en_US-amy-low"}`)
	ttsPostStreamHandler(r)(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "text query parameter is required") {
		t.Fatalf("expected error message in body, got %q", w.Body.String())
	}
}

func TestTTSPostStreamHandler_Mp3Rejected(t *testing.T) {
	r := map[string]TTSRequestStore{}
	c, w := newTestContext("POST", "/api/tts/stream", `{"text":"hello","outputFormat":"mp3"}`)
	ttsPostStreamHandler(r)(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "mp3") {
		t.Fatalf("expected mp3 rejection message, got %q", w.Body.String())
	}
}

func TestTTSPostStreamHandler_ValidRequest(t *testing.T) {
	r := map[string]TTSRequestStore{}
	c, w := newTestContext("POST", "/api/tts/stream", `{"text":"hello"}`)
	ttsPostStreamHandler(r)(c)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	streamId, ok := result["streamId"].(string)
	if !ok || streamId == "" {
		t.Fatalf("expected non-empty streamId, got %v", result["streamId"])
	}
	if _, exists := r[streamId]; !exists {
		t.Fatalf("expected streamId %q in store", streamId)
	}
}

func TestTTSPostStreamHandler_StreamIdIsUnique(t *testing.T) {
	r := map[string]TTSRequestStore{}

	c1, w1 := newTestContext("POST", "/api/tts/stream", `{"text":"hello"}`)
	ttsPostStreamHandler(r)(c1)

	c2, w2 := newTestContext("POST", "/api/tts/stream", `{"text":"hello"}`)
	ttsPostStreamHandler(r)(c2)

	var r1, r2 map[string]interface{}
	json.Unmarshal(w1.Body.Bytes(), &r1)
	json.Unmarshal(w2.Body.Bytes(), &r2)

	id1 := r1["streamId"].(string)
	id2 := r2["streamId"].(string)
	if id1 == id2 {
		t.Fatalf("expected unique stream IDs, both were %q", id1)
	}
}

func TestTTSGetStreamHandler_NotFound(t *testing.T) {
	voices := Voices{}
	r := map[string]TTSRequestStore{}
	c, w := newTestContext("GET", "/api/tts/stream/unknown-id", "")
	c.Params = gin.Params{{Key: "streamId", Value: "unknown-id"}}
	ttsGetStreamHandler(&voices, r)(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Stream not found") {
		t.Fatalf("expected 'Stream not found', got %q", w.Body.String())
	}
}

func TestTTSGetStreamHandler_ExpiredStream(t *testing.T) {
	voices := Voices{}
	r := map[string]TTSRequestStore{
		"expired-id": {
			Request: TTSRequestInput{Text: "hello"},
			Expires: time.Now().Add(-1 * time.Second),
		},
	}
	c, w := newTestContext("GET", "/api/tts/stream/expired-id", "")
	c.Params = gin.Params{{Key: "streamId", Value: "expired-id"}}
	ttsGetStreamHandler(&voices, r)(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if _, exists := r["expired-id"]; exists {
		t.Fatal("expected expired entry to be deleted from store")
	}
}
