package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/walterfan/english-agent/internal/agent"
	"github.com/walterfan/english-agent/internal/api"
	"github.com/walterfan/english-agent/internal/config"
	"github.com/walterfan/english-agent/internal/logger"
	"github.com/walterfan/english-agent/internal/storage"
)

func main() {
	// 1. Load .env
	if err := godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist, we might be using real env vars
		// But let's log it just in case
		// log.Println("No .env file found")
	}

	// 2. Initialize
	config.Init()
	logger.Init()
	storage.Init("english_agent.db")
	defer storage.Close()

	// 3. Parse flags
	explainCmd := flag.NewFlagSet("explain", flag.ExitOnError)
	simplifyCmd := flag.NewFlagSet("simplify", flag.ExitOnError)
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("expected 'start', 'explain' or 'simplify' subcommands")
		os.Exit(1)
	}

	ctx := context.Background()

	switch os.Args[1] {
	case "start":
		startCmd.Parse(os.Args[2:])
		runServer(ctx)
	case "explain":
		a, err := agent.NewAgent(ctx)
		if err != nil {
			log.Fatalf("failed to create agent: %v", err)
		}
		explainCmd.Parse(os.Args[2:])
		text := strings.Join(explainCmd.Args(), " ")
		if text == "" {
			fmt.Println("please provide text to explain")
			return
		}
		runAgent(ctx, a, text, "Explain the meaning and list useful vocabulary")
	case "simplify":
		a, err := agent.NewAgent(ctx)
		if err != nil {
			log.Fatalf("failed to create agent: %v", err)
		}
		simplifyCmd.Parse(os.Args[2:])
		text := strings.Join(simplifyCmd.Args(), " ")
		if text == "" {
			fmt.Println("please provide text to simplify")
			return
		}
		runAgent(ctx, a, text, "Rewrite in simpler English for a beginner")
	default:
		fmt.Println("expected 'start', 'explain' or 'simplify' subcommands")
		os.Exit(1)
	}
}

func runAgent(ctx context.Context, a *agent.Agent, text, task string) {
	fmt.Printf("🤔 Thinking... (Task: %s)\n", task)
	result, err := a.Run(ctx, text, task)
	if err != nil {
		log.Fatalf("agent error: %v", err)
	}
	fmt.Println("\n🤖 Response:")
	fmt.Println(result)
}

func runServer(ctx context.Context) {
	s, err := api.NewServer(ctx)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	fmt.Println("🚀 Starting Gin server...")
	fmt.Println("👉 Run 'streamlit run web/app.py' in another terminal to start UI")
	if err := s.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
