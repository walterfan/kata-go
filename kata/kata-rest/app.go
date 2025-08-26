package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// GUI Components
type App struct {
	window           fyne.Window
	config           *Config
	collectionSelect *widget.Select
	requestSelect    *widget.Select
	methodSelect     *widget.Select
	urlEntry         *widget.Entry
	headersEntry     *widget.Entry
	parametersEntry  *widget.Entry
	bodyEntry        *widget.Entry
	responseEntry    *widget.Entry
	statusLabel      *widget.Label
	sendButton       *widget.Button
	currentRequest   *Request
}

func (app *App) updateCollectionSelect() {
	collectionNames := []string{}
	for _, collection := range app.config.Collections {
		collectionNames = append(collectionNames, collection.Name)
	}
	app.collectionSelect.Options = collectionNames
	app.collectionSelect.Refresh()
	app.requestSelect.Options = []string{}
	app.requestSelect.Refresh()
	app.requestSelect.Disable()
}

func (app *App) onCollectionSelected(collectionName string) {
	// Find the selected collection
	var selectedCollection *Collection
	for _, collection := range app.config.Collections {
		if collection.Name == collectionName {
			selectedCollection = &collection
			break
		}
	}

	if selectedCollection == nil {
		return
	}

	// Update request selector
	requestNames := []string{}
	for _, request := range selectedCollection.Requests {
		requestNames = append(requestNames, request.Name)
	}

	app.requestSelect.Options = requestNames
	app.requestSelect.Enable()
	app.requestSelect.ClearSelected()
	app.requestSelect.Refresh()
}

func (app *App) onRequestSelected(requestName string) {
	// Find the selected request
	var selectedRequest *Request
	for _, collection := range app.config.Collections {
		if collection.Name == app.collectionSelect.Selected {
			for _, request := range collection.Requests {
				if request.Name == requestName {
					selectedRequest = &request
					break
				}
			}
			break
		}
	}

	if selectedRequest == nil {
		return
	}

	app.currentRequest = selectedRequest
	app.loadRequestIntoForm(selectedRequest)
}

func (app *App) loadRequestIntoForm(request *Request) {
	app.methodSelect.SetSelected(request.Method)
	app.urlEntry.SetText(request.URL)

	// Load headers
	headerLines := []string{}
	for key, value := range request.Headers {
		headerLines = append(headerLines, fmt.Sprintf("%s: %s", key, value))
	}
	app.headersEntry.SetText(strings.Join(headerLines, "\n"))

	// Load parameters
	paramPairs := []string{}
	for key, value := range request.Parameters {
		paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", key, value))
	}
	app.parametersEntry.SetText(strings.Join(paramPairs, "&"))

	// Load body
	app.bodyEntry.SetText(request.Body)

	app.statusLabel.SetText(fmt.Sprintf("Loaded: %s", request.Name))
}

func (app *App) clearRequestForm() {
	app.methodSelect.SetSelected("GET")
	app.urlEntry.SetText("")
	app.headersEntry.SetText("")
	app.parametersEntry.SetText("")
	app.bodyEntry.SetText("")
	app.responseEntry.SetText("")
	app.statusLabel.SetText("New request")
	app.currentRequest = nil
}
