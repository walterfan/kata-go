package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"

	"github.com/walterfan/llm-agent-go/internal/llm"
)

// Request structure for incoming JSON
type PromptRequest struct {
	Name           string `json:"name"`
	SystemPrompt   string `json:"system_prompt"`
	UserPrompt     string `json:"user_prompt"`
	OutputLanguage string `json:"output_language,omitempty"`
}

type APIRequest struct {
	Prompt   PromptRequest `json:"prompt"`
	Stream   bool          `json:"stream"`
	CodePath string        `json:"code_path"`
}

// Response structure for outgoing JSON
type APIResponse struct {
	Answer string `json:"answer"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins
	},
}

func handleWebSocket(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("Failed to upgrade connection")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	logger.Info("handleWebSocket")

	// Read initial prompt data from WebSocket
	var req APIRequest
	err = conn.ReadJSON(&req)
	if err != nil {
		logger.Error("Invalid request")
		conn.WriteJSON(gin.H{"error": "Invalid request"})
		return
	}
	logger.Sugar().Infof("received request: %v", req)
	// Read code file
	code, err := ioutil.ReadFile(req.CodePath)
	if err != nil {
		conn.WriteJSON(gin.H{"error": "Failed to read code file"})
		return
	}

	// Build final prompt
	selectedPrompt, err := getPromptConfigByName(req.Prompt.Name)
	//append if it not exist in user_prompt: ```\n\nnote: please use {{output_language}} to output

	if !strings.Contains(selectedPrompt.UserPrompt, "{{output_language}}") {
		selectedPrompt.UserPrompt += "\n\nnote: please use {{output_language}} to output"
	}
	//and repleace {{output_language}} with req.OutputLanguage
	selectedPrompt.UserPrompt = strings.ReplaceAll(selectedPrompt.UserPrompt, "{{output_language}}", req.Prompt.OutputLanguage)

	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	promptText := selectedPrompt.UserPrompt
	promptText = strings.Replace(promptText, "{{code}}", string(code), -1)
	if req.Prompt.OutputLanguage != "" {
		promptText = strings.Replace(promptText, "{{output_language}}", req.Prompt.OutputLanguage, 1)
	}

	// Stream response via WebSocket
	err = llm.AskLLMWithStream(selectedPrompt.SystemPrompt, promptText, func(chunk string) {
		_ = conn.WriteMessage(websocket.TextMessage, []byte(chunk))
	})

	if err != nil {
		_ = conn.WriteJSON(gin.H{"error": "LLM processing failed"})
	}
}

// processLLM handles the LLM request logic
func processLLM(c *gin.Context) {
	var req APIRequest
	var err error
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Build prompt text
	var promptText string
	if req.Prompt.Name != "" {
		selectedPrompt, err := getPromptConfigByName(req.Prompt.Name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		req.Prompt.SystemPrompt = selectedPrompt.SystemPrompt
		req.Prompt.UserPrompt = selectedPrompt.UserPrompt
	}
	promptText = req.Prompt.UserPrompt

	// Read code from file
	if req.CodePath != "" && strings.Contains(promptText, "{{code}}") {
		code, err := ioutil.ReadFile(req.CodePath)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read code file"})
			return
		}
		promptText = strings.Replace(promptText, "{{code}}", string(code), -1)
	}

	if !strings.Contains(promptText, "{{output_language}}") {
		promptText += "\n\nnote: please use {{output_language}} to output"
	}

	if req.Prompt.OutputLanguage != "" {
		promptText = strings.Replace(promptText, "{{output_language}}", req.Prompt.OutputLanguage, 1)
	} else {
		promptText = strings.Replace(promptText, "{{output_language}}", "English", 1)
	}

	// Call LLM
	var answer string

	if req.Stream {
		err = llm.AskLLMWithStream(req.Prompt.SystemPrompt, promptText, func(chunk string) {
			answer += chunk
		})
	} else {
		answer, err = llm.AskLLM(req.Prompt.SystemPrompt, promptText)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM processing failed"})
		return
	}

	// Return simplified response without usage
	c.JSON(http.StatusOK, APIResponse{
		Answer: answer,
	})
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start a web server to handle LLM requests via HTTP",
	Run: func(cmd *cobra.Command, args []string) {

		r := gin.Default()

		r.GET("/api/v1/stream", handleWebSocket)

		// Define POST endpoint
		r.POST("/api/v1/process", processLLM)

		r.GET("/api/v1/commands", func(c *gin.Context) {
			var promptConfigs []PromptConfig

			for _, name := range cachedPromptList {
				item := cachedPromptMap[name].(map[string]interface{})
				promptConfigs = append(promptConfigs, PromptConfig{
					Name:         name,
					Description:  item["description"].(string),
					SystemPrompt: item["system_prompt"].(string),
					UserPrompt:   item["user_prompt"].(string),
					Tags:         item["tags"].(string),
				})
			}

			c.JSON(http.StatusOK, promptConfigs)
		})
		// Fallback middleware: serve static files only if no API route matched
		r.Use(func(c *gin.Context) {

			// Try to serve static file
			filepath := "./web" + c.Request.URL.Path
			if _, err := os.Stat(filepath); err == nil {
				c.File(filepath)
				c.Abort()
			} else {
				c.Next()
			}
		})
		// Start server on user-defined port
		addr := fmt.Sprintf(":%s", port)
		fmt.Printf("Starting web server on %s\n", addr)
		if err := r.Run(addr); err != nil {
			panic(err)
		}
	},
}

func init() {
	if err := InitConfig(); err != nil {
		panic(err)
	}

	webCmd.Flags().StringVarP(&port, "port", "p", "8080", "Specify custom port for the web server")
	rootCmd.AddCommand(webCmd)
}
