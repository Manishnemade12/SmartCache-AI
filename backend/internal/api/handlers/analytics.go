package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smartcache-ai/backend/internal/analytics"
)

// AnalyticsHandler handles GET /api/analytics
type AnalyticsHandler struct {
	tracker *analytics.Tracker
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(t *analytics.Tracker) *AnalyticsHandler {
	return &AnalyticsHandler{tracker: t}
}

// Handle returns current system metrics
func (h *AnalyticsHandler) Handle(c *gin.Context) {
	metrics := h.tracker.GetMetrics(c.Request.Context())
	c.JSON(http.StatusOK, metrics)
}
