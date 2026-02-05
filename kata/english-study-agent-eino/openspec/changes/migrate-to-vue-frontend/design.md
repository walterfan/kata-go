# Design: Migrate Frontend from Streamlit to Vue.js 3 + TypeScript

**Change ID**: `migrate-to-vue-frontend`  
**Status**: Draft  
**Last Updated**: 2026-01-02

---

## Context

This design doc covers the technical decisions, architecture, and implementation strategy for migrating the existing Streamlit-based frontend to a Vue.js 3 + TypeScript SPA.

### Current State
- **Frontend**: Streamlit (`web/app.py`, ~1339 lines) with custom HTML components for TTS.
- **Backend**: Gin (Go) API with REST endpoints + SSE streaming (`/api/chat/stream`).
- **Features**: RSS feed reader, article loading, URL fetching, AI actions (explain, summarize, translate, etc.), TTS, i18n (EN/CN), link extraction.
- **Limitations**: Streamlit's rerun model, iframe hacks for browser APIs, mixed Python/JS logic.

### Goals
1. Replicate all existing features in Vue.js with feature parity.
2. Improve UX with smooth transitions, optimistic updates, and native browser API integration.
3. Maintain backward compatibility with Go backend API (no breaking changes).
4. Enable future extensibility (e.g., user profiles, learning plans).

---

## Technology Stack

### Frontend Stack
- **Framework**: Vue.js 3.5+ (Composition API)
- **Language**: TypeScript 5+
- **Build Tool**: Vite 5+ (fast HMR, optimized builds)
- **State Management**: Pinia (Vue 3 recommended store)
- **Routing**: Vue Router 4 (history mode)
- **UI Library**: **Element Plus** (Material Design-inspired, TypeScript support, comprehensive components)
- **HTTP Client**: Axios (familiar, interceptors for error handling)
- **SSE Handling**: Native `EventSource` API
- **i18n**: Vue I18n 9+ (supports Composition API)
- **TTS**: Native Web Speech API (`SpeechSynthesis`)
- **Styling**: **Tailwind CSS** (utility-first, highly customizable, small bundle size)
  - Element Plus + Tailwind can coexist (Element Plus for components, Tailwind for custom styling)
- **Testing**: **Vitest** (unit tests) + **Vue Test Utils** (component tests)

### Build & Deploy
- **Package Manager**: `pnpm` (fast, disk-efficient)
- **Linting**: ESLint + Vue ESLint Plugin + Prettier
- **Type Checking**: `vue-tsc` (Vue TypeScript compiler)
- **Deployment**: **Go backend serves the Vue SPA**
  - Embed built Vue SPA into Go binary using `embed` package
  - Serve static files from `/assets`
  - Fallback to `index.html` for Vue Router history mode
  - Single binary deployment for simplicity

---

## Architecture

### Project Structure

```
english-study-agent-eino/
├── backend/                      # Renamed from root structure
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── agent/
│   │   ├── api/
│   │   ├── config/
│   │   ├── logger/
│   │   ├── rss/
│   │   └── storage/
│   ├── config.yaml
│   ├── go.mod
│   └── Makefile
├── frontend/                     # NEW Vue.js app
│   ├── public/
│   │   └── favicon.ico
│   ├── src/
│   │   ├── assets/               # Static assets (images, fonts)
│   │   ├── components/           # Reusable Vue components
│   │   │   ├── ArticleLoader.vue
│   │   │   ├── ActionButtons.vue
│   │   │   ├── RssFeedList.vue
│   │   │   ├── UrlFetcher.vue
│   │   │   ├── ResultDisplay.vue
│   │   │   └── LanguageToggle.vue
│   │   ├── composables/          # Composition API hooks
│   │   │   ├── useAgent.ts       # AI agent API calls
│   │   │   ├── useRss.ts         # RSS feed API calls
│   │   │   ├── useTTS.ts         # Text-to-Speech logic
│   │   │   └── useStreaming.ts   # SSE streaming logic
│   │   ├── i18n/                 # Translations
│   │   │   ├── index.ts
│   │   │   ├── en.ts
│   │   │   └── zh.ts
│   │   ├── router/               # Vue Router config
│   │   │   └── index.ts
│   │   ├── stores/               # Pinia stores
│   │   │   ├── app.ts            # Global app state (lang, mode)
│   │   │   ├── content.ts        # Current text, articles
│   │   │   ├── rss.ts            # RSS feeds & sources
│   │   │   └── results.ts        # AI action results
│   │   ├── types/                # TypeScript types
│   │   │   ├── api.ts
│   │   │   ├── rss.ts
│   │   │   └── agent.ts
│   │   ├── utils/                # Utilities
│   │   │   ├── api.ts            # Axios instance + interceptors
│   │   │   ├── linkExtractor.ts  # Extract links from HTML/Markdown
│   │   │   └── markdown.ts       # Markdown parsing helpers
│   │   ├── views/                # Top-level pages
│   │   │   ├── HomeView.vue      # Main learning interface
│   │   │   └── SettingsView.vue  # RSS feed management
│   │   ├── App.vue               # Root component
│   │   └── main.ts               # Entry point
│   ├── index.html
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── package.json
│   └── .eslintrc.cjs
├── openspec/
├── README.md
└── start.sh                      # Updated to build Vue + start Go server
```

**Note**: We could keep the current structure and put Vue in `web/` or `frontend/`. For clarity, I propose renaming Go code to `backend/` and Vue to `frontend/`. This is optional.

---

## Component Hierarchy

### Main Views

1. **HomeView** (Main Learning Interface)
   - **Sidebar** (Collapsible)
     - `LanguageToggle` (EN/CN)
     - `ArticleLoader` (RSS feed selection, article list)
     - `UrlFetcher` (URL input + fetch button)
     - `ActionButtons` (Explain, Summarize, Translate, etc.)
   - **Main Pane**
     - **Header**: Title + proverb
     - **CurrentTextDisplay** (Expander with article content + link extraction)
     - **ResultDisplay** (AI action results with streaming support)
     - **TTS Controls** (Read/Stop toggle button)

2. **SettingsView** (RSS Feed Management)
   - **CustomFeedList** (CRUD interface for custom RSS feeds)
   - **DefaultFeedList** (Read-only list of default feeds from `config.yaml`)

### Reusable Components

- `ActionButtons.vue`: Grid of action buttons (Explain, Summarize, etc.)
- `ArticleLoader.vue`: RSS source dropdown + article list + preview
- `RssFeedList.vue`: Display list of RSS articles with metadata
- `UrlFetcher.vue`: URL input + "Load Article" button
- `ResultDisplay.vue`: Handles streaming results, markdown rendering, TTS integration
- `LanguageToggle.vue`: EN/CN language switcher
- `LinkList.vue`: Displays extracted links with "Fetch" buttons

---

## State Management (Pinia Stores)

### 1. `appStore` (Global App State)
```typescript
interface AppState {
  language: 'en' | 'zh';
  inputMode: 'article' | 'text' | 'url';
  sidebarCollapsed: boolean;
  streamingEnabled: boolean;
  loading: boolean;
}
```

### 2. `contentStore` (Current Content)
```typescript
interface ContentState {
  currentText: string;
  fetchedUrl: string;
  selectedArticle: Article | null;
  extractedLinks: Array<{ url: string; text: string }>;
}
```

### 3. `rssStore` (RSS Feeds)
```typescript
interface RssState {
  sources: RssSource[];
  articles: Article[];
  customFeeds: CustomFeed[];
  selectedSource: string;
}
```

### 4. `resultsStore` (AI Action Results)
```typescript
interface ResultsState {
  currentResult: string;
  streamingTask: string | null;
  isStreaming: boolean;
  history: Array<{ task: string; result: string; timestamp: Date }>;
}
```

---

## API Integration

### HTTP Client Setup (Axios)

```typescript
// src/utils/api.ts
import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Error interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle errors globally (e.g., show toast notification)
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export default api;
```

### SSE Streaming (EventSource)

```typescript
// src/composables/useStreaming.ts
import { ref } from 'vue';

export function useStreaming() {
  const content = ref('');
  const isStreaming = ref(false);
  const error = ref<string | null>(null);

  function startStream(url: string, payload: any) {
    content.value = '';
    error.value = null;
    isStreaming.value = true;

    // POST request to get streaming endpoint
    fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
      .then((response) => {
        const reader = response.body?.getReader();
        const decoder = new TextDecoder();

        function read() {
          reader?.read().then(({ done, value }) => {
            if (done) {
              isStreaming.value = false;
              return;
            }

            const chunk = decoder.decode(value, { stream: true });
            const lines = chunk.split('\n');

            for (const line of lines) {
              if (line.startsWith('data:')) {
                let data = line.slice(5); // Remove "data:"
                if (data.startsWith(' ')) data = data.slice(1); // Remove SSE space
                if (data) {
                  data = data.replace(/\\n/g, '\n'); // Unescape newlines
                  content.value += data;
                }
              } else if (line.startsWith('event: done')) {
                isStreaming.value = false;
                return;
              } else if (line.startsWith('event: error')) {
                error.value = 'Streaming error occurred';
                isStreaming.value = false;
                return;
              }
            }

            read(); // Continue reading
          });
        }

        read();
      })
      .catch((err) => {
        error.value = err.message;
        isStreaming.value = false;
      });
  }

  function stopStream() {
    isStreaming.value = false;
  }

  return { content, isStreaming, error, startStream, stopStream };
}
```

### TTS Integration (Web Speech API)

```typescript
// src/composables/useTTS.ts
import { ref } from 'vue';

export function useTTS() {
  const isSpeaking = ref(false);
  const utterance = ref<SpeechSynthesisUtterance | null>(null);

  function speak(text: string, rate = 0.9) {
    if (isSpeaking.value) {
      stop();
    }

    utterance.value = new SpeechSynthesisUtterance(text);
    utterance.value.lang = 'en-US';
    utterance.value.rate = rate;

    // Try to use an English voice
    const voices = window.speechSynthesis.getVoices();
    const englishVoice = voices.find((v) => v.lang.startsWith('en'));
    if (englishVoice) {
      utterance.value.voice = englishVoice;
    }

    utterance.value.onstart = () => {
      isSpeaking.value = true;
    };

    utterance.value.onend = () => {
      isSpeaking.value = false;
    };

    utterance.value.onerror = () => {
      isSpeaking.value = false;
    };

    window.speechSynthesis.speak(utterance.value);
  }

  function stop() {
    window.speechSynthesis.cancel();
    isSpeaking.value = false;
  }

  function toggle(text: string, rate = 0.9) {
    if (isSpeaking.value) {
      stop();
    } else {
      speak(text, rate);
    }
  }

  return { isSpeaking, speak, stop, toggle };
}
```

---

## Routing

### Routes

```typescript
// src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import SettingsView from '@/views/SettingsView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/settings',
      name: 'settings',
      component: SettingsView,
    },
  ],
});

export default router;
```

### Go Backend Changes for History Mode

To support Vue Router's history mode, the Go backend needs to serve `index.html` for all non-API routes:

```go
// internal/api/server.go
func (s *Server) setupRoutes() {
	// API routes
	api := s.router.Group("/api")
	{
		// ... existing routes
	}

	// Serve static files (embedded or from disk)
	s.router.Static("/assets", "./frontend/dist/assets")
	s.router.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	// Fallback to index.html for Vue Router history mode
	s.router.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})
}
```

---

## i18n (Internationalization)

### Setup

```typescript
// src/i18n/index.ts
import { createI18n } from 'vue-i18n';
import en from './en';
import zh from './zh';

const i18n = createI18n({
  legacy: false, // Use Composition API
  locale: localStorage.getItem('language') || 'en',
  fallbackLocale: 'en',
  messages: {
    en,
    zh,
  },
});

export default i18n;
```

### Translation Files

```typescript
// src/i18n/en.ts
export default {
  title: '📚 English Agent',
  subtitle: 'AI-powered English learning',
  main_title: 'The limits of my language mean the limits of my world.',
  main_subtitle: '— Ludwig Wittgenstein | AI-powered English learning',
  // ... (copy all translations from Streamlit app.py)
};

// src/i18n/zh.ts (similar structure)
```

### Usage in Components

```vue
<template>
  <h1>{{ t('main_title') }}</h1>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n';
const { t } = useI18n();
</script>
```

---

## Link Extraction

### Utility Function

```typescript
// src/utils/linkExtractor.ts
export interface ExtractedLink {
  url: string;
  text: string;
}

export function extractLinks(html: string, baseUrl = ''): ExtractedLink[] {
  const links: ExtractedLink[] = [];
  const seen = new Set<string>();

  // HTML <a> tags
  const htmlRegex = /<a\s+(?:[^>]*?\s+)?href=["'](.*?)(?=["'])[^>]*?>(.*?)<\/a>/gi;
  let match;
  while ((match = htmlRegex.exec(html)) !== null) {
    const url = resolveUrl(match[1], baseUrl);
    const text = stripHtml(match[2]);
    if (url && !seen.has(url)) {
      seen.add(url);
      links.push({ url, text: text || url });
    }
  }

  // Markdown [text](url)
  const mdRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
  while ((match = mdRegex.exec(html)) !== null) {
    const url = resolveUrl(match[2], baseUrl);
    const text = match[1];
    if (url && !seen.has(url)) {
      seen.add(url);
      links.push({ url, text });
    }
  }

  return links;
}

function resolveUrl(url: string, baseUrl: string): string {
  try {
    return new URL(url, baseUrl || window.location.origin).href;
  } catch {
    return url;
  }
}

function stripHtml(html: string): string {
  return html.replace(/<[^>]*>/g, '').trim();
}
```

---

## Build & Deployment

### Vite Config

```typescript
// frontend/vite.config.ts
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: false,
  },
  css: {
    postcss: {
      plugins: [
        require('tailwindcss'),
        require('autoprefixer'),
      ],
    },
  },
});
```

### Tailwind Config

```javascript
// frontend/tailwind.config.js
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {},
  },
  plugins: [],
  // Don't conflict with Element Plus
  corePlugins: {
    preflight: false, // Disable Tailwind's base styles to avoid conflicts with Element Plus
  },
};
```

### Makefile Updates

```makefile
# Existing Go targets
.PHONY: build run clean

build:
	go build -o bin/english-agent cmd/main.go

run: build
	./bin/english-agent

clean:
	rm -rf bin/

# NEW: Frontend targets
.PHONY: frontend-install frontend-dev frontend-build

frontend-install:
	cd frontend && pnpm install

frontend-dev:
	cd frontend && pnpm dev

frontend-build:
	cd frontend && pnpm build

# Combined targets
.PHONY: dev-all build-all

dev-all:
	# Run backend and frontend in parallel
	make -j2 run frontend-dev

build-all: frontend-build build
	# Build frontend first, then Go binary with embedded assets
```

### Embedding Vue App in Go Binary (Optional)

```go
// cmd/main.go
package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/english-agent/internal/api"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
	// ... (existing initialization)

	server, err := api.NewServer(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Serve embedded frontend
	distFS, _ := fs.Sub(frontendFS, "frontend/dist")
	server.Router().StaticFS("/assets", http.FS(distFS))
	server.Router().NoRoute(func(c *gin.Context) {
		data, _ := frontendFS.ReadFile("frontend/dist/index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	server.Run()
}
```

---

## Migration Strategy

### Phase 1: Scaffold & Setup (No features)
- Scaffold Vue project with Vite, TypeScript, Element Plus, Pinia, Vue Router, Vue I18n.
- Configure build, linting, and basic routing.
- Serve Vue app from Go backend (proxy API during dev).

### Phase 2: Core UI Components (No AI yet)
- Implement `LanguageToggle`, `ArticleLoader`, `UrlFetcher`, `ActionButtons`, `ResultDisplay`.
- Integrate i18n (copy translations from Streamlit).
- Wire up RSS feed loading (GET `/api/rss-sources`, `/api/feeds`).

### Phase 3: AI Actions (Non-streaming)
- Implement API calls for all actions (explain, summarize, translate, etc.) using non-streaming endpoint (`/api/chat`).
- Display results in `ResultDisplay` component.
- Test all actions with existing backend.

### Phase 4: Streaming Support
- Implement `useStreaming` composable with SSE via `fetch` + `ReadableStream`.
- Add streaming toggle in UI.
- Test with `/api/chat/stream` endpoint.

### Phase 5: TTS & Link Extraction
- Implement `useTTS` composable with Web Speech API.
- Add "Read" toggle button in `ResultDisplay`.
- Implement link extraction utility and `LinkList` component.

### Phase 6: Custom RSS Feed Management
- Implement Settings view with CRUD for custom feeds.
- Wire up API calls (GET/POST/PUT/DELETE `/api/custom-feeds`).

### Phase 7: Polish & Testing
- Manual testing of all features.
- Fix edge cases, styling, responsive layout.
- Unit tests for composables and utilities.
- Update documentation.

### Phase 8: Cutover
- Remove `web/app.py` and `web/requirements.txt`.
- Update `README.md` and `start.sh`.
- Optionally embed Vue SPA into Go binary.

---

## Risks & Mitigations

| Risk | Impact | Likelihood | Mitigation |
|------|--------|-----------|------------|
| SSE streaming doesn't work correctly in Vue | High | Low | Test early with existing backend; use proven `fetch` + `ReadableStream` pattern. |
| Element Plus doesn't meet UI needs | Medium | Low | Evaluate alternative libraries (PrimeVue) in Phase 1; Element Plus is mature and comprehensive. |
| Team unfamiliar with Vue.js 3 Composition API | Medium | Medium | Provide training, use clear examples, leverage TypeScript for autocomplete. |
| Go static file serving breaks API routes | High | Low | Test routing carefully; use `/api` prefix for all backend routes, serve Vue SPA with `NoRoute` fallback. |
| Embedding frontend in Go binary increases binary size | Low | High | Accept trade-off for single-binary deployment; alternatively, serve separately. |
| Missing Streamlit features during migration | High | Low | Strict feature parity checklist; thorough manual testing before cutover. |

---

## Design Decisions (Finalized)

1. ✅ **UI Library**: Element Plus (Material Design-inspired, comprehensive components)
2. ✅ **Build Tool**: Vite (fast HMR, optimized builds)
3. ✅ **State Management**: Pinia (Vue 3 recommended store)
4. ✅ **Testing**: Vitest (unit tests) + Vue Test Utils (component tests)
5. ✅ **Deployment**: Go backend serves Vue SPA (embedded in binary)
6. ✅ **Styling**: Tailwind CSS (utility-first, works alongside Element Plus)

## Open Questions

1. **Rollout**: Immediate cutover or run Streamlit and Vue in parallel during testing?
2. **Accessibility**: Should we prioritize ARIA labels, keyboard navigation, and screen reader support in Phase 1 or Phase 2?
3. **Mobile**: Should we optimize for mobile/tablet layouts in Phase 2 (layout) or defer to Phase 9 (polish)?
4. **E2E Testing**: Should we add E2E tests (Playwright/Cypress) in a future phase, or skip for now?

---

## Success Criteria (Repeated from Proposal)

1. ✅ All existing Streamlit features replicated in Vue.js.
2. ✅ SSE streaming works correctly (no space/HTML escaping issues).
3. ✅ TTS works natively (no iframes).
4. ✅ i18n toggle works (EN/CN).
5. ✅ Link extraction + fetch buttons work.
6. ✅ Custom RSS feed CRUD works.
7. ✅ Build process documented and automated.
8. ✅ Go backend serves Vue SPA (optional embedding).
9. ✅ No API regressions.
10. ✅ Manual testing checklist completed.

---

## Next Steps

See `tasks.md` for detailed implementation tasks.

