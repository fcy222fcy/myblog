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

# 取响应状态码，连不上时返回空串。
# --insecure 是有意的：这里只验服务可用性，不让证书续期状态阻塞部署。
http_status() {
  curl --insecure --silent --output /dev/null --write-out '%{http_code}' --max-time 10 "$1" 2>/dev/null || true
}

# 冒烟必须走 443 并显式断言 200：
#   80 端口对一切请求无条件 return 301 到 https（Host: 127.0.0.1 也一样），
#   而 curl --fail 只在状态码 >=400 时才失败，301 会被当成成功 —— 那样两个
#   探针实际只证明了 80 在监听，后端起不来也照样放行。
#   断言 200 还能顺带拦住 404/502 这类"能连上但服务不对"的情况。
smoke_check() {
  attempt=1
  while [ "$attempt" -le "$HEALTH_ATTEMPTS" ]; do
    if container_ready gin-blog-mysql \
      && container_ready gin-blog-redis \
      && container_ready gin-blog-backend \
      && container_ready gin-blog-nginx; then
      published=$(compose port nginx 443 2>/dev/null | sed -n '1p')
      port=${published##*:}
      if [ -n "$port" ] \
        && [ "$(http_status "https://127.0.0.1:$port/")" = '200' ] \
        && [ "$(http_status "https://127.0.0.1:$port/api/v1/articles")" = '200' ]; then
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
