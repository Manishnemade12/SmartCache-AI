package worker

import "time"

// Job states
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// Job represents an async summarization task
type Job struct {
	ID          string    `json:"job_id"`
	Input       string    `json:"input"`
	Status      string    `json:"status"`
	Summary     string    `json:"summary,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Error       string    `json:"error,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	DurationMs  int64     `json:"duration_ms,omitempty"`
}
