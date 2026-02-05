# Design: English Learning AI Agent (Eino Implementation)

## Context
This project builds an AI agent for English learners using **Go** and the **CloudWeGo/Eino** framework. The agent helps users improve English by reading RSS feeds, extracting learning content, and providing AI-assisted explanations.

Target users: Non-native software engineers/developers who want to improve their English efficiently (5-10 mins/day).

## Goals / Non-Goals

### Goals
- Provide a quick daily workflow (5-10 minutes) via CLI or Web UI
- Extract actionable learning items: 3 useful phrases + 1 sentence structure
- Support AI-powered simplification and explanation using Eino Chains
- **MVP Timeline**: 
  - Day 1: CLI + Basic Agent (Explain/Simplify)
  - Day 2: RSS Input + Storage
  - Day 3: Web UI (Streamlit + Gin)

### Non-Goals
- Complex Single Page Application (React/Vue) - use Streamlit for rapid UI
- Complex multi-agent orchestration (start with single agent + tools)

## Decisions

### Technology Stack
- **Framework**: [Eino](https://github.com/cloudwego/eino) (CloudWeGo)
- **HTTP Server**: [Gin](https://github.com/gin-gonic/gin) - to serve API for Streamlit
- **Web UI**: [Streamlit](https://streamlit.io/) (Python) - consumes Gin API
- **Language**: Go (Backend), Python (Frontend/UI only)
- **LLM**: OpenAI (gpt-4o-mini or gpt-4) via Eino's `openai` component
- **Storage**: SQLite (local)
- **RSS**: `gofeed` library
- **Logging**: [Zap](https://github.com/uber-go/zap) - High performance logging
- **Configuration**: [Viper](https://github.com/spf13/viper) - Configuration management

### Agent Persona (System Prompt)
The agent will use this core persona in Eino's `prompt.SystemMessage`:
```text
You are an English learning assistant for non-native software engineers.
You explain English clearly, simply, and practically.
You avoid complex grammar terms.
You always give examples.
```

### Eino Architecture Pattern
We will use Eino's **Chain** and **Agent** primitives:

1.  **Simple Chain** (for direct tasks like Simplification/Explanation):
    *   **Input**: Text + Task (e.g., "Explain for a junior developer")
    *   **Components**: `PromptTemplate` → `ChatModel` (LLM)
    *   **Output**: Simplified text or explanation

2.  **Agent with Tools** (for complex workflows):
    *   **Tools**:
        *   `VocabularyTool`: Extracts phrases and structures (can use simple heuristics or LLM sub-chain).
        *   `RSSReaderTool`: Fetches content from feeds (optional, can be outside agent loop initially).
    *   **Reasoning**: The agent decides whether to just explain text or extract vocabulary based on user input.

### Data Storage
- **SQLite** for:
  - `articles`: Content and metadata
  - `learning_items`: Extracted phrases/structures
  - `history`: User interaction log

## Architecture Diagram

```
User (Browser)        User (Terminal)
      ↓                     ↓
[ Streamlit App ]      [ CLI App ]
(Python Frontend)           │
      ↓ (HTTP/JSON)         │
      └──────────┐          │
                 ↓          ↓
           [ Gin HTTP Server ]
           (Go API Layer)
                 ↓
        [ Eino Chain / Agent ]
                 ↓
┌─────────────────┐
│  Prompt Template│ (Injects Persona + Task)
└────────┬────────┘
         ↓
┌─────────────────┐
│    LLM Model    │ (OpenAI / Claude)
└────────┬────────┘
         │
    ┌────▼────┐
    │  Tools  │ (Optional)
    ├─────────┤
    │ Vocab   │
    │ Grammar │
    └────┬────┘
         ↓
   Structured Output
         ↓
   SQLite Storage
```

## Risks / Trade-offs

### Risk: Complexity of Tools
- **Mitigation**: Start with a simple **Chain** for the MVP (Input -> LLM -> Output) and add **Tools** only for distinct operations like "Save to DB" or "Fetch RSS" if needed.

### Risk: LLM Cost
- **Mitigation**: Use `gpt-4o-mini` for bulk operations (extraction) and `gpt-4` only for complex explanations. Implement caching.

## Open Questions
- Should the RSS fetching be a "Tool" called by the agent, or a separate CLI command that feeds data *into* the agent? 
  - **Decision**: Keep RSS fetching as a separate CLI command (`english-agent list/read`) that feeds text into the Eino Chain. This is simpler and more predictable for a daily workflow.
