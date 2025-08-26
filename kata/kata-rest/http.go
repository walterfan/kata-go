package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/go-resty/resty/v2"
)

func (app *App) formatRequestBody() {
	body := app.bodyEntry.Text
	if body == "" {
		return
	}

	var prettyJSON bytes.Buffer
	if json.Valid([]byte(body)) {
		err := json.Indent(&prettyJSON, []byte(body), "", "  ")
		if err == nil {
			app.bodyEntry.SetText(prettyJSON.String())
			app.statusLabel.SetText("JSON formatted")
		}
	} else {
		app.statusLabel.SetText("Invalid JSON format")
	}
}

func (app *App) sendRequest() {
	method := app.methodSelect.Selected
	requestURL := app.urlEntry.Text
	body := app.bodyEntry.Text

	if requestURL == "" {
		app.statusLabel.SetText("❌ Please enter a URL")
		return
	}

	app.statusLabel.SetText("🔄 Sending request...")
	app.sendButton.Disable()

	go func() {
		client := resty.New()

		// Parse headers
		headers := app.parseHeaders()

		// Parse parameters and add to URL
		finalURL := app.buildURLWithParameters(requestURL)

		// Create request
		req := client.R()
		for key, value := range headers {
			req.SetHeader(key, value)
		}

		if body != "" {
			req.SetBody(body)
		}

		var resp *resty.Response
		var err error

		// Send request based on method
		switch method {
		case "GET":
			resp, err = req.Get(finalURL)
		case "POST":
			resp, err = req.Post(finalURL)
		case "PUT":
			resp, err = req.Put(finalURL)
		case "DELETE":
			resp, err = req.Delete(finalURL)
		case "PATCH":
			resp, err = req.Patch(finalURL)
		case "HEAD":
			resp, err = req.Head(finalURL)
		case "OPTIONS":
			resp, err = req.Options(finalURL)
		default:
			err = fmt.Errorf("unsupported method: %s", method)
		}

		// Update UI in main thread
		app.window.Canvas().SetOnTypedKey(func(ke *fyne.KeyEvent) {
			app.handleResponse(resp, err)
		})

		// Trigger UI update
		app.handleResponse(resp, err)
		app.sendButton.Enable()
	}()
}

func (app *App) parseHeaders() map[string]string {
	headers := make(map[string]string)
	headerText := app.headersEntry.Text

	if headerText == "" {
		return headers
	}

	lines := strings.Split(headerText, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	return headers
}

func (app *App) buildURLWithParameters(baseURL string) string {
	paramText := app.parametersEntry.Text
	if paramText == "" {
		return baseURL
	}

	// Parse existing URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	// Parse parameters
	values := parsedURL.Query()

	// Add new parameters
	pairs := strings.Split(paramText, "&")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			values.Add(key, value)
		}
	}

	parsedURL.RawQuery = values.Encode()
	return parsedURL.String()
}

func (app *App) handleResponse(resp *resty.Response, err error) {
	if err != nil {
		app.responseEntry.SetText(fmt.Sprintf("❌ Error: %v", err))
		app.statusLabel.SetText("❌ Request failed")
		return
	}

	// Format status
	statusText := fmt.Sprintf("✅ %d %s | %v", resp.StatusCode(), resp.Status(), resp.Time())
	app.statusLabel.SetText(statusText)

	// Format response body
	responseBody := resp.Body()
	if len(responseBody) == 0 {
		app.responseEntry.SetText("(Empty response)")
		return
	}

	// Try to pretty print JSON
	var prettyJSON bytes.Buffer
	if json.Valid(responseBody) {
		err := json.Indent(&prettyJSON, responseBody, "", "  ")
		if err == nil {
			app.responseEntry.SetText(prettyJSON.String())
			return
		}
	}

	// Fallback to plain text
	app.responseEntry.SetText(string(responseBody))
}
