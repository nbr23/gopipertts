package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8080")
	preloadVoices := getEnv("PRELOAD_VOICES", "")

	if os.Getenv("DEBUG") == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	voices := getAvailableVoices(VOICES_JSON_PATH)
	loadVoicesDetails()
	ensureVoices(strings.Split(preloadVoices, ","), &voices)
	requestsMap := initTTSRequestsStore()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.GET("/", homeHandler)
	r.GET("/api/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.GET("/api/voices", voicesHandler(&voices))
	r.POST("/api/tts", ttsHandler(&voices))
	r.GET("/api/tts", ttsHandler(&voices))
	r.POST("/api/tts/stream", ttsPostStreamHandler(requestsMap))
	r.GET("/api/tts/stream/:streamId", ttsGetStreamHandler(&voices, requestsMap))

	fmt.Printf("Listening and serving HTTP on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
