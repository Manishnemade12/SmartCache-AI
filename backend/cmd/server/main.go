package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartcache-ai/backend/config"
	"github.com/smartcache-ai/backend/internal/ai"
	"github.com/smartcache-ai/backend/internal/analytics"
	"github.com/smartcache-ai/backend/internal/api/handlers"
	"github.com/smartcache-ai/backend/internal/cache"
	"github.com/smartcache-ai/backend/internal/services"
	"github.com/smartcache-ai/backend/internal/worker"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize Valkey/Redis client
	cacheClient, err := cache.NewClient(config.C.RedisURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Valkey: %v", err)
	}
	log.Println("✅ Connected to Valkey/Redis")

	// Initialize Gemini AI client
	aiClient, err := ai.NewClient(config.C.GeminiKey)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Gemini client: %v", err)
	}
	log.Println("✅ Gemini AI client initialized")

	// Initialize analytics tracker
	tracker := analytics.New(cacheClient)

	// Initialize processor
	processor := services.New(cacheClient, aiClient, tracker, config.C.CacheTTL)

	// Initialize and start worker pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := worker.NewPool(config.C.WorkerCount, cacheClient, processor)
	pool.Start(ctx)

	// Set up Gin router
	r := gin.Default()

	// CORS middleware for React frontend
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Initialize handlers
	submitHandler := handlers.NewSubmitHandler(cacheClient, processor, tracker)
	statusHandler := handlers.NewStatusHandler(cacheClient)
	analyticsHandler := handlers.NewAnalyticsHandler(tracker)

	// Register routes
	api := r.Group("/api")
	{
		api.POST("/submit", submitHandler.Handle)
		api.GET("/status/:job_id", statusHandler.Handle)
		api.GET("/analytics", analyticsHandler.Handle)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
		})
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + config.C.Port,
		Handler: r,
	}

	log.Printf("🚀 SmartCache AI backend running on port %s", config.C.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")
	cancel() // Stop workers

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited cleanly")
}
