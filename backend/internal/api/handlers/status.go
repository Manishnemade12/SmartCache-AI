package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcache-ai/backend/internal/cache"
	"github.com/smartcache-ai/backend/internal/worker"
)

// StatusHandler handles GET /api/status/:job_id
type StatusHandler struct {
	cache *cache.Client
}

// NewStatusHandler creates a new status handler
func NewStatusHandler(c *cache.Client) *StatusHandler {
	return &StatusHandler{cache: c}
}

// Handle returns the current job status
func (h *StatusHandler) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	jobID := c.Param("job_id")

	if jobID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_id is required"})
		return
	}

	raw, err := h.cache.GetJob(ctx, jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	var job worker.Job
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse job"})
		return
	}

	resp := gin.H{
		"job_id":     job.ID,
		"status":     job.Status,
		"created_at": job.CreatedAt,
	}

	if job.Status == worker.StatusCompleted {
		resp["summary"] = job.Summary
		resp["tags"] = job.Tags
		resp["duration_ms"] = job.DurationMs
		resp["completed_at"] = job.CompletedAt
	}

	if job.Status == worker.StatusFailed {
		resp["error"] = job.Error
	}

	c.JSON(http.StatusOK, resp)
}
