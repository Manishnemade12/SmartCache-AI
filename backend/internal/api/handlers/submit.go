package handlers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/smartcache-ai/backend/internal/analytics"
	"github.com/smartcache-ai/backend/internal/cache"
	"github.com/smartcache-ai/backend/internal/services"
	"github.com/smartcache-ai/backend/internal/worker"
)

// SubmitHandler handles POST /api/submit
type SubmitHandler struct {
	cache     *cache.Client
	processor *services.Processor
	analytics *analytics.Tracker
}

// NewSubmitHandler creates a new submit handler
func NewSubmitHandler(c *cache.Client, p *services.Processor, a *analytics.Tracker) *SubmitHandler {
	return &SubmitHandler{cache: c, processor: p, analytics: a}
}

type submitRequest struct {
	Input string `json:"input" binding:"required"`
}

// Handle processes the submit request
func (h *SubmitHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	h.analytics.TrackRequest(ctx)

	var req submitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input field is required"})
		return
	}

	input := strings.TrimSpace(req.Input)
	if input == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input cannot be empty"})
		return
	}

	// Generate deterministic hash key
	hash := hashInput(input)

	// Check cache first (HIT path)
	if raw, err := h.cache.GetSummary(ctx, hash); err == nil {
		h.analytics.TrackCacheHit(ctx)

		var cached map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &cached); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"job_id":  hash,
				"status":  "completed",
				"cached":  true,
				"summary": cached["summary"],
				"tags":    cached["tags"],
			})
			return
		}
	}

	// MISS path — create job and enqueue
	h.analytics.TrackCacheMiss(ctx)

	jobID := uuid.New().String()
	job := worker.Job{
		ID:        jobID,
		Input:     input,
		Status:    worker.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := h.cache.SetJob(ctx, jobID, job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create job"})
		return
	}

	if err := h.cache.PushQueue(ctx, jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue job"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"job_id": jobID,
		"status": worker.StatusPending,
		"cached": false,
	})
}

// hashInput generates a stable SHA-256 hash for deduplication
func hashInput(input string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(input))))
	return fmt.Sprintf("%x", sum)[:16] // 16-char prefix for readability
}
