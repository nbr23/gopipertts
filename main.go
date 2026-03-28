package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8080")
	preloadVoices := getEnv("PRELOAD_VOICES", "")

	debug := os.Getenv("DEBUG")
	if debug == "1" || debug == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	voices := getAvailableVoices(VOICES_JSON_PATH)
	loadVoicesDetails()
	ensureVoices(strings.Split(preloadVoices, ","), &voices)
	requestsMap := initTTSRequestsStore()

	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/api/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Use(gin.Logger())
	r.GET("/", homeHandler)
	r.GET("/api/voices", voicesHandler(&voices))
	r.POST("/api/tts", ttsHandler(&voices))
	r.GET("/api/tts", ttsHandler(&voices))
	r.POST("/api/tts/stream", ttsPostStreamHandler(requestsMap))
	r.GET("/api/tts/stream/:streamId", ttsGetStreamHandler(&voices, requestsMap))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		fmt.Printf("Listening and serving HTTP on :%s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	killAllProcesses()
	log.Println("Server exited")
}
