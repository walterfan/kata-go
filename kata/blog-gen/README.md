# Blog Generator

A Go program that generates technical blogs using OpenAI's API with predefined tools, including weather information integration.

## Features

- **OpenAI API Integration**: Uses OpenAI's API to generate blog content
- **Weather Tool**: Integrates weather information for today and tomorrow
- **Streaming Support**: Supports both streaming and non-streaming modes
- **Template-based**: Uses Jinja2-like templates for blog structure
- **Tool Server**: Local tool server for weather API calls

## Prerequisites

- Go 1.24.0 or later
- OpenAI API key

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Configuration

Create a `.env` file in the project root with the following variables:

```env
LLM_BASE_URL=https://api.openai.com
LLM_API_KEY=your_openai_api_key_here
LLM_MODEL=gpt-4o
LLM_STREAM=false
WEATHER_API_URL=optional_external_weather_api_url
```

### Weather API Format

The program supports weather APIs that return data in the following format (like Amap/高德地图 API):

```json
{
    "status": "1",
    "count": "1",
    "info": "OK",
    "infocode": "10000",
    "lives": [
        {
            "province": "安徽",
            "city": "合肥市",
            "adcode": "340100",
            "weather": "大雨",
            "temperature": "25",
            "winddirection": "北",
            "windpower": "≤3",
            "humidity": "97",
            "reporttime": "2025-07-31 22:01:22",
            "temperature_float": "25.0",
            "humidity_float": "97.0"
        }
    ]
}
```

## Usage

### Basic Usage

```bash
go run main.go
```

This will generate a blog with the default idea "用 WebRTC 和 Pion 打造一款网络录音机" and Beijing weather.

### Custom Options

```bash
# With custom idea
go run main.go --idea "Building a WebRTC-based Network Recorder"

# With custom location for weather
go run main.go --location "Shanghai"

# With specific OpenAI model
go run main.go --model "gpt-4o-mini"

# Combine options
go run main.go --idea "AI in Software Development" --location "New York" --model "gpt-4o"
```

### Command Line Options

- `-i, --idea`: Today's technical inspiration/title
- `-l, --location`: City for weather information (default: Beijing)
- `-m, --model`: OpenAI model name (e.g., gpt-4o)

## How It Works

1. **Template Rendering**: The program reads a Markdown template and renders it with basic data
2. **Tool Server**: Starts a local HTTP server on port 8080 to handle weather API calls
3. **OpenAI API Call**: Sends the rendered template to OpenAI with tool definitions
4. **Weather Integration**: When OpenAI requests weather information, it calls the local tool server
5. **Content Generation**: OpenAI generates the final blog content with weather information
6. **File Output**: Saves the generated blog to a file named `blog-YYYY-MM-DD.md`

## Blog Structure

The generated blog includes:

- Technical idea discussion (what, why, how, example, summary, reference)
- Daily GitHub hot project recommendation
- Daily software development best practice
- Daily LeetCode algorithm practice
- Daily design pattern practice
- Daily English quotes recitation
- Weather information for today and tomorrow

## Weather Tool

The program includes a built-in weather tool that:

- Provides mock weather data when no external API is configured
- Can integrate with external weather APIs via `WEATHER_API_URL` (supports Amap/高德地图 API format)
- Returns weather information including temperature, humidity, wind direction, and weather conditions
- Is accessible via HTTP POST to `http://localhost:8080/tool/get_weather`

## Streaming Mode

Enable streaming mode by setting `LLM_STREAM=true` in your `.env` file. This will:

- Show real-time content generation in the console
- Provide immediate feedback during API calls
- Display content as it's being generated

## Error Handling

The program includes comprehensive error handling for:

- Missing environment variables
- API request failures
- Template rendering errors
- File I/O operations
- Tool server startup issues

## Example Output

The program generates a blog file named `blog-2024-01-15.md` with content like:

```markdown
# my blog at 2024-01-15

## "用 WebRTC 和 Pion 打造一款网络录音机"

### what
WebRTC and Pion can be used to create a network recorder...

### why
Network recording is essential for debugging and monitoring...

### how
Using Pion's WebRTC implementation...

### example
Here's a simple example...

### summary
WebRTC and Pion provide powerful tools...

### reference
- Pion WebRTC documentation
- WebRTC specification

## Daily recommendation of 1 github hot project
...

## Daily recommendation of 1 best practice in software development and AI applications
...

## Daily practice of 1 leetcode algorithm question by one of the following languages: go/java/python/typescript/rust
...

## Daily practice of 1 classic design pattern by one of the following languages: go/java/python/typescript/rust
...

## Daily recitation of 10 English quotes
...
```

## Development

### Project Structure

```
blog-gen/
├── cmd/
│   └── root.go          # CLI command implementation
├── internal/
│   ├── llm_service.go   # OpenAI API integration
│   └── tool_weather.go  # Weather tool server
├── templates/
│   └── blog.md.teml     # Blog template
├── main.go              # Entry point
├── go.mod               # Dependencies
└── README.md           # This file
```

### Adding New Tools

To add new tools to the OpenAI API calls:

1. Define the tool in the `tools` array in `llm_service.go`
2. Implement the tool handler in a new file in `internal/`
3. Update the tool server to handle the new endpoint

## License

This project is open source and available under the MIT License. 