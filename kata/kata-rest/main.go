package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

// Config structures moved to config.go

// GUI Components moved to app.go

func main() {
	a := app.New()
	a.SetIcon(theme.ComputerIcon())

	myApp := &App{}
	myApp.window = a.NewWindow("🚀 GoPostman - REST API Client")
	myApp.window.SetIcon(theme.ComputerIcon())

	// Load configuration
	myApp.loadConfig()

	// Setup UI
	myApp.setupUI()

	// Set window properties
	myApp.window.Resize(fyne.NewSize(1000, 700))
	myApp.window.CenterOnScreen()
	myApp.window.ShowAndRun()
}

// UI functions moved to ui.go

// All functions moved to separate modules:
// - Config structures and loading → config.go
// - App struct and logic → app.go
// - UI creation and layout → ui.go
// - HTTP request handling → http.go
