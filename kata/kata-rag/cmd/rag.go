package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sashabaranov/go-openai"

	"rag-assistant/internal/embedder"
	"rag-assistant/internal/loader"
	"rag-assistant/internal/vectorstore"
)

func main() {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			panic("Error loading .env file: " + err.Error())
		}
	}

	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	config := openai.DefaultConfig(os.Getenv("OPENAI_API_KEY"))
	config.BaseURL = os.Getenv("OPENAI_BASE_URL")
	openaiClient := openai.NewClientWithConfig(config)

	// Check if documents table is already populated
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM documents").Scan(&count)
	if err != nil {
		panic(err)
	}

	if count == 0 {
		// 1. Load code files
		files, err := loader.LoadCodeFiles(os.Getenv("LOCAL_REPO_PATH"))
		if err != nil {
			panic(err)
		}

		// 2. Embed and store in pgvector
		for _, f := range files {
			emb, err := embedder.GetEmbedding(ctx, f.Content)
			if err != nil {
				panic(err)
			}
			err = vectorstore.InsertEmbedding(ctx, db, f.Path, f.Content, emb)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println("✅ Code indexed.")
	} else {
		fmt.Println("⏭️  Documents already exist in the database. Skipping indexing.")
	}

	// 3. RAG loop (simple CLI)
	//for {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nAsk: ")
	question, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	question = strings.TrimSpace(question)

	embQ, _ := embedder.GetEmbedding(ctx, question)
	results, _ := vectorstore.SearchSimilar(ctx, db, embQ, 5)

	var sb strings.Builder
	sb.WriteString("You are a code assistant. Use the following context to answer:\n")
	for i, r := range results {
		sb.WriteString(fmt.Sprintf("[%d]: %s\n", i+1, r))
	}
	sb.WriteString(fmt.Sprintf("\nQuestion: %s", question))

	fmt.Printf("----\nPrompt: %s\n------\n", sb.String())
	resp, err := openaiClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: os.Getenv("OPENAI_MODEL"),
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are a coding assistant."},
			{Role: "user", Content: sb.String()},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("\n--- Answer ---")
	fmt.Println(resp.Choices[0].Message.Content)

	//}
}
