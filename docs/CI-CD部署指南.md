# CI/CD 部署指南

## 流水线行为

- Pull Request 和 `main`：运行 Go 检查、双前端构建、部署脚本测试、生产 Compose 校验和两个生产镜像构建。
- `main` CI 成功：发布 `sha-<完整提交 SHA>` 与 `latest` 镜像到 GHCR。
- `production` 批准：服务器部署不可变 SHA 标签，并在冒烟失败时恢复上一个成功标签。

Pull Request 不具备包写入权限，也无法读取生产 Environment 的密钥。

## 一次性服务器准备

1. 安装 Docker Engine 和 Docker Compose v2。
2. 创建专用部署账号，并允许它使用 Docker。
3. 创建部署目录，例如 `/opt/gin-blog`，并让部署账号拥有该目录。
4. 将项目的 `.env.example` 复制为 `/opt/gin-blog/.env`，替换所有生产密码和密钥。
5. 确保网关端口和 SSH 端口符合服务器防火墙策略。

示例：

```sh
sudo install -d -o deploy -g deploy -m 750 /opt/gin-blog
sudo install -o deploy -g deploy -m 600 .env.example /opt/gin-blog/.env
sudo -u deploy docker compose version
```

生产 `.env` 只保存在服务器。不要把它复制回仓库或保存为 GitHub Secret。

## GitHub production Environment

在仓库 Settings → Environments 中创建名为 `production` 的 Environment：

1. 将 Deployment branches 限制为 `main`。
2. 设置 Required reviewers；如团队条件允许，同时禁止发起者自行批准。
3. 配置以下 Environment Secrets：

   - `SSH_HOST`
   - `SSH_PORT`
   - `SSH_USER`
   - `SSH_PRIVATE_KEY`
   - `SSH_KNOWN_HOSTS`
   - `GHCR_USERNAME`
   - `GHCR_TOKEN`

4. 配置 Environment Variable：`DEPLOY_PATH=/opt/gin-blog`。

`GHCR_TOKEN` 使用 classic PAT，并且只授予 `read:packages`。如果 GHCR 包继承私有仓库权限，令牌所属账号还必须对仓库具有读取权限。

`SSH_KNOWN_HOSTS` 必须从可信管理通道核验服务器主机公钥后生成。可以在受信任的管理员电脑上取得候选记录，与服务器控制台显示的 Ed25519 指纹核对无误后再保存；不要在 Actions 工作流中临时运行 `ssh-keyscan` 并直接信任结果。

## 首次发布

合并到 `main` 后：

1. 等待 `CI` 工作流通过。
2. 等待 `Release and Deploy` 发布两个 GHCR 镜像。
3. 在 Actions 中检查目标提交 SHA，然后批准 `production`。
4. 等待部署任务完成并检查线上页面。

首次部署没有旧版本可自动恢复。批准前必须确认服务器 `.env`、目录权限、Docker、GHCR 读取令牌和端口均已准备。

## 验证

在 GitHub 确认 `Release and Deploy` 工作流成功。在服务器部署目录运行：

```sh
cd /opt/gin-blog
docker compose --env-file .env -f docker-compose.prod.yml ps
cat .release
curl --fail --show-error http://127.0.0.1/
curl --fail --show-error http://127.0.0.1/api/v1/articles
```

验收标准：

- MySQL、Redis、后端和 Nginx 四个服务处于运行状态。
- 首页及 `/api/v1/articles` 返回 HTTP 2xx。
- `.release` 内容等于本次 `sha-<完整提交 SHA>`。

若 `GATEWAY_PORT` 不是 `80`，在两个 `curl` 地址中使用实际端口。

## 手动回滚

自动回滚失败或需要主动恢复历史版本时，在服务器部署目录执行：

```sh
cd /opt/gin-blog
previous=sha-填写要恢复的完整提交SHA
RELEASE_TAG="$previous" docker compose --env-file .env -f docker-compose.prod.yml pull backend nginx
RELEASE_TAG="$previous" docker compose --env-file .env -f docker-compose.prod.yml up -d backend nginx
printf '%s\n' "$previous" > .release
```

随后重新检查容器状态、首页和文章 API。应用镜像回滚不会逆转数据库；包含删除列、重命名列或不可逆数据变换的版本必须先备份，并使用单独迁移方案。

## 常见失败边界

- CI 失败：不会发布镜像，也不会请求生产批准。
- GHCR 发布失败：不会运行部署任务。
- SSH、登录或拉取失败：工作流失败；拉取失败不会切换正在运行的容器。
- 新版本冒烟失败：脚本尝试恢复 `.release` 记录的上一个应用镜像。
- 首次部署失败：由于没有历史 `.release`，无法自动回滚，工作流会明确失败。

## 证据边界

本地测试和镜像构建成功只证明配置与产物可以生成。只有 GitHub Environment 获批、真实服务器工作流成功、容器状态正常且线上 HTTP 检查通过，才能说明生产部署完成。
