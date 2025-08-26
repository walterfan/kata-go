#!/bin/bash

echo "Blog Generator API Diagnosis"
echo "==========================="

# Check if .env file exists
if [ ! -f .env ]; then
    echo "❌ .env file not found"
    echo "Please create one with: cp env.example .env"
    exit 1
fi

echo "✅ .env file found"

# Check if API key is set
API_KEY=$(grep LLM_API_KEY .env | cut -d'=' -f2)
if [ -z "$API_KEY" ] || [ "$API_KEY" = "your_openai_api_key_here" ]; then
    echo "❌ LLM_API_KEY not set or still has default value"
    echo "Please edit .env and set your actual OpenAI API key"
    exit 1
fi

echo "✅ LLM_API_KEY is set"

# Check other environment variables
BASE_URL=$(grep LLM_BASE_URL .env | cut -d'=' -f2)
MODEL=$(grep LLM_MODEL .env | cut -d'=' -f2)

echo "Base URL: ${BASE_URL:-https://api.openai.com}"
echo "Model: ${MODEL:-gpt-4o}"

# Test the actual API call
echo ""
echo "Testing API connection..."
echo "This will make a real API call to OpenAI"

# Run with a simple test
go run main.go --idea "Test" --location "Hefei" 2>&1 | head -20

echo ""
echo "Diagnosis complete!" 