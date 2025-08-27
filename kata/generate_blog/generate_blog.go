package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// ----------- Config -------------
const blogTemplate = `# %s

> %s

##  “%s”
### what
%s
### why
%s
### how
%s
### example
%s
### summary
%s
### reference
%s

## Daily recommendation of 1 github hot project with a detailed explanation
%s

## Daily recommendation of 1 best practice in software development and AI applications
%s

## Daily practice of 1 common library by golang
%s

## Daily practice of 1 classic design pattern  by golang
%s

## Daily recitation of 10 English sentences with a detailed explanation in daily work
%s
`

// ----------- Utils -------------

func fetchGitHubTrending() (string, error) {
	// Simple crawler: use GitHub trending API proxy
	resp, err := http.Get("https://ghapi.huchen.dev/repositories?since=daily")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var repos []map[string]interface{}
	if err := json.Unmarshal(body, &repos); err != nil {
		return "", err
	}

	if len(repos) == 0 {
		return "No trending projects found.", nil
	}
	repo := repos[0]
	return fmt.Sprintf("[%s](%s): %s",
		repo["name"],
		repo["url"],
		repo["description"],
	), nil
}

func fetchBestPractice() string {
	// TODO: Replace with real curated source later
	return "Always write unit tests before committing production code. It reduces regression risk."
}

func callOpenAI(prompt string) (string, error) {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("LLM_API_KEY environment variable is not set")
	}

	client := openai.NewClient(apiKey)
	ctx := context.Background()

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gpt-4o-mini" // default fallback
	}

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a technical writer."},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}

// ----------- Main -------------

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: bloggen \"Your daily idea here\"")
		return
	}
	idea := os.Args[1]
	today := time.Now().Format("2006-01-02")

	// Call AI for different parts
	quote, _ := callOpenAI("Give me one famous English quote with Chinese explanation.")
	what, _ := callOpenAI("Explain WHAT about: " + idea + " in Chinese.")
	why, _ := callOpenAI("Explain WHY about: " + idea + " in Chinese.")
	how, _ := callOpenAI("Explain HOW to implement: " + idea + " in Chinese.")
	example, _ := callOpenAI("Give an EXAMPLE about: " + idea + " in Chinese with code snippet.")
	summary, _ := callOpenAI("Summarize: " + idea + " in Chinese.")
	reference, _ := callOpenAI("List 3 reference links about: " + idea)

	// External data
	githubProj, _ := fetchGitHubTrending()
	bestPractice := fetchBestPractice()
	golangLib, _ := callOpenAI("介绍一个 Go 常用库，带代码示例，用中文。")
	designPattern, _ := callOpenAI("用 Go 演示一个经典设计模式，带代码示例，用中文。")
	englishSentences, _ := callOpenAI("Give me 10 English sentences for daily work, with detailed explanation in Chinese.")

	// Fill template
	content := fmt.Sprintf(blogTemplate,
		"每日博客 "+today,
		quote,
		idea,
		what,
		why,
		how,
		example,
		summary,
		reference,
		githubProj,
		bestPractice,
		golangLib,
		designPattern,
		englishSentences,
	)

	// Save file
	fileName := "blog-" + today + ".md"
	if err := ioutil.WriteFile(fileName, []byte(content), 0644); err != nil {
		panic(err)
	}

	fmt.Println("✅ Blog generated:", fileName)
}
