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
