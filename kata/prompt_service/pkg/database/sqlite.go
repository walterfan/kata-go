package database

import (
	"log"
	"github.com/walterfan/prompt-service/pkg/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("prompt.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	if err := DB.AutoMigrate(&models.Prompt{}); err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	InitData()
}

func InitData() {
	// Check if data already exists
	var count int64
	DB.Model(&models.Prompt{}).Count(&count)

	if count > 0 {
		log.Println("Sample data already exists. Skipping initialization.")
		return
	}

	prompts := []models.Prompt{
		{
			Name:         "Explain Goroutines",
			Description:  "Describe how goroutines work in Go.",
			SystemPrompt: "You are a Go language expert.",
			UserPrompt:   "What is a goroutine and how does it differ from a thread?",
			Tags:         "concurrency,golang",
		},
		{
			Name:         "Go Interface Usage",
			Description:  "Explain interfaces in Go with examples.",
			SystemPrompt: "You are a Go language expert.",
			UserPrompt:   "How do you define and use interfaces in Go?",
			Tags:         "interface,golang",
		},
		{
			Name:         "Explain Context Package",
			Description:  "Explain the purpose and usage of the context package in Go.",
			SystemPrompt: "You are a Go language expert.",
			UserPrompt:   "Why is the context package important in Go applications?",
			Tags:         "context,golang",
		},
		{
			Name:         "Go Modules Introduction",
			Description:  "Introduction to Go modules for dependency management.",
			SystemPrompt: "You are a Go language expert.",
			UserPrompt:   "What are Go modules and how do they help manage dependencies?",
			Tags:         "modules,dependency,golang",
		},
	}

	result := DB.Create(&prompts)
	if result.Error != nil {
		log.Fatal("Failed to insert initial data: ", result.Error)
	}

	log.Println("Initialized database with sample Go-related prompts.")
}