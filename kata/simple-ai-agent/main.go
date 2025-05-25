package main

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/joho/godotenv"
	"context"
	"go.uber.org/zap"
	"encoding/json"
	"fmt"
	"os"
)

// Tool definitions for OpenAI API
var searchFileTool = openai.ChatCompletionToolParam{
	Function: openai.FunctionDefinitionParam{
		Name:        "SearchFileAndRead",
		Description: openai.String("Searches for a file by name in a directory and returns its content."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"fileName": map[string]interface{}{
					"type":        "string",
					"description": "The name of the file to search for (e.g., input.txt).",
				},
				"directory": map[string]interface{}{
					"type":        "string",
					"description": "The directory to search in (e.g., .).",
				},
			},
			"required": []string{"fileName", "directory"},
		},
	},
}

var drawImageTool = openai.ChatCompletionToolParam{
	Function: openai.FunctionDefinitionParam{
		Name:        "DrawImage",
		Description: openai.String("Generates a mindmap in DOT format from a list of topics."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{

				"script": map[string]interface{}{
					"type":        "string",
					"description": "The image script as plantuml script format.",
				},
				"imageType": map[string]interface{}{
					"type":        "string",
					"description": "The image type to generate (e.g., uml or mindmap).",
				},
				"outputPath": map[string]interface{}{
					"type":        "string",
					"description": "The image file save path (e.g., /tmp/image.png).",
				},
			},
			"required": []string{"imageType", "script"},
		},
	},
}

var getWeatherTool = openai.ChatCompletionToolParam{

	Function: openai.FunctionDefinitionParam{
		Name:        "get_weather",
		Description: openai.String("Get weather at the given location"),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]string{
					"type": "string",
				},
			},
			"required": []string{"location"},
		},
	},

}

func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Flushes buffer, if any


	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using environment variables")
	}
	// defaults to os.LookupEnv("OPENAI_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)

	ctx := context.Background()

	question := "What is the weather in New York City?"

	print("> ")
	println(question)

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolParam{
			getWeatherTool,
		},
		Seed:  openai.Int(0),
		Model: os.Getenv("OPENAI_MODEL"),
	}

	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	toolCalls := completion.Choices[0].Message.ToolCalls

	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}

	// If there is a was a function call, continue the conversation
	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			// Extract the location from the function call arguments
			var args map[string]interface{}
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				panic(err)
			}
			location := args["location"].(string)

			// Simulate getting weather data
			weatherData := getWeather(location)

			// Print the weather data
			fmt.Printf("Weather in %s: %s\n", location, weatherData)

			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
	}

	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}

	println(completion.Choices[0].Message.Content)
}
