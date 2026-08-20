# GitHub Actions CI/CD Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build CI for the Go backend, two Vue applications, and production images, then publish immutable GHCR images and deploy them to a Linux Docker Compose server after production approval.

**Architecture:** CI runs on Pull Requests and `main` pushes with read-only permissions. A successful push CI run triggers a privileged workflow that publishes two SHA-tagged images, waits on the protected `production` Environment, then copies a production Compose file and rollback-capable deployment script to the server over pinned-host SSH.

**Tech Stack:** GitHub Actions, Go 1.25, Node.js 22, npm, Docker Buildx, GHCR, POSIX shell, Docker Compose v2, SSH/SCP

**Spec:** `docs/superpowers/specs/2026-08-20-github-actions-cicd-design.md`

## Global Constraints

- Pull Requests never receive package-write permission or production Environment secrets.
- Deployment uses `sha-<full commit SHA>` tags; `latest` is informational only.
- Production deployment requires approval through GitHub Environment `production`.
- Production database, JWT, email, and administrator secrets remain only in `${DEPLOY_PATH}/.env`.
- SSH host identity is supplied by the pre-verified `SSH_KNOWN_HOSTS` secret; do not call `ssh-keyscan` in the workflow.
- Only one production deployment may run at a time, and an in-progress deployment is never cancelled.
- Application rollback changes only `backend` and `nginx`; destructive database migrations remain outside automation.
- Backend startup runs GORM `AutoMigrate`; every automatically deployed schema change must remain backward-compatible with the previous application image.
- Do not modify or stage unrelated business-code changes, `.env*` files, `.deploy/`, image archives, logs, or local build artifacts.
- Use the current official action majors verified on 2026-08-20: `actions/checkout@v7`, `actions/setup-go@v7`, `actions/setup-node@v7`, `docker/setup-buildx-action@v4`, `docker/login-action@v4`, and `docker/build-push-action@v7`.

## File Map

- `.github/workflows/ci.yml`: unprivileged validation for Go, both frontends, deployment shell tests, Compose expansion, and two image builds.
- `.github/workflows/release-deploy.yml`: trusted `main` release, GHCR publishing, protected production deployment, SSH setup, file transfer, and remote execution.
- `deploy/docker-compose.prod.yml`: production services using immutable GHCR image tags and persistent data.
- `deploy/deploy.sh`: remote image pull, service update, bounded smoke checks, version recording, and application rollback.
- `deploy/deploy_test.sh`: isolated fake-Docker/fake-curl tests for successful deployment and failed-smoke rollback.
- `docs/CI-CD部署指南.md`: one-time GitHub/server setup, release operation, verification, and manual rollback.

---

### Task 1: Establish the unprivileged CI workflow

**Files:**
- Modify: `.github/workflows/ci.yml`

**Interfaces:**
- Consumes: `go.mod`, `go.sum`, `blog-web/package-lock.json`, `blog-admin/package-lock.json`, root `Dockerfile`, `docker/nginx/Dockerfile`.
- Produces: a workflow named exactly `CI`, which `release-deploy.yml` references through `workflow_run.workflows: [CI]`.

- [ ] **Step 1: Record the expected CI contract before editing**

Run:

```powershell
$ci = Get-Content .github/workflows/ci.yml -Raw
@('name: CI','permissions:','contents: read','go test -race','blog-web/package-lock.json','blog-admin/package-lock.json','docker/nginx/Dockerfile','deploy/deploy_test.sh') |
  ForEach-Object { if ($ci -notmatch [regex]::Escape($_)) { Write-Error "Missing CI contract: $_" } }
```

Expected: the command fails because the current draft does not run the deployment script tests or build the unified Nginx image.

- [ ] **Step 2: Replace the CI workflow with the complete configuration**

Use this exact structure:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

concurrency:
  group: ci-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

permissions:
  contents: read

jobs:
  backend:
    name: Backend (Go)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache-dependency-path: go.sum
      - name: Verify modules
        run: go mod verify
      - name: Check formatting
        shell: bash
        run: |
          unformatted="$(gofmt -l cmd internal pkg)"
          test -z "$unformatted" || { printf 'Unformatted Go files:\n%s\n' "$unformatted"; exit 1; }
      - name: Vet
        run: go vet ./cmd/... ./internal/... ./pkg/...
      - name: Build
        run: go build ./cmd/... ./internal/... ./pkg/...
      - name: Test
        run: go test -race -timeout 180s -coverprofile=coverage.out ./cmd/... ./internal/... ./pkg/...
      - name: Upload coverage
        if: always()
        uses: actions/upload-artifact@v6
        with:
          name: backend-coverage
          path: coverage.out
          if-no-files-found: ignore

  frontend-web:
    name: Frontend (blog-web)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: blog-web
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-node@v7
        with:
          node-version: '22'
          cache: npm
          cache-dependency-path: blog-web/package-lock.json
      - run: npm ci
      - run: npm run build

  frontend-admin:
    name: Frontend (blog-admin)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: blog-admin
    steps:
      - uses: actions/checkout@v7
      - uses: actions/setup-node@v7
        with:
          node-version: '22'
          cache: npm
          cache-dependency-path: blog-admin/package-lock.json
      - run: npm ci
      - run: npm run build

  deployment-config:
    name: Deployment configuration
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - name: Test deployment script
        run: sh deploy/deploy_test.sh
      - name: Validate production Compose
        env:
          RELEASE_TAG: sha-0000000000000000000000000000000000000000
          DB_ROOT_PASSWORD: ci-root-password
          DB_NAME: blog
          DB_USER: blog_user
          DB_PASSWORD: ci-db-password
          JWT_SECRET: ci-only-secret-with-at-least-32-characters
          BLOGGER_PASSWORD: ci-admin-password
        shell: bash
        run: |
          cp .env.example deploy/.env
          trap 'rm -f deploy/.env' EXIT
          docker compose -f deploy/docker-compose.prod.yml config --quiet

  images:
    name: Production images
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - uses: docker/setup-buildx-action@v4
      - name: Build backend image
        uses: docker/build-push-action@v7
        with:
          context: .
          file: Dockerfile
          push: false
          tags: gin-blog-backend:ci
          cache-from: type=gha,scope=ci-backend
          cache-to: type=gha,mode=max,scope=ci-backend
      - name: Build gateway image
        uses: docker/build-push-action@v7
        with:
          context: .
          file: docker/nginx/Dockerfile
          push: false
          tags: gin-blog-nginx:ci
          cache-from: type=gha,scope=ci-nginx
          cache-to: type=gha,mode=max,scope=ci-nginx
```

- [ ] **Step 3: Re-run the CI contract check**

Run the PowerShell command from Step 1.

Expected: exit code 0 and no `Missing CI contract` errors.

- [ ] **Step 4: Check whitespace and the isolated workflow diff**

Run:

```powershell
git diff --check -- .github/workflows/ci.yml
git diff -- .github/workflows/ci.yml
```

Expected: no whitespace errors; the diff contains only CI workflow changes.

- [ ] **Step 5: Commit the CI workflow**

```powershell
git add -- .github/workflows/ci.yml
git commit -m "ci: validate backend frontend and images"
```

### Task 2: Add production Compose configuration

**Files:**
- Create: `deploy/docker-compose.prod.yml`

**Interfaces:**
- Consumes: `${RELEASE_TAG}`, server `.env`, GHCR images `ghcr.io/fcy222fcy/blog-backend` and `ghcr.io/fcy222fcy/blog-nginx`.
- Produces: Compose services named `mysql`, `redis`, `backend`, and `nginx`, plus fixed container names consumed by `deploy.sh` health checks.

- [ ] **Step 1: Prove the production Compose file is absent**

Run:

```powershell
if (Test-Path deploy/docker-compose.prod.yml) { throw 'Expected production Compose file to be absent before implementation' }
```

Expected: exit code 0.

- [ ] **Step 2: Create `deploy/docker-compose.prod.yml`**

```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: gin-blog-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_ROOT_PASSWORD:?DB_ROOT_PASSWORD is required}
      MYSQL_DATABASE: ${DB_NAME:?DB_NAME is required}
      MYSQL_USER: ${DB_USER:?DB_USER is required}
      MYSQL_PASSWORD: ${DB_PASSWORD:?DB_PASSWORD is required}
      TZ: Asia/Shanghai
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_unicode_ci
      - --default-time-zone=+08:00
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD-SHELL", "mysqladmin ping -h localhost -uroot -p$$MYSQL_ROOT_PASSWORD --silent"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 30s
    networks: [gin-blog]

  redis:
    image: redis:7-alpine
    container_name: gin-blog-redis
    restart: unless-stopped
    command: >-
      sh -c 'if [ -n "$$REDIS_PASSWORD" ]; then
        redis-server --appendonly yes --requirepass "$$REDIS_PASSWORD";
      else
        redis-server --appendonly yes;
      fi'
    environment:
      REDIS_PASSWORD: ${REDIS_PASSWORD:-}
      TZ: Asia/Shanghai
    volumes:
      - redis_data:/data
    healthcheck:
      test: >-
        sh -c 'if [ -n "$$REDIS_PASSWORD" ]; then
          redis-cli -a "$$REDIS_PASSWORD" ping;
        else
          redis-cli ping;
        fi'
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 10s
    networks: [gin-blog]

  backend:
    image: ghcr.io/fcy222fcy/blog-backend:${RELEASE_TAG:?RELEASE_TAG is required}
    container_name: gin-blog-backend
    restart: unless-stopped
    env_file: [.env]
    environment:
      SERVER_PORT: ${SERVER_PORT:-9090}
      DB_HOST: mysql
      REDIS_HOST: redis
      APP_INIT_SQL_DIR: /app/scripts
      APP_UPLOAD_DIR: ${APP_UPLOAD_DIR:-/app/uploads}
      TZ: Asia/Shanghai
    expose: ['9090']
    volumes:
      - ./uploads:/app/uploads
      - ./logs/backend:/app/logs
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks: [gin-blog]

  nginx:
    image: ghcr.io/fcy222fcy/blog-nginx:${RELEASE_TAG:?RELEASE_TAG is required}
    container_name: gin-blog-nginx
    restart: unless-stopped
    ports:
      - "${GATEWAY_PORT:-80}:80"
    depends_on: [backend]
    networks: [gin-blog]

networks:
  gin-blog:
    driver: bridge

volumes:
  mysql_data:
  redis_data:
```

- [ ] **Step 3: Verify missing release tags fail closed**

Run in PowerShell with the server `.env` excluded from output:

```powershell
Copy-Item .env.example deploy/.env
try {
$env:RELEASE_TAG = ''
docker compose --env-file .env.example -f deploy/docker-compose.prod.yml config --quiet
if ($LASTEXITCODE -eq 0) { throw 'Compose unexpectedly accepted an empty RELEASE_TAG' }
} finally {
  Remove-Item deploy/.env -Force
}
```

Expected: non-zero Compose result mentioning `RELEASE_TAG is required`; the PowerShell assertion itself completes successfully.

- [ ] **Step 4: Verify the complete Compose model expands**

Run:

```powershell
Copy-Item .env.example deploy/.env
try {
$env:RELEASE_TAG = 'sha-0000000000000000000000000000000000000000'
docker compose --env-file .env.example -f deploy/docker-compose.prod.yml config --quiet
if ($LASTEXITCODE -ne 0) { throw 'Production Compose validation failed' }
} finally {
  Remove-Item deploy/.env -Force
}
```

Expected: exit code 0.

- [ ] **Step 5: Commit the production Compose file**

```powershell
git add -- deploy/docker-compose.prod.yml
git commit -m "deploy: add production compose configuration"
```

### Task 3: Implement and test rollback-capable deployment

**Files:**
- Create: `deploy/deploy_test.sh`
- Create: `deploy/deploy.sh`

**Interfaces:**
- Consumes: required `RELEASE_TAG`, optional `DEPLOY_PATH`, `HEALTH_ATTEMPTS`, and `HEALTH_INTERVAL`; `${DEPLOY_PATH}/.env`; Task 2 Compose services and container names.
- Produces: `${DEPLOY_PATH}/.release` containing the last successful SHA tag; exit 0 only after service and HTTP checks succeed.

- [ ] **Step 1: Write `deploy/deploy_test.sh` first**

The test creates mock `docker`, `curl`, and `sleep` commands. It must cover a successful release and a failed smoke check that restores `sha-old` while returning failure.

```sh
#!/usr/bin/env sh
set -eu

root_dir=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
sandbox=$(mktemp -d)
trap 'rm -rf "$sandbox"' EXIT HUP INT TERM
mkdir -p "$sandbox/bin" "$sandbox/app"
cp "$root_dir/deploy/deploy.sh" "$sandbox/app/deploy.sh"
cp "$root_dir/deploy/docker-compose.prod.yml" "$sandbox/app/docker-compose.prod.yml"
printf '%s\n' 'DB_ROOT_PASSWORD=test' 'DB_NAME=blog' 'DB_USER=blog' 'DB_PASSWORD=test' > "$sandbox/app/.env"

cat > "$sandbox/bin/docker" <<'MOCK_DOCKER'
#!/usr/bin/env sh
printf 'tag=%s docker %s\n' "${RELEASE_TAG:-unset}" "$*" >> "$MOCK_LOG"
if [ "${1:-}" = inspect ]; then
  printf '%s\n' healthy
fi
if [ "${1:-}" = compose ] && [ "${2:-}" = --env-file ] && [ "${6:-}" = port ]; then
  printf '%s\n' '0.0.0.0:8080'
fi
MOCK_DOCKER

cat > "$sandbox/bin/curl" <<'MOCK_CURL'
#!/usr/bin/env sh
printf 'curl %s\n' "$*" >> "$MOCK_LOG"
[ "${MOCK_CURL_FAIL:-0}" -eq 0 ]
MOCK_CURL

cat > "$sandbox/bin/sleep" <<'MOCK_SLEEP'
#!/usr/bin/env sh
exit 0
MOCK_SLEEP

chmod +x "$sandbox/bin/docker" "$sandbox/bin/curl" "$sandbox/bin/sleep" "$sandbox/app/deploy.sh"
export MOCK_LOG="$sandbox/mock.log"

printf '%s\n' sha-old > "$sandbox/app/.release"
PATH="$sandbox/bin:$PATH" DEPLOY_PATH="$sandbox/app" RELEASE_TAG=sha-new HEALTH_ATTEMPTS=1 HEALTH_INTERVAL=0 "$sandbox/app/deploy.sh"
test "$(cat "$sandbox/app/.release")" = sha-new
grep -q 'tag=sha-new docker compose' "$MOCK_LOG"

: > "$MOCK_LOG"
printf '%s\n' sha-old > "$sandbox/app/.release"
if MOCK_CURL_FAIL=1 PATH="$sandbox/bin:$PATH" DEPLOY_PATH="$sandbox/app" RELEASE_TAG=sha-bad HEALTH_ATTEMPTS=1 HEALTH_INTERVAL=0 "$sandbox/app/deploy.sh"; then
  echo 'failed smoke check unexpectedly returned success' >&2
  exit 1
fi
test "$(cat "$sandbox/app/.release")" = sha-old
grep -q 'tag=sha-old docker compose' "$MOCK_LOG"
printf '%s\n' 'deploy tests passed'
```

- [ ] **Step 2: Run the test and verify it fails before implementation**

Run:

```powershell
sh deploy/deploy_test.sh
```

Expected: non-zero exit because `deploy/deploy.sh` does not exist.

- [ ] **Step 3: Create `deploy/deploy.sh`**

```sh
#!/usr/bin/env sh
set -eu

: "${RELEASE_TAG:?RELEASE_TAG is required}"
DEPLOY_PATH=${DEPLOY_PATH:-$(pwd)}
HEALTH_ATTEMPTS=${HEALTH_ATTEMPTS:-12}
HEALTH_INTERVAL=${HEALTH_INTERVAL:-5}
compose_file="$DEPLOY_PATH/docker-compose.prod.yml"
env_file="$DEPLOY_PATH/.env"
release_file="$DEPLOY_PATH/.release"

cd "$DEPLOY_PATH"
test -f "$compose_file" || { echo "missing $compose_file" >&2; exit 1; }
test -f "$env_file" || { echo "missing $env_file" >&2; exit 1; }

previous_tag=''
if [ -f "$release_file" ]; then
  previous_tag=$(sed -n '1p' "$release_file")
fi

compose() {
  docker compose --env-file "$env_file" -f "$compose_file" "$@"
}

container_ready() {
  status=$(docker inspect --format '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$1" 2>/dev/null || true)
  [ "$status" = healthy ] || [ "$status" = running ]
}

smoke_check() {
  attempt=1
  while [ "$attempt" -le "$HEALTH_ATTEMPTS" ]; do
    if container_ready gin-blog-mysql \
      && container_ready gin-blog-redis \
      && container_ready gin-blog-backend \
      && container_ready gin-blog-nginx; then
      published=$(compose port nginx 80 2>/dev/null | sed -n '1p')
      port=${published##*:}
      if [ -n "$port" ] \
        && curl --fail --silent --show-error --max-time 10 "http://127.0.0.1:$port/" >/dev/null \
        && curl --fail --silent --show-error --max-time 10 "http://127.0.0.1:$port/api/v1/articles" >/dev/null; then
        return 0
      fi
    fi
    attempt=$((attempt + 1))
    sleep "$HEALTH_INTERVAL"
  done
  return 1
}

rollback() {
  if [ -z "$previous_tag" ]; then
    echo 'deployment failed and no previous successful release is recorded' >&2
    return 1
  fi
  echo "deployment failed; restoring $previous_tag" >&2
  RELEASE_TAG=$previous_tag
  export RELEASE_TAG
  compose pull backend nginx || true
  compose up -d backend nginx || true
  smoke_check || echo 'rollback started but health verification also failed' >&2
  return 1
}

export RELEASE_TAG
compose pull backend nginx
if ! compose up -d mysql redis backend nginx; then
  rollback
  exit 1
fi
if ! smoke_check; then
  rollback
  exit 1
fi

printf '%s\n' "$RELEASE_TAG" > "$release_file.tmp"
mv "$release_file.tmp" "$release_file"
docker image prune -f >/dev/null 2>&1 || true
echo "deployment succeeded: $RELEASE_TAG"
```

- [ ] **Step 4: Run syntax and behavior tests**

Run:

```powershell
sh -n deploy/deploy.sh
sh -n deploy/deploy_test.sh
sh deploy/deploy_test.sh
```

Expected: both syntax checks exit 0 and the final line is `deploy tests passed`.

- [ ] **Step 5: Mark scripts executable and commit**

Run:

```powershell
git add -- deploy/deploy.sh deploy/deploy_test.sh
git update-index --chmod=+x deploy/deploy.sh deploy/deploy_test.sh
git commit -m "deploy: add verified rollout and rollback script"
```

### Task 4: Publish images and deploy the exact successful commit

**Files:**
- Create: `.github/workflows/release-deploy.yml`

**Interfaces:**
- Consumes: workflow named `CI`; `workflow_run.head_sha`; Task 2 and Task 3 files; Environment secrets `SSH_HOST`, `SSH_PORT`, `SSH_USER`, `SSH_PRIVATE_KEY`, `SSH_KNOWN_HOSTS`, `GHCR_USERNAME`, `GHCR_TOKEN`; Environment variable `DEPLOY_PATH`.
- Produces: GHCR images tagged `sha-<full SHA>` and `latest`; a protected production deployment invoking `deploy.sh` with the same full SHA tag.

- [ ] **Step 1: Write a failing release-workflow contract check**

Run before creating the file:

```powershell
if (Test-Path .github/workflows/release-deploy.yml) { throw 'Expected release workflow to be absent before implementation' }
```

Expected: exit code 0.

- [ ] **Step 2: Create `.github/workflows/release-deploy.yml`**

```yaml
name: Release and Deploy

on:
  workflow_run:
    workflows: [CI]
    types: [completed]

permissions:
  contents: read

jobs:
  release:
    name: Publish production images
    if: >-
      github.event.workflow_run.conclusion == 'success' &&
      github.event.workflow_run.event == 'push' &&
      github.event.workflow_run.head_branch == 'main'
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    outputs:
      release-tag: ${{ steps.release.outputs.tag }}
    steps:
      - uses: actions/checkout@v7
        with:
          ref: ${{ github.event.workflow_run.head_sha }}
      - id: release
        name: Select immutable release tag
        shell: bash
        run: echo "tag=sha-${{ github.event.workflow_run.head_sha }}" >> "$GITHUB_OUTPUT"
      - uses: docker/setup-buildx-action@v4
      - uses: docker/login-action@v4
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Publish backend
        uses: docker/build-push-action@v7
        with:
          context: .
          file: Dockerfile
          push: true
          tags: |
            ghcr.io/fcy222fcy/blog-backend:${{ steps.release.outputs.tag }}
            ghcr.io/fcy222fcy/blog-backend:latest
          labels: org.opencontainers.image.source=https://github.com/${{ github.repository }}
          cache-from: type=gha,scope=release-backend
          cache-to: type=gha,mode=max,scope=release-backend
      - name: Publish gateway
        uses: docker/build-push-action@v7
        with:
          context: .
          file: docker/nginx/Dockerfile
          push: true
          tags: |
            ghcr.io/fcy222fcy/blog-nginx:${{ steps.release.outputs.tag }}
            ghcr.io/fcy222fcy/blog-nginx:latest
          labels: org.opencontainers.image.source=https://github.com/${{ github.repository }}
          cache-from: type=gha,scope=release-nginx
          cache-to: type=gha,mode=max,scope=release-nginx

  deploy:
    name: Deploy production
    needs: release
    runs-on: ubuntu-latest
    concurrency:
      group: production
      cancel-in-progress: false
    environment:
      name: production
    steps:
      - uses: actions/checkout@v7
        with:
          ref: ${{ github.event.workflow_run.head_sha }}
      - name: Validate deployment settings
        shell: bash
        env:
          SSH_HOST: ${{ secrets.SSH_HOST }}
          SSH_PORT: ${{ secrets.SSH_PORT }}
          SSH_USER: ${{ secrets.SSH_USER }}
          DEPLOY_PATH: ${{ vars.DEPLOY_PATH }}
        run: |
          [[ "$SSH_HOST" =~ ^[A-Za-z0-9._:-]+$ ]]
          [[ "$SSH_PORT" =~ ^[0-9]+$ ]]
          [[ "$SSH_USER" =~ ^[A-Za-z0-9._-]+$ ]]
          [[ "$DEPLOY_PATH" =~ ^/[A-Za-z0-9._/-]+$ ]]
      - name: Configure pinned SSH host
        shell: bash
        env:
          SSH_PRIVATE_KEY: ${{ secrets.SSH_PRIVATE_KEY }}
          SSH_KNOWN_HOSTS: ${{ secrets.SSH_KNOWN_HOSTS }}
        run: |
          install -d -m 700 "$HOME/.ssh"
          install -m 600 /dev/null "$HOME/.ssh/id_ed25519"
          printf '%s\n' "$SSH_PRIVATE_KEY" > "$HOME/.ssh/id_ed25519"
          printf '%s\n' "$SSH_KNOWN_HOSTS" > "$HOME/.ssh/known_hosts"
      - name: Copy deployment files
        shell: bash
        env:
          SSH_HOST: ${{ secrets.SSH_HOST }}
          SSH_PORT: ${{ secrets.SSH_PORT }}
          SSH_USER: ${{ secrets.SSH_USER }}
          DEPLOY_PATH: ${{ vars.DEPLOY_PATH }}
        run: |
          target="$SSH_USER@$SSH_HOST"
          ssh -p "$SSH_PORT" "$target" "mkdir -p '$DEPLOY_PATH'"
          scp -P "$SSH_PORT" deploy/docker-compose.prod.yml deploy/deploy.sh "$target:$DEPLOY_PATH/"
      - name: Login server to GHCR and deploy
        shell: bash
        env:
          SSH_HOST: ${{ secrets.SSH_HOST }}
          SSH_PORT: ${{ secrets.SSH_PORT }}
          SSH_USER: ${{ secrets.SSH_USER }}
          DEPLOY_PATH: ${{ vars.DEPLOY_PATH }}
          GHCR_USERNAME: ${{ secrets.GHCR_USERNAME }}
          GHCR_TOKEN: ${{ secrets.GHCR_TOKEN }}
          RELEASE_TAG: ${{ needs.release.outputs.release-tag }}
        run: |
          target="$SSH_USER@$SSH_HOST"
          printf '%s' "$GHCR_TOKEN" | ssh -p "$SSH_PORT" "$target" "docker login ghcr.io -u '$GHCR_USERNAME' --password-stdin"
          ssh -p "$SSH_PORT" "$target" "chmod 700 '$DEPLOY_PATH/deploy.sh' && DEPLOY_PATH='$DEPLOY_PATH' RELEASE_TAG='$RELEASE_TAG' '$DEPLOY_PATH/deploy.sh'"
```

- [ ] **Step 3: Verify privileged-trigger filters and exact-SHA wiring**

Run:

```powershell
$release = Get-Content .github/workflows/release-deploy.yml -Raw
@("workflow_run.event == 'push'","workflow_run.head_branch == 'main'",'workflow_run.head_sha','environment:','name: production','cancel-in-progress: false','SSH_KNOWN_HOSTS') |
  ForEach-Object {
    if ($release -notmatch [regex]::Escape($_)) { Write-Error "Missing release contract: $_" }
  }
if ($release -match 'ssh-keyscan') { Write-Error 'Runtime host-key trust is forbidden' }
```

Expected: exit code 0, no missing-contract error, and no `ssh-keyscan` usage.

- [ ] **Step 4: Review permissions and workflow diff**

Run:

```powershell
git diff --check -- .github/workflows/release-deploy.yml
git diff -- .github/workflows/release-deploy.yml
```

Expected: package write permission appears only on the `release` Job; the `deploy` Job alone references the `production` Environment and its secrets.

- [ ] **Step 5: Commit the release workflow**

```powershell
git add -- .github/workflows/release-deploy.yml
git commit -m "ci: publish and deploy production images"
```

### Task 5: Document setup, operation, rollback, and proof boundaries

**Files:**
- Create: `docs/CI-CD部署指南.md`
- Modify: `docs/superpowers/specs/2026-08-20-github-actions-cicd-design.md` (the already-made `SSH_KNOWN_HOSTS` security clarification only)

**Interfaces:**
- Consumes: exact secret/variable names and behavior from Tasks 1–4.
- Produces: a user-facing runbook that can configure GitHub and the Linux server without reading workflow internals.

- [ ] **Step 1: Write a documentation contract check**

Run before creating the guide:

```powershell
if (Test-Path 'docs/CI-CD部署指南.md') { throw 'Expected deployment guide to be absent before implementation' }
```

Expected: exit code 0.

- [ ] **Step 2: Create the guide with these exact sections and values**

Write `docs/CI-CD部署指南.md` with:

```markdown
# CI/CD 部署指南

## 流水线行为

- PR/main：运行 Go 检查、双前端构建、部署脚本测试、Compose 校验和两个镜像构建。
- main CI 成功：发布 `sha-<完整提交 SHA>` 与 `latest` 镜像到 GHCR。
- production 批准：服务器部署 SHA 标签，并在冒烟失败时恢复上一个成功标签。

## 一次性服务器准备

1. 安装 Docker Engine 和 Docker Compose v2。
2. 创建专用部署账号，并允许它使用 Docker。
3. 创建部署目录，例如 `/opt/gin-blog`，并让部署账号拥有该目录。
4. 将 `.env.example` 复制为 `/opt/gin-blog/.env`，替换所有生产密码和密钥。
5. 确保网关端口和 SSH 端口符合服务器防火墙策略。

## GitHub production Environment

创建名为 `production` 的 Environment，限制为 `main`，并设置 Required reviewers。

Environment Secrets：

- `SSH_HOST`
- `SSH_PORT`
- `SSH_USER`
- `SSH_PRIVATE_KEY`
- `SSH_KNOWN_HOSTS`（从可信管理通道核验并复制服务器公钥记录，不能在流水线中临时扫描）
- `GHCR_USERNAME`
- `GHCR_TOKEN`（classic PAT，仅 `read:packages`）

Environment Variable：

- `DEPLOY_PATH=/opt/gin-blog`

## 首次发布

合并到 `main` 后等待 CI 与镜像发布完成，在 Actions 中批准 `production`。首次部署没有旧版本可自动恢复，批准前应确认服务器 `.env`、目录权限、Docker 与端口均已准备。

## 验证

- Actions 的 `Release and Deploy` 工作流成功。
- `docker compose --env-file .env -f docker-compose.prod.yml ps` 显示四个服务运行。
- 首页和 `/api/v1/articles` 返回 2xx。
- `.release` 内容等于本次 `sha-<完整提交 SHA>`。

## 手动回滚

在服务器部署目录执行：

```sh
previous=sha-填写要恢复的完整提交SHA
RELEASE_TAG="$previous" docker compose --env-file .env -f docker-compose.prod.yml pull backend nginx
RELEASE_TAG="$previous" docker compose --env-file .env -f docker-compose.prod.yml up -d backend nginx
printf '%s\n' "$previous" > .release
```

应用镜像回滚不会逆转数据库。包含破坏性数据库变更的版本必须先备份，并使用单独迁移方案。

## 证据边界

本地测试和镜像构建成功只证明配置与产物可生成；只有 GitHub Environment 获批、真实服务器工作流成功、线上 HTTP 检查通过，才能说明生产部署完成。
```

- [ ] **Step 3: Verify documentation matches workflow names**

Run:

```powershell
$guide = Get-Content 'docs/CI-CD部署指南.md' -Raw
@('SSH_KNOWN_HOSTS','GHCR_TOKEN','DEPLOY_PATH','sha-<完整提交 SHA>','/api/v1/articles','read:packages','数据库') |
  ForEach-Object { if ($guide -notmatch [regex]::Escape($_)) { Write-Error "Guide missing: $_" } }
```

Expected: exit code 0.

- [ ] **Step 4: Run the full local verification set**

Use workspace-local Go caches on Windows so verification does not write to the restricted user cache:

```powershell
$env:GOCACHE = Join-Path (Get-Location) '.gocache'
$env:GOMODCACHE = Join-Path (Get-Location) '.gomodcache'
$env:GOPATH = Join-Path (Get-Location) '.gopath'
go mod verify
$unformatted = gofmt -l cmd internal pkg
if ($unformatted) { $unformatted; throw 'gofmt check failed' }
go vet ./cmd/... ./internal/... ./pkg/...
go build ./cmd/... ./internal/... ./pkg/...
go test -race -timeout 180s ./cmd/... ./internal/... ./pkg/...
sh -n deploy/deploy.sh
sh -n deploy/deploy_test.sh
sh deploy/deploy_test.sh
$env:RELEASE_TAG = 'sha-0000000000000000000000000000000000000000'
Copy-Item .env.example deploy/.env
try {
  docker compose --env-file .env.example -f deploy/docker-compose.prod.yml config --quiet
  if ($LASTEXITCODE -ne 0) { throw 'Production Compose validation failed' }
} finally {
  Remove-Item deploy/.env -Force
}
npm ci --prefix blog-web
npm run build --prefix blog-web
npm ci --prefix blog-admin
npm run build --prefix blog-admin
docker build -t gin-blog-backend:verify -f Dockerfile .
docker build -t gin-blog-nginx:verify -f docker/nginx/Dockerfile .
git diff --check
```

Expected: every command exits 0. Treat frontend chunk-size warnings as warnings, not build failures. If Docker is unavailable locally, report both image builds as unverified instead of claiming they pass.

- [ ] **Step 5: Check scope before the final commit**

Run:

```powershell
git status --short
git diff --name-only HEAD
```

Expected for this task: only `docs/CI-CD部署指南.md` and the approved `SSH_KNOWN_HOSTS` clarification remain. Existing unrelated dirty files may still appear in `git status`, but must not appear in the staged set.

- [ ] **Step 6: Commit documentation only**

```powershell
git add -- docs/CI-CD部署指南.md docs/superpowers/specs/2026-08-20-github-actions-cicd-design.md
git commit -m "docs: add CI/CD deployment runbook"
```

- [ ] **Step 7: Report local proof separately from live deployment**

Final handoff must list:

- Files and commits created.
- Exact local checks that exited 0, with warnings separated from failures.
- GitHub Environment Secrets/Variable the user must configure.
- Remaining live acceptance: merge/push workflow run, production approval, server pull, container health, and real HTTP checks.
- No claim that production is deployed until the live workflow evidence exists.
