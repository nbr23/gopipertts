package main

import (
	"embed"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:embed static/index.html
var staticFiles embed.FS

type TTSRequestInput struct {
	Text         string  `json:"text"`
	Voice        string  `json:"voice"`
	Speaker      string  `json:"speaker"`
	Speed        float64 `json:"speed"`
	OutputFormat string  `json:"outputFormat"`
}

func homeHandler(c *gin.Context) {
	html, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		c.String(http.StatusInternalServerError, "Error reading index.html")
		return
	}
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, string(html))
}

func voicesHandler(voices *Voices) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, voices)
	}
}

func writeWavStreamHttpHeaders(c *gin.Context, sampleRate int, channels int, bitsPerSample int) error {
	c.Header("Content-Type", "audio/wav")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	header := generateWAVHeader(sampleRate, channels, bitsPerSample)
	if _, err := c.Writer.Write(header); err != nil {
		return err
	}
	c.Writer.Flush()
	return nil
}

func getTTSStrParameter(c *gin.Context, postValue string, key string, defaultValue string) string {
	value := postValue
	if value == "" {
		value = c.Query(key)
	}
	if value == "" {
		value = defaultValue
	}
	return value
}

func ttsGetStreamHandler(voices *Voices, r *TTSRequestsStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamId := c.Param("streamId")
		ttsRequest, ok := r.get(streamId)
		if !ok {
			c.String(http.StatusNotFound, "Stream not found")
			return
		}
		if ttsRequest.Expires.Before(time.Now()) {
			r.delete(streamId)
			c.String(http.StatusNotFound, "Stream not found")
			return
		}
		piperToAudioStream(c, ttsRequest.Request, voices)
	}
}

func ttsPostStreamHandler(r *TTSRequestsStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		streamId := uuid.New().String()

		ttsRequestInput, err := getTTSRequestInput(c)
		if err != nil {
			c.String(http.StatusBadRequest, "Invalid request")
			return
		}

		if ttsRequestInput.Text == "" {
			c.String(http.StatusBadRequest, "text query parameter is required")
			return
		}
		if ttsRequestInput.OutputFormat == "mp3" {
			c.String(http.StatusBadRequest, "outputFormat 'mp3' is not supported on stream endpoints")
			return
		}
		entry := TTSRequestStore{
			Request: ttsRequestInput,
			Expires: time.Now().Add(time.Duration(STREAM_EXPIRATION_MINUTES) * time.Minute),
		}
		r.set(streamId, entry)
		c.JSON(http.StatusOK, gin.H{
			"streamId": streamId,
			"expires":  entry.Expires,
		})
	}
}

func getTTSFloatParameter(c *gin.Context, postValue float64, key string, defaultValue float64) float64 {
	value := postValue
	if value == 0 {
		parsedValue, err := strconv.ParseFloat(c.Query(key), 64)
		if err == nil {
			value = parsedValue
		}
	}
	if value == 0 {
		value = defaultValue
	}
	return value
}

func getTTSRequestInput(c *gin.Context) (TTSRequestInput, error) {
	var ttsRequestInput TTSRequestInput

	if c.Request.Method == "POST" {
		if err := c.ShouldBindJSON(&ttsRequestInput); err != nil {
			return TTSRequestInput{}, err
		}
	}

	ttsRequestInput.Voice = getTTSStrParameter(c, ttsRequestInput.Voice, "voice", "en_US-amy-low")
	ttsRequestInput.Speaker = getTTSStrParameter(c, ttsRequestInput.Speaker, "speaker", "")
	ttsRequestInput.Speed = getTTSFloatParameter(c, ttsRequestInput.Speed, "speed", 1.0)
	ttsRequestInput.Text = getTTSStrParameter(c, ttsRequestInput.Text, "text", "")
	ttsRequestInput.OutputFormat = getTTSStrParameter(c, ttsRequestInput.OutputFormat, "outputFormat", "wav")

	return ttsRequestInput, nil
}

func piperToAudioStream(c *gin.Context, ttsRequestInput TTSRequestInput, voices *Voices) {
	if ttsRequestInput.Text == "" {
		c.String(http.StatusBadRequest, "text query parameter is required")
		return
	}

	if ttsRequestInput.OutputFormat != "wav" && ttsRequestInput.OutputFormat != "mp3" {
		c.String(http.StatusBadRequest, "invalid outputFormat, must be 'wav' or 'mp3'")
		return
	}

	voice, err := getVoiceDetails(voices, ttsRequestInput.Voice)
	if err != nil {
		c.String(http.StatusBadRequest, "Voice not found")
		return
	}
	speaker, ok := voice.SpeakerIdMap[ttsRequestInput.Speaker]
	if !ok {
		speaker = 0
	}

	sampleRate := voice.Audio.SampleRate
	lengthScale := 1.0
	if ttsRequestInput.Speed > 0 {
		lengthScale = 1.0 / ttsRequestInput.Speed
	}
	channels := 1
	bitsPerSample := 16

	if ttsRequestInput.OutputFormat == "mp3" {
		if err := streamTTSAsMp3(c, ttsRequestInput.Voice, speaker, ttsRequestInput.Text, sampleRate, lengthScale); err != nil {
			log.Printf("Error streaming MP3 TTS: %v", err)
		}
		return
	}

	err = writeWavStreamHttpHeaders(c, sampleRate, channels, bitsPerSample)
	if err != nil {
		log.Printf("error writting http headers: %v", err)
		c.String(http.StatusInternalServerError, "Error streaming TTS")
		return
	}

	err = streamTTS(c, ttsRequestInput.Voice, speaker, ttsRequestInput.Text, sampleRate, channels, bitsPerSample, lengthScale)
	if err != nil {
		log.Printf("Error streaming TTS: %v", err)
		c.String(http.StatusInternalServerError, "Error streaming TTS")
	}
}

func ttsHandler(voices *Voices) gin.HandlerFunc {
	return func(c *gin.Context) {
		ttsRequestInput, err := getTTSRequestInput(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request, JSON body required"})
			return
		}
		piperToAudioStream(c, ttsRequestInput, voices)
	}
}
