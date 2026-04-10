package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Client wraps the Gemini generative AI client
type Client struct {
	model *genai.GenerativeModel
}

// SummarizeResult holds the AI output
type SummarizeResult struct {
	Summary string   `json:"summary"`
	Tags    []string `json:"tags"`
}

// NewClient initializes a new Gemini client
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()
	c, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := c.GenerativeModel("gemini-1.5-flash")
	model.SetTemperature(0.3)

	return &Client{model: model}, nil
}

// Summarize sends text to Gemini and returns a summary + tags
func (c *Client) Summarize(ctx context.Context, text string) (*SummarizeResult, error) {
	prompt := BuildPrompt(text)

	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini")
	}

	raw := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// Clean markdown code blocks if present
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var result SummarizeResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w (raw: %s)", err, raw)
	}

	return &result, nil
}
