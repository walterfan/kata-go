# LLM Function Call Example

This is a Go implementation of an LLM function calling example that queries weather information using the Amap (高德) Weather API.

## Features

- Uses OpenAI's function calling capability to determine when to fetch weather data
- Integrates with Amap Weather API for real weather information in Chinese cities
- Demonstrates the tool/function calling flow in LLM applications

## Prerequisites

1. An OpenAI API key (or compatible LLM API)
2. An Amap (高德) LBS API key from [https://lbs.amap.com/](https://lbs.amap.com/)

## Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` and add your API keys:
   ```bash
   LLM_API_KEY=your_openai_api_key
   LLM_BASE_URL=https://api.openai.com/v1
   LLM_MODEL=gpt-3.5-turbo
   LBS_API_KEY=your_amap_api_key
   ```

3. Install dependencies:
   ```bash
   make tidy
   ```

## Usage

### Build and Run

```bash
# Build the binary
make build

# Run with default city (Hefei)
make run

# Run with a specific city
make run CITY=Shanghai
./llm-function-call Shanghai

# Run with Chinese city names
make run-hefei
make run-wuhu
```

### Command Line

```bash
# Using go run
go run . Hefei
go run . 合肥

# Using the built binary
./llm-function-call 芜湖
./llm-function-call -city="New York"
```

## How It Works

1. The user asks about the weather in a specific city
2. The LLM analyzes the request and decides to call the `get_weather` function
3. The program fetches real weather data from Amap API using the city code
4. The weather data is sent back to the LLM
5. The LLM generates a natural language response

## Supported Cities

The following cities have predefined city codes for the Amap API:

| City Name | City Code |
|-----------|-----------|
| 合肥市/合肥/HEFEI | 340100 |
| 芜湖市/芜湖/WUHU | 340200 |

For other cities, it defaults to Hefei (340100). You can extend the `cityCodes` map in `main.go` to support more cities.

## API Reference

- [Amap Weather API Documentation](https://lbs.amap.com/api/webservice/guide/api-advanced/weatherinfo)
- [OpenAI Function Calling](https://platform.openai.com/docs/guides/function-calling)

## Python Version

The original Python version is available in `llm_function_call.py` in this directory.

