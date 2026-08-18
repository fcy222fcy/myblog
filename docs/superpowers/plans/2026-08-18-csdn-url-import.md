# CSDN URL Import Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an authenticated single-URL CSDN importer that safely creates an editable local draft with source metadata and best-effort local images.

**Architecture:** A dedicated `article_import` API delegates to a focused import service. The service validates and fetches public resources, parses CSDN HTML into sanitized Markdown, localizes images, prevents duplicate source URLs, and calls the existing article service to create a draft. The Vue admin adds a modal that collects URL and category, then navigates to the created draft.

**Tech Stack:** Go 1.25, Gin, GORM, `golang.org/x/net/html`, Vue 3, Element Plus, Axios, Vite.

**Spec:** `docs/superpowers/specs/2026-08-18-csdn-url-import-design.md`

## Global Constraints

- Only public single-article URLs matching `blog.csdn.net/<user>/article/details/<numeric-id>` are accepted.
- Imports always create `draft` articles and never publish automatically.
- No CSDN cookies, credentials, headless browser, generic extraction, bulk import, comment import, or statistics import.
- Page body limit: 5 MiB. Converted Markdown limit: 500,000 bytes.
- Image limits: 30 images, 5 MiB each, 50 MiB total; image failures are warnings.
- Every outbound request rejects loopback, private, unspecified, multicast, and link-local IP addresses, including redirects.
- Existing untracked workspace files must not be staged or changed.

---

### Task 1: Safe remote fetching and CSDN URL normalization

**Files:**
- Create: `internal/service/article_import_http.go`
- Test: `internal/service/article_import_http_test.go`

**Interfaces:**
- Produces: `NormalizeCSDNArticleURL(raw string) (string, error)`.
- Produces: `isPublicIP(ip net.IP) bool`.
- Produces: `RemoteFetcher` with `Fetch(ctx context.Context, rawURL string, maxBytes int64) (*FetchedResource, error)`.
- Produces: `NewSafeRemoteFetcher(timeout time.Duration) RemoteFetcher`.

- [ ] **Step 1: Write URL and IP safety tests**

```go
func TestNormalizeCSDNArticleURL(t *testing.T) {
    got, err := NormalizeCSDNArticleURL("https://blog.csdn.net/demo/article/details/123?spm=x#part")
    require.NoError(t, err)
    assert.Equal(t, "https://blog.csdn.net/demo/article/details/123", got)
}

func TestNormalizeCSDNArticleURLRejectsOtherHosts(t *testing.T) {
    _, err := NormalizeCSDNArticleURL("http://127.0.0.1/article/details/123")
    require.Error(t, err)
}

func TestIsPublicIP(t *testing.T) {
    assert.False(t, isPublicIP(net.ParseIP("127.0.0.1")))
    assert.False(t, isPublicIP(net.ParseIP("169.254.1.1")))
    assert.False(t, isPublicIP(net.ParseIP("10.0.0.1")))
    assert.True(t, isPublicIP(net.ParseIP("1.1.1.1")))
}
```

- [ ] **Step 2: Run the focused tests and verify they fail**

Run: `go test ./internal/service -run 'NormalizeCSDN|PublicIP' -v`

Expected: FAIL because the functions do not exist.

- [ ] **Step 3: Implement normalization and a bounded public-network HTTP client**

```go
type FetchedResource struct {
    Body        []byte
    ContentType string
    FinalURL    string
}

type RemoteFetcher interface {
    Fetch(ctx context.Context, rawURL string, maxBytes int64) (*FetchedResource, error)
}
```

Use a custom `http.Transport.DialContext` that resolves the requested host, rejects every non-public resolved IP, and dials a public result. Set `CheckRedirect` to a maximum of three redirects and revalidate each redirect through the same transport. Read through `io.LimitReader(maxBytes+1)` and reject bodies over the limit.

- [ ] **Step 4: Run the focused tests**

Run: `go test ./internal/service -run 'NormalizeCSDN|PublicIP' -v`

Expected: PASS.

- [ ] **Step 5: Commit the task**

```bash
git add internal/service/article_import_http.go internal/service/article_import_http_test.go
git commit -m "feat: add safe CSDN resource fetcher"
```

### Task 2: CSDN HTML parsing and sanitized Markdown conversion

**Files:**
- Create: `internal/service/article_import_parser.go`
- Test: `internal/service/article_import_parser_test.go`

**Interfaces:**
- Consumes: normalized CSDN URL from Task 1.
- Produces: `ParsedImportedArticle { Title, Content, Summary, Cover string; PublishedAt *time.Time }`.
- Produces: `ParseCSDNArticle(page []byte, sourceURL string) (*ParsedImportedArticle, error)`.

- [ ] **Step 1: Write parser fixture tests**

```go
func TestParseCSDNArticle(t *testing.T) {
    page := []byte(`<html><head><meta property="og:title" content="Go 导入测试"><meta name="description" content="摘要"><meta property="article:published_time" content="2025-01-02T03:04:05+08:00"></head><body><div id="content_views"><h2>标题</h2><p>正文 <strong>加粗</strong></p><pre><code class="language-go">fmt.Println("ok")</code></pre><img src="https://i-blog.csdnimg.cn/test.png" onerror="alert(1)"><script>alert(1)</script></div></body></html>`)
    got, err := ParseCSDNArticle(page, "https://blog.csdn.net/demo/article/details/123")
    require.NoError(t, err)
    assert.Equal(t, "Go 导入测试", got.Title)
    assert.Contains(t, got.Content, "## 标题")
    assert.Contains(t, got.Content, "```go")
    assert.NotContains(t, got.Content, "script")
    assert.NotContains(t, got.Content, "onerror")
}
```

Add cases for fallback `h1.title-article`, missing body, tables, ordered/unordered lists, relative image URLs, and `javascript:` links.

- [ ] **Step 2: Run the parser tests and verify they fail**

Run: `go test ./internal/service -run 'ParseCSDNArticle' -v`

Expected: FAIL because the parser does not exist.

- [ ] **Step 3: Implement DOM lookup and Markdown rendering**

Use `html.Parse`, recursive node search helpers, metadata priority lists, and one recursive Markdown renderer. Render only known-safe structural elements; skip `script`, `style`, `iframe`, `button`, `svg`, and `noscript`. Resolve relative links against `sourceURL`, reject non-HTTP(S) image URLs and `javascript:` links, normalize whitespace, and cap summary to 200 runes.

- [ ] **Step 4: Run parser tests**

Run: `go test ./internal/service -run 'ParseCSDNArticle' -v`

Expected: PASS.

- [ ] **Step 5: Commit the task**

```bash
git add go.mod go.sum internal/service/article_import_parser.go internal/service/article_import_parser_test.go
git commit -m "feat: parse CSDN articles into Markdown"
```

### Task 3: Source metadata and duplicate lookup

**Files:**
- Modify: `internal/model/entity/article.go`
- Modify: `internal/model/dto/request/article.go`
- Modify: `internal/model/dto/response/article.go`
- Modify: `internal/repository/article_interface.go`
- Modify: `internal/repository/article_repository.go`
- Modify: `internal/service/article_service.go`
- Modify: `internal/service/article_service_test.go`
- Modify: `pkg/errors/code.go`
- Create: `scripts/migrations/003_add_article_import_source.sql`

**Interfaces:**
- Produces: nullable `Article.SourceURL *string`, `SourcePlatform string`, `SourcePublishedAt *time.Time`.
- Produces: `ArticleRepository.FindBySourceURL(sourceURL string) (*entity.Article, error)`.
- Produces: `CodeArticleSourceExists = 4009`.
- Allows internal callers to set source fields on `CreateArticleRequest` using `json:"-"`.

- [ ] **Step 1: Add failing article-service metadata test**

```go
func TestCreateArticlePersistsSourceMetadata(t *testing.T) {
    source := "https://blog.csdn.net/demo/article/details/123"
    published := time.Date(2025, 1, 2, 3, 4, 5, 0, time.FixedZone("CST", 8*3600))
    req := &request.CreateArticleRequest{Title: "导入", Content: "正文", CategoryID: 1, Status: entity.ArticleStatusDraft, SourceURL: &source, SourcePlatform: "csdn", SourcePublishedAt: &published}
    id, err := svc.CreateArticle(req)
    require.NoError(t, err)
    assert.Equal(t, &source, repo.articles[id].SourceURL)
}
```

Update the existing mock repository with `FindBySourceURL`.

- [ ] **Step 2: Run the focused test and verify it fails**

Run: `go test ./internal/service -run 'CreateArticlePersistsSourceMetadata' -v`

Expected: FAIL because the fields and repository method do not exist.

- [ ] **Step 3: Implement model, repository, DTO, error code, and migration changes**

Use `gorm:"type:varchar(1000);uniqueIndex"` on `*string` source URL. `FindBySourceURL` must return `(nil, nil)` on `gorm.ErrRecordNotFound`. Copy source fields in `CreateArticle`, and return them from `AdminArticleDetailResponse`. Increase create-content binding to `max=500000`.

- [ ] **Step 4: Run article service tests**

Run: `go test ./internal/service -run 'Article|CreateArticle' -v`

Expected: PASS.

- [ ] **Step 5: Commit the task**

```bash
git add internal/model internal/repository internal/service/article_service.go internal/service/article_service_test.go pkg/errors/code.go scripts/migrations/003_add_article_import_source.sql
git commit -m "feat: persist imported article source metadata"
```

### Task 4: Image localization and import orchestration

**Files:**
- Create: `internal/service/article_import_interface.go`
- Create: `internal/service/article_import_service.go`
- Create: `internal/service/article_import_images.go`
- Test: `internal/service/article_import_service_test.go`
- Test: `internal/service/article_import_images_test.go`
- Create: `internal/model/dto/request/article_import.go`
- Create: `internal/model/dto/response/article_import.go`

**Interfaces:**
- Consumes: `RemoteFetcher`, `ParseCSDNArticle`, `ArticleRepository`, `CategoryRepository`, and `ArticleService`.
- Produces: `ArticleImportService.Import(ctx context.Context, req *request.ArticleImportRequest) (*response.ArticleImportResponse, error)`.
- Produces: `ArticleImportResponse { ID uint; Warnings []string }`.

- [ ] **Step 1: Write import-service and image tests**

```go
func TestArticleImportCreatesDraft(t *testing.T) {
    result, err := importer.Import(context.Background(), &request.ArticleImportRequest{URL: "https://blog.csdn.net/demo/article/details/123", CategoryID: 1})
    require.NoError(t, err)
    created := articleRepo.articles[result.ID]
    assert.Equal(t, entity.ArticleStatusDraft, created.Status)
    assert.Equal(t, "csdn", created.SourcePlatform)
}

func TestArticleImportRejectsDuplicateSource(t *testing.T) {
    _, err := importer.Import(context.Background(), &request.ArticleImportRequest{URL: existingURL, CategoryID: 1})
    var bizErr *bizerrors.BizError
    require.ErrorAs(t, err, &bizErr)
    assert.Equal(t, bizerrors.CodeArticleSourceExists, bizErr.Code)
}
```

Add image cases for valid PNG replacement, non-image response, body over 5 MiB, more than 30 images, failed image preserving the remote URL, and cleanup when article creation fails.

- [ ] **Step 2: Run focused tests and verify they fail**

Run: `go test ./internal/service -run 'ArticleImport|LocalizeImportedImages' -v`

Expected: FAIL because the importer and localizer do not exist.

- [ ] **Step 3: Implement image localization**

Extract Markdown image destinations, de-duplicate URLs, fetch through `RemoteFetcher`, validate `image/jpeg`, `image/png`, `image/gif`, or `image/webp`, write unique files beneath `<uploadDir>/article`, and return updated Markdown, warnings, and absolute saved paths. On individual failure, preserve the original destination.

- [ ] **Step 4: Implement import orchestration**

Validate the category exists before network work. Normalize and duplicate-check the source URL, fetch the page with a 5 MiB limit, parse, reject content over 500,000 bytes, localize images, then call `ArticleService.CreateArticle` with `Status: draft` and source fields. Remove saved files if article creation fails.

- [ ] **Step 5: Run focused and full service tests**

Run: `go test ./internal/service -v`

Expected: PASS.

- [ ] **Step 6: Commit the task**

```bash
git add internal/service/article_import* internal/model/dto/request/article_import.go internal/model/dto/response/article_import.go
git commit -m "feat: import CSDN articles as local drafts"
```

### Task 5: Authenticated import API and dependency wiring

**Files:**
- Create: `internal/api/v1/article_import/controller.go`
- Create: `internal/api/v1/article_import/routes.go`
- Test: `internal/api/v1/article_import/controller_test.go`
- Modify: `internal/api/routes.go`
- Modify: `internal/app/app.go`

**Interfaces:**
- Consumes: `ArticleImportService` from Task 4.
- Produces: authenticated `POST /api/v1/admin/article-imports`.

- [ ] **Step 1: Write controller validation tests**

```go
func TestImportRejectsInvalidBody(t *testing.T) {
    req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/article-imports", strings.NewReader(`{"url":"","category_id":0}`))
    req.Header.Set("Content-Type", "application/json")
    recorder := httptest.NewRecorder()
    router.ServeHTTP(recorder, req)
    assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
```

Add success response and service business-error cases. Route-level authentication is verified using the real middleware setup or a focused route test without an Authorization header.

- [ ] **Step 2: Run API tests and verify they fail**

Run: `go test ./internal/api/v1/article_import -v`

Expected: FAIL because the package does not exist.

- [ ] **Step 3: Implement controller and protected route**

Bind `ArticleImportRequest` with required URL and category ID, call `Import(ctx.Request.Context(), &req)`, map business errors with `response.BizError`, and return `response.Success`. Register the route under the existing JWT middleware.

- [ ] **Step 4: Wire the service**

Construct one safe fetcher with a 20-second timeout, construct the import service with `config.App.UploadDir`, pass it into `api.NewRouter`, add `articleImportController`, and register its routes in `Router.Setup`.

- [ ] **Step 5: Run API and full backend tests**

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 6: Commit the task**

```bash
git add internal/api/v1/article_import internal/api/routes.go internal/app/app.go
git commit -m "feat: expose authenticated article import API"
```

### Task 6: Vue admin URL import modal

**Files:**
- Modify: `blog-admin/src/api/article.js`
- Modify: `blog-admin/src/views/article/Edit.vue`

**Interfaces:**
- Consumes: `POST /admin/article-imports` returning `{ code, data: { id, warnings } }`.
- Produces: “URL 导入” modal on the new-article editor.

- [ ] **Step 1: Add the API function and a longer per-request timeout**

```js
export const importArticleFromUrl = (data) => request.post('/admin/article-imports', data, { timeout: 60000 })
```

- [ ] **Step 2: Add modal state and validation**

Add `showImportModal`, `importUrl`, `importCategoryId`, and `importing`. Show the entry button only for a new article. If the current new article has unsaved title/content, confirm replacement before opening. Require an HTTP(S) URL and a selected category.

- [ ] **Step 3: Submit and navigate**

```js
const result = await importArticleFromUrl({
  url: importUrl.value.trim(),
  category_id: Number(importCategoryId.value)
})
const warnings = result.data?.warnings || []
ElMessage.success(warnings.length ? `导入成功，${warnings.length} 项内容需要检查` : '导入成功，已创建草稿')
router.replace(`/articles/edit/${result.data.id}`)
```

Disable controls while importing and rely on the existing Axios interceptor for server error messages.

- [ ] **Step 4: Add responsive styles**

Keep the dialog within `min(520px, calc(100vw - 32px))`, allow long URLs and warning text to wrap, and stack actions below 520px.

- [ ] **Step 5: Build the admin application**

Run: `npm run build`

Working directory: `blog-admin`

Expected: build succeeds with no Vue compiler errors.

- [ ] **Step 6: Commit the task**

```bash
git add blog-admin/src/api/article.js blog-admin/src/views/article/Edit.vue
git commit -m "feat: add CSDN URL import dialog"
```

### Task 7: End-to-end verification and documentation alignment

**Files:**
- Modify only files required to fix failures found by the checks below.

**Interfaces:**
- Consumes: all previous task outputs.
- Produces: verified backend and frontend build artifacts without checking generated `dist` files into Git.

- [ ] **Step 1: Format and test backend**

Run: `gofmt -w internal/api/v1/article_import internal/model internal/repository internal/service pkg/errors`

Run: `go test ./...`

Expected: PASS.

- [ ] **Step 2: Build both frontends**

Run in `blog-admin`: `npm run build`

Run in `blog-web`: `npm run build`

Expected: both builds succeed; size warnings are recorded separately and are not treated as failures.

- [ ] **Step 3: Inspect the final diff and repository status**

Run: `git status --short`

Run: `git diff --check`

Run: `git diff --stat HEAD~1..HEAD`

Expected: no whitespace errors; unrelated pre-existing untracked files remain untouched.

- [ ] **Step 4: Perform a live public-page probe when network access is available**

Use one public CSDN article URL through the authenticated local API or a focused Go test harness. Verify the created row is a draft, content is non-empty, and localized image URLs begin with `/uploads/article/`. If external access is unavailable, report this acceptance item as not runtime-verified rather than claiming it passed.

- [ ] **Step 5: Record the verification boundary**

In the delivery message, list each command that passed, preserve build warnings separately from failures, and state explicitly whether a live CSDN page was imported through the running application. If any verification fix was required, include it in the final implementation commit by staging only its exact file path.
