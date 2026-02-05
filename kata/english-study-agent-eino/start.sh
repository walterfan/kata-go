#!/bin/bash

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤖 Starting English Learning Agent...${NC}"

# Function to cleanup background processes on exit
cleanup() {
    echo -e "\n${BLUE}Shutting down...${NC}"
    if [ -n "$GO_PID" ]; then
        kill $GO_PID 2>/dev/null
    fi
    exit 0
}

trap cleanup SIGINT SIGTERM

# Check for .env or config
if [ ! -f "config.yaml" ]; then
    echo "⚠️  config.yaml not found. Please create one with your API keys."
    exit 1
fi

# 1. Start Backend
echo -e "${GREEN}[1/2] Starting Go Backend (Gin + Eino)...${NC}"
go run cmd/main.go start &
GO_PID=$!

# Wait for backend to be ready (simple sleep for now, could check port)
echo "Waiting for backend to initialize..."
sleep 2

# 2. Start Frontend
echo -e "${GREEN}[2/2] Starting Streamlit UI...${NC}"
if command -v streamlit &> /dev/null; then
    streamlit run web/app.py
else
    echo "❌ Streamlit not found. Please run: pip install -r web/requirements.txt"
    kill $GO_PID
    exit 1
fi

# Wait for background process
wait $GO_PID

