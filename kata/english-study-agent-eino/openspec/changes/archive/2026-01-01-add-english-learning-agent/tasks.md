## 1. Project Setup (Eino & Go)
- [x] 1.1 Initialize Go module: `go mod init github.com/yourname/english-agent`
- [x] 1.2 Install Core Libs: `go get github.com/cloudwego/eino github.com/gin-gonic/gin`
- [x] 1.3 Install Utilities: `go get go.uber.org/zap` (Log), `github.com/spf13/viper` (Config), `github.com/mmcdole/gofeed` (RSS), `modernc.org/sqlite` (DB)
- [x] 1.4 Create project structure: `cmd/`, `internal/agent/`, `internal/api/`, `internal/rss/`, `internal/storage/`, `internal/logger/`, `web/`
- [x] 1.5 Setup Logger: Create `internal/logger/logger.go` with Zap configuration.
- [x] 1.6 Setup Config: Create `internal/config/config.go` with Viper to load `config.yaml`.

## 2. Basic Eino Agent Implementation (Day 1 Goal)
- [x] 2.1 Define Agent Persona: Create `internal/agent/prompts.go` with the specific System Prompt.
- [x] 2.2 Implement Basic Chain: Create `internal/agent/chain.go` using `chain.NewChain(prompt, llm)`.
- [x] 2.3 Implement Task Support: Update prompt template to accept `{{.text}}` and `{{.task}}`.
- [x] 2.4 Create CLI "Explain" Command: Wire up `main.go` to invoke the chain with user input.
- [x] 2.5 Verify "Simplify" and "Explain" tasks work with `gpt-4o-mini`.

## 3. Web UI & API Layer (Day 1-2 Goal)
- [x] 3.1 Implement Gin Server: Create `internal/api/server.go` to serve Agent requests using `gin.Default()`.
- [x] 3.2 Define API Routes: `POST /api/chat`, `POST /api/explain`, `GET /api/feeds`.
- [x] 3.3 Create Streamlit App: Create `web/app.py` with `st.chat_message` interface.
- [x] 3.4 Wire Frontend to Backend: Implement Python logic to call `localhost:8080/api/...`.
- [x] 3.5 Add "Start" Command: Update CLI to start both Gin server and Streamlit (via subprocess or instruction).

## 4. RSS Feed Integration
- [x] 4.1 Implement RSS Fetcher: `internal/rss/fetcher.go` using `gofeed`.
- [x] 4.2 Implement Article Parser: Clean HTML content from feeds.
- [x] 4.3 Update API: Add endpoint to list headlines and fetch article body.
- [x] 4.4 Update Streamlit: Add sidebar for RSS feed selection and article reading.

## 5. Advanced Agent Capabilities (Tools)
- [x] 5.1 Define Tool Interface: `internal/agent/tools/vocabulary.go` implementing Eino Tool interface.
- [x] 5.2 Implement Vocabulary Extraction Logic: Can be heuristic-based or a sub-chain.
- [x] 5.3 Update Agent to use Tools: Refactor `chain.NewChain` to `agent.NewAgent(llm, agent.WithTools(...))`.
- [x] 5.4 Update UI/CLI: Support "Extraction" actions in both interfaces.

## 6. Storage & Persistence (Day 2 Goal)
- [x] 6.1 Setup SQLite DB: `internal/storage/db.go` with schema for Articles and LearningItems.
- [x] 6.2 Implement "Save" functionality: Store extracted vocabulary to DB.
- [x] 6.3 Implement "Review" UI: Add Streamlit page to list saved phrases/structures.
- [x] 6.4 Implement Cache: Cache LLM responses in SQLite to save costs.

## 7. Polish & Release
- [x] 7.1 Add Configuration: Ensure `config.yaml` is fully mapped in Viper (API Keys, Feeds, Port).
- [x] 7.2 Add Interactive Mode: A simple REPL or TUI (Bubbletea optional) for the daily workflow.
- [x] 7.3 Write README with installation (Go+Python) and usage instructions.
