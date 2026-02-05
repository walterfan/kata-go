package agent

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/walterfan/english-agent/internal/config"
	"github.com/walterfan/english-agent/internal/storage"
)

type Agent struct {
	model *openai.ChatModel
}

func NewAgent(ctx context.Context) (*Agent, error) {
	cfg := config.Get()

	// Configure HTTP client for TLS skipping if needed
	var httpClient *http.Client
	if cfg.AI.SkipTLS {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	// Temperature as pointer
	temperature := float32(0.7)

	// LLM Config
	llmCfg := &openai.ChatModelConfig{
		APIKey:      cfg.AI.APIKey,
		Model:       cfg.AI.Model,
		BaseURL:     cfg.AI.BaseURL,
		Temperature: &temperature,
	}

	if httpClient != nil {
		llmCfg.HTTPClient = httpClient
	}

	llm, err := openai.NewChatModel(ctx, llmCfg)
	if err != nil {
		return nil, err
	}

	return &Agent{
		model: llm,
	}, nil
}

func (a *Agent) Run(ctx context.Context, text, task string) (string, error) {
	// Check Cache
	inputStr := fmt.Sprintf("Text:\n%s\n\nTask: %s", text, task)
	hash := storage.HashRequest(inputStr)
	if cached, ok := storage.GetCache(hash); ok {
		return cached, nil
	}

	// Build messages
	messages := []*schema.Message{
		schema.SystemMessage(SystemPersona),
		schema.UserMessage(fmt.Sprintf(`Text:
%s

Task: %s

Rules:
- Be concise
- Use plain English
- Give examples
`, text, task)),
	}

	// Call the model
	response, err := a.model.Generate(ctx, messages)
	if err != nil {
		return "", err
	}

	result := response.Content

	// Save to Cache
	storage.SetCache(hash, result)

	return result, nil
}

// BuildMessages creates the message array for the LLM
func (a *Agent) BuildMessages(text, task string) []*schema.Message {
	return []*schema.Message{
		schema.SystemMessage(SystemPersona),
		schema.UserMessage(fmt.Sprintf(`Text:
%s

Task: %s

Rules:
- Be concise
- Use plain English
- Give examples
`, text, task)),
	}
}

// RunStream streams the LLM response through a channel
func (a *Agent) RunStream(ctx context.Context, text, task string) (<-chan string, <-chan error) {
	contentCh := make(chan string, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(contentCh)
		defer close(errCh)

		// Check Cache first
		inputStr := fmt.Sprintf("Text:\n%s\n\nTask: %s", text, task)
		hash := storage.HashRequest(inputStr)
		if cached, ok := storage.GetCache(hash); ok {
			contentCh <- cached
			return
		}

		// Build messages
		messages := a.BuildMessages(text, task)

		// Stream from the model
		stream, err := a.model.Stream(ctx, messages)
		if err != nil {
			errCh <- err
			return
		}

		var fullContent string
		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				chunk, err := stream.Recv()
				if err != nil {
					// EOF is expected when stream ends
					if err.Error() == "EOF" || fullContent != "" {
						// Save to cache
						if fullContent != "" {
							storage.SetCache(hash, fullContent)
						}
						return
					}
					errCh <- err
					return
				}

				if chunk != nil && chunk.Content != "" {
					contentCh <- chunk.Content
					fullContent += chunk.Content
				}
			}
		}
	}()

	return contentCh, errCh
}
