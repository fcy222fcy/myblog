<p align="center">
  <h1 align="center">✍️ Gin Blog — 自研博客系统</h1>
</p>

<p align="center">
  <a href="https://github.com/fcy222fcy/myblog"><img src="https://img.shields.io/github/stars/fcy222fcy/myblog?style=flat-square" alt="Stars"></a>
  <a href="https://github.com/fcy222fcy/myblog/forks"><img src="https://img.shields.io/github/forks/fcy222fcy/myblog?style=flat-square" alt="Forks"></a>
  <a href="https://github.com/fcy222fcy/myblog/issues"><img src="https://img.shields.io/github/issues/fcy222fcy/myblog?style=flat-square" alt="Issues"></a>
  <a href="https://github.com/fcy222fcy/myblog/actions"><img src="https://img.shields.io/github/actions/workflow/status/fcy222fcy/myblog/ci.yml?style=flat-square&logo=githubactions&label=CI" alt="CI"></a>
  <a href="https://github.com/fcy222fcy/myblog/blob/main/LICENSE"><img src="https://img.shields.io/github/license/fcy222fcy/myblog?style=flat-square" alt="License"></a>
  <img src="https://img.shields.io/github/languages/top/fcy222fcy/myblog?style=flat-square" alt="Language">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square" alt="PRs Welcome">
</p>

<p align="center">
  🌐 在线体验：<a href="https://blog.fcyan.xyz">blog.fcyan.xyz</a>
  &nbsp;·&nbsp;
  🚀 从零自研的 <strong>Go + Gin + GORM + Vue 3</strong> 现代化博客系统
</p>

<p align="center">
  前后端分离 · 双前端架构 · 每日一问 · 评论系统 · CI/CD 自动发布
</p>

---

## 📌 项目简介

一个**完全从零自研**的个人博客系统：Go 后端 + Vue 3 双前端（前台展示 / 后台管理），Docker 编排部署，GitHub Actions 多分支 CI/CD 自动发布。支持文章管理、每日一问、评论互动、全文搜索、RSS 订阅等完整博客能力，并具备从 CSDN 等平台**批量迁移历史文章**的实用工具链。

> 本项目是个人技术博客的工程化实践，沉淀了从架构设计、分层实现到自动化部署的全栈经验，欢迎 Star、Fork 与 PR。

---

## ✨ 功能特性

### 📝 内容管理
- **文章**：Markdown 写作、草稿/定时发布、分类/标签、SEO 字段、阅读时长与浏览量统计
- **CSDN 迁移工具链**：HTML → Markdown 转换、图片批量本地化、保留原始发布时间与标签
- **每日一问**：日历排期、未来日期预排、答案渐变遮罩、浏览计数

### 💬 互动体系
- **评论系统**：游客 / 登录用户 / 博主三级身份，多级回复、热度排序、点赞防刷
- **智能头像**：QQ 邮箱 → Gravatar → 首字母 SVG 三级兜底
- **博主标识**：强制使用服务端配置值，前端不可伪造

### 🔐 安全与可观测
- **认证**：JWT 令牌 + bcrypt 密码加密
- **审计**：全量操作审计日志，请求体敏感字段自动脱敏
- **防护**：评论点赞 IP 唯一约束防刷、UA 解析

### 🛠️ 工程化
- **架构**：Controller → Service → Repository 三层分层
- **CI/CD**：push main 自动构建 → 发布 GHCR 镜像 → 生产环境人工审批部署 → 失败自动回滚
- **可观测**：Zap 结构化日志、请求追踪

### 📡 其他
- 全文搜索（正文命中高亮）、时间归档、RSS + Sitemap
- 深浅色主题、响应式布局（桌面 / 平板 / 手机）
- 站点运行时长、访问统计

---

## 🏗️ 技术栈

### 后端

<p>
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Gin-1.12-00ADD8?style=flat-square" alt="Gin">
  <img src="https://img.shields.io/badge/GORM-1.31-20B2AA?style=flat-square" alt="GORM">
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat-square&logo=mysql&logoColor=white" alt="MySQL">
  <img src="https://img.shields.io/badge/Redis-7-DC382D?style=flat-square&logo=redis&logoColor=white" alt="Redis">
  <img src="https://img.shields.io/badge/JWT-v5-000000?style=flat-square&logo=json-web-tokens&logoColor=white" alt="JWT">
  <img src="https://img.shields.io/badge/Zap-1.27-00ADD8?style=flat-square" alt="Zap">
  <img src="https://img.shields.io/badge/Viper-1.19-00ADD8?style=flat-square" alt="Viper">
</p>

### 前台前端 (blog-web)

<p>
  <img src="https://img.shields.io/badge/Vue-3.4-4FC08D?style=flat-square&logo=vue.js&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/Vite-5.2-646CFF?style=flat-square&logo=vite&logoColor=white" alt="Vite">
  <img src="https://img.shields.io/badge/Pinia-2.1-F7D336?style=flat-square" alt="Pinia">
  <img src="https://img.shields.io/badge/Vue_Router-4.3-4FC08D?style=flat-square&logo=vue.js&logoColor=white" alt="Vue Router">
  <img src="https://img.shields.io/badge/TailwindCSS-3.4-38B2AC?style=flat-square&logo=tailwind-css&logoColor=white" alt="TailwindCSS">
  <img src="https://img.shields.io/badge/Axios-1.7-5A29E4?style=flat-square&logo=axios&logoColor=white" alt="Axios">
</p>

### 后台前端 (blog-admin)

<p>
  <img src="https://img.shields.io/badge/Vue-3.4-4FC08D?style=flat-square&logo=vue.js&logoColor=white" alt="Vue">
  <img src="https://img.shields.io/badge/Element_Plus-2.7-409EFF?style=flat-square&logo=element&logoColor=white" alt="Element Plus">
  <img src="https://img.shields.io/badge/Markdown_Editor-md--editor--v3-519ABA?style=flat-square&logo=markdown&logoColor=white" alt="Markdown Editor">
  <img src="https://img.shields.io/badge/Pinia-2.1-F7D336?style=flat-square" alt="Pinia">
</p>

---

## 🖼️ 界面预览

| 前台首页 | 文章详情 |
|---|---|
| 每日一问 + 文章卡片列表，支持深浅色主题 | Markdown 渲染 + 评论互动 |

| 后台管理 | 文章编辑 |
|---|---|
| 数据看板 + 全模块管理 | Markdown 写作 + 分类标签 |

> 💡 更多细节可直接访问线上示例：<https://blog.fcyan.xyz>

---

## 📁 项目结构

```
.
├── 📦 blog-web/               # 前台展示前端 (Vue 3 + TailwindCSS)
├── 📦 blog-admin/             # 后台管理前端 (Vue 3 + Element Plus)
├── ⚙️ internal/               # Go 后端核心代码
│   ├── api/v1/             # API 控制器 & 路由（按模块划分）
│   ├── middleware/         # 中间件（认证/审计/日志/恢复）
│   ├── model/              # 数据模型（Entity / DTO）
│   ├── repository/         # 数据访问层
│   ├── service/            # 业务逻辑层
│   └── router/             # 路由注册
├── 📦 pkg/                   # 公共工具包 (JWT / 日志 / 加密 / 头像)
├── 🐳 docker/                # Docker 构建配置（nginx 等）
├── 🐳 docker-compose.yml     # Docker 编排（MySQL+Redis+后端+网关）
├── ⚙️ deploy/                # 生产部署脚本（含健康检查与回滚）
├── 🔧 .github/workflows/     # CI/CD 流水线
├── 🧪 test/                  # 测试（unit / integration / e2e）
├── 📖 docs/                  # 文档（API / 设计 / 规范）
└── 🔧 .env.example           # 环境变量示例
```

> 💡 后端采用经典的 **Controller → Service → Repository** 三层架构，模块边界清晰，便于扩展与测试。

---

## 🚀 快速开始

### 方式一：Docker 一键部署 ⭐（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/fcy222fcy/myblog.git
cd myblog

# 2. 配置环境变量
cp .env.example .env
# 编辑 .env，修改数据库密码与 JWT_SECRET（生产环境务必修改）

# 3. 一键启动全部服务
docker-compose up -d
```

启动完成后访问：

| 入口 | 地址 |
|------|------|
| 🌐 前台 | http://localhost |
| 🔧 后台 | http://localhost:8081 |
| ⚙️ API | http://localhost:9090 |

> 默认网关端口由 `GATEWAY_PORT` 控制（默认 80），API 端口由 `SERVER_PORT` 控制（默认 9090）。

### 方式二：本地开发

#### 前置要求

- [Go](https://go.dev/) ≥ 1.25
- [Node.js](https://nodejs.org/) ≥ 18
- [MySQL](https://www.mysql.com/) ≥ 8.0
- [Redis](https://redis.io/) ≥ 7（可选，缺省时自动降级）

#### 启动后端

```bash
# 1. 配置 .env 中的数据库连接
cp .env.example .env

# 2. 安装依赖
go mod download

# 3. 启动服务（首次启动自动迁移数据库表）
go run ./internal/app
# 默认监听 :9090
```

#### 启动前台

```bash
cd blog-web
npm install
npm run dev    # http://localhost:5173
```

#### 启动后台

```bash
cd blog-admin
npm install
npm run dev    # http://localhost:5174
```

---

## ⚙️ 配置说明

核心环境变量参考 [.env.example](.env.example)：

| 变量 | 说明 |
|------|------|
| `DB_*` | MySQL 数据库连接信息 |
| `REDIS_*` | Redis 连接配置 |
| `JWT_SECRET` | JWT 签名密钥（**生产环境务必修改**） |
| `JWT_EXPIRE_HOURS` | Token 过期时间，默认 7 天 |
| `EMAIL_*` | SMTP 邮件配置（评论通知） |
| `CORS_ORIGINS` | 允许跨域的域名，逗号分隔 |
| `GATEWAY_PORT` / `SERVER_PORT` | 网关 / API 端口 |
| `LOG_*` | 日志文件配置（路径/大小/备份/保留） |

> 🔐 **博主身份**：通过 `blogger.user_id` 服务端配置标识，前端不可伪造，杜绝冒充评论。

---

## 🔄 CI/CD 流水线

项目内置完整的 GitHub Actions 发布流水线：

```
push main ──► CI 校验 ──► 构建镜像 ──► 发布 GHCR ──► 生产审批 ──► 部署服务器
                │            │                          │             │
              gofmt/vet    backend     ghcr.io/       Review       健康检查
              测试+覆盖率  双前端      blog-*         approvals    失败自动回滚
```

### 流水线组成

| 工作流 | 触发时机 | 内容 |
|--------|----------|------|
| `ci.yml` | push main / PR | 后端（gofmt、go vet、测试、覆盖率）+ 双前端构建 + 部署配置校验 + 镜像构建 |
| `release-deploy.yml` | CI 成功 | 发布不可变 tag 镜像（`sha-<commit>`）→ SSH 部署到服务器 → 健康检查 → 失败自动回滚 |

### 关键设计
- **不可变发布**：每个提交对应唯一镜像 tag（`sha-<commit>`），支持精准回滚到任意历史版本
- **人工审批**：生产环境部署需在 GitHub 上批准，防止误发布
- **自动回滚**：部署后健康检查失败自动恢复上一个成功版本

---

## 🧪 测试

```bash
# 运行全部测试
go test ./...

# 按层次分类
go test ./internal/service/...    # 业务逻辑测试
go test ./internal/repository/... # 数据访问测试
go test ./test/unit/...           # 单元测试
go test ./test/integration/...    # 集成测试
go test ./test/e2e/...            # 端到端测试

# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

> CI 流水线会自动执行测试并上传覆盖率报告（backend-coverage artifact）。

---

## 📖 文档

- [API 接口文档](docs/README.md) · [OpenAPI 规范](docs/api.yaml)
- [CI/CD 部署指南](docs/CI-CD部署指南.md)
- [后端代码开发规范](docs/后端代码开发规范.md)
- [项目设计文档](docs/设计文档.md)
- [搜索功能设计](docs/搜索功能设计文档.md)
- [用户头像获取设计](docs/用户头像获取功能设计.md)
- [模块调用关系图](docs/backend-modules-flow.md)

---

## 🤝 参与贡献

任何形式的贡献都非常欢迎！👏

1. Fork 本仓库
2. 创建特性分支：`git checkout -b feature/your-feature`
3. 提交改动：`git commit -m 'feat: add some feature'`
4. 推送到分支：`git push origin feature/your-feature`
5. 提交 Pull Request

> 💡 提交信息建议遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范；代码请通过 `gofmt` 与 `go vet` 检查。

---

## 🗺️ Roadmap

- [x] 核心博客能力（文章/分类/标签/评论）
- [x] 每日一问 + 日历排期
- [x] Docker 部署 + CI/CD 自动发布
- [x] 从 CSDN 批量迁移历史文章
- [ ] 站内短链接 / 代码片段分享
- [ ] 文章数据统计报表增强
- [ ] 基于 AI 的摘要与标签推荐

---

## 📄 许可证

MIT © [fcy222fcy](https://github.com/fcy222fcy)

---

<p align="center">
  <samp>
    如果这个项目对你有帮助的话，不妨点个 ⭐ Star 支持一下~
  </samp>
</p>
