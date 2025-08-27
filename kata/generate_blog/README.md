# Blog Generator Example

This is an example Go project that demonstrates how to generate daily technical blog content using OpenAI's API.

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Create a `.env` file** in the project root with the following variables:
   ```env
   LLM_API_KEY=your_openai_api_key_here
   LLM_MODEL=gpt-4o-mini
   LLM_TEMPERATURE=0.7
   ```

3. **Run the blog generator:**
   ```bash
   go run generate_blog.go "Your daily idea here"
   ```

## Environment Variables

- `LLM_API_KEY`: Your OpenAI API key (required)
- `LLM_MODEL`: The model to use (defaults to "gpt-4o-mini")
- `LLM_TEMPERATURE`: Controls randomness in responses (defaults to 0.7)

## Features

- Generates structured blog content with sections for what, why, how, examples, and summaries
- Fetches trending GitHub projects
- Includes daily best practices and Go library examples
- Provides English sentences with Chinese explanations for daily work
- Uses streaming responses for better user experience

## Output

The program generates a markdown file named `blog-YYYY-MM-DD.md` with the complete blog content.
