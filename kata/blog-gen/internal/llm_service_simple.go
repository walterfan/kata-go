package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Simple structures for basic chat
type SimpleChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SimpleChatChoice struct {
	Index        int               `json:"index"`
	Message      SimpleChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

type SimpleChatResponse struct {
	Choices []SimpleChatChoice `json:"choices"`
}

// Simple LLM service for testing
func (s *LlmService) SimpleAsk(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	client := resty.New()

	messages := []SimpleChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := map[string]interface{}{
		"model":    s.cfg.Model,
		"messages": messages,
		"stream":   false,
	}

	// Validate API key before making request
	if s.cfg.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set. Please set LLM_API_KEY in your .env file")
	}

	apiURL := s.cfg.BaseURL + "/v1/chat/completions"
	logrus.Infof("Making API request to: %s", apiURL)
	logrus.Infof("Using model: %s", s.cfg.Model)

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&SimpleChatResponse{}).
		SetDebug(true).
		Post(apiURL)

	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	logrus.Infof("Response status: %d", resp.StatusCode())

	if resp.StatusCode() != http.StatusOK {
		// Try to get response body for better error information
		body := resp.Body()
		return "", fmt.Errorf("API request failed with status: %d, response: %s", resp.StatusCode(), string(body))
	}

	result := resp.Result().(*SimpleChatResponse)
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	return result.Choices[0].Message.Content, nil
}

// Simple version of FullToolAwareChatFlow for debugging
func (s *LlmService) SimpleFullToolAwareChatFlow(ctx context.Context, systemPrompt, userPrompt, location string) (string, error) {
	// For now, just add weather info to the prompt instead of using tools
	enhancedPrompt := fmt.Sprintf("%s\n\n注意：请在生成的博客中包含 %s 的天气信息。", userPrompt, location)

	return s.SimpleAsk(ctx, systemPrompt, enhancedPrompt)
}
