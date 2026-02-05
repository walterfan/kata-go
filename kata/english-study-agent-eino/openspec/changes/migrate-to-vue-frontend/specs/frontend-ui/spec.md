# Frontend UI Capability

## ADDED Requirements

### Requirement: Modern SPA Architecture
The system SHALL provide a modern single-page application (SPA) frontend built with Vue.js 3 and TypeScript for improved performance and user experience.

#### Scenario: Fast initial load
- **WHEN** a user visits the application
- **THEN** the frontend SHALL load within 2 seconds on a standard broadband connection
- **AND** display a loading indicator during initialization

#### Scenario: Smooth interactions
- **WHEN** a user performs any action (click, input, scroll)
- **THEN** the UI SHALL respond within 100ms (perceived performance)
- **AND** use transitions and animations for state changes

#### Scenario: Client-side routing
- **WHEN** a user navigates between pages (Home, Settings)
- **THEN** the browser SHALL not reload the entire page
- **AND** the URL SHALL update to reflect the current route
- **AND** the browser back/forward buttons SHALL work correctly

### Requirement: Responsive Layout
The frontend SHALL provide a responsive layout that adapts to different screen sizes and devices.

#### Scenario: Desktop layout
- **WHEN** viewed on a desktop screen (>1024px width)
- **THEN** the UI SHALL display a two-column layout with collapsible sidebar and main content pane
- **AND** the sidebar SHALL default to expanded

#### Scenario: Tablet layout
- **WHEN** viewed on a tablet screen (768-1024px width)
- **THEN** the sidebar SHALL default to collapsed
- **AND** all features SHALL remain accessible

#### Scenario: Mobile layout
- **WHEN** viewed on a mobile screen (<768px width)
- **THEN** the UI SHALL adapt to single-column layout
- **AND** navigation SHALL use a slide-out drawer or bottom tabs

#### Scenario: Touch support
- **WHEN** used on touch devices
- **THEN** all interactive elements SHALL have adequate touch targets (min 44x44px)
- **AND** support touch gestures where appropriate (swipe, pinch-to-zoom)

### Requirement: Internationalization (i18n)
The frontend SHALL support multiple languages with seamless switching and persistence.

#### Scenario: Language selection
- **WHEN** a user selects a language from the language toggle
- **THEN** all UI text SHALL update immediately without page reload
- **AND** the selected language SHALL be saved to browser localStorage
- **AND** the language preference SHALL persist across sessions

#### Scenario: Supported languages
- **WHEN** the application initializes
- **THEN** it SHALL support at least English (en) and Chinese (zh) locales
- **AND** allow adding new locales through translation files

#### Scenario: Fallback translations
- **WHEN** a translation key is missing in the selected language
- **THEN** the UI SHALL display the English (fallback) translation
- **AND** log a warning in the console for debugging

### Requirement: Component-Based Architecture
The frontend SHALL use a modular component-based architecture for maintainability and reusability.

#### Scenario: Reusable components
- **WHEN** developing UI features
- **THEN** common UI elements SHALL be extracted into reusable components
- **AND** components SHALL accept props for configuration
- **AND** components SHALL emit events for parent communication

#### Scenario: Component isolation
- **WHEN** a component is modified
- **THEN** the change SHALL not affect other unrelated components
- **AND** styles SHALL be scoped to the component (no global CSS pollution)

#### Scenario: Type safety
- **WHEN** developing components
- **THEN** all props, events, and internal state SHALL have TypeScript types
- **AND** the TypeScript compiler SHALL catch type errors before runtime

### Requirement: State Management
The frontend SHALL use centralized state management for predictable data flow and debugging.

#### Scenario: Global state store
- **WHEN** managing application-wide state (language, user preferences, current content)
- **THEN** the UI SHALL use Pinia stores
- **AND** all state mutations SHALL go through store actions (not direct mutation)

#### Scenario: Reactive updates
- **WHEN** state changes in a store
- **THEN** all components depending on that state SHALL automatically re-render
- **AND** only affected components SHALL update (no full page re-render)

#### Scenario: DevTools integration
- **WHEN** developing or debugging
- **THEN** the Vue DevTools browser extension SHALL display all store state
- **AND** allow time-travel debugging (undo/redo state changes)

### Requirement: HTTP API Integration
The frontend SHALL communicate with the Go backend via RESTful HTTP APIs and handle errors gracefully.

#### Scenario: API base URL configuration
- **WHEN** the application initializes
- **THEN** the API base URL SHALL be configurable via environment variables
- **AND** default to `http://localhost:8080/api` for local development

#### Scenario: Request/response handling
- **WHEN** making API requests
- **THEN** the frontend SHALL use Axios with a configured instance
- **AND** include appropriate headers (Content-Type, Accept)
- **AND** handle timeouts (default 30 seconds)

#### Scenario: Global error handling
- **WHEN** an API request fails
- **THEN** the frontend SHALL display a user-friendly error message (using Element Plus `ElMessage`)
- **AND** log the error details to the browser console
- **AND** allow retrying the request if applicable

#### Scenario: Loading indicators
- **WHEN** an API request is in progress
- **THEN** the UI SHALL display a loading indicator (spinner, skeleton, progress bar)
- **AND** disable interactive elements to prevent duplicate requests

### Requirement: Server-Sent Events (SSE) Streaming
The frontend SHALL support real-time streaming of AI responses via SSE for improved user experience.

#### Scenario: Native SSE handling
- **WHEN** streaming mode is enabled and an AI action is triggered
- **THEN** the frontend SHALL use the Fetch API with `ReadableStream` to consume SSE
- **AND** parse SSE format: `event: {type}\ndata: {content}\n\n`
- **AND** handle events: `message` (content chunk), `done` (stream complete), `error` (stream failure)

#### Scenario: Incremental content display
- **WHEN** receiving SSE chunks
- **THEN** the UI SHALL append each chunk to the result display in real-time
- **AND** auto-scroll to the bottom to show new content
- **AND** preserve formatting (newlines, spaces)

#### Scenario: Stream cancellation
- **WHEN** a user navigates away or clicks "Stop"
- **THEN** the frontend SHALL abort the ongoing stream
- **AND** clean up resources (close reader, reset state)

#### Scenario: Streaming toggle
- **WHEN** a user toggles streaming mode on/off
- **THEN** the setting SHALL be saved to localStorage
- **AND** all subsequent AI actions SHALL use the selected mode (streaming or non-streaming)

### Requirement: Text-to-Speech (TTS) Integration
The frontend SHALL support browser-native text-to-speech for reading English content aloud.

#### Scenario: Native Web Speech API
- **WHEN** a user clicks the "Read" button
- **THEN** the frontend SHALL use the browser's `SpeechSynthesis` API
- **AND** select an English voice (preferring `en-US` locale)
- **AND** set speech rate to 0.9x (slightly slower for learners)

#### Scenario: TTS controls
- **WHEN** TTS is playing
- **THEN** the "Read" button SHALL change to "Stop"
- **AND** clicking "Stop" SHALL cancel the speech immediately
- **AND** the button state SHALL sync with `isSpeaking` reactive state

#### Scenario: Sentence-level TTS
- **WHEN** displaying sentence extraction results
- **THEN** each sentence SHALL have its own "Read" button
- **AND** clicking a sentence button SHALL read only that sentence

#### Scenario: Browser compatibility
- **WHEN** TTS is used on different browsers
- **THEN** it SHALL work on Chrome, Edge, Safari, and Firefox
- **AND** display a warning if `SpeechSynthesis` is not supported

### Requirement: Link Extraction and Navigation
The frontend SHALL automatically extract links from article content and provide quick navigation to linked articles.

#### Scenario: Link detection
- **WHEN** article content is loaded (from RSS or URL)
- **THEN** the frontend SHALL extract all links from HTML `<a>` tags and Markdown `[text](url)` format
- **AND** deduplicate links (same URL appears only once)
- **AND** resolve relative URLs to absolute URLs using the article's base URL

#### Scenario: Link display
- **WHEN** links are extracted
- **THEN** the UI SHALL display a "🔗 Links Found" section below the article content
- **AND** show each link with its text and a "📥 Fetch" button

#### Scenario: Link navigation
- **WHEN** a user clicks a "Fetch" button
- **THEN** the frontend SHALL call the URL fetch API with the link's URL
- **AND** replace the current content with the fetched article
- **AND** update the extracted links list with links from the new article

### Requirement: Custom RSS Feed Management
The frontend SHALL provide a user interface for managing custom RSS feeds with CRUD operations.

#### Scenario: View custom feeds
- **WHEN** a user navigates to the Settings page
- **THEN** the UI SHALL display a table of custom RSS feeds
- **AND** show columns: title, URL, category, actions

#### Scenario: Add custom feed
- **WHEN** a user clicks "Add Feed"
- **THEN** a dialog SHALL open with form fields: title (required), URL (required, validated), category (optional)
- **AND** clicking "Save" SHALL call `POST /api/custom-feeds`
- **AND** the new feed SHALL appear in the list and RSS source dropdown

#### Scenario: Edit custom feed
- **WHEN** a user clicks "Edit" on a feed
- **THEN** a dialog SHALL open pre-filled with the feed's current data
- **AND** clicking "Save" SHALL call `PUT /api/custom-feeds/:id`
- **AND** the feed SHALL update in the list

#### Scenario: Delete custom feed
- **WHEN** a user clicks "Delete" on a feed
- **THEN** a confirmation dialog SHALL appear
- **AND** confirming SHALL call `DELETE /api/custom-feeds/:id`
- **AND** the feed SHALL be removed from the list and RSS source dropdown

#### Scenario: Validation
- **WHEN** adding or editing a feed
- **THEN** the UI SHALL validate:
  - Title is not empty
  - URL is a valid HTTP(S) URL
  - URL does not already exist in the list
- **AND** display validation errors next to the relevant field

### Requirement: Build and Deployment
The frontend SHALL be built as a static SPA and served from the Go backend for single-binary deployment.

#### Scenario: Production build
- **WHEN** building the frontend for production
- **THEN** Vite SHALL bundle all assets (JS, CSS, images) into an optimized `dist/` folder
- **AND** generate a manifest for cache busting
- **AND** minify and tree-shake code to reduce bundle size

#### Scenario: Integration with Go backend
- **WHEN** the Go server starts
- **THEN** it SHALL serve static files from `frontend/dist/assets`
- **AND** serve `index.html` for all non-API routes (Vue Router history mode)
- **AND** optionally embed the frontend assets into the Go binary using `//go:embed`

#### Scenario: Development mode
- **WHEN** developing locally
- **THEN** the Vite dev server SHALL run on port 5173 with hot module replacement (HMR)
- **AND** proxy `/api` requests to the Go backend on port 8080
- **AND** allow parallel development of frontend and backend

#### Scenario: Automated build
- **WHEN** running `make build-all`
- **THEN** the Makefile SHALL:
  - Install frontend dependencies (`pnpm install`)
  - Build the Vue app (`pnpm build`)
  - Build the Go binary with embedded frontend (optional)
- **AND** output a single executable in `bin/`

### Requirement: User Experience Enhancements
The frontend SHALL provide a polished and intuitive user experience with helpful guidance and feedback.

#### Scenario: Welcome guide
- **WHEN** a new user opens the application
- **THEN** the UI SHALL display a welcome message with quick start instructions
- **AND** provide example inputs or a guided tour (optional)

#### Scenario: Empty states
- **WHEN** no content is loaded
- **THEN** the UI SHALL display a friendly message (e.g., "Load an article to get started")
- **AND** suggest actions (e.g., "Click Refresh to load articles")

#### Scenario: Error states
- **WHEN** an error occurs (network failure, API error, empty content)
- **THEN** the UI SHALL display a clear error message with context
- **AND** suggest next steps (e.g., "Check your connection and try again")

#### Scenario: Success feedback
- **WHEN** an action succeeds (article loaded, feed added)
- **THEN** the UI SHALL display a success message (toast notification)
- **AND** automatically dismiss the message after 3 seconds

### Requirement: Accessibility
The frontend SHALL follow accessibility best practices to support users with disabilities.

#### Scenario: Keyboard navigation
- **WHEN** a user navigates using only the keyboard
- **THEN** all interactive elements SHALL be reachable via Tab key
- **AND** the current focus SHALL be visually indicated
- **AND** Enter or Space key SHALL activate buttons

#### Scenario: ARIA labels
- **WHEN** interactive elements are rendered
- **THEN** they SHALL have appropriate ARIA labels (aria-label, aria-describedby)
- **AND** screen readers SHALL announce the element's purpose

#### Scenario: Contrast and readability
- **WHEN** displaying text and UI elements
- **THEN** color contrast SHALL meet WCAG AA standards (minimum 4.5:1 for normal text)
- **AND** font sizes SHALL be at least 14px (or adjustable via browser zoom)

---

## Related Capabilities

- **ai-assistant**: AI actions (explain, summarize, translate, etc.) are triggered from the frontend and consume the backend API.
- **content-extractor**: Sentence and vocabulary extraction results are displayed in the frontend with markdown rendering.
- **rss-feed-reader**: RSS feeds and articles are fetched and displayed in the frontend's Article Loader component.

