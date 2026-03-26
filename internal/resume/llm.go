package resume

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/yourusername/autoreach-backend/internal/config"
	"github.com/yourusername/autoreach-backend/pkg/logger"
)

const resumePrompt = `You are a precise resume parser. Extract the following structured data from the resume text below.

Extract:
- name (string)
- email (string)
- skills (array of strings)
- experience (array of objects with: company_name, role, start_date, end_date, work_done)
- projects (array of objects with: name, description, tech_stack, url)
- education (array of objects with: degree, institution, year)

Rules:
1. Return ONLY valid JSON. No explanation, no markdown fences, no extra text.
2. If a field is missing or unclear, use an empty string or empty array.
3. For dates, use format "YYYY-MM" if available, otherwise use whatever is provided.
4. Keep skill names concise (e.g., "Go", "PostgreSQL", not long descriptions).
5. For work_done and descriptions, keep them concise but informative.

Resume Text:
---
%s
---

Return ONLY the JSON object.`

// CallOpenAI sends cleaned resume text to GPT and returns the structured JSON string.
func CallOpenAI(ctx context.Context, cleanedText string) (string, error) {
	apiKey := config.AppConfig.OpenAIKey
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not configured")
	}

	client := openai.NewClient(apiKey)
	prompt := fmt.Sprintf(resumePrompt, cleanedText)

	var lastErr error
	maxRetries := 2

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			logger.Warn(fmt.Sprintf("OpenAI retry attempt %d", attempt))
		}

		resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a JSON-only resume parser. Output nothing but valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
			MaxTokens:   4096,
		})
		if err != nil {
			lastErr = fmt.Errorf("OpenAI API error: %w", err)
			continue
		}

		if len(resp.Choices) == 0 {
			lastErr = fmt.Errorf("OpenAI returned empty response")
			continue
		}

		content := resp.Choices[0].Message.Content
		logger.Info(fmt.Sprintf("OpenAI response received (%d chars)", len(content)))

		return content, nil
	}

	return "", fmt.Errorf("OpenAI failed after %d attempts: %w", maxRetries, lastErr)
}
