# Tech Stack Reference: Vue.js 3 Frontend

**Change ID**: `migrate-to-vue-frontend`  
**Status**: ✅ Finalized & Validated

---

## Core Stack

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| **Framework** | Vue.js | 3.5+ | Reactive UI framework with Composition API |
| **Language** | TypeScript | 5+ | Type-safe JavaScript with excellent IDE support |
| **Build Tool** | Vite | 5+ | Fast dev server, optimized production builds |
| **Package Manager** | pnpm | Latest | Fast, disk-efficient package management |

---

## Frontend Libraries

| Category | Library | Purpose |
|----------|---------|---------|
| **State Management** | Pinia | Centralized state store (Vue 3 recommended) |
| **Routing** | Vue Router 4 | Client-side routing with history mode |
| **UI Components** | Element Plus | Material Design-inspired component library |
| **Styling** | Tailwind CSS | Utility-first CSS framework |
| **HTTP Client** | Axios | Promise-based HTTP client with interceptors |
| **i18n** | Vue I18n 9+ | Internationalization with Composition API support |

---

## Development Tools

| Category | Tool | Purpose |
|----------|------|---------|
| **Testing** | Vitest | Fast unit test runner (Vite-native) |
| **Component Testing** | Vue Test Utils | Vue component testing utilities |
| **Linting** | ESLint + Vue ESLint Plugin | Code quality and style enforcement |
| **Formatting** | Prettier | Consistent code formatting |
| **Type Checking** | vue-tsc | Vue TypeScript compiler |

---

## Integration with Go Backend

- **Deployment**: Go backend serves the built Vue SPA
- **Embedding**: Vue `dist/` folder embedded in Go binary using `//go:embed`
- **Routing**: Go serves `index.html` for all non-API routes (Vue Router history mode)
- **API**: Axios calls Go backend at `/api/*` endpoints
- **Streaming**: Native `fetch` API with `ReadableStream` for SSE

---

## Key Design Patterns

1. **Composition API**: Use `<script setup>` syntax for components
2. **Composables**: Extract reusable logic into `composables/` (e.g., `useAgent`, `useTTS`, `useStreaming`)
3. **Pinia Stores**: One store per domain (`app`, `content`, `rss`, `results`)
4. **TypeScript Types**: All props, events, state, and API responses are typed
5. **Scoped Styles**: Tailwind utilities + Element Plus components (preflight disabled to avoid conflicts)

---

## File Structure

```
frontend/
├── src/
│   ├── components/        # Reusable Vue components
│   ├── composables/       # Composition API logic (useAgent, useTTS, etc.)
│   ├── i18n/              # Translation files (en.ts, zh.ts)
│   ├── router/            # Vue Router config
│   ├── stores/            # Pinia stores (app, content, rss, results)
│   ├── types/             # TypeScript type definitions
│   ├── utils/             # Utilities (API client, link extractor)
│   ├── views/             # Top-level pages (HomeView, SettingsView)
│   ├── styles/            # Global styles (Tailwind imports)
│   ├── App.vue            # Root component
│   └── main.ts            # Entry point
├── public/                # Static assets
├── index.html             # HTML template
├── vite.config.ts         # Vite configuration
├── tailwind.config.js     # Tailwind configuration
├── tsconfig.json          # TypeScript configuration
├── package.json           # Dependencies
└── .eslintrc.cjs          # ESLint configuration
```

---

## Development Commands

```bash
# Install dependencies
make frontend-install
# or: cd frontend && pnpm install

# Run dev server (with HMR)
make frontend-dev
# or: cd frontend && pnpm dev

# Build for production
make frontend-build
# or: cd frontend && pnpm build

# Run tests
cd frontend && pnpm test

# Lint code
cd frontend && pnpm lint

# Type check
cd frontend && pnpm type-check
```

---

## Production Build Output

```
frontend/dist/
├── index.html           # Main HTML entry point
├── assets/
│   ├── index-[hash].js  # Main JS bundle (minified, tree-shaken)
│   ├── index-[hash].css # Main CSS bundle (Tailwind + Element Plus)
│   └── [other assets]   # Fonts, images, icons
└── favicon.ico
```

**Go Backend Integration**:
- Embed `frontend/dist/*` using `//go:embed`
- Serve static files from `/assets`
- Fallback to `index.html` for all non-API routes

---

## Notes

- **Tailwind Preflight Disabled**: To avoid conflicts with Element Plus base styles
- **Vite Proxy**: Dev server proxies `/api` requests to `http://localhost:8080` (Go backend)
- **SSE Streaming**: Use native `fetch` API (not EventSource) for better control
- **TTS**: Native Web Speech API (no iframe workarounds)
- **Testing**: Vitest for unit tests, Vue Test Utils for component tests (E2E tests optional)

---

## Next Steps

1. **Review** the proposal documents (`proposal.md`, `design.md`, `tasks.md`)
2. **Apply** the change: `/openspec-apply migrate-to-vue-frontend`
3. **Start implementation** with Phase 1 (Project Scaffold & Setup)

---

## References

- [Vue 3 Documentation](https://vuejs.org/)
- [Vite Documentation](https://vitejs.dev/)
- [Element Plus Documentation](https://element-plus.org/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Vitest Documentation](https://vitest.dev/)

