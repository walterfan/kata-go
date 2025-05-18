package main

import (
	"github.com/joho/godotenv"
	"log"
	"github.com/walterfan/llm-agent-go/cmd"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	cmd.Execute()
}
