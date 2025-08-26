package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (app *App) setupUI() {
	// Create main sections
	leftPanel := app.createLeftPanel()
	rightPanel := app.createRightPanel()

	// Create horizontal split
	mainContainer := container.NewHSplit(leftPanel, rightPanel)
	mainContainer.SetOffset(0.3) // 30% for left panel, 70% for right panel

	// Set main content
	app.window.SetContent(mainContainer)
}

func (app *App) createLeftPanel() *fyne.Container {
	// Title
	title := widget.NewLabelWithStyle("📚 Collections", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Collection selector
	collectionNames := []string{}
	for _, collection := range app.config.Collections {
		collectionNames = append(collectionNames, collection.Name)
	}

	app.collectionSelect = widget.NewSelect(collectionNames, app.onCollectionSelected)
	app.collectionSelect.PlaceHolder = "Select a collection..."

	// Request selector
	app.requestSelect = widget.NewSelect([]string{}, app.onRequestSelected)
	app.requestSelect.PlaceHolder = "Select a request..."
	app.requestSelect.Disable()

	// Buttons
	refreshBtn := widget.NewButton("🔄 Refresh Config", func() {
		app.loadConfig()
		app.updateCollectionSelect()
	})
	refreshBtn.Importance = widget.MediumImportance

	newRequestBtn := widget.NewButton("➕ New Request", func() {
		app.clearRequestForm()
	})
	newRequestBtn.Importance = widget.LowImportance

	// Layout
	return container.NewVBox(
		title,
		widget.NewSeparator(),
		widget.NewLabel("Collection:"),
		app.collectionSelect,
		widget.NewLabel("Request:"),
		app.requestSelect,
		widget.NewSeparator(),
		refreshBtn,
		newRequestBtn,
		layout.NewSpacer(),
	)
}

func (app *App) createRightPanel() *container.Split {
	// Request section
	requestSection := app.createRequestSection()

	// Response section
	responseSection := app.createResponseSection()

	// Create vertical split for request/response
	return container.NewVSplit(requestSection, responseSection)
}

func (app *App) createRequestSection() *fyne.Container {
	// Title
	title := widget.NewLabelWithStyle("🌐 Request", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Method and URL row
	app.methodSelect = widget.NewSelect([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}, nil)
	app.methodSelect.SetSelected("GET")
	app.methodSelect.Resize(fyne.NewSize(100, 0))

	app.urlEntry = widget.NewEntry()
	app.urlEntry.SetPlaceHolder("Enter request URL...")

	app.sendButton = widget.NewButton("🚀 Send", app.sendRequest)
	app.sendButton.Importance = widget.HighImportance

	urlContainer := container.NewBorder(nil, nil, app.methodSelect, app.sendButton, app.urlEntry)

	// Parameters section
	paramLabel := widget.NewLabelWithStyle("📋 Parameters", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	app.parametersEntry = widget.NewEntry()
	app.parametersEntry.SetPlaceHolder("key1=value1&key2=value2")
	app.parametersEntry.MultiLine = true

	// Headers section
	headerLabel := widget.NewLabelWithStyle("🔖 Headers", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	app.headersEntry = widget.NewEntry()
	app.headersEntry.SetPlaceHolder("Content-Type: application/json\nAuthorization: Bearer token")
	app.headersEntry.MultiLine = true

	// Body section
	bodyLabel := widget.NewLabelWithStyle("📝 Body", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	app.bodyEntry = widget.NewEntry()
	app.bodyEntry.SetPlaceHolder("Request body (JSON, XML, etc.)")
	app.bodyEntry.MultiLine = true
	app.bodyEntry.Wrapping = fyne.TextWrapWord

	// Format body button
	formatBtn := widget.NewButton("🎨 Format JSON", app.formatRequestBody)
	formatBtn.Importance = widget.LowImportance

	bodyContainer := container.NewBorder(nil, formatBtn, nil, nil, app.bodyEntry)

	// Layout
	return container.NewVBox(
		title,
		widget.NewSeparator(),
		urlContainer,
		widget.NewSeparator(),
		paramLabel,
		container.NewWithoutLayout(app.parametersEntry),
		widget.NewSeparator(),
		headerLabel,
		container.NewWithoutLayout(app.headersEntry),
		widget.NewSeparator(),
		bodyLabel,
		bodyContainer,
	)
}

func (app *App) createResponseSection() *fyne.Container {
	// Title and status
	title := widget.NewLabelWithStyle("📨 Response", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	app.statusLabel = widget.NewLabel("Ready")
	app.statusLabel.TextStyle = fyne.TextStyle{Italic: true}

	titleContainer := container.NewBorder(nil, nil, title, app.statusLabel)

	// Response content
	app.responseEntry = widget.NewEntry()
	app.responseEntry.SetPlaceHolder("Response will appear here...")
	app.responseEntry.MultiLine = true
	app.responseEntry.Wrapping = fyne.TextWrapWord

	// Response buttons
	copyBtn := widget.NewButton("📋 Copy", func() {
		app.window.Clipboard().SetContent(app.responseEntry.Text)
		app.statusLabel.SetText("Response copied to clipboard")
	})

	clearBtn := widget.NewButton("🗑️ Clear", func() {
		app.responseEntry.SetText("")
		app.statusLabel.SetText("Response cleared")
	})

	buttonContainer := container.NewHBox(copyBtn, clearBtn)

	// Layout
	return container.NewVBox(
		titleContainer,
		widget.NewSeparator(),
		app.responseEntry,
		buttonContainer,
	)
}
