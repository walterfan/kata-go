package main

import (
	"encoding/json"
	"fmt"
	"log"
)

// Weather API response structure (same as in tool_weather.go)
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

func main() {
	// Sample weather API response
	sampleResponse := `{
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
	}`

	var weatherResp WeatherResponse
	if err := json.Unmarshal([]byte(sampleResponse), &weatherResp); err != nil {
		log.Fatalf("Failed to parse weather response: %v", err)
	}

	fmt.Println("Weather API Response Parsing Test")
	fmt.Println("==================================")
	fmt.Printf("Status: %s\n", weatherResp.Status)
	fmt.Printf("Info: %s\n", weatherResp.Info)

	if len(weatherResp.Lives) > 0 {
		live := weatherResp.Lives[0]
		fmt.Printf("Location: %s, %s\n", live.Province, live.City)
		fmt.Printf("Weather: %s\n", live.Weather)
		fmt.Printf("Temperature: %s°C\n", live.Temperature)
		fmt.Printf("Humidity: %s%%\n", live.Humidity)
		fmt.Printf("Wind: %s %s\n", live.WindDirection, live.WindPower)
		fmt.Printf("Report Time: %s\n", live.ReportTime)
	}

	fmt.Println("\nTest completed successfully!")
}
