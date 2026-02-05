# Summary: Migrate Frontend from Streamlit to Vue.js 3 + TypeScript

**Change ID**: `migrate-to-vue-frontend`  
**Status**: ✅ Ready for Review  
**Validation**: Passed `openspec validate --strict`

---

## Quick Overview

This proposal outlines the migration of the existing Streamlit-based frontend to a modern Vue.js 3 + TypeScript single-page application (SPA), while maintaining full feature parity and backward compatibility with the Go backend API.

---

## Key Benefits

1. **Performance**: 
   - No Python runtime overhead
   - Client-side rendering with code splitting
   - Optimistic UI updates
   - Faster perceived performance with smooth transitions

2. **User Experience**:
   - Native browser API integration (no iframe hacks for TTS)
   - Smooth streaming with proper SSE handling
   - Better mobile responsiveness
   - Modern, polished UI with Element Plus components

3. **Developer Experience**:
   - Clear frontend/backend separation
   - TypeScript type safety (catch errors before runtime)
   - Hot Module Replacement (HMR) for instant feedback
   - Component-based architecture for maintainability

4. **Deployment**:
   - Single-binary deployment (Vue SPA embedded in Go binary)
   - No Python dependencies
   - Static assets served from CDN (optional)

---

## Technology Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| **Framework** | Vue.js 3.5+ | Composition API, excellent TypeScript support, mature ecosystem |
| **Language** | TypeScript 5+ | Type safety, better IDE support, fewer runtime errors |
| **Build Tool** | Vite 5+ | Fast HMR, optimized builds, great DX |
| **State Management** | Pinia | Vue 3 recommended store, cleaner API than Vuex |
| **Routing** | Vue Router 4 | History mode for SEO-friendly URLs |
| **UI Library** | Element Plus | Comprehensive components, TypeScript support, Material Design-inspired |
| **Styling** | Tailwind CSS | Utility-first, highly customizable, small bundle size |
| **HTTP Client** | Axios | Familiar API, interceptors for error handling |
| **i18n** | Vue I18n 9+ | Composition API support, SSR-ready |
| **Testing** | Vitest + Vue Test Utils | Fast unit/component tests, compatible with Vite |

---

## What's Changing

### ✅ Features Retained (100% Parity)

- RSS feed loading and article selection
- URL fetching for arbitrary articles
- Manual text input
- All AI actions: Explain, Summarize, Translate, Refine, Extract Sentences, Extract Vocabulary
- Streaming mode (SSE) for real-time LLM responses
- Text-to-Speech (TTS) for articles and sentences
- Link extraction with "Fetch" buttons
- Custom RSS feed management (CRUD)
- Internationalization (English/Chinese toggle)
- Collapsible sidebar

### 🆕 Improvements

- **Native SSE handling** (no more `iter_lines` workarounds, proper space handling)
- **Native TTS** (no more `components.v1.html` iframe hacks)
- **Client-side routing** (smooth page transitions, browser back/forward support)
- **Responsive layout** (better mobile/tablet support)
- **Markdown rendering** (proper display of formatted content)
- **Optimistic updates** (faster perceived performance)
- **Better error handling** (toast notifications, retry logic)

### 🗑️ Removed

- Streamlit Python files (`web/app.py`, `web/requirements.txt`)
- Python runtime dependency for frontend

### 🔄 Backend Changes (Minimal)

- Serve static files from `frontend/dist/`
- Add `NoRoute` handler for Vue Router history mode
- Optionally embed Vue SPA in Go binary (`//go:embed`)
- No breaking API changes

---

## Implementation Plan

### Phased Approach (10 Phases, 75 Tasks)

1. **Phase 1**: Project scaffold & setup (tooling, dependencies, config)
2. **Phase 2**: Core layout & i18n (UI shell, language toggle)
3. **Phase 3**: RSS feed loading (API integration)
4. **Phase 4**: URL fetcher & text input (other input modes)
5. **Phase 5**: AI actions (non-streaming)
6. **Phase 6**: Streaming support (SSE)
7. **Phase 7**: TTS & link extraction (enhanced UX)
8. **Phase 8**: Custom RSS feed management (settings)
9. **Phase 9**: Polish & testing (QA, bug fixes)
10. **Phase 10**: Backend integration & deployment (cutover)

**Estimated Effort**: 37-59 hours (1-2 developers, ~1-2 weeks)

---

## File Structure (After Migration)

```
english-study-agent-eino/
├── backend/                      # (Optional: rename from root)
│   ├── cmd/main.go
│   ├── internal/{agent,api,config,logger,rss,storage}/
│   ├── config.yaml
│   └── Makefile
├── frontend/                     # NEW Vue.js 3 + TypeScript SPA
│   ├── src/
│   │   ├── components/           # Reusable Vue components
│   │   ├── composables/          # Composition API hooks (useAgent, useTTS, etc.)
│   │   ├── i18n/                 # Translations (en.ts, zh.ts)
│   │   ├── router/               # Vue Router config
│   │   ├── stores/               # Pinia stores (app, content, rss, results)
│   │   ├── types/                # TypeScript types
│   │   ├── utils/                # Utilities (API client, link extractor)
│   │   ├── views/                # Top-level pages (HomeView, SettingsView)
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── vite.config.ts
│   ├── package.json
│   └── tsconfig.json
├── openspec/
├── README.md
└── start.sh                      # Updated to build Vue + start Go
```

---

## Success Criteria

✅ All existing features work (feature parity checklist completed)  
✅ SSE streaming works without formatting issues  
✅ TTS works natively (no iframes)  
✅ i18n toggle works (EN/CN)  
✅ Link extraction + fetch buttons work  
✅ Custom RSS feed CRUD works  
✅ Build process documented (`make build-all`)  
✅ Go backend serves Vue SPA (optional embedding)  
✅ No API regressions  
✅ Manual testing completed (all actions + edge cases)

---

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| SSE streaming doesn't work | Test early with existing backend; use proven `fetch` + `ReadableStream` pattern |
| Team unfamiliar with Vue | Provide training, use TypeScript for autocomplete, leverage docs |
| Missing features during migration | Strict feature parity checklist; thorough testing before cutover |
| Build complexity | Document setup clearly; automate with Makefile |

---

## ✅ Decisions Made

1. **UI Library**: Element Plus (comprehensive, TypeScript support)
2. **Build Tool**: Vite (fast, modern)
3. **State Management**: Pinia (Vue 3 recommended)
4. **Testing**: Vitest + Vue Test Utils (unit & component tests)
5. **Deployment**: Go backend serves Vue SPA (embedded in binary)
6. **Styling**: Tailwind CSS (utility-first, works with Element Plus)

## Remaining Questions (Optional)

1. **Rollout Strategy**: Immediate cutover or parallel deployment for testing?
2. **Accessibility Priority**: Phase 1 (core) or Phase 2 (polish)?
3. **Mobile Optimization**: Phase 2 (layout) or Phase 9 (polish)?
4. **E2E Tests**: Add Playwright/Cypress later or skip?

---

## Files Created

```
openspec/changes/migrate-to-vue-frontend/
├── proposal.md           ✅ Created (Why, What, Impact, Alternatives)
├── design.md             ✅ Created (Architecture, tech stack, component hierarchy, API integration)
├── tasks.md              ✅ Created (10 phases, 75 tasks, dependencies)
└── specs/
    └── frontend-ui/
        └── spec.md       ✅ Created (11 requirements, 47 scenarios)
```

**Validation**: ✅ Passed `openspec validate migrate-to-vue-frontend --strict`

---

## Next Steps

1. **Review this proposal** and decide on open questions.
2. **Approve** the proposal (or request changes).
3. **Apply** with `/openspec-apply migrate-to-vue-frontend` when ready to start implementation.

---

## Questions?

- Want to see code samples for specific components?
- Need clarification on any technical decisions?
- Want to adjust the tech stack (e.g., use PrimeVue instead of Element Plus)?
- Want to discuss the migration timeline or resource allocation?

Let me know and I'll provide more details! 🚀

