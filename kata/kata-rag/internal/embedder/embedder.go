package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

/*
	curl -X POST http://127.0.0.1:11434/api/embeddings \
	  -H "Content-Type: application/json" \
	  -d '{
	    "model": "nomic-embed-text",
	    "prompt": "your text"
	  }'
*/
func GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	url := os.Getenv("EMBED_URL")
	payload := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
	}{
		Model:  os.Getenv("EMBED_MODEL"),
		Prompt: text,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Embedding, nil
}
func VectorToPGArray(vector []float32) string {
	strValues := make([]string, len(vector))
	for i, v := range vector {
		strValues[i] = fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(strValues, ","))
}
