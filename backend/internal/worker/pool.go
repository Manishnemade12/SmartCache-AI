package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/smartcache-ai/backend/internal/cache"
	"github.com/smartcache-ai/backend/internal/services"
)

// Pool manages a set of background goroutine workers
type Pool struct {
	size      int
	cache     *cache.Client
	processor *services.Processor
}

// NewPool creates a new worker pool
func NewPool(size int, c *cache.Client, p *services.Processor) *Pool {
	return &Pool{size: size, cache: c, processor: p}
}

// Start launches N goroutine workers that continuously process jobs
func (p *Pool) Start(ctx context.Context) {
	log.Printf("🚀 Starting %d workers...", p.size)
	for i := 0; i < p.size; i++ {
		go p.runWorker(ctx, i+1)
	}
}

// runWorker is the main loop for a single worker goroutine
func (p *Pool) runWorker(ctx context.Context, id int) {
	log.Printf("Worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Worker %d shutting down", id)
			return
		default:
			jobID, err := p.cache.PopQueue(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("Worker %d: queue error: %v", id, err)
				time.Sleep(1 * time.Second)
				continue
			}
			p.processJob(ctx, id, jobID)
		}
	}
}

// processJob handles a single job from the queue
func (p *Pool) processJob(ctx context.Context, workerID int, jobID string) {
	log.Printf("Worker %d: processing job %s", workerID, jobID)

	// Fetch job from Valkey
	raw, err := p.cache.GetJob(ctx, jobID)
	if err != nil {
		log.Printf("Worker %d: failed to get job %s: %v", workerID, jobID, err)
		return
	}

	var job Job
	if err := json.Unmarshal([]byte(raw), &job); err != nil {
		log.Printf("Worker %d: failed to unmarshal job %s: %v", workerID, jobID, err)
		return
	}

	// Mark as processing
	job.Status = StatusProcessing
	if err := p.cache.SetJob(ctx, jobID, job); err != nil {
		log.Printf("Worker %d: failed to update job status: %v", workerID, err)
	}

	// Process (fetch URL if needed + AI call) using processor
	result, err := p.processor.Process(ctx, jobID, job.Input)
	if err != nil {
		log.Printf("Worker %d: job %s failed: %v", workerID, jobID, err)
		job.Status = StatusFailed
		job.Error = fmt.Sprintf("Processing error: %v", err)
		_ = p.cache.SetJob(ctx, jobID, job)
		return
	}

	// Update job with AI results
	job.Summary = result.Summary
	job.Tags = result.Tags
	job.Status = StatusCompleted
	job.CompletedAt = time.Now()
	job.DurationMs = result.DurationMs

	if err := p.cache.SetJob(ctx, jobID, job); err != nil {
		log.Printf("Worker %d: failed to save completed job: %v", workerID, err)
	}

	log.Printf("Worker %d: job %s completed in %dms", workerID, jobID, result.DurationMs)
}
