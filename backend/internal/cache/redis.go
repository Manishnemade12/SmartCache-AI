package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	KeyPrefixSummary   = "summary:"
	KeyPrefixJob       = "job:"
	KeyJobQueue        = "job_queue"
	KeyMetricRequests  = "metrics:total_requests"
	KeyMetricCacheHits = "metrics:cache_hits"
	KeyMetricCacheMiss = "metrics:cache_misses"
	KeyMetricFailed    = "metrics:failed"
	KeyMetricTotalTime = "metrics:total_time_ms"
)

// Client wraps the Redis/Valkey client
type Client struct {
	rdb *redis.Client
}

// NewClient initializes the Valkey/Redis connection
func NewClient(url string) (*Client, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("invalid redis URL: %w", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Valkey/Redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// GetSummary retrieves a cached summary by hash key
func (c *Client) GetSummary(ctx context.Context, hash string) (string, error) {
	return c.rdb.Get(ctx, KeyPrefixSummary+hash).Result()
}

// SetSummary stores a summary with TTL
func (c *Client) SetSummary(ctx context.Context, hash string, data interface{}, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, KeyPrefixSummary+hash, b, ttl).Err()
}

// GetJob retrieves a job by ID
func (c *Client) GetJob(ctx context.Context, jobID string) (string, error) {
	return c.rdb.Get(ctx, KeyPrefixJob+jobID).Result()
}

// SetJob stores a job (no expiry — kept until manually cleaned)
func (c *Client) SetJob(ctx context.Context, jobID string, job interface{}) error {
	b, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, KeyPrefixJob+jobID, b, 24*time.Hour).Err()
}

// PushQueue pushes a job ID onto the queue list
func (c *Client) PushQueue(ctx context.Context, jobID string) error {
	return c.rdb.RPush(ctx, KeyJobQueue, jobID).Err()
}

// PopQueue blocks and pops a job ID from the queue (timeout 0 = block forever)
func (c *Client) PopQueue(ctx context.Context) (string, error) {
	result, err := c.rdb.BLPop(ctx, 0, KeyJobQueue).Result()
	if err != nil {
		return "", err
	}
	if len(result) < 2 {
		return "", fmt.Errorf("unexpected BLPop response")
	}
	return result[1], nil
}

// QueueSize returns current number of jobs in queue
func (c *Client) QueueSize(ctx context.Context) (int64, error) {
	return c.rdb.LLen(ctx, KeyJobQueue).Result()
}

// IncrMetric atomically increments a metric counter
func (c *Client) IncrMetric(ctx context.Context, key string) {
	c.rdb.Incr(ctx, key)
}

// AddMetricTime adds processing time to the total
func (c *Client) AddMetricTime(ctx context.Context, ms int64) {
	c.rdb.IncrBy(ctx, KeyMetricTotalTime, ms)
}

// GetMetricInt retrieves an integer metric value
func (c *Client) GetMetricInt(ctx context.Context, key string) int64 {
	val, _ := c.rdb.Get(ctx, key).Int64()
	return val
}
