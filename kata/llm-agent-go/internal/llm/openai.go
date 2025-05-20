package llm

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ChatRequest struct {
	Model       string      `json:"model"`
	Messages    []ChatEntry `json:"messages"`
	Stream      bool        `json:"stream"`
	Temperature float64     `json:"temperature,omitempty"` // Default: 1.0
	//TopP        float64     `json:"top_p,omitempty"`       // Default: 1.0
	//MaxTokens   int         `json:"max_tokens,omitempty"`  // Default: 4096
}

type ChatEntry struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatEntry `json:"message"`
	} `json:"choices"`
}

func createClient() (*http.Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: transport}, nil
}

func buildChatRequest(systemPrompt, userPrompt, model string, stream bool, temperature float64) (*ChatRequest, error) {
	return &ChatRequest{
		Model:       model,
		Stream:      stream,
		Messages: []ChatEntry{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: temperature,
	}, nil
}

func loadBaseConfig() (string, string, string, float64, error) {
	baseUrl := os.Getenv("LLM_BASE_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	temperatureStr := os.Getenv("LLM_TEMPERATURE")
	if temperatureStr == "" {
		temperatureStr = "1.0"
	}
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		temperature = 1.0
	}

	return baseUrl, apiKey, model, temperature, nil
}
func AskLLM(systemPrompt string, userPrompt string) (string, error) {
	baseUrl, apiKey, model, temperature, err := loadBaseConfig()
	if err != nil {
		return "", err
	}

	req, err := buildChatRequest(systemPrompt, userPrompt, model, false, temperature)
	if err != nil {
		return "", err
	}

	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client, err := createClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		log.Printf("Decode error: %v", err)
		return "", err
	}

	return out.Choices[0].Message.Content, nil
}

func AskLLMWithStream(systemPrompt string, userPrompt string, processChunk func(string)) error {
	baseUrl, apiKey, model, temperature, err := loadBaseConfig()
	if err != nil {
		return err
	}

	req, err := buildChatRequest(systemPrompt, userPrompt, model, true, temperature)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client, err := createClient()
	if err != nil {
		return err
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		trimmedLine := bytes.TrimSpace(line)
		if !bytes.HasPrefix(trimmedLine, []byte("data: ")) {
			continue
		}

		data := trimmedLine[6:]
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) {
			continue
		}

		var chunk map[string]interface{}
		if err := json.Unmarshal(data, &chunk); err != nil {
			log.Printf("JSON decode error: %v (raw data: %s)", err, data)
			continue
		}

		if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if delta, ok := choice["delta"].(map[string]interface{}); ok {
					if content, ok := delta["content"].(string); ok && content != "" {
						processChunk(content)
					}
				}
			}
		}
	}

	return nil
}