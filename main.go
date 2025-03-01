package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8080")

	if os.Getenv("DEBUG") == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	voices := getAvailableVoices(VOICES_JSON_PATH)
	loadVoicesDetails()

	r.GET("/", homeHandler)
	r.GET("/api/voices", voicesHandler(&voices))
	r.POST("/api/tts", ttsHandler(&voices))

	fmt.Printf("Listening and serving HTTP on :%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
