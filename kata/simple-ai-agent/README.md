# Building an AI Agent with Golang and OpenAI API

## Project Structure

```
├── main.go
├── agent/
│   └── agent.go
├── tools/
│   ├── search.go
│   └── mindmap.go
├── go.mod
└── go.sum
```

## Step 1: Set Up the Project

```bash
mkdir file-ai-agent
cd file-ai-agent
go mod init github.com/yourusername/file-ai-agent
go get github.com/sashabaranov/go-openai
```

## Step 2: Implement the Tools

### tools/search.go

```go
package tools

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

// SearchResult represents a file search result
type SearchResult struct {
    FilePath string
    Content  string
    Size     int64
}

// SearchFileAndRead searches for files containing the keyword and returns their content
func SearchFileAndRead(directory, keyword string, maxResults int) ([]SearchResult, error) {
    results := []SearchResult{}
    
    err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Skip directories
        if info.IsDir() {
            return nil
        }
        
        // Check if it's a text file (you can add more extensions)
        ext := strings.ToLower(filepath.Ext(path))
        if ext != ".txt" && ext != ".md" && ext != ".go" && ext != ".py" && ext != ".js" && ext != ".html" && ext != ".csv" {
            return nil
        }
        
        // Check file size (avoid reading very large files)
        if info.Size() > 10*1024*1024 { // 10MB limit
            return nil
        }
        
        // Read file content
        content, err := ioutil.ReadFile(path)
        if err != nil {
            return nil
        }
        
        // Check if content contains the keyword
        if strings.Contains(strings.ToLower(string(content)), strings.ToLower(keyword)) {
            results = append(results, SearchResult{
                FilePath: path,
                Content:  string(content),
                Size:     info.Size(),
            })
            
            // Stop if we reached the maximum number of results
            if len(results) >= maxResults {
                return filepath.SkipAll
            }
        }
        
        return nil
    })
    
    return results, err
}
```

### tools/mindmap.go

```go
package tools

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "strings"
)

// MindMapNode represents a node in the mind map
type MindMapNode struct {
    Title    string         `json:"title"`
    Children []*MindMapNode `json:"children,omitempty"`
}

// DrawMindmap generates a mind map from the provided structure
// It saves the mind map as a Markdown file and optionally renders it
func DrawMindmap(rootNode *MindMapNode, outputPath string) error {
    // Convert the mind map to Markdown format
    markdown := convertToMarkdown(rootNode, 0)
    
    // Save to file
    err := os.WriteFile(outputPath, []byte(markdown), 0644)
    if err != nil {
        return fmt.Errorf("failed to save mind map: %v", err)
    }
    
    fmt.Printf("Mind map saved to %s\n", outputPath)
    
    // If Markmap CLI is installed, render the mind map (optional)
    if hasCommand("markmap") {
        cmd := exec.Command("markmap", outputPath)
        err = cmd.Run()
        if err != nil {
            return fmt.Errorf("failed to render mind map: %v", err)
        }
        fmt.Println("Mind map rendered with Markmap")
    } else {
        fmt.Println("To visualize the mind map, install Markmap: npm install -g markmap-cli")
        fmt.Println("Then run: markmap " + outputPath)
    }
    
    return nil
}

// Convert the mind map structure to Markdown
func convertToMarkdown(node *MindMapNode, level int) string {
    if node == nil {
        return ""
    }
    
    indent := strings.Repeat("#", level+1)
    result := fmt.Sprintf("%s %s\n\n", indent, node.Title)
    
    for _, child := range node.Children {
        result += convertToMarkdown(child, level+1)
    }
    
    return result
}

// Check if a command is available
func hasCommand(cmd string) bool {
    _, err := exec.LookPath(cmd)
    return err == nil
}
```

## Step 3: Implement the AI Agent

### agent/agent.go

```go
package agent

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "strings"

    "github.com/sashabaranov/go-openai"
    "github.com/yourusername/file-ai-agent/tools"
)

type Agent struct {
    client     *openai.Client
    systemPrompt string
}

func NewAgent(apiKey string) *Agent {
    client := openai.NewClient(apiKey)
    
    systemPrompt := `You are an AI assistant that helps users search through files, summarize content, and create mind maps.
You have access to the following tools:
1. SearchFileAndRead(directory, keyword, maxResults) - Searches for files containing the keyword and returns their content.
2. DrawMindmap(rootNode, outputPath) - Generates a mind map from the provided structure and saves it.

Always respond in a helpful, concise manner. When creating mind maps, organize information hierarchically.`

    return &Agent{
        client:     client,
        systemPrompt: systemPrompt,
    }
}

func (a *Agent) ProcessQuery(query string) (string, error) {
    messages := []openai.ChatCompletionMessage{
        {
            Role:    openai.ChatMessageRoleSystem,
            Content: a.systemPrompt,
        },
        {
            Role:    openai.ChatMessageRoleUser,
            Content: query,
        },
    }

    resp, err := a.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model:    openai.GPT4,
            Messages: messages,
            Temperature: 0.2,
            MaxTokens: 1000,
        },
    )
    if err != nil {
        return "", err
    }

    aiResponse := resp.Choices[0].Message.Content
    
    // Check if the response contains a tool call
    if strings.Contains(aiResponse, "SearchFileAndRead") {
        return a.handleSearchFileCall(aiResponse, messages)
    } else if strings.Contains(aiResponse, "DrawMindmap") {
        return a.handleDrawMindmapCall(aiResponse)
    }
    
    return aiResponse, nil
}

func (a *Agent) handleSearchFileCall(aiResponse string, messages []openai.ChatCompletionMessage) (string, error) {
    // Extract parameters from the response
    startIdx := strings.Index(aiResponse, "SearchFileAndRead(")
    if startIdx == -1 {
        return aiResponse, nil
    }
    
    startIdx += len("SearchFileAndRead(")
    endIdx := strings.Index(aiResponse[startIdx:], ")")
    if endIdx == -1 {
        return aiResponse, nil
    }
    
    paramsStr := aiResponse[startIdx : startIdx+endIdx]
    paramsParts := strings.Split(paramsStr, ",")
    if len(paramsParts) < 3 {
        return "", errors.New("invalid search parameters")
    }
    
    directory := strings.Trim(paramsParts[0], " \"'")
    keyword := strings.Trim(paramsParts[1], " \"'")
    maxResults := 5 // Default value
    
    // Execute the search
    results, err := tools.SearchFileAndRead(directory, keyword, maxResults)
    if err != nil {
        return "", fmt.Errorf("search failed: %v", err)
    }
    
    // Format search results
    var resultContent string
    if len(results) == 0 {
        resultContent = "No files found containing the keyword: " + keyword
    } else {
        resultContent = fmt.Sprintf("Found %d files containing '%s':\n\n", len(results), keyword)
        for i, result := range results {
            resultContent += fmt.Sprintf("File %d: %s (%.2f KB)\n", i+1, result.FilePath, float64(result.Size)/1024)
            
            // Add file content preview (first 1000 chars)
            preview := result.Content
            if len(preview) > 1000 {
                preview = preview[:1000] + "... (content truncated)"
            }
            resultContent += "Content:\n" + preview + "\n\n"
        }
    }
    
    // Ask the AI to summarize the search results
    messages = append(messages, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleAssistant,
        Content: aiResponse,
    })
    
    messages = append(messages, openai.ChatCompletionMessage{
        Role:    openai.ChatMessageRoleUser,
        Content: "Here are the search results. Please summarize them:\n\n" + resultContent,
    })
    
    resp, err := a.client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model:    openai.GPT4,
            Messages: messages,
            Temperature: 0.2,
            MaxTokens: 1500,
        },
    )
    if err != nil {
        return "", err
    }
    
    return resp.Choices[0].Message.Content, nil
}

func (a *Agent) handleDrawMindmapCall(aiResponse string) (string, error) {
    // Extract parameters from the response
    startIdx := strings.Index(aiResponse, "DrawMindmap(")
    if startIdx == -1 {
        return aiResponse, nil
    }
    
    startIdx += len("DrawMindmap(")
    endIdx := strings.Index(aiResponse[startIdx:], ")")
    if endIdx == -1 {
        return aiResponse, nil
    }
    
    paramsStr := aiResponse[startIdx : startIdx+endIdx]
    
    // The first parameter should be a JSON object representing the mind map structure
    // Let's find where the JSON starts and ends
    jsonStartIdx := strings.Index(paramsStr, "{")
    if jsonStartIdx == -1 {
        return "", errors.New("invalid mindmap structure: missing JSON object")
    }
    
    // Find the matching closing brace by counting braces
    jsonEndIdx := -1
    braceCount := 0
    for i, char := range paramsStr[jsonStartIdx:] {
        if char == '{' {
            braceCount++
        } else if char == '}' {
            braceCount--
            if braceCount == 0 {
                jsonEndIdx = jsonStartIdx + i + 1
                break
            }
        }
    }
    
    if jsonEndIdx == -1 {
        return "", errors.New("invalid mindmap structure: unbalanced braces")
    }
    
    // Extract the JSON string and parse it
    jsonStr := paramsStr[jsonStartIdx:jsonEndIdx]
    var rootNode tools.MindMapNode
    err := json.Unmarshal([]byte(jsonStr), &rootNode)
    if err != nil {
        return "", fmt.Errorf("failed to parse mindmap structure: %v", err)
    }
    
    // Extract the output path
    outputPathParts := strings.Split(paramsStr[jsonEndIdx:], ",")
    if len(outputPathParts) < 1 {
        return "", errors.New("missing output path parameter")
    }
    
    outputPath := strings.Trim(outputPathParts[len(outputPathParts)-1], " \"'")
    if !strings.HasSuffix(outputPath, ".md") {
        outputPath += ".md"
    }
    
    // Generate the mind map
    err = tools.DrawMindmap(&rootNode, outputPath)
    if err != nil {
        return "", fmt.Errorf("failed to generate mindmap: %v", err)
    }
    
    return fmt.Sprintf("Mind map has been created and saved to %s. %s", 
        outputPath, 
        aiResponse), nil
}
```

## Step 4: Implement the Main Program

### main.go

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/yourusername/file-ai-agent/agent"
)

func main() {
    // Get OpenAI API key from environment variable
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("Please set the OPENAI_API_KEY environment variable")
        return
    }

    fmt.Println("Welcome to the AI File Assistant!")
    fmt.Println("I can help you search files, summarize content, and create mind maps.")
    fmt.Println("Type 'exit' to quit.")
    fmt.Println()

    // Create an agent
    fileAgent := agent.NewAgent(apiKey)
    
    // Start the interaction loop
    scanner := bufio.NewScanner(os.Stdin)
    for {
        fmt.Print("> ")
        if !scanner.Scan() {
            break
        }
        
        query := scanner.Text()
        if strings.ToLower(query) == "exit" {
            break
        }
        
        response, err := fileAgent.ProcessQuery(query)
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            continue
        }
        
        fmt.Println("\nAI Assistant:")
        fmt.Println(response)
        fmt.Println()
    }
}
```

## How to Use the Program

1. Set up your OpenAI API key as an environment variable:
   ```bash
   export OPENAI_API_KEY=your_api_key_here
   ```

2. Build and run the program:
   ```bash
   go build
   ./file-ai-agent
   ```

3. Example interactions:

   ```
   > Search for "golang" in my projects folder
   
   > Summarize the file search results and create a mindmap
   ```

## Notes on Mind Map Visualization

For mind map visualization, the program saves the mind map as a Markdown file. You can use the [markmap-cli](https://www.npmjs.com/package/markmap-cli) tool to render it:

```bash
npm install -g markmap-cli
markmap mindmap.md
```

This will generate an HTML file and open it in your browser.

## Enhancements You Could Add

1. Add support for more file formats (PDF, Word, etc.)
2. Implement more complex search criteria
3. Add a web interface using a framework like Gin
4. Add caching to avoid re-processing the same files
5. Implement pagination for search results
6. Add error handling and retries for API calls

This implementation provides a solid foundation for an AI-powered file assistant that can search, summarize, and visualize content in mind maps.