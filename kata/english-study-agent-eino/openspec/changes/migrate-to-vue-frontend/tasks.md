# Tasks: Migrate Frontend from Streamlit to Vue.js 3 + TypeScript

**Change ID**: `migrate-to-vue-frontend`  
**Status**: Draft  
**Last Updated**: 2026-01-02

---

## Implementation Phases

### Phase 1: Project Scaffold & Setup (Foundation)

**Goal**: Set up Vue 3 + TypeScript project with all tooling and dependencies.

- [ ] **T1.1**: Scaffold Vue 3 project with Vite
  - Run `pnpm create vite@latest frontend -- --template vue-ts`
  - Verify dev server runs (`pnpm dev`)
- [ ] **T1.2**: Install core dependencies
  - `pnpm add vue-router@4 pinia vue-i18n@9 axios element-plus @element-plus/icons-vue`
  - `pnpm add -D sass @types/node tailwindcss postcss autoprefixer`
  - Run `npx tailwindcss init -p` to generate Tailwind config
- [ ] **T1.3**: Configure Vite
  - Set up path alias (`@` -> `src`)
  - Configure dev proxy (`/api` -> `http://localhost:8080`)
  - Set build output directory (`dist`)
  - Configure Tailwind CSS in PostCSS pipeline
- [ ] **T1.4**: Configure TypeScript
  - Update `tsconfig.json` with strict mode and path mappings
  - Add `tsconfig.node.json` for Vite config
- [ ] **T1.5**: Set up ESLint + Prettier
  - Install `eslint`, `@vue/eslint-config-typescript`, `prettier`
  - Create `.eslintrc.cjs` and `.prettierrc`
- [ ] **T1.6**: Create project structure
  - Create folders: `src/{components,composables,i18n,router,stores,types,utils,views,assets}`
  - Add placeholder `index.ts` files
  - Create `src/styles/main.css` with Tailwind imports:
    ```css
    @tailwind base;
    @tailwind components;
    @tailwind utilities;
    ```
- [ ] **T1.7**: Set up Pinia store
  - Create `src/stores/index.ts` with Pinia initialization
  - Create placeholder stores: `app.ts`, `content.ts`, `rss.ts`, `results.ts`
- [ ] **T1.8**: Set up Vue Router
  - Create `src/router/index.ts` with routes: `/` (HomeView), `/settings` (SettingsView)
  - Create placeholder view files
- [ ] **T1.9**: Set up Vue I18n
  - Create `src/i18n/{index.ts,en.ts,zh.ts}`
  - Copy translations from `web/app.py` (EN/CN dictionaries)
- [ ] **T1.10**: Configure Axios instance
  - Create `src/utils/api.ts` with base URL, timeout, error interceptor
- [ ] **T1.11**: Update Makefile
  - Add targets: `frontend-install`, `frontend-dev`, `frontend-build`, `dev-all`
- [ ] **T1.12**: Test dev setup
  - Run `make frontend-dev` and verify Vue app loads
  - Verify API proxy works (test with `/api/feeds` call)

**Validation**: Dev server runs, hot reload works, API proxy forwards requests to Go backend.

---

### Phase 2: Core Layout & i18n (UI Shell)

**Goal**: Implement main layout, navigation, and language switching.

- [ ] **T2.1**: Create `App.vue`
  - Set up root layout with Element Plus global styles
  - Include `<RouterView>` for page routing
- [ ] **T2.2**: Implement `LanguageToggle.vue` component
  - Dropdown or button to switch between EN/CN
  - Update `appStore.language` and `i18n.global.locale`
  - Persist language choice in `localStorage`
- [ ] **T2.3**: Create `HomeView.vue` layout
  - Use Element Plus `<el-container>` with sidebar and main pane
  - Add collapsible sidebar (toggle button)
  - Display title + proverb in header
- [ ] **T2.4**: Create `SettingsView.vue` placeholder
  - Add navigation link in sidebar
  - Display "Settings" title
- [ ] **T2.5**: Implement `appStore` (Pinia)
  - State: `language`, `inputMode`, `sidebarCollapsed`, `streamingEnabled`, `loading`
  - Actions: `setLanguage()`, `setInputMode()`, `toggleSidebar()`, `toggleStreaming()`
- [ ] **T2.6**: Test i18n switching
  - Verify all UI text updates when language toggles
  - Check `localStorage` persistence across page reloads

**Validation**: Layout renders correctly, sidebar toggles, language switching works, navigation between Home/Settings works.

---

### Phase 3: RSS Feed Loading (Article Mode)

**Goal**: Fetch and display RSS articles from backend.

- [ ] **T3.1**: Define TypeScript types
  - Create `src/types/rss.ts`: `RssSource`, `Article`, `CustomFeed`
  - Create `src/types/api.ts`: `ApiResponse`, `ApiError`
- [ ] **T3.2**: Implement `rssStore` (Pinia)
  - State: `sources[]`, `articles[]`, `customFeeds[]`, `selectedSource`
  - Actions: `fetchSources()`, `fetchArticles(source)`, `fetchCustomFeeds()`
- [ ] **T3.3**: Create `RssFeedList.vue` component
  - Display list of articles with title, source, published date
  - Click to select article
  - Show "No articles" message if empty
- [ ] **T3.4**: Create `ArticleLoader.vue` component
  - Dropdown to select RSS source
  - "Refresh Headlines" button
  - Display `<RssFeedList>` with loaded articles
  - Article preview pane
- [ ] **T3.5**: Integrate API calls in `rssStore`
  - `GET /api/rss-sources` -> update `sources`
  - `GET /api/feeds?source={source}` -> update `articles`
  - Handle errors (show toast notification via Element Plus `ElMessage`)
- [ ] **T3.6**: Wire up article selection
  - Update `contentStore.currentText` and `contentStore.selectedArticle` when user clicks article
  - Display article content in main pane
- [ ] **T3.7**: Test RSS loading
  - Verify articles load from backend
  - Check error handling for network failures
  - Verify article selection updates UI

**Validation**: RSS sources load, articles display, user can select an article and see content.

---

### Phase 4: URL Fetcher & Text Input (Other Input Modes)

**Goal**: Support URL fetching and manual text input.

- [ ] **T4.1**: Implement `contentStore` (Pinia)
  - State: `currentText`, `fetchedUrl`, `selectedArticle`, `extractedLinks`
  - Actions: `setCurrentText()`, `fetchFromUrl(url)`, `clearContent()`
- [ ] **T4.2**: Create `UrlFetcher.vue` component
  - URL input field (Element Plus `<el-input>`)
  - "Load Article" button
  - Loading spinner during fetch
  - Error message display
- [ ] **T4.3**: Integrate `POST /api/fetch-url` in `contentStore`
  - Call API with `{ url }`
  - Update `currentText`, `fetchedUrl`, and `extractedLinks`
  - Handle errors (empty content, network failure)
- [ ] **T4.4**: Create `TextInput.vue` component
  - Textarea for manual text input
  - Word count display
  - "Clear" button
- [ ] **T4.5**: Add input mode toggle in sidebar
  - Radio buttons or tabs: "Article", "URL", "Text Input"
  - Update `appStore.inputMode`
  - Show/hide components based on mode
- [ ] **T4.6**: Test URL fetching
  - Verify URLs load (HTML and Markdown)
  - Check error handling (404, empty content, timeouts)
- [ ] **T4.7**: Test text input
  - Verify manual text updates `contentStore.currentText`

**Validation**: User can fetch articles from URLs, enter text manually, and switch between input modes.

---

### Phase 5: AI Actions (Non-Streaming)

**Goal**: Implement all AI action buttons with non-streaming API.

- [ ] **T5.1**: Define action types
  - Create `src/types/agent.ts`: `AgentAction`, `AgentTask`
  - List all actions: `explain`, `summarize`, `translate`, `refine`, `extract_sentences`, `extract_vocabulary`
- [ ] **T5.2**: Create `ActionButtons.vue` component
  - Grid of buttons (Element Plus `<el-button>`)
  - Icons for each action (using `@element-plus/icons-vue`)
  - Emit action click events
- [ ] **T5.3**: Implement `resultsStore` (Pinia)
  - State: `currentResult`, `streamingTask`, `isStreaming`, `history[]`
  - Actions: `executeAction(text, task)`, `clearResults()`, `addToHistory()`
- [ ] **T5.4**: Create API service in `src/composables/useAgent.ts`
  - Function `callAgent(text, task)` -> `POST /api/chat`
  - Return `Promise<string>` with result
  - Handle errors
- [ ] **T5.5**: Wire up action buttons
  - On button click, call `resultsStore.executeAction()`
  - Show loading indicator (`appStore.loading = true`)
  - Display result in `ResultDisplay` component
- [ ] **T5.6**: Create `ResultDisplay.vue` component
  - Display result text with markdown rendering (use `markdown-it` or `marked`)
  - Show loading spinner during API call
  - "Clear Results" button
- [ ] **T5.7**: Test all actions
  - Verify each action calls correct API
  - Check markdown rendering (especially for `extract_sentences`, `extract_vocabulary`)
  - Test error handling (LLM timeout, invalid input)

**Validation**: All AI actions work correctly, results display, markdown renders properly.

---

### Phase 6: Streaming Support (SSE)

**Goal**: Enable streaming mode for real-time LLM responses.

- [ ] **T6.1**: Create `useStreaming` composable
  - Function `startStream(url, payload)` using `fetch` + `ReadableStream`
  - Parse SSE format: `event: message\ndata: {content}`
  - Handle `event: done` and `event: error`
  - Return reactive `content`, `isStreaming`, `error`
- [ ] **T6.2**: Add streaming toggle in UI
  - Checkbox or switch in sidebar: "⚡ Streaming Mode"
  - Update `appStore.streamingEnabled`
- [ ] **T6.3**: Update `resultsStore` to use streaming
  - If `streamingEnabled`, call `startStream('/api/chat/stream', { text, task })`
  - Update `currentResult` incrementally as chunks arrive
  - Show typing cursor animation during streaming
- [ ] **T6.4**: Update `ResultDisplay` to show streaming
  - Display `currentResult` with auto-scroll to bottom
  - Show typing cursor (e.g., blinking `|`)
  - Disable "Clear" button during streaming
- [ ] **T6.5**: Test streaming
  - Verify SSE chunks are parsed correctly (no missing spaces)
  - Check newline handling (`\n` -> `\n`)
  - Test cancellation (user navigates away or clicks "Stop")
- [ ] **T6.6**: Add "Stop Streaming" button
  - Button to cancel ongoing stream
  - Call `stopStream()` in `useStreaming` composable

**Validation**: Streaming mode works, no formatting issues, cancellation works, toggle between streaming/non-streaming.

---

### Phase 7: TTS & Link Extraction (Enhanced UX)

**Goal**: Add text-to-speech and link navigation features.

- [ ] **T7.1**: Create `useTTS` composable
  - Function `speak(text, rate)` using `SpeechSynthesisUtterance`
  - Function `stop()` using `window.speechSynthesis.cancel()`
  - Function `toggle(text, rate)` to start/stop
  - Reactive `isSpeaking` state
- [ ] **T7.2**: Add TTS controls to `ResultDisplay`
  - Toggle button: "🔊 Read" / "⏸️ Stop"
  - Button state updates with `isSpeaking`
  - For `extract_sentences`: individual "Read" buttons per sentence
- [ ] **T7.3**: Create `src/utils/linkExtractor.ts`
  - Function `extractLinks(html, baseUrl)` -> `ExtractedLink[]`
  - Handle HTML `<a>` tags and Markdown `[text](url)`
  - Deduplicate links, resolve relative URLs
- [ ] **T7.4**: Create `LinkList.vue` component
  - Display list of extracted links with text and URL
  - "📥 Fetch" button per link
  - On click, call `contentStore.fetchFromUrl(link.url)`
- [ ] **T7.5**: Integrate link extraction in `contentStore`
  - When `currentText` updates, call `extractLinks(currentText, fetchedUrl)`
  - Update `extractedLinks` state
- [ ] **T7.6**: Display `LinkList` in `HomeView`
  - Show below "Current Text" expander
  - Only visible if `extractedLinks.length > 0`
- [ ] **T7.7**: Test TTS
  - Verify speech starts/stops correctly
  - Check English voice selection
  - Test on different browsers (Chrome, Firefox, Safari)
- [ ] **T7.8**: Test link extraction
  - Verify links are extracted from RSS content and URL-fetched articles
  - Check relative URL resolution
  - Test "Fetch" button updates content

**Validation**: TTS works natively, link extraction and navigation work, no browser console errors.

---

### Phase 8: Custom RSS Feed Management (Settings)

**Goal**: Implement CRUD interface for custom RSS feeds.

- [ ] **T8.1**: Update `rssStore` with custom feed actions
  - `fetchCustomFeeds()` -> `GET /api/custom-feeds`
  - `addCustomFeed(feed)` -> `POST /api/custom-feeds`
  - `updateCustomFeed(id, feed)` -> `PUT /api/custom-feeds/:id`
  - `deleteCustomFeed(id)` -> `DELETE /api/custom-feeds/:id`
- [ ] **T8.2**: Create `CustomFeedList.vue` component
  - Table displaying custom feeds (title, URL, category)
  - "Edit" and "Delete" buttons per row
  - "Add Feed" button to open dialog
- [ ] **T8.3**: Create `CustomFeedDialog.vue` component
  - Form with fields: `title`, `url`, `category`
  - Validation: URL format, required fields
  - "Save" and "Cancel" buttons
- [ ] **T8.4**: Create `DefaultFeedList.vue` component
  - Read-only table of default feeds from `config.yaml`
  - Display title, URL, category
- [ ] **T8.5**: Integrate into `SettingsView`
  - Display `<DefaultFeedList>` and `<CustomFeedList>` in tabs or sections
  - Wire up CRUD actions to `rssStore`
- [ ] **T8.6**: Test CRUD operations
  - Add a custom feed -> verify it appears in list and RSS source dropdown
  - Edit a feed -> verify changes persist
  - Delete a feed -> verify it's removed from list
  - Test validation (invalid URL, empty title)

**Validation**: Users can manage custom RSS feeds, changes persist in SQLite database.

---

### Phase 9: Polish & Testing (Quality Assurance)

**Goal**: Fix bugs, improve styling, and test edge cases.

- [ ] **T9.1**: Manual testing checklist
  - Test all input modes (Article, URL, Text Input)
  - Test all AI actions (Explain, Summarize, Translate, Refine, Sentences, Vocabulary)
  - Test streaming toggle on/off
  - Test TTS for full article and individual sentences
  - Test link extraction and fetch
  - Test RSS feed management (CRUD)
  - Test language toggle (EN/CN)
  - Test sidebar collapse/expand
  - Test error handling (network failure, empty content, API errors)
- [ ] **T9.2**: Responsive layout fixes
  - Test on mobile/tablet screen sizes
  - Adjust sidebar behavior (auto-collapse on small screens)
  - Ensure buttons don't wrap, text is readable
- [ ] **T9.3**: Accessibility improvements (optional)
  - Add ARIA labels to buttons and inputs
  - Ensure keyboard navigation works (Tab, Enter)
  - Test with screen reader (VoiceOver, NVDA)
- [ ] **T9.4**: Unit tests for utilities (Vitest)
  - Test `linkExtractor.ts` with various HTML/Markdown inputs
  - Test API error handling in `api.ts`
  - Test i18n helper functions
- [ ] **T9.5**: Component tests (Vue Test Utils)
  - Test `LanguageToggle` switches locale
  - Test `ActionButtons` emits correct events
  - Test `ResultDisplay` renders markdown
  - Test `RssFeedList` displays articles correctly
- [ ] **T9.6**: Styling polish
  - Ensure consistent spacing, colors, typography
  - Match Element Plus theme or customize as needed
  - Add loading skeletons for better perceived performance
- [ ] **T9.7**: Fix any remaining bugs from manual testing
- [ ] **T9.8**: Performance optimization
  - Check bundle size (`pnpm build`, analyze with `vite-plugin-visualizer`)
  - Lazy-load heavy components (e.g., Settings view)
  - Optimize images/assets

**Validation**: All features work, no critical bugs, UI is polished and responsive.

---

### Phase 10: Backend Integration & Deployment (Cutover)

**Goal**: Integrate Vue app with Go backend and retire Streamlit.

- [ ] **T10.1**: Build Vue app for production
  - Run `make frontend-build`
  - Verify `frontend/dist/` contains static assets
- [ ] **T10.2**: Update Go backend to serve Vue SPA
  - Modify `internal/api/server.go`:
    - Serve static files from `frontend/dist/assets`
    - Add `NoRoute` handler to serve `index.html` for Vue Router history mode
  - Test locally: `make build-all && ./bin/english-agent`
- [ ] **T10.3**: Test full stack integration
  - Access app at `http://localhost:8080`
  - Verify all features work (API calls, routing, SSE streaming)
  - Check browser console for errors
- [ ] **T10.4**: Embed Vue app in Go binary
  - Use `//go:embed frontend/dist/*` in `cmd/main.go`
  - Update `server.go` to serve embedded FS
  - Test single-binary deployment
  - Verify all routes work with embedded assets
- [ ] **T10.5**: Update documentation
  - Update `README.md`:
    - Remove Streamlit references
    - Add Vue.js setup instructions (`make frontend-install`, `make dev-all`)
    - Document build process (`make build-all`)
  - Update `start.sh` to build Vue and start Go server
- [ ] **T10.6**: Remove Streamlit files
  - Delete `web/app.py`
  - Delete `web/requirements.txt`
  - Remove Python dependencies from CI/CD (if any)
- [ ] **T10.7**: Final smoke test
  - Run `make build-all && ./bin/english-agent`
  - Test all features end-to-end
  - Verify no regressions
- [ ] **T10.8** (Optional): Update OpenSpec `project.md`
  - Document Vue.js conventions (component structure, naming, styling)
  - Add frontend testing strategy

**Validation**: Vue app is fully integrated with Go backend, single command to build and run, Streamlit is removed, documentation is updated.

---

## Summary of Tasks

| Phase | Tasks | Dependencies |
|-------|-------|--------------|
| Phase 1 | T1.1 - T1.12 (12 tasks) | None |
| Phase 2 | T2.1 - T2.6 (6 tasks) | Phase 1 |
| Phase 3 | T3.1 - T3.7 (7 tasks) | Phase 2 |
| Phase 4 | T4.1 - T4.7 (7 tasks) | Phase 2, Phase 3 |
| Phase 5 | T5.1 - T5.7 (7 tasks) | Phase 4 |
| Phase 6 | T6.1 - T6.6 (6 tasks) | Phase 5 |
| Phase 7 | T7.1 - T7.8 (8 tasks) | Phase 5 |
| Phase 8 | T8.1 - T8.6 (6 tasks) | Phase 3 |
| Phase 9 | T9.1 - T9.8 (8 tasks) | Phase 6, Phase 7, Phase 8 |
| Phase 10 | T10.1 - T10.8 (8 tasks) | Phase 9 |
| **Total** | **75 tasks** | |

---

## Parallelization Opportunities

- **Phase 3, 4, 5** can be partially parallelized (different components, minimal overlap).
- **Phase 6, 7, 8** can be developed in parallel (streaming, TTS, settings are independent).
- **Phase 9** testing should be done after all features are complete.

---

## Estimated Effort

- **Phase 1**: 4-6 hours (setup, tooling)
- **Phase 2**: 2-4 hours (layout, i18n)
- **Phase 3**: 4-6 hours (RSS loading, API integration)
- **Phase 4**: 2-4 hours (URL fetch, text input)
- **Phase 5**: 6-8 hours (AI actions, markdown rendering)
- **Phase 6**: 4-6 hours (SSE streaming, debugging)
- **Phase 7**: 4-6 hours (TTS, link extraction)
- **Phase 8**: 3-5 hours (CRUD for custom feeds)
- **Phase 9**: 6-10 hours (testing, bug fixes, polish)
- **Phase 10**: 2-4 hours (integration, documentation)

**Total**: 37-59 hours (assuming 1-2 developers)

---

## Next Steps

1. Review and approve this task breakdown.
2. Create spec deltas for `frontend-ui` capability (see `specs/frontend-ui/spec.md`).
3. Validate proposal with `openspec validate migrate-to-vue-frontend --strict`.
4. Begin Phase 1 after approval.

