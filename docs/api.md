# 博客系统 API 接口文档

## 通用说明

### 基础信息

| 项目 | 说明 |
|------|------|
| 基础 URL | `http://localhost:8080/api/v1` |
| 数据格式 | JSON |
| 字符编码 | UTF-8 |
| 认证方式 | JWT Token（管理员接口需要；评论接口可选，用于识别博主身份） |

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

**说明**：`code = 0` 表示成功，非 0 表示失败（参考错误码表）。

### 错误码说明

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| 0 | 200 | 成功 |
| 40000+ | 400 | 请求参数错误（具体见细分码） |
| 40100 | 401 | 未授权（未登录或 Token 无效） |
| 40300 | 403 | 禁止访问（权限不足） |
| 40400 | 404 | 资源不存在 |
| 50000 | 500 | 服务器内部错误 |
| 10001 | 400 | 用户名或密码错误 |
| 10002 | 400 | 旧密码错误 |
| 20001 | 400 | 评论内容不能为空 |
| 20002 | 400 | 文章不存在 |
| 20003 | 400 | 评论不存在 |

### 分页参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 当前页码 |
| page_size | int | 否 | 10 | 每页数量（最大 50） |

### 分页响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 一、认证模块

### 1.1 管理员登录

**请求**

```
POST /auth/login
```

**请求体**

```json
{
  "username": "admin",
  "password": "123456"
}
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": 1719000000
  }
}
```

### 1.2 管理员注册（可选）

**请求**

```
POST /auth/register
```

**请求体**

```json
{
  "username": "admin",
  "password": "123456",
  "email": "admin@example.com"
}
```

**说明**：通常用于首次部署，生产环境建议关闭或限制访问。

### 1.3 修改密码

**请求**

```
PUT /auth/password
```

**请求头**

```
Authorization: Bearer <token>
```

**请求体**

```json
{
  "old_password": "123456",
  "new_password": "new_password"
}
```

**响应**

```json
{
  "code": 0,
  "message": "密码修改成功",
  "data": null
}
```

---

## 二、用户模块

### 2.1 获取公开博主信息（前台）

**请求**

```
GET /user/info
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "nickname": "Fu Chengyan",
    "avatar": "https://example.com/avatar.jpg",
    "email": "admin@example.com",
    "description": "Go 开发者",
    "social_links": [
      {"name": "GitHub", "url": "https://github.com/xxx"},
      {"name": "Twitter", "url": "https://twitter.com/xxx"}
    ]
  }
}
```

### 2.2 获取当前登录用户信息（后台）

**请求**

```
GET /user/profile
```

**请求头**

```
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "Fu Chengyan",
    "avatar": "https://example.com/avatar.jpg",
    "email": "admin@example.com",
    "description": "Go 开发者",
    "social_links": []
  }
}
```

### 2.3 更新当前用户信息（后台）

**请求**

```
PUT /user/profile
```

**请求头**

```
Authorization: Bearer <token>
```

**请求体**

```json
{
  "nickname": "新昵称",
  "avatar": "https://example.com/new-avatar.jpg",
  "email": "new@example.com",
  "description": "新个人简介",
  "social_links": [
    {"name": "GitHub", "url": "https://github.com/xxx"}
  ]
}
```

---

## 三、文章模块（前台）

### 3.1 获取文章列表

**请求**

```
GET /articles
```

**查询参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| category_id | int | 否 | 分类 ID |
| tag_id | int | 否 | 标签 ID |
| status | string | 否 | 固定传 `published`（前台仅返回已发布） |

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "我的博客主题已开源",
        "slug": "open-source-theme",
        "summary": "基于 Hugo Theme Stack 打造的开箱即用博客模板...",
        "cover": "https://example.com/cover.jpg",
        "category": {
          "id": 1,
          "name": "搭建网站",
          "slug": "build",
          "description": "网站搭建、博客部署相关"
        },
        "tags": [
          {"id": 1, "name": "Hugo", "slug": "hugo"},
          {"id": 2, "name": "开源", "slug": "open-source"}
        ],
        "view_count": 252,
        "like_count": 42,
        "comment_count": 8,
        "created_at": "2026-04-24T10:00:00Z",
        "updated_at": "2026-04-24T10:00:00Z"
      }
    ],
    "total": 6,
    "page": 1,
    "page_size": 10
  }
}
```

### 3.2 搜索文章（带 search_snippet 摘要）

**请求**

```
GET /articles/search
```

**查询参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| keyword | string | 是 | 搜索关键词，多关键词以空格分隔 |
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "title": "我的博客主题已开源",
        "slug": "open-source-theme",
        "cover": "https://example.com/cover.jpg",
        "category": {
          "id": 1,
          "name": "搭建网站"
        },
        "tags": [
          {"id": 1, "name": "Hugo"}
        ],
        "search_snippet": "...主题已<strong>开源</strong>，基于 Hugo Theme Stack 打造，支持自定义主题颜色...",
        "view_count": 252,
        "created_at": "2026-04-24T10:00:00Z"
      }
    ],
    "total": 3,
    "page": 1,
    "page_size": 10
  }
}
```

**说明**：`search_snippet` 已清洗 Markdown，并将匹配关键词用 `<strong>` 包裹高亮，前后保留上下文。

### 3.3 获取文章详情

**请求**

```
GET /articles/:slug
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "我的博客主题已开源",
    "slug": "open-source-theme",
    "content": "# 标题\n\n这里是文章正文...",
    "summary": "基于 Hugo Theme Stack 打造的开箱即用博客模板...",
    "cover": "https://example.com/cover.jpg",
    "category": {
      "id": 1,
      "name": "搭建网站",
      "slug": "build",
      "description": "网站搭建、博客部署相关"
    },
    "tags": [
      {"id": 1, "name": "Hugo", "slug": "hugo"}
    ],
    "view_count": 252,
    "like_count": 42,
    "comment_count": 8,
    "reading_time": 3,
    "created_at": "2026-04-24T10:00:00Z",
    "updated_at": "2026-04-24T10:00:00Z"
  }
}
```

### 3.4 文章点赞

**请求**

```
POST /articles/:id/like
```

**响应**

```json
{
  "code": 0,
  "message": "点赞成功",
  "data": {
    "like_count": 43
  }
}
```

### 3.5 获取文章归档

**请求**

```
GET /articles/archives
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "year": 2026,
      "articles": [
        {
          "id": 1,
          "title": "我的博客主题已开源",
          "slug": "open-source-theme",
          "created_at": "2026-04-24T10:00:00Z"
        }
      ]
    }
  ]
}
```

---

## 四、文章管理（后台）

### 4.1 获取文章列表（后台）

**请求**

```
GET /admin/articles
```

**请求头**

```
Authorization: Bearer <token>
```

**查询参数**

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| page | int | 否 | 页码 |
| page_size | int | 否 | 每页数量 |
| category_id | int | 否 | 分类 ID |
| status | string | 否 | 状态：published/draft/空(全部) |
| keyword | string | 否 | 搜索关键词(标题) |

### 4.2 获取文章详情（后台）

**请求**

```
GET /admin/articles/:id
```

**请求头**

```
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "我的博客主题已开源",
    "slug": "open-source-theme",
    "content": "# 标题\n\n这里是文章正文...",
    "summary": "摘要",
    "cover": "https://example.com/cover.jpg",
    "category_id": 1,
    "tag_ids": [1, 2],
    "status": "published",
    "is_top": false,
    "created_at": "2026-04-24T10:00:00Z"
  }
}
```

### 4.3 创建文章

```
POST /admin/articles
Authorization: Bearer <token>
```

**请求体**

```json
{
  "title": "新文章标题",
  "content": "# 标题\n\n文章正文...",
  "summary": "文章摘要（可选，留空自动截取）",
  "cover": "https://example.com/cover.jpg",
  "category_id": 1,
  "tag_ids": [1, 2],
  "status": "draft",
  "is_top": false
}
```

### 4.4 更新文章

```
PUT /admin/articles/:id
Authorization: Bearer <token>
```

### 4.5 删除文章

```
DELETE /admin/articles/:id
Authorization: Bearer <token>
```

### 4.6 批量删除文章

```
POST /admin/articles/batch-delete
Authorization: Bearer <token>

{ "ids": [1, 2, 3] }
```

---

## 五、分类模块

### 5.1 获取分类列表（前台）

**请求**

```
GET /categories
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "搭建网站",
      "slug": "build",
      "description": "网站搭建、博客部署相关",
      "article_count": 3,
      "sort_order": 1
    }
  ]
}
```

### 5.2 获取分类详情（前台）

**请求**

```
GET /categories/:id
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "搭建网站",
    "slug": "build",
    "description": "网站搭建、博客部署相关",
    "article_count": 3
  }
}
```

### 5.3 创建分类（后台）

```
POST /admin/categories
Authorization: Bearer <token>

{
  "name": "新分类",
  "slug": "new-category",
  "description": "分类描述",
  "sort_order": 0
}
```

### 5.4 更新分类（后台）

```
PUT /admin/categories/:id
Authorization: Bearer <token>

{
  "name": "更新后的分类名",
  "slug": "updated-slug",
  "description": "更新后的描述",
  "sort_order": 1
}
```

### 5.5 删除分类（后台）

```
DELETE /admin/categories/:id
Authorization: Bearer <token>
```

---

## 六、标签模块

### 6.1 获取标签列表（前台）

```
GET /tags
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": [
    { "id": 1, "name": "Hugo", "slug": "hugo", "article_count": 3 }
  ]
}
```

### 6.2 创建标签（后台）

```
POST /admin/tags
Authorization: Bearer <token>

{ "name": "新标签", "slug": "new-tag" }
```

### 6.3 更新标签（后台）

```
PUT /admin/tags/:id
Authorization: Bearer <token>

{ "name": "更新后的标签名", "slug": "updated-tag" }
```

### 6.4 删除标签（后台）

```
DELETE /admin/tags/:id
Authorization: Bearer <token>
```

---

## 七、评论模块

### 7.1 获取文章评论列表（前台，支持 OptionalAuth）

**请求**

```
GET /comments/article/:articleId
```

**查询参数**

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码 |
| page_size | int | 否 | 10 | 每页数量 |
| sort_by | string | 否 | hot | 排序方式：`hot`(热度，点赞×2+回复数，同分按时间倒序) / `time_desc`(时间倒序) / `time_asc`(时间正序) |

**请求头（可选）**

```
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "content": "很棒的文章，学到了很多！",
        "nickname": "张三",
        "email": "zhangsan@example.com",
        "website": "https://zhangsan.com",
        "avatar": "https://thirdqq.qlogo.cn/g?b=qq&nk=10000&s=100",
        "is_admin": false,
        "reply_to": null,
        "like_count": 5,
        "user_id": null,
        "status": "approved",
        "created_at": "2026-04-24T12:00:00Z",
        "replies": [
          {
            "id": 2,
            "content": "谢谢支持！",
            "nickname": "博主",
            "avatar": "https://.../avatar.jpg",
            "is_admin": true,
            "reply_to": {
              "id": 1,
              "nickname": "张三"
            },
            "like_count": 2,
            "user_id": 1,
            "created_at": "2026-04-24T13:00:00Z"
          }
        ]
      }
    ],
    "total": 8,
    "page": 1,
    "page_size": 10
  }
}
```

**关键说明**：
- `avatar` 获取优先级：QQ邮箱→提取QQ号走 qlogo.cn → 其他邮箱走 Gravatar → 空字符串（前端兜底首字母头像）
- `is_admin`：仅当用户通过登录按钮登录 **且** 其 UserID 与配置 `blogger.user_id` 一致时为 true（防冒充）
- `reply_to`：回复某条评论时，给出被回复人的 id 和 nickname；根评论为 null
- 排序在后端全量根评论计算后再分页，保证多页数据一致性

### 7.2 提交评论（前台，OptionalAuth）

**请求**

```
POST /comments
Authorization: Bearer <token>  # 可选
```

**请求体**

```json
{
  "article_id": 1,
  "content": "评论内容",
  "nickname": "张三",
  "email": "zhangsan@example.com",
  "website": "https://zhangsan.com",
  "parent_id": 0
}
```

**说明**：
- 若携带了有效登录 Token，会关联 `user_id` 并从用户表填充昵称、邮箱、头像（游客填写的会被覆盖）
- 游客必须至少填 nickname 和 email

### 7.3 评论点赞（前台）

```
POST /comments/:id/like
```

### 7.4 获取评论列表（后台）

```
GET /admin/comments
Authorization: Bearer <token>
```

**查询参数**：page / page_size / status(pending/approved/rejected) / article_id

### 7.5 审核评论（后台）

```
PUT /admin/comments/:id/status
Authorization: Bearer <token>

{ "status": "approved" }
```

### 7.6 删除评论（后台）

```
DELETE /admin/comments/:id
Authorization: Bearer <token>
```

### 7.7 批量删除评论（后台）

```
POST /admin/comments/batch-delete
Authorization: Bearer <token>

{ "ids": [1, 2, 3] }
```

---

## 八、每日一问模块

### 8.1 获取最新问题（前台）

```
GET /daily-questions/latest
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "question": "什么是最好的编程语言？",
    "answer": "没有最好的编程语言，只有最适合的...",
    "date": "2026-07-11",
    "like_count": 42,
    "comment_count": 8
  }
}
```

### 8.2 获取全部已发布问题（前台）

```
GET /daily-questions/all
```

### 8.3 按日期获取问题（前台）

```
GET /daily-questions/date/:date
```

路径参数 `date` 格式：`2026-07-11`

### 8.4 上一天 / 下一天（前台）

```
GET /daily-questions/previous/:date
GET /daily-questions/next/:date
```

### 8.5 每日一问点赞（前台）

```
POST /daily-questions/:id/like
```

### 8.6 获取问题列表（后台）

```
GET /admin/daily-questions
Authorization: Bearer <token>
```

查询参数：page / page_size / status(0禁用,1启用) / keyword

### 8.7 创建问题（后台）

```
POST /admin/daily-questions
Authorization: Bearer <token>

{
  "question": "新问题？",
  "answer": "问题答案...",
  "date": "2026-07-12",
  "status": 1
}
```

### 8.8 更新问题（后台）

```
PUT /admin/daily-questions/:id
Authorization: Bearer <token>
```

### 8.9 启用/禁用问题（后台）

```
PUT /admin/daily-questions/:id/status
Authorization: Bearer <token>

{ "status": 1 }
```

### 8.10 删除问题（后台）

```
DELETE /admin/daily-questions/:id
Authorization: Bearer <token>
```

---

## 九、媒体模块（后台，仅上传/列表/删除，无媒体库页面）

### 9.1 上传文件

```
POST /media/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| file | file | 是 | 支持 JPG、PNG、GIF、WebP、PDF 等 |

**响应**

```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "filename": "1783049034191.jpg",
    "url": "http://localhost:8080/uploads/1783049034191.jpg",
    "size": 102400,
    "mime_type": "image/jpeg"
  }
}
```

### 9.2 获取媒体列表

```
GET /media
Authorization: Bearer <token>
```

### 9.3 删除媒体文件

```
DELETE /media/:filename
Authorization: Bearer <token>
```

说明：根据文件名从 `uploads/` 目录物理删除。

---

## 十、关于页面

### 10.1 获取关于页面内容（前台）

```
GET /about
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "content": "# 关于我\n\n这里是关于页面 Markdown 内容...",
    "updated_at": "2026-07-10T00:00:00Z"
  }
}
```

### 10.2 更新关于页面（后台）

```
PUT /admin/about
Authorization: Bearer <token>

{
  "content": "# 关于我\n\n更新后的内容..."
}
```

---

## 十一、仪表盘（后台）

### 11.1 获取仪表盘统计

```
GET /admin/dashboard/stats
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "article_count": 20,
    "published_count": 18,
    "draft_count": 2,
    "total_views": 12500,
    "today_views": 36,
    "comment_count": 256,
    "pending_count": 3
  }
}
```

### 11.2 获取最近文章

```
GET /admin/dashboard/recent-articles?limit=5
Authorization: Bearer <token>
```

查询参数 `limit`：1 ~ 20，默认 5。

---

## 十二、审计日志（后台）

### 12.1 获取审计日志列表

```
GET /admin/audit-logs
Authorization: Bearer <token>
```

**查询参数**

| 参数 | 类型 | 说明 |
|------|------|------|
| page | int | 页码 |
| page_size | int | 每页数量 |
| action | string | 操作类型过滤（登录/创建文章/删除评论 等） |
| module | string | 模块过滤（auth/article/comment 等） |
| keyword | string | 操作者/详情 搜索 |

**说明**：仅记录管理员写操作（创建/更新/删除/登录等），读操作不记录。

---

## 十三、RSS / Sitemap（公开，无需 Token）

### 13.1 RSS 订阅

```
GET /rss.xml   或  GET /rss
```

### 13.2 Sitemap

```
GET /sitemap.xml
```

---

## 附录

### 1. 状态码枚举

**文章状态**
- `published` - 已发布
- `draft` - 草稿

**评论状态**
- `pending` - 待审核
- `approved` - 已通过
- `rejected` - 已拒绝

**每日一问状态**
- `0` - 禁用
- `1` - 启用

**评论排序方式**
- `hot` - 热度优先（点赞数×2 + 所有层级回复数，同分按时间倒序）
- `time_desc` - 时间倒序（最新在前）
- `time_asc` - 时间正序（最早在前）

### 2. 头像获取规则

1. 邮箱为 QQ 邮箱（含 `@vip.qq.com`）→ 提取 QQ 号，走 `https://thirdqq.qlogo.cn/g?b=qq&nk={QQ号}&s=100`
2. 其他邮箱 → 走 Gravatar
3. 无邮箱或 Gravatar 无头像 → 后端返回空字符串，前端走 `ui-avatars.com` 生成首字母头像兜底

### 3. 媒体类型

- `image` - 图片
- `document` - 文档

### 4. 分页默认值

- 默认页码：1
- 默认每页数量：10
- 最大每页数量：50

### 5. 模块变更记录

| 版本 | 变更内容 |
|------|----------|
| v2026-07-11 | 移除娱乐模块(entertainment)、移除友链模块(link)；新增评论排序 sort_by；新增文章 /articles/search；补充用户/关于/审计日志/RSS/Sitemap 接口 |
| v2026-07-10 | 评论响应补充 is_admin、reply_to、like_count；新增 OptionalAuth；新增修改密码 PUT /auth/password |
