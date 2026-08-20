# GitHub Actions CI/CD 设计

## 1. 目标

为 Gin 博客项目建立一条可追溯、可人工把关、可回滚应用版本的 CI/CD 流水线：

- Pull Request 和 `main` 分支推送必须通过后端、双前端和容器构建检查。
- `main` 通过 CI 后，将后端与统一 Nginx 网关镜像发布到 GitHub Container Registry（GHCR）。
- 生产部署进入 GitHub `production` Environment，等待 Required reviewers 人工批准。
- 批准后通过 SSH 在 Linux 服务器使用 Docker Compose 拉取指定提交版本并更新服务。
- 业务密钥只保存在生产服务器，Pull Request 不接触发布权限或生产密钥。

## 2. 当前项目边界

项目由以下可部署单元组成：

- Go 1.25 后端：根目录 `Dockerfile`，运行 `cmd/server`。
- Vue 3 前台：`blog-web`，由 `docker/nginx/Dockerfile` 构建。
- Vue 3 后台：`blog-admin`，由 `docker/nginx/Dockerfile` 构建并挂载到 `/admin/`。
- Nginx 统一网关：提供两个前端，并代理 `/api/` 和 `/uploads/` 到后端。
- MySQL 8、Redis 7：由生产 Compose 管理并使用命名卷持久化。

本次不修改后端业务逻辑、前端页面、数据库表结构或域名/TLS 终止方式。TLS 可以由服务器现有反向代理或后续独立任务处理。

## 3. 总体架构

### 3.1 CI

`.github/workflows/ci.yml` 在以下事件运行：

- 向 `main` 发起或更新 Pull Request。
- 向 `main` 推送提交。

CI 使用最小 `contents: read` 权限，按独立 Job 并行执行：

1. 后端：检查 `gofmt`、执行 `go vet`、编译并运行 Go 测试。
2. 前台：`npm ci` 后执行 `npm run build`。
3. 后台：`npm ci` 后执行 `npm run build`。
4. 容器：分别构建根目录后端镜像和 `docker/nginx/Dockerfile` 统一网关镜像，但不推送。

同一分支或 Pull Request 的新提交取消旧 CI，避免浪费执行时间。

### 3.2 镜像发布

独立发布部署工作流只在 `main` 对应的 CI 成功后运行。发布 Job：

- 使用 `packages: write` 和仓库自带的短期 `GITHUB_TOKEN` 登录 GHCR。
- 构建并推送两个镜像：
  - `ghcr.io/fcy222fcy/blog-backend`
  - `ghcr.io/fcy222fcy/blog-nginx`
- 每个镜像写入两个标签：不可变的 `sha-<完整提交 SHA>` 与便于查看的 `latest`。
- 部署始终使用不可变 SHA 标签，不使用 `latest`，避免标签漂移。
- Release 与 deploy Job 显式使用触发 CI 的 `head_sha` 检出和标记代码，避免工作流排队期间 `main` 前进后部署错版本。

### 3.3 生产部署

Deploy Job 依赖镜像发布成功，并绑定 GitHub `production` Environment。仓库管理员在 GitHub 设置 Required reviewers，从而实现人工批准。

服务器不保留完整源码，只保留：

- 生产 `docker-compose.prod.yml`。
- 部署脚本。
- 只存在于服务器的 `.env`。
- 当前成功版本记录文件。
- 上传、日志与 Docker 命名卷数据。

Actions 将本次 Compose 文件和部署脚本复制到 `DEPLOY_PATH`，然后通过 SSH 以 `RELEASE_TAG=sha-<完整提交 SHA>` 执行部署。服务器使用只具备 GHCR `read:packages` 权限的令牌登录并拉取镜像。

## 4. 配置与密钥

### 4.1 GitHub Environment Secrets

以下值配置在 GitHub `production` Environment，不写入仓库：

- `SSH_HOST`：生产服务器地址。
- `SSH_PORT`：SSH 端口。
- `SSH_USER`：专用部署账号。
- `SSH_PRIVATE_KEY`：对应部署账号的私钥。
- `SSH_KNOWN_HOSTS`：预先核验的生产服务器 SSH 主机公钥记录，用于固定服务器身份；工作流不得使用运行时 `ssh-keyscan` 结果替代它。
- `GHCR_USERNAME`：服务器拉取镜像所用 GitHub 用户名。
- `GHCR_TOKEN`：仅授予 `read:packages` 的令牌。

### 4.2 GitHub Environment Variable

- `DEPLOY_PATH`：部署账号可写的绝对目录；文档示例为 `/opt/gin-blog`。

### 4.3 服务器 `.env`

数据库密码、JWT 密钥、邮件凭据、管理员初始密码和应用配置只保存在 `${DEPLOY_PATH}/.env`。工作流不得打印、覆盖或上传该文件。

## 5. 部署数据流

1. 开发者提交 Pull Request，CI 并行验证后端、前端和容器。
2. 合并到 `main` 后再次运行 CI。
3. CI 成功事件触发发布部署工作流，并锁定该次 CI 的提交 SHA。
4. 发布 Job 将两个 SHA 标签镜像推送到 GHCR。
5. Deploy Job 等待 `production` Environment 审批。
6. 审批后，Actions 上传部署配置并通过 SSH 调用部署脚本。
7. 脚本保存当前版本，拉取新镜像，以新 SHA 标签更新 `backend` 和 `nginx`。
8. 脚本验证容器状态、网关首页和 `/api/v1/articles`。
9. 验证成功后记录新版本并清理未使用的旧镜像；验证失败则恢复上一个镜像标签并重新启动。

## 6. 失败与回滚

- CI 任一 Job 失败：不发布镜像，不进入生产部署。
- 任一镜像构建或推送失败：Deploy Job 不运行。
- SSH、GHCR 登录或镜像拉取失败：现有容器保持当前版本，工作流失败。
- 新容器启动或冒烟检查失败：部署脚本使用先前成功标签恢复 `backend` 与 `nginx`，再次执行 Compose 更新，并让工作流以失败结束以保留告警信号。
- 首次部署没有先前成功版本时，自动回滚不可用；脚本必须明确报告该状态并保持失败，不能伪报成功。
- 同一时间只允许一个 production 部署，后续部署排队，不能取消正在进行的生产部署。

应用回滚不等于数据库回滚。后端启动会执行 GORM `AutoMigrate`，因此数据库变更必须保持向后兼容。删除列、重命名列、不可逆数据变换等破坏性迁移不纳入自动部署，必须单独备份、评审和人工执行。

## 7. 健康检查与成功标准

部署脚本在有限重试窗口内检查：

1. `docker compose ps` 中 MySQL、Redis、后端与 Nginx 没有退出或不健康状态。
2. `http://127.0.0.1:${GATEWAY_PORT}/` 返回 HTTP 2xx 或 3xx。
3. `http://127.0.0.1:${GATEWAY_PORT}/api/v1/articles` 返回 HTTP 2xx。

CI/CD 配置完成的本地证据包括：

- GitHub Actions YAML 可解析。
- 生产 Compose 使用示例变量可以成功展开配置。
- Go 的格式、静态检查、编译和测试按 CI 同等命令执行。
- 两个前端均以锁文件安装依赖并成功构建。
- 后端与统一网关两个生产镜像均成功构建。

真正的“已上线”必须以 GitHub production 审批、真实服务器工作流成功和线上 HTTP 检查为证。仅本地配置和构建通过不能表述为生产部署已验证。

## 8. 服务器前置条件

- Linux 服务器已安装 Docker Engine 与 Docker Compose v2 插件。
- 部署账号可使用 Docker，并对 `DEPLOY_PATH`、上传目录和日志目录具有读写权限。
- 防火墙已开放实际网关端口；SSH 仅开放给需要的来源。
- `${DEPLOY_PATH}/.env` 已依据 `.env.example` 设置生产值，敏感默认值均已替换。
- GitHub `production` Environment 已启用 Required reviewers。

## 9. 文件变更范围

- 修改 `.github/workflows/ci.yml`：完善 CI 检查与两个镜像构建验证。
- 新增 `.github/workflows/release-deploy.yml`：发布 GHCR 镜像并经审批部署。
- 新增 `deploy/docker-compose.prod.yml`：生产环境仅拉取镜像的编排。
- 新增 `deploy/deploy.sh`：版本切换、检查与应用回滚。
- 新增 `docs/CI-CD部署指南.md`：GitHub 与服务器的一次性设置、发布和回滚操作。
- 修改 `.gitignore`（仅在需要时）：确保可提交部署模板，同时继续忽略服务器密钥和本地部署产物。

现有业务代码、`.env`、`.env.production`、`.deploy/` 本地产物及其他未提交改动均不在本次范围内。
