package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/siddontang/go/log"
)

type QuestionForm struct {
	Question string `form:"question"`
	Answer   string
}

type DeepSeekRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type LlmConfig struct {
	ApiKey string `json:"apiKey"`
	Model  string `json:"model"`
	BaseUrl string `json:"baseUrl"`
}

func NewLlmConfig() *LlmConfig {
	return &LlmConfig{
		ApiKey: os.Getenv("LLM_API_KEY"),
		Model: os.Getenv("LLM_MODEL"),
		BaseUrl: os.Getenv("LLM_BASE_URL"),
	}

}

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			panic("Error loading .env file: " + err.Error())
		}
	}

	// Set up Gin router
	r := gin.Default()
	r.Use(BasicAuthMiddleware())
	r.SetFuncMap(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	})
	r.LoadHTMLGlob("templates/*")
	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.POST("/ask", func(c *gin.Context) {
		var form QuestionForm
		if err := c.ShouldBind(&form); err != nil {
			c.HTML(400, "index.html", gin.H{"error": "Invalid request"})
			return
		}

		// Call DeepSeek
		ctx := context.Background()
		llmConfig := NewLlmConfig()
		answer, err := callDeepSeek(ctx, llmConfig, form.Question)
		if err != nil {
			log.Errorf("Error calling DeepSeek: %s, %s, %s, %v", baseUrl, llmModel, apiKey, err)
			c.HTML(500, "index.html", gin.H{"error": "Error generating response"})
			return
		}

		form.Answer = answer

		c.HTML(200, "index.html", gin.H{"form": form})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func callDeepSeek(ctx context.Context, llmConfig *LlmConfig, question string) (string, error) {
	requestBody := DeepSeekRequest{
		Model: llmConfig.Model,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: question},
		},
		Stream: false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", llmConfig.BaseUrl + "/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + llmConfig.ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API returned non-OK status: %d", resp.StatusCode)
	}

	var response DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("No choices in response from DeepSeek API")
	}

	return response.Choices[0].Message.Content, nil
}

func BasicAuthMiddleware() gin.HandlerFunc {
	expectedUser := os.Getenv("BASIC_AUTH_USER")
	expectedPass := os.Getenv("BASIC_AUTH_PASS")

	log.Infof("basic auth : %s:%s", expectedUser, expectedPass)

	return gin.BasicAuth(gin.Accounts{
		expectedUser: expectedPass,
	})
}
