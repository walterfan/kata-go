.PHONY: all build run ui clean test deps help

# Variables
BINARY_NAME=english-agent
GO_FILES=$(shell find . -name '*.go')

all: build

## 🏗️ Build
build: ## Build the Go binary
	@echo "Building backend..."
	@go build -o bin/$(BINARY_NAME) cmd/main.go

## 🚀 Run
run: ## Run the backend server
	@echo "Starting backend..."
	@go run cmd/main.go start

ui: ## Run the Streamlit UI
	@echo "Starting frontend..."
	@streamlit run web/app.py

start: ## Start both backend and frontend (requires split terminal or background)
	@./start.sh

explain: ## CLI: Explain text (usage: make explain text="your text")
	@go run cmd/main.go explain "$(text)"

simplify: ## CLI: Simplify text (usage: make simplify text="your text")
	@go run cmd/main.go simplify "$(text)"

## 🧹 Maintenance
clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f english_agent.db

test: ## Run Go tests
	@go test ./...

deps: ## Install dependencies (Go + Python)
	@echo "Installing Go dependencies..."
	@go mod tidy
	@echo "Installing Python dependencies..."
	@pip install -r web/requirements.txt

## ❓ Help
help: ## Show this help message
	@echo "Usage: make [target]"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

