# Normal Routing and Reading Experience Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Serve the public and admin Vue applications at stable history URLs and enrich article reading with TOC, code tools, and resilient images.

**Architecture:** Keep the existing two-SPA Nginx gateway. Move both routers to HTML5 history using their Vite base path, and isolate Markdown presentation transforms in a pure utility consumed by the article view.

**Tech Stack:** Vue 3, Vue Router 4, Vite 5, Marked 12, Highlight.js 11, Node test runner, Nginx.

---

### Task 1: Article rendering utility

**Files:**
- Create: `blog-web/src/utils/articleEnhancements.js`
- Test: `blog-web/src/utils/articleEnhancements.test.js`

- [ ] Run the new test and confirm it fails because the utility is missing.
- [ ] Implement heading anchors and TOC extraction.
- [ ] Implement highlighted code wrappers with language metadata and numbered lines.
- [ ] Implement lazy image markup with a failure placeholder.
- [ ] Run the focused test and confirm all cases pass.

### Task 2: Article view integration

**Files:**
- Modify: `blog-web/src/views/Article.vue`
- Modify: `blog-web/src/assets/styles/main.css`

- [ ] Replace direct Marked rendering with the tested utility.
- [ ] Render a Hugo Theme Stack-style sticky TOC beside the article and track the active heading without changing the progress ring.
- [ ] Add delegated copy and captured image-error behavior.
- [ ] Add responsive light/dark styles using existing design tokens.

### Task 3: History routing and gateway ownership

**Files:**
- Modify: `blog-web/src/router/index.js`
- Modify: `blog-admin/src/router/index.js`
- Modify: `docker/nginx/nginx.conf`

- [ ] Use `createWebHistory(import.meta.env.BASE_URL)` in both applications.
- [ ] Redirect exact `/admin` requests to `/admin/`.
- [ ] Preserve separate admin and public SPA fallbacks.

### Task 4: Verification

- [ ] Run article utility and existing frontend utility tests.
- [ ] Build `blog-web` and `blog-admin` for production.
- [ ] Inspect generated admin asset URLs and Nginx routing rules.
- [ ] Run page-level smoke checks when the local services are available.
