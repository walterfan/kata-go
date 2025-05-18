package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ChatRequest struct {
	Model    string      `json:"model"`
	Messages []ChatEntry `json:"messages"`
	Stream   bool        `json:"stream"`
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

func AskLLM(prompt string) (string, error) {
	baseUrl := os.Getenv("LLM_BASE_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	req := ChatRequest{
		Model: model,
		Stream: false,
		Messages: []ChatEntry{
			{Role: "user", Content: prompt},
		},
	}
	body, _ := json.Marshal(req)

	httpReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
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

func AskLLMWithStream(prompt string, processChunk func(string)) error {
	baseUrl := os.Getenv("LLM_BASE_URL")
	apiKey := os.Getenv("LLM_API_KEY")
	model := os.Getenv("LLM_MODEL")

	req := ChatRequest{
		Model:  model,
		Stream: true, // 🔥 Enable stream mode
		Messages: []ChatEntry{
			{Role: "user", Content: prompt},
		},
	}
	body, _ := json.Marshal(req)

	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", baseUrl), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Stream response line-by-line
	reader := bufio.NewReader(resp.Body)
	var contentBuilder []byte

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

		data := trimmedLine[6:] // Remove "data: " prefix
		if len(data) == 0 || bytes.Equal(data, []byte("[DONE]")) {
			continue
		}

		// Parse the chunk
		var chunk map[string]interface{}
		if err := json.Unmarshal(data, &chunk); err != nil {
			log.Printf("JSON decode error: %v (raw data: %s)", err, data)
			continue
		}

		// Extract delta content
		if choices, ok := chunk["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if delta, ok := choice["delta"].(map[string]interface{}); ok {
					if content, ok := delta["content"].(string); ok && content != "" {
						contentBuilder = append(contentBuilder, content...)
						processChunk(content) // Callback for handling each received chunk
					}
				}
			}
		}
	}

	return nil
}