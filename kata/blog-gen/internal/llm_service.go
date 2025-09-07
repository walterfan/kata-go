package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type LlmConfig struct {
	BaseURL string
	APIKey  string
	Model   string
	Stream  bool
}

type ChatMessage struct {
	Role       string             `json:"role"`
	Content    string             `json:"content"`
	ToolCalls  []ToolFunctionCall `json:"tool_calls,omitempty"`
	ToolCallId string             `json:"tool_call_id,omitempty"`
	Name       string             `json:"name,omitempty"`
}

type ToolFunctionCall struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
}

type StreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role         string            `json:"role,omitempty"`
			Content      string            `json:"content,omitempty"`
			FunctionCall *ToolFunctionCall `json:"function_call,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason,omitempty"`
	} `json:"choices"`
}

type LlmService struct {
	cfg *LlmConfig
}

func NewLlmService(cfg *LlmConfig) *LlmService {
	return &LlmService{cfg: cfg}
}

func LoadLlmConfigFromEnv() *LlmConfig {
	cfg := &LlmConfig{
		BaseURL: getEnvOrDefault("LLM_BASE_URL", "https://api.openai.com"),
		APIKey:  os.Getenv("LLM_API_KEY"),
		Model:   getEnvOrDefault("LLM_MODEL", "gpt-4o"),
		Stream:  getEnvBoolOrDefault("LLM_STREAM", false),
	}

	// Validate required configuration
	if cfg.APIKey == "" {
		logrus.Warn("LLM_API_KEY is not set. Please set your OpenAI API key in the .env file.")
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

func (s *LlmService) Ask(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if s.cfg.Stream {
		return s.askWithStream(ctx, systemPrompt, userPrompt)
	}
	return s.askWithoutStream(ctx, systemPrompt, userPrompt)
}

func (s *LlmService) askWithoutStream(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	client := resty.New()

	messages := []ChatMessage{
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

	apiURL := s.cfg.BaseURL + "/chat/completions"
	logrus.Debugf("Making API request to: %s", apiURL)

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&ChatResponse{}).
		Post(apiURL)

	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		// Try to get response body for better error information
		body := resp.Body()
		return "", fmt.Errorf("API request failed with status: %d, response: %s", resp.StatusCode(), string(body))
	}

	result := resp.Result().(*ChatResponse)
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	choice := result.Choices[0]

	// Handle tool calls
	if len(choice.Message.ToolCalls) > 0 {
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Function.Name == "get_weather" {
				return s.handleWeatherToolCall(ctx, messages, &toolCall)
			}
		}
	}

	return choice.Message.Content, nil
}

func (s *LlmService) askWithStream(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	client := resty.New()

	messages := []ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	req := map[string]interface{}{
		"model":    s.cfg.Model,
		"messages": messages,
		"tools": []map[string]interface{}{
			{
				"type": "function",
				"function": map[string]interface{}{
					"name":        "get_weather",
					"description": "获取今日和明日的天气信息",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type":        "string",
								"description": "地理位置，如 Hefei, Beijing",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		"stream": true,
	}

	// Validate API key before making request
	if s.cfg.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key is not set. Please set LLM_API_KEY in your .env file")
	}

	apiURL := s.cfg.BaseURL + "/v1/chat/completions"
	logrus.Debugf("Making streaming API request to: %s", apiURL)

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetDebug(true).
		Post(apiURL)

	if err != nil {
		return "", fmt.Errorf("streaming API request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		// Try to get response body for better error information
		body := resp.Body()
		return "", fmt.Errorf("streaming API request failed with status: %d, response: %s", resp.StatusCode(), string(body))
	}

	var content strings.Builder
	// removed unused toolCall assembly for streaming

	reader := bytes.NewReader(resp.Body())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			logrus.Warnf("Failed to parse stream response: %v", err)
			continue
		}

		if len(streamResp.Choices) > 0 {
			choice := streamResp.Choices[0]

			if choice.Delta.Content != "" {
				content.WriteString(choice.Delta.Content)
				fmt.Print(choice.Delta.Content) // Print to console for real-time display
			}

			// streaming tool call chunks are ignored for now
		}
	}

	// no streaming tool-call handling; return accumulated content
	return content.String(), nil
}

func (s *LlmService) handleWeatherToolCall(ctx context.Context, messages []ChatMessage, toolCall *ToolFunctionCall) (string, error) {
	// Parse tool call arguments
	var args struct {
		Location string `json:"location"`
	}
	if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse tool call arguments: %w", err)
	}

	// Call weather API
	weatherText, err := s.callWeatherAPI(args.Location)
	if err != nil {
		return "", fmt.Errorf("weather API call failed: %w", err)
	}

	// Add tool call and result to messages
	messages = append(messages, ChatMessage{
		Role:      "assistant",
		Content:   "",
		ToolCalls: []ToolFunctionCall{*toolCall},
	})
	messages = append(messages, ChatMessage{
		Role:       "tool",
		ToolCallId: toolCall.Id,
		Content:    weatherText,
	})

	// Make final call to get complete response
	return s.makeFinalCall(ctx, messages)
}

func (s *LlmService) callWeatherAPI(location string) (string, error) {
	payload := map[string]string{"location": location}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:8080/tool/get_weather", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("weather tool server request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read weather response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse weather response: %w", err)
	}

	return fmt.Sprintf("Weather in %s: %s", result["location"], result["weather"]), nil
}

func (s *LlmService) makeFinalCall(ctx context.Context, messages []ChatMessage) (string, error) {
	client := resty.New()

	req := map[string]interface{}{
		"model":    s.cfg.Model,
		"messages": messages,
		"stream":   s.cfg.Stream,
	}

	if s.cfg.Stream {
		return s.makeFinalStreamCall(ctx, client, req)
	}

	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&ChatResponse{}).
		SetDebug(true).
		Post(s.cfg.BaseURL + "/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("final API call failed: %w", err)
	}

	result := resp.Result().(*ChatResponse)
	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("final response missing content")
}

func (s *LlmService) makeFinalStreamCall(ctx context.Context, client *resty.Client, req map[string]interface{}) (string, error) {
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.cfg.APIKey).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Post(s.cfg.BaseURL + "/v1/chat/completions")

	if err != nil {
		return "", fmt.Errorf("final streaming API call failed: %w", err)
	}

	var content strings.Builder
	reader := bytes.NewReader(resp.Body())
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
			content.WriteString(streamResp.Choices[0].Delta.Content)
			fmt.Print(streamResp.Choices[0].Delta.Content)
		}
	}

	return content.String(), nil
}

func (s *LlmService) FullToolAwareChatFlow(ctx context.Context, systemPrompt, userPrompt, location string) (string, error) {
	// Set default location if not provided
	if location == "" {
		location = "Hefei"
	}

	// Start tool server if not already running
	StartToolServer()

	// Wait a moment for the server to start
	time.Sleep(100 * time.Millisecond)

	return s.Ask(ctx, systemPrompt, userPrompt)
}
