package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/smartcache-ai/backend/internal/ai"
	"github.com/smartcache-ai/backend/internal/analytics"
	"github.com/smartcache-ai/backend/internal/cache"
)

// ProcessResult is the output of AI processing. Worker uses this to update job state.
type ProcessResult struct {
	Summary    string
	Tags       []string
	DurationMs int64
}

// Processor orchestrates: fetch content → call AI → cache result
type Processor struct {
	cache     *cache.Client
	aiClient  *ai.Client
	analytics *analytics.Tracker
	cacheTTL  time.Duration
}

// New creates a new Processor
func New(c *cache.Client, a *ai.Client, tracker *analytics.Tracker, cacheTTL int) *Processor {
	return &Processor{
		cache:     c,
		aiClient:  a,
		analytics: tracker,
		cacheTTL:  time.Duration(cacheTTL) * time.Second,
	}
}

// Process handles one job: fetch (if URL), call AI, cache result. Returns ProcessResult.
func (p *Processor) Process(ctx context.Context, jobID, input string) (*ProcessResult, error) {
	start := time.Now()

	// Resolve text content
	text, err := resolveInput(input)
	if err != nil {
		p.analytics.TrackFailure(ctx)
		return nil, fmt.Errorf("failed to resolve input: %w", err)
	}

	// Call Gemini AI
	result, err := p.aiClient.Summarize(ctx, text)
	if err != nil {
		p.analytics.TrackFailure(ctx)
		return nil, fmt.Errorf("AI summarization failed: %w", err)
	}

	elapsed := time.Since(start).Milliseconds()
	p.analytics.TrackProcessingTime(ctx, elapsed)

	// Cache the result keyed by job ID
	if err := p.cache.SetSummary(ctx, jobID, result, p.cacheTTL); err != nil {
		fmt.Printf("Warning: failed to cache result for job %s: %v\n", jobID, err)
	}

	return &ProcessResult{
		Summary:    result.Summary,
		Tags:       result.Tags,
		DurationMs: elapsed,
	}, nil
}

// GetCachedResult retrieves a previously cached result
func (p *Processor) GetCachedResult(ctx context.Context, hash string) (*ai.SummarizeResult, error) {
	raw, err := p.cache.GetSummary(ctx, hash)
	if err != nil {
		return nil, err
	}
	var result ai.SummarizeResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// resolveInput returns raw text from either plain text or a URL
func resolveInput(input string) (string, error) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "http://") || strings.HasPrefix(input, "https://") {
		return fetchURL(input)
	}
	return input, nil
}

// fetchURL fetches text content from a URL
func fetchURL(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 50_000))
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	text := stripHTML(string(body))
	if len(text) > 8000 {
		text = text[:8000]
	}

	return text, nil
}

// stripHTML removes HTML tags from content
func stripHTML(html string) string {
	result := strings.Builder{}
	inTag := false
	for _, r := range html {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
			result.WriteRune(' ')
		} else if !inTag {
			result.WriteRune(r)
		}
	}
	parts := strings.Fields(result.String())
	return strings.Join(parts, " ")
}
