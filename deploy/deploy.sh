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
