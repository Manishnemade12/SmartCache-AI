package ai

import "fmt"

// BuildPrompt creates the summarization prompt
func BuildPrompt(text string) string {
	return fmt.Sprintf(`You are a concise summarization assistant.

Summarize the following text in 2-3 sentences and extract 2-4 relevant tags.

Return ONLY valid JSON in this exact format (no markdown, no code blocks):
{
  "summary": "Your 2-3 sentence summary here.",
  "tags": ["tag1", "tag2", "tag3"]
}

TEXT TO SUMMARIZE:
%s`, text)
}
