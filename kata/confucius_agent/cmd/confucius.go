
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"

)

// Define custom types for components similar to Python's LangChain
type Document struct {
	PageContent string
	Metadata    map[string]interface{}
}

type Message struct {
	Role    string
	Content string
}

type ChatHistory struct {
	Messages []Message
}

type VectorStore struct {
	Documents []Document
	// In a real implementation, this would include embeddings and vector storage
}

// LazyLlmAgent represents the main agent structure
type LazyLlmAgent struct {
	SystemPrompt string
	SessionID    string
	MaxTokens    int
	ModelName    string
	ApiKey       string
	ApiBase      string
	Store        map[string]*ChatHistory
	VectorStore  *VectorStore
	OpenAIClient *openai.Client
}

// NewLazyLlmAgent creates a new instance of LazyLlmAgent
func NewLazyLlmAgent(systemPrompt, sessionID string, maxTokens int) *LazyLlmAgent {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		log.Fatal("LLM_API_KEY not set in environment")
	}

	// Set OpenAI API environment
	os.Setenv("OPENAI_API_KEY", apiKey)
	os.Setenv("USER_AGENT", "waltertest")

	modelName := os.Getenv("LLM_MODEL")
	apiBase := os.Getenv("LLM_BASE_URL")

	// Initialize OpenAI client
	config := openai.DefaultConfig(apiKey)
	if apiBase != "" {
		config.BaseURL = apiBase
	}
	client := openai.NewClientWithConfig(config)

	return &LazyLlmAgent{
		SystemPrompt: systemPrompt,
		SessionID:    sessionID,
		MaxTokens:    maxTokens,
		ModelName:    modelName,
		ApiKey:       apiKey,
		ApiBase:      apiBase,
		Store:        make(map[string]*ChatHistory),
		VectorStore:  &VectorStore{Documents: []Document{}},
		OpenAIClient: client,
	}
}

// LoadDocument loads a document into the vector store
func (a *LazyLlmAgent) LoadDocument(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// In a real implementation, you'd split the text and generate embeddings
	// For now, we'll just store the raw text as a document
	a.VectorStore.Documents = append(a.VectorStore.Documents, Document{
		PageContent: string(content),
		Metadata:    map[string]interface{}{"source": filePath},
	})

	fmt.Printf("Loaded %s to vector store successfully\n", filePath)
	return nil
}

// GetSessionHistory gets or creates a chat history for a session
func (a *LazyLlmAgent) GetSessionHistory(sessionID string) *ChatHistory {
	if _, exists := a.Store[sessionID]; !exists {
		a.Store[sessionID] = &ChatHistory{Messages: []Message{}}
	}
	return a.Store[sessionID]
}

// FormatDocs formats documents for context
func (a *LazyLlmAgent) FormatDocs(docs []Document) string {
	var contents []string
	for _, doc := range docs {
		contents = append(contents, doc.PageContent)
	}
	return strings.Join(contents, "\n\n")
}

// SimpleRetriever simulates document retrieval based on query
// In a real implementation, this would use vector similarity search
func (a *LazyLlmAgent) SimpleRetriever(query string) []Document {
	// For demonstration, return all documents
	// In a real implementation, this would find relevant documents using embeddings
	return a.VectorStore.Documents
}

// ProcessInput processes user input and returns the agent's response
func (a *LazyLlmAgent) ProcessInput(userInput string) (string, error) {
	// Get session history
	history := a.GetSessionHistory(a.SessionID)
	
	// Retrieve relevant context
	retrievedDocs := a.SimpleRetriever(userInput)
	contextDocs := a.FormatDocs(retrievedDocs)

	// Build messages for OpenAI API
	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: fmt.Sprintf(`You are an assistant for question-answering tasks.
			Use the following pieces of retrieved context to answer the question. If you don't know the answer, just say that you don't know.
			Context: %s
			
			%s`, contextDocs, a.SystemPrompt),
		},
	}
	
	// Add history messages
	for _, msg := range history.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add current user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	})

	ctx := context.Background()
	// Create the completion request
	resp, err := a.OpenAIClient.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:     a.ModelName,
			Messages:  messages,
			MaxTokens: a.MaxTokens,
			Stream:    true,
		},
	)
	
	if err != nil {
		return "", fmt.Errorf("completion error: %v", err)
	}
	defer resp.Close()
	
	// Process the streaming response
	var fullResponse strings.Builder
	for {
		response, err := resp.Recv()
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return fullResponse.String(), fmt.Errorf("stream error: %v", err)
		}
		
		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content
			fmt.Print(content)
			fullResponse.WriteString(content)
		}
	}
	fmt.Println() // New line after response

	// Update history
	history.Messages = append(history.Messages, Message{Role: openai.ChatMessageRoleUser, Content: userInput})
	history.Messages = append(history.Messages, Message{Role: openai.ChatMessageRoleAssistant, Content: fullResponse.String()})
	
	return fullResponse.String(), nil
}

// RunInteractive runs the agent in interactive mode
func (a *LazyLlmAgent) RunInteractive() {
	fmt.Printf("Starting interactive session with ID: %s\n", a.SessionID)
	fmt.Println("Type 'exit' to end the session")
	
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You:> ")
		if !scanner.Scan() {
			break
		}
		
		userInput := scanner.Text()
		if strings.ToLower(userInput) == "exit" {
			break
		}
		if userInput == "" {
			continue
		}
		
		_, err := a.ProcessInput(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func runAgent() {
	app := &cli.App{
		Name:  "LazyLlmAgent",
		Usage: "A conversational AI agent",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "role",
				Aliases: []string{"r"},
				Value:   "你现在扮演孔子的角色，尽量按照孔子的风格回复，不要出现'子曰'",
				Usage:   "System prompt for the agent",
			},
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Value:   "waltertest",
				Usage:   "Session ID for conversation history",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "analects.txt",
				Usage:   "File path for document loading",
			},
		},
		Action: func(c *cli.Context) error {
			agent := NewLazyLlmAgent(
				c.String("role"),
				c.String("session"),
				4096,
			)
			
			err := agent.LoadDocument(c.String("file"))
			if err != nil {
				return fmt.Errorf("failed to load document: %v", err)
			}
			
			agent.RunInteractive()
			return nil
		},
	}
	
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var confuciusCmd = &cobra.Command{
	Use:   "confucius",
	Short: "Confucius say",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAgent()
	},
}

func init() {
	rootCmd.AddCommand(confuciusCmd)
}
