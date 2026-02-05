# English Learning Agent 🤖

A powerful AI agent built with **Go (Eino)** and **Streamlit** to help you improve your English daily.

## Features

- **RSS Feed Reader**: Automatically fetches headlines from BBC and VOA Learning English.
- **AI Explanations**: Ask the agent to simplify text or explain meanings.
- **Vocabulary Extraction**: Extract useful phrases and sentence structures.
- **Review System**: Save and review your learning items.
- **Dual Interface**: Use via CLI or a beautiful Web UI.

## Tech Stack

- **Backend**: Go, CloudWeGo/Eino, Gin, Zap, Viper, SQLite
- **Frontend**: Python, Streamlit
- **AI**: OpenAI (GPT-4o-mini)

## Installation

### Prerequisites

- Go 1.21+
- Python 3.10+
- OpenAI API Key

### Backend Setup

1. Clone the repo
2. Install Go dependencies:
   ```bash
   go mod tidy
   ```
3. Configure `config.yaml`:
   ```yaml
   ai:
     provider: openai
     api_key: "YOUR_OPENAI_API_KEY"
     model: gpt-4o-mini
   ```

### Frontend Setup

1. Install Python dependencies:
   ```bash
   pip install -r web/requirements.txt
   ```

## Usage

### 1. Start the Backend Server

```bash
go run cmd/main.go start
```
Server runs at `http://localhost:8080`.

### 2. Start the Web UI

Open a new terminal:
```bash
streamlit run web/app.py
```

### 3. CLI Mode

You can also use the agent directly from the terminal:

```bash
# Explain text
go run cmd/main.go explain "The company is rolling out the feature."

# Simplify text
go run cmd/main.go simplify "The implementation of the new policy facilitated considerable improvements."
```

## Daily Workflow

1. Open the Web UI.
2. Click **Refresh Headlines** in the sidebar.
3. Click **Study this** on an interesting article.
4. Select **Extract Vocabulary** or ask "Simplify this".
5. Review your saved items in the **Review Items** tab.

## License

MIT

