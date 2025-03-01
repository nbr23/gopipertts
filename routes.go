package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TTSRequest struct {
	Text    string  `json:"text"`
	Voice   string  `json:"voice"`
	Speaker string  `json:"speaker"`
	Speed   float64 `json:"speed"`
}

func homeHandler(c *gin.Context) {
	html := INDEX_HTML
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

func voicesHandler(voices *Voices) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, voices)
	}
}

func writeWavStreamHttpHeaders(c *gin.Context, sampleRate int, channels int, bitsPerSample int) error {
	c.Header("Content-Type", "audio/wav")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	header := generateWAVHeader(sampleRate, channels, bitsPerSample)
	if _, err := c.Writer.Write(header); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}

func ttsHandler(voices *Voices) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ttsRequest TTSRequest

		voiceName := c.Query("voice")
		if voiceName == "" {
			voiceName = "en_US-amy-low"
		}
		speakerName := c.Query("speaker")

		speed, err := strconv.ParseFloat(c.Query("speed"), 64)
		if err != nil {
			speed = 1.0
		}

		if err := c.ShouldBindJSON(&ttsRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if ttsRequest.Text == "" {
			c.String(http.StatusBadRequest, "text query parameter is required")
			return
		}
		if ttsRequest.Voice != "" {
			voiceName = ttsRequest.Voice
		}
		if ttsRequest.Speaker != "" {
			speakerName = ttsRequest.Speaker
		}
		if ttsRequest.Speed != 0 {
			speed = ttsRequest.Speed
		}

		voice, err := getVoiceDetails(voices, voiceName)
		if err != nil {
			c.String(http.StatusBadRequest, "Voice not found")
			return
		}
		speaker, ok := voice.SpeakerIdMap[speakerName]
		if !ok {
			speaker = 0
		}

		sampleRate := int(float64(voice.Audio.SampleRate) * speed)
		channels := 1
		bitsPerSample := 16

		err = writeWavStreamHttpHeaders(c, sampleRate, channels, bitsPerSample)
		if err != nil {
			log.Printf("error writting http headers: %v", err)
			c.String(http.StatusInternalServerError, "Error streaming TTS")
			return
		}

		err = streamTTS(c, voiceName, speaker, ttsRequest.Text, int(sampleRate), channels, bitsPerSample)
		if err != nil {
			log.Printf("Error streaming TTS: %v", err)
			c.String(http.StatusInternalServerError, "Error streaming TTS")
		}
	}
}
