# Proposal: Migrate Frontend from Streamlit to Vue.js 3 + TypeScript

**Change ID**: `migrate-to-vue-frontend`  
**Status**: Draft  
**Author**: AI Assistant  
**Date**: 2026-01-02

---

## Why

**Problem**: The current Streamlit-based frontend is functional but has several limitations:

1. **Performance**: Streamlit reruns the entire Python script on every interaction, causing unnecessary overhead and latency.
2. **User Experience**: Limited control over UI/UX patterns (e.g., custom layouts, animations, granular state management).
3. **SSE Streaming Complexity**: Streamlit's component model requires workarounds (e.g., `components.v1.html` for TTS, manual SSE parsing) that are fragile.
4. **Development Workflow**: Python-based frontend mixing backend logic makes it harder to separate concerns and scale the team.
5. **Browser Feature Access**: Limited native browser API integration (Web Speech API requires iframe hacks).
6. **Deployment Flexibility**: Streamlit requires a Python runtime; a static SPA can be served from CDN or embedded in the Go binary.

**Opportunity**: Migrating to Vue.js 3 + TypeScript will provide:
- **Modern SPA architecture**: Component-based, reactive, with first-class TypeScript support.
- **Better UX**: Smooth transitions, optimistic updates, native SSE handling via `EventSource`.
- **Maintainability**: Clear separation of frontend (Vue) and backend (Go API), easier testing and debugging.
- **Ecosystem**: Rich UI libraries, build tools (Vite), and community support.
- **Performance**: Client-side rendering, code splitting, and minimal backend coupling.

---

## What Changes

### High-Level Scope

**Migrate** the existing Streamlit UI (`web/app.py`) to a **Vue.js 3 + TypeScript SPA** while **retaining** the Go backend API and all existing functionality.

### Affected Components

1. **Frontend (NEW)**:
   - Scaffold a new Vue 3 + TypeScript project using Vite.
   - Implement all existing UI features: RSS feed reader, article loading, URL fetching, AI actions, TTS, i18n (EN/CN), link extraction.
   - Use native `EventSource` API for SSE streaming from `/api/chat/stream`.
   - Integrate a UI component library (e.g., Element Plus or PrimeVue) for rapid development.
   - State management with Pinia (Vue 3 recommended store).

2. **Backend (MINIMAL CHANGES)**:
   - **No breaking changes** to existing API endpoints.
   - Optionally embed the built Vue SPA into the Go binary (using `embed` package) for single-binary deployment.
   - Add a fallback route to serve `index.html` for Vue Router history mode.

3. **Build & Deploy**:
   - Add `package.json`, `vite.config.ts`, `tsconfig.json` to the project root or `web/` directory.
   - Update `Makefile` and `start.sh` to build Vue app and serve it from Go.
   - Remove `web/app.py` and `web/requirements.txt` after migration is complete.

4. **Documentation**:
   - Update `README.md` with new frontend tech stack and development instructions.
   - Optionally update OpenSpec `project.md` with Vue conventions.

### Out of Scope

- **No new features**: This is a **like-for-like migration**. New features (e.g., user profiles, learning plans from `upgrade-to-intelligent-agent`) will come later.
- **No backend refactoring**: The Go API remains unchanged (unless minor CORS or static file serving tweaks are needed).
- **No database changes**: SQLite schema and storage logic remain as-is.

---

## Impact

### User Impact
- **Positive**:
  - Faster, more responsive UI with smooth transitions.
  - Better mobile responsiveness (if UI library supports it).
  - More reliable TTS and link navigation (native browser APIs).
- **Neutral**:
  - No visible feature changes; users should see the same functionality.
- **Negative** (Mitigation):
  - Potential short-term bugs during migration → **Mitigate with thorough testing and phased rollout**.

### Developer Impact
- **Positive**:
  - Clear frontend/backend separation enables parallel development.
  - TypeScript improves code quality and reduces runtime errors.
  - Easier to onboard frontend developers (Vue.js is popular).
- **Negative**:
  - Team needs to learn Vue.js 3 + TypeScript if unfamiliar → **Mitigate with documentation and training**.
  - Increased build tooling complexity (Node.js, Vite, `npm`/`pnpm`) → **Mitigate with clear setup scripts**.

### Technical Debt
- **Removes**: Python runtime dependency for frontend, Streamlit workarounds (iframe TTS, manual SSE parsing).
- **Adds**: Node.js build toolchain, frontend state management complexity.
- **Net**: Positive reduction in long-term technical debt.

---

## Alternatives Considered

1. **Keep Streamlit, optimize Python code**:
   - Pros: No migration cost, team already familiar.
   - Cons: Doesn't solve UX, performance, or deployment flexibility issues.

2. **Use React instead of Vue.js**:
   - Pros: Larger ecosystem, more developers available.
   - Cons: More boilerplate (e.g., Redux), Vue 3 Composition API is cleaner for this use case.

3. **Use Svelte or Alpine.js**:
   - Pros: Smaller bundle size, simpler syntax.
   - Cons: Smaller ecosystem, fewer UI libraries, less TypeScript maturity.

4. **Server-side rendering (SSR) with Next.js or Nuxt.js**:
   - Pros: SEO benefits, faster initial load.
   - Cons: Overkill for this use case (private learning tool, not public website).

**Decision**: Vue.js 3 + TypeScript strikes the best balance of ecosystem maturity, developer ergonomics, and technical fit.

---

## Success Criteria

1. ✅ All existing Streamlit features are replicated in Vue.js with feature parity.
2. ✅ SSE streaming works reliably with `EventSource` (no missing spaces or HTML escaping issues).
3. ✅ TTS (Web Speech API) works without iframe hacks.
4. ✅ i18n (English/Chinese) toggle works with Vue I18n.
5. ✅ Link extraction and "Fetch" buttons work correctly.
6. ✅ Custom RSS feed management (CRUD) works via API calls.
7. ✅ Build process is documented and automated (`make frontend`, `make run`).
8. ✅ Go backend can serve the built Vue SPA (optionally embedded).
9. ✅ No regressions in existing backend API behavior.
10. ✅ Manual testing checklist completed (all actions + edge cases).

---

## Next Steps

1. **Review this proposal** with the team and approve/reject/iterate.
2. **Scaffold the design doc** (`design.md`) with architecture decisions, component hierarchy, routing, state management, and API integration patterns.
3. **Create spec deltas** for `frontend-ui` (new capability) and updates to `ai-assistant` (streaming behavior clarification).
4. **Break down tasks** in `tasks.md` into phased implementation (scaffold → core features → polish → cutover).
5. **Validate** with `openspec validate migrate-to-vue-frontend --strict`.

---

## Design Decisions (Finalized)

✅ **UI Library**: Element Plus (comprehensive components, TypeScript support)  
✅ **Build Tool**: Vite (fast HMR, optimized builds)  
✅ **State Management**: Pinia (Vue 3 recommended store)  
✅ **Testing**: Vitest + Vue Test Utils (unit & component tests)  
✅ **Deployment**: Go backend serves Vue SPA (embedded in binary)  
✅ **Styling**: Tailwind CSS (utility-first, works alongside Element Plus)

## Remaining Questions

- **Rollout**: Immediate cutover or parallel deployment during testing?
- **Accessibility**: Priority in Phase 1 or defer to Phase 2?
- **Mobile**: Optimize in Phase 2 (layout) or Phase 9 (polish)?
- **E2E Tests**: Add Playwright/Cypress in future or skip?

