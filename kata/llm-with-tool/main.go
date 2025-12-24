package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// City code mapping for Amap Weather API
var cityCodes = map[string]string{
	"合肥市":   "340100",
	"合肥":    "340100",
	"HEFEI": "340100",
	"长沙市":   "430100",
	"长沙":    "430100",
	"CHANGSHA":  "430100",
}

// WeatherResponse represents the Amap weather API response
type WeatherResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Lives    []struct {
		Province         string `json:"province"`
		City             string `json:"city"`
		Adcode           string `json:"adcode"`
		Weather          string `json:"weather"`
		Temperature      string `json:"temperature"`
		WindDirection    string `json:"winddirection"`
		WindPower        string `json:"windpower"`
		Humidity         string `json:"humidity"`
		ReportTime       string `json:"reporttime"`
		TemperatureFloat string `json:"temperature_float"`
		HumidityFloat    string `json:"humidity_float"`
	} `json:"lives"`
}

// getCityCode returns the city code for the given city name
func getCityCode(cityName string) string {
	cityName = strings.TrimSpace(strings.ToUpper(cityName))
	if code, ok := cityCodes[cityName]; ok {
		return code
	}
	// Default to Hefei
	return "340100"
}

// getWeatherTool defines the weather function tool for OpenAI
var getWeatherTool = openai.ChatCompletionToolParam{
	Function: openai.FunctionDefinitionParam{
		Name:        "get_weather",
		Description: openai.String("Get weather of a location. User must supply a location first."),
		Parameters: openai.FunctionParameters{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g., 合肥市",
				},
			},
			"required": []string{"location"},
		},
	},
}

// fetchWeather fetches weather data from Amap API
func fetchWeather(cityCode string) (*WeatherResponse, error) {
	apiKey := os.Getenv("LBS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LBS_API_KEY environment variable is not set")
	}

	baseURL := os.Getenv("LBS_BASE_URL") + "/v3/weather/weatherInfo"
	params := url.Values{}
	params.Add("city", cityCode)
	params.Add("key", apiKey)

	fullURL := baseURL + "?" + params.Encode()

	// Create HTTP client with TLS config (skip verification like the Python version)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &weatherResp, nil
}

// createOpenAIClient creates an OpenAI client with custom configuration
func createOpenAIClient() openai.Client {
	apiKey := os.Getenv("LLM_API_KEY")
	baseURL := os.Getenv("LLM_BASE_URL")

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}

	return openai.NewClient(opts...)
}

func main() {
	// Parse command line arguments
	city := flag.String("city", "Hefei", "City name to get weather for")
	flag.Parse()

	// Check if there's a positional argument
	if flag.NArg() > 0 {
		*city = flag.Arg(0)
	}

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	ctx := context.Background()
	client := createOpenAIClient()

	// Initial user message
	userMessage := fmt.Sprintf("please recommend today's dressing according today's weather in %s in Chinese?", *city)
	fmt.Printf("User> %s\n", userMessage)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(userMessage),
	}

	// First call: Model decides whether to use a tool
	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    os.Getenv("LLM_MODEL"),
		Messages: messages,
		Tools: []openai.ChatCompletionToolParam{
			getWeatherTool,
		},
	})
	if err != nil {
		fmt.Printf("Error calling OpenAI: %v\n", err)
		return
	}

	message := completion.Choices[0].Message
	toolCalls := message.ToolCalls

	// Check if a tool was called
	if len(toolCalls) == 0 {
		// No tool called; model provided a direct answer
		fmt.Printf("Model> %s\n", message.Content)
		return
	}

	// Add assistant message with tool calls to the conversation
	messages = append(messages, message.ToParam())

	// Process tool calls
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			fmt.Printf("Model> %s(%s)\n", toolCall.Function.Name, toolCall.Function.Arguments)

			// Parse the function arguments
			var args struct {
				Location string `json:"location"`
			}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				fmt.Printf("Error parsing arguments: %v\n", err)
				return
			}

			// Get the city code and fetch weather
			cityCode := getCityCode(args.Location)
			weatherResp, err := fetchWeather(cityCode)
			if err != nil {
				fmt.Printf("Error fetching weather: %v\n", err)
				return
			}

			// Format the weather response
			if len(weatherResp.Lives) == 0 {
				fmt.Println("No weather data available")
				return
			}

			todayWeather := weatherResp.Lives[0]
			toolResponse := fmt.Sprintf("%s is %s, %s°C, at %s",
				todayWeather.City,
				todayWeather.Weather,
				todayWeather.Temperature,
				todayWeather.ReportTime,
			)
			fmt.Printf("Tool> %s\n", toolResponse)
			// Add tool response to messages
			messages = append(messages, openai.ToolMessage(toolResponse, toolCall.ID))
		}
	}

	// Final call: Model generates natural language answer
	finalCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    os.Getenv("LLM_MODEL"),
		Messages: messages,
		Tools: []openai.ChatCompletionToolParam{
			getWeatherTool,
		},
	})
	if err != nil {
		fmt.Printf("Error calling OpenAI: %v\n", err)
		return
	}

	fmt.Printf("Model> %s\n", finalCompletion.Choices[0].Message.Content)
}
