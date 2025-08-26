package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var toolServerStarted bool

// Weather API response structure
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

func StartToolServer() {
	if toolServerStarted {
		return
	}

	r := gin.Default()
	r.POST("/tool/get_weather", func(c *gin.Context) {
		var req struct {
			Location string `json:"location"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		weatherInfo, err := getWeatherInfo(req.Location)
		if err != nil {
			logrus.Errorf("Failed to get weather for %s: %v", req.Location, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch weather"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"location": req.Location,
			"weather":  weatherInfo,
		})
	})

	go func() {
		if err := r.Run(":8080"); err != nil {
			logrus.Errorf("Failed to start tool server: %v", err)
		}
	}()

	toolServerStarted = true
	logrus.Info("Tool server started on :8080")
}

func getWeatherInfo(location string) (string, error) {
	// Try to get weather from environment variable first
	weatherAPIURL := os.Getenv("WEATHER_API_URL")
	if weatherAPIURL == "" {
		// Fallback to a mock weather service
		return getMockWeatherInfo(location), nil
	}

	// Build the API URL with location parameter
	baseURL := weatherAPIURL
	if baseURL[len(baseURL)-1] != '=' {
		baseURL += "&city="
	} else {
		baseURL += "city="
	}

	// Encode the location for URL
	encodedLocation := url.QueryEscape(location)
	fullURL := baseURL + encodedLocation

	// Make request to external weather API
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch weather from API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read weather response: %w", err)
	}

	// Parse the weather response
	var weatherResp WeatherResponse
	if err := json.Unmarshal(body, &weatherResp); err != nil {
		return "", fmt.Errorf("failed to parse weather response: %w", err)
	}

	// Check if the API call was successful
	if weatherResp.Status != "1" {
		return "", fmt.Errorf("weather API returned error status: %s, info: %s", weatherResp.Status, weatherResp.Info)
	}

	// Format the weather information
	if len(weatherResp.Lives) > 0 {
		live := weatherResp.Lives[0]
		weatherInfo := fmt.Sprintf("Location: %s, %s\nWeather: %s\nTemperature: %s°C\nHumidity: %s%%\nWind: %s %s\nReport Time: %s",
			live.Province, live.City, live.Weather, live.Temperature, live.Humidity, live.WindDirection, live.WindPower, live.ReportTime)
		return weatherInfo, nil
	}

	return "Weather information not available", nil
}

func getMockWeatherInfo(location string) string {
	now := time.Now()

	// Generate mock weather data in the same format as the real API
	mockResponse := WeatherResponse{
		Status:   "1",
		Count:    "1",
		Info:     "OK",
		Infocode: "10000",
		Lives: []struct {
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
		}{
			{
				Province:         "安徽",
				City:             location,
				Adcode:           "340100",
				Weather:          "晴",
				Temperature:      "25",
				WindDirection:    "北",
				WindPower:        "≤3",
				Humidity:         "65",
				ReportTime:       now.Format("2006-01-02 15:04:05"),
				TemperatureFloat: "25.0",
				HumidityFloat:    "65.0",
			},
		},
	}

	jsonData, _ := json.Marshal(mockResponse)
	return string(jsonData)
}
