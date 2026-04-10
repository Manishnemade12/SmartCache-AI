package analytics

import (
	"context"

	"github.com/smartcache-ai/backend/internal/cache"
)

// Tracker handles observability metrics using Valkey
type Tracker struct {
	cache *cache.Client
}

// New creates a new analytics tracker
func New(c *cache.Client) *Tracker {
	return &Tracker{cache: c}
}

// TrackRequest increments the total request counter
func (t *Tracker) TrackRequest(ctx context.Context) {
	t.cache.IncrMetric(ctx, cache.KeyMetricRequests)
}

// TrackCacheHit increments the cache hit counter
func (t *Tracker) TrackCacheHit(ctx context.Context) {
	t.cache.IncrMetric(ctx, cache.KeyMetricCacheHits)
}

// TrackCacheMiss increments the cache miss counter
func (t *Tracker) TrackCacheMiss(ctx context.Context) {
	t.cache.IncrMetric(ctx, cache.KeyMetricCacheMiss)
}

// TrackFailure increments the failure counter
func (t *Tracker) TrackFailure(ctx context.Context) {
	t.cache.IncrMetric(ctx, cache.KeyMetricFailed)
}

// TrackProcessingTime adds the processing time in ms
func (t *Tracker) TrackProcessingTime(ctx context.Context, ms int64) {
	t.cache.AddMetricTime(ctx, ms)
}

// Metrics holds all analytics values
type Metrics struct {
	TotalRequests      int64   `json:"total_requests"`
	CacheHits          int64   `json:"cache_hits"`
	CacheMisses        int64   `json:"cache_misses"`
	FailedJobs         int64   `json:"failed_jobs"`
	QueueSize          int64   `json:"queue_size"`
	AvgProcessingMs    float64 `json:"avg_processing_time_ms"`
}

// GetMetrics retrieves all current analytics
func (t *Tracker) GetMetrics(ctx context.Context) Metrics {
	total := t.cache.GetMetricInt(ctx, cache.KeyMetricRequests)
	hits := t.cache.GetMetricInt(ctx, cache.KeyMetricCacheHits)
	misses := t.cache.GetMetricInt(ctx, cache.KeyMetricCacheMiss)
	failed := t.cache.GetMetricInt(ctx, cache.KeyMetricFailed)
	totalTime := t.cache.GetMetricInt(ctx, cache.KeyMetricTotalTime)
	queueSize, _ := t.cache.QueueSize(ctx)

	var avgTime float64
	processed := misses // each cache miss triggers processing
	if processed > 0 {
		avgTime = float64(totalTime) / float64(processed)
	}

	return Metrics{
		TotalRequests:   total,
		CacheHits:       hits,
		CacheMisses:     misses,
		FailedJobs:      failed,
		QueueSize:       queueSize,
		AvgProcessingMs: avgTime,
	}
}
