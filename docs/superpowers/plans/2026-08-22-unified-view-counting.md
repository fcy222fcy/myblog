# Unified View Counting Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make article and daily-question views concurrency-safe daily unique views with privacy-preserving anonymous identity and explicit frontend recording.

**Architecture:** A generic `content_view_events` repository owns the unique event insert and target counter update in one transaction. A dedicated public view endpoint derives an HMAC visitor key, while all content GET endpoints become side-effect free and the Vue frontend records only genuinely displayed content.

**Tech Stack:** Go 1.25, Gin, GORM, MySQL/SQLite tests, Vue 3, Axios, Node test runner

**Spec:** `docs/superpowers/specs/2026-08-22-unified-view-counting-design.md`

## Global Constraints

- One visitor, one content item, one Asia/Shanghai calendar day counts at most once.
- Never persist raw anonymous IDs or raw IP addresses in the new event table.
- A content event and its counter increment must commit or roll back together.
- GET content endpoints have no view-counting side effects.
- Preserve existing cumulative `view_count` values and unrelated worktree changes.
- Do not commit or push without explicit authorization.

---

### Task 1: Atomic content-view repository

**Files:**
- Create: `internal/model/entity/content_view_event.go`
- Create: `internal/repository/content_view_repository.go`
- Create: `internal/repository/content_view_repository_test.go`
- Modify: `pkg/database/database.go`

**Interfaces:**
- Produces: `RecordView(contentType string, contentID uint, visitorKey string, viewDate time.Time) (counted bool, viewCount int64, err error)`
- Produces: `CountByDate(contentType string, viewDate time.Time) (int64, error)`

- [ ] Write SQLite tests creating published/draft articles and daily questions. Assert first record returns `(true, 1)`, same key/day returns `(false, 1)`, next day returns `(true, 2)`, invalid targets leave zero events, and 20 concurrent identical calls produce one event and one increment.
- [ ] Run `go test ./internal/repository -run ContentView -count=1` and confirm failure because the repository/entity do not exist.
- [ ] Add the entity with unique index `(content_type, content_id, visitor_key, view_date)` and implement transactional insert using `clause.OnConflict{DoNothing:true}` followed by a guarded atomic counter update.
- [ ] Register the entity in `Database.AutoMigrate` and rerun the targeted tests until green.

### Task 2: Identity derivation and view service/API

**Files:**
- Create: `pkg/visitor/identity.go`
- Create: `pkg/visitor/identity_test.go`
- Create: `internal/service/content_view_service.go`
- Create: `internal/service/content_view_service_test.go`
- Create: `internal/api/v1/content_view/controller.go`
- Create: `internal/api/v1/content_view/routes.go`
- Modify: `internal/app/app.go`
- Modify: `internal/api/routes.go`
- Modify: `pkg/config/config.go`
- Modify: `.env.example`
- Modify: `docker-compose.yml`
- Test: `internal/api/trusted_proxy_test.go`

**Interfaces:**
- Produces: `visitor.DeriveKey(secret, anonymousID, clientIP string) (string, error)`
- Produces: `ContentViewService.Record(contentType string, contentID uint, visitorKey string) (*ContentViewResult, error)`
- Consumes: repository `RecordView` and `CountByDate`.

- [ ] Write tests proving valid UUID-like IDs are stable, raw values are absent from derived keys, invalid/missing IDs fall back to IP, and empty identity fails.
- [ ] Write service tests for allowed types, invalid IDs, repository error propagation and Asia/Shanghai date selection.
- [ ] Run the focused tests and confirm the missing APIs fail.
- [ ] Implement the identity helper, service, request DTO `{content_type, content_id}`, controller response `{counted, view_count}`, route `POST /views`, dependency wiring, trusted-proxy configuration and dashboard today-count migration to `content_view_events`.
- [ ] Run focused tests until green.

### Task 3: Remove GET side effects and fix daily-question payload

**Files:**
- Modify: `internal/service/article_interface.go`
- Modify: `internal/service/article_service.go`
- Modify: `internal/api/v1/article/controller.go`
- Modify: `internal/service/article_service_test.go`
- Modify: `internal/service/daily_question_interface.go`
- Modify: `internal/service/daily_question_service.go`
- Create: `internal/service/daily_question_service_test.go`

**Interfaces:**
- Changes: `GetArticleDetail(slug string)` no longer accepts client IP.
- Changes: `GetAllPublishedQuestions() ([]*response.DailyQuestionResponse, error)` returns complete published records.

- [ ] Add tests proving repeated article GETs never call a view repository and repeated latest/date/all daily reads do not alter `view_count`.
- [ ] Run focused service tests and confirm failure against the current mutating behavior.
- [ ] Remove article `HasVisited`/increment branches and daily-question increments; convert the public all-list to full responses.
- [ ] Run focused service tests until green.

### Task 4: Frontend visitor identity and explicit recording

**Files:**
- Create: `blog-web/src/utils/visitorIdentity.js`
- Create: `blog-web/src/utils/visitorIdentity.test.js`
- Create: `blog-web/src/utils/viewRecorder.js`
- Create: `blog-web/src/utils/viewRecorder.test.js`
- Create: `blog-web/src/api/view.js`
- Modify: `blog-web/src/api/request.js`
- Modify: `internal/middleware/auth.go`
- Test: `internal/middleware/cors_test.go`
- Modify: `blog-web/src/stores/article.js`
- Modify: `blog-web/src/views/Article.vue`
- Modify: `blog-web/src/views/DailyQuestion.vue`

**Interfaces:**
- Produces: `getOrCreateVisitorID(storage, cryptoProvider)`.
- Produces: `createPageViewRecorder(recordFn)` which suppresses duplicate content keys for one page lifecycle.
- Consumes: `POST /views` response `{counted, view_count}`.

- [ ] Write Node tests proving visitor ID reuse/fallback generation and per-page duplicate suppression.
- [ ] Run `node --test src/utils/visitorIdentity.test.js src/utils/viewRecorder.test.js` and confirm missing modules fail.
- [ ] Implement `X-Visitor-ID` header injection and the view API.
- [ ] Record articles only after successful detail load; update `currentArticle.view_count` from the server response.
- [ ] Replace the daily N+1 detail loading with the full `/all` payload and an `IntersectionObserver`; record each visible question once and apply the returned count.
- [ ] Run Node tests and `npm run build` until green.

### Task 5: Final verification and handoff

**Files:**
- Verify all changed files; no new production files beyond Tasks 1-4.

- [ ] Run targeted repository, service and visitor tests with repository-local Go caches.
- [ ] Run `go test ./... -count=1` with repository-local caches.
- [ ] Run the two Node unit-test files and `npm run build` in `blog-web`.
- [ ] Run `gofmt` on changed Go files and `git diff --check`.
- [ ] Inspect `git status --short` and confirm `outputs/` and all unrelated existing changes were untouched.
- [ ] Report code-level verification separately from pending live MySQL/Redis/browser acceptance.
