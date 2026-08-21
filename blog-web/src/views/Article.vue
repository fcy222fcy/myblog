<template>
  <div class="page-view">
    <span class="back-btn" role="button" tabindex="0" @click="goBack" @keydown.enter.prevent="goBack">← 返回列表</span>

    <!-- 骨架屏加载 -->
    <ArticleDetailSkeleton v-if="articleStore.loading" />

    <!-- 错误状态 -->
    <ErrorState v-else-if="articleStore.error" :message="articleStore.error" @retry="loadArticle" />

    <!-- 文章内容 -->
    <template v-else-if="articleStore.currentArticle">
      <article class="post-detail">
        <!-- 全宽封面图：有图才渲染，失败自动隐藏 -->
        <div
          v-if="hasCover(articleStore.currentArticle)"
          class="post-cover-wrap"
        >
          <img
            :src="resolveCover(articleStore.currentArticle.cover)"
            :alt="articleStore.currentArticle.title"
            class="post-cover-img"
            @error="onCoverError"
          >
        </div>

        <!-- 分类 + 标签：同一行横着排，# 前缀 + 主页胶囊尺寸 -->
        <div
          v-if="(articleStore.currentArticle.category && articleStore.currentArticle.category.name) || (articleStore.currentArticle.tags && articleStore.currentArticle.tags.length)"
          class="post-pills-row"
        >
          <router-link
            v-if="articleStore.currentArticle.category && articleStore.currentArticle.category.name"
            :to="'/category/' + articleStore.currentArticle.category.id"
            class="post-category-pill"
            :title="'查看「' + articleStore.currentArticle.category.name + '」分类下的文章'"
          >
            {{ articleStore.currentArticle.category.name }}
          </router-link>
          <router-link
            v-for="t in (articleStore.currentArticle.tags || [])"
            :key="t.id"
            :to="'/tag/' + t.id"
            class="post-tag-pill"
            :title="'查看「' + t.name + '」标签下的文章'"
          >#{{ t.name }}</router-link>
        </div>

        <h1 class="post-title">{{ articleStore.currentArticle.title }}</h1>
        <div class="post-meta">
          <span class="meta-item">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
            {{ formatDateTime(articleStore.currentArticle.created_at) }}
          </span>
          <span class="meta-item" v-if="articleStore.currentArticle.updated_at && articleStore.currentArticle.updated_at !== articleStore.currentArticle.created_at">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
            更新于 {{ formatDateTime(articleStore.currentArticle.updated_at) }}
          </span>
          <span class="meta-item">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
            时长 {{ articleStore.currentArticle.reading_time || 1 }} 分钟
          </span>
          <span class="meta-item">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
            浏览量 {{ articleStore.currentArticle.view_count }}
          </span>
        </div>
        <div class="post-content" v-html="renderedContent"></div>
      </article>

      <!-- 评论区 -->
      <div id="comment-section">
        <CommentSection ref="commentSectionRef" />
      </div>
    </template>

    <!-- 文章不存在 -->
    <div v-else class="empty-state">
      <p>文章不存在</p>
      <router-link to="/" class="back-home-link">返回首页</router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useArticleStore } from '../stores/article'
import { marked } from 'marked'
import CommentSection from '../components/comment/CommentSection.vue'
import ArticleDetailSkeleton from '../components/common/ArticleDetailSkeleton.vue'
import ErrorState from '../components/common/ErrorState.vue'
import { updateMetaTags, addStructuredData, resetMetaTags } from '../utils/seo'
import { formatDate } from '../utils/date'

const formatDateTime = (d) => formatDate(d, { withTime: true })

// 封面：有图才渲染，URL 兼容，加载失败自动隐藏
const hasCover = (article) => !!(article && article.cover && typeof article.cover === 'string' && article.cover.trim())
const resolveCover = (url) => {
  if (!url) return ''
  const s = url.trim()
  if (/^https?:\/\//i.test(s)) return s
  if (s.startsWith('//')) return location.protocol + s
  return s.startsWith('/') ? s : '/' + s
}
const onCoverError = (e) => {
  const img = e.target
  if (!img) return
  img.style.display = 'none'
  const wrap = img.closest('.post-cover-wrap')
  if (wrap) wrap.style.display = 'none'
}

const route = useRoute()
const router = useRouter()
const articleStore = useArticleStore()

// 「返回列表」：返回进入详情前的来源页（首页/归档/分类/标签/搜索），
// 若直接打开详情（无站内来源）则回归档列表，避免误跳主页
const LIST_ROUTE_NAMES = ['Home', 'Archives', 'Category', 'Tag', 'Search']
const goBack = () => {
  const prevName = router.__prevRouteName
  const prevPath = router.__prevRoutePath
  if (prevName && (LIST_ROUTE_NAMES.includes(prevName) || prevPath === '/')) {
    router.back()
  } else {
    router.push('/archives')
  }
}

// 简单 HTML 净化：移除事件属性和危险标签
const sanitizeHtml = (html) => {
  // 移除 on* 事件属性（onerror, onclick, onload 等）
  let clean = html.replace(/\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)/gi, '')
  // 移除 javascript: 协议
  clean = clean.replace(/href\s*=\s*(?:"javascript:[^"]*"|'javascript:[^']*')/gi, '')
  clean = clean.replace(/src\s*=\s*(?:"javascript:[^"]*"|'javascript:[^']*')/gi, '')
  // 移除 <script> 和 <iframe> 标签
  clean = clean.replace(/<script[\s\S]*?<\/script>/gi, '')
  clean = clean.replace(/<iframe[\s\S]*?<\/iframe>/gi, '')
  return clean
}

const renderedContent = computed(() => {
  if (!articleStore.currentArticle?.content) return ''
  let html = sanitizeHtml(marked(articleStore.currentArticle.content))
  html = html.replace(/<table>/g, '<div class="table-wrapper"><table>')
  html = html.replace(/<\/table>/g, '</table></div>')
  return html
})

// 计算阅读时长（中文约 400 字/分钟，英文约 200 词/分钟）
const readingTime = computed(() => {
  if (!articleStore.currentArticle?.content) return 1
  const content = articleStore.currentArticle.content
  // 去除 markdown 标记
  let cleaned = content
  const replacements = ['#', '*', '`', '>', '[', ']', '(', ')', '!', '-']
  for (const r of replacements) {
    cleaned = cleaned.split(r).join('')
  }
  cleaned = cleaned.trim()

  if (!cleaned) return 1

  // 计算字符数（使用 spread 正确处理 UTF-16）
  const runeCount = [...cleaned].length
  // 计算单词数
  const wordCount = cleaned.split(/\s+/).filter(w => w).length

  // 混合计算：与后端保持一致
  const minutes = Math.floor(runeCount / 400) + Math.floor(wordCount / 200)
  return Math.max(1, minutes)
})

const loadArticle = () => {
  const slug = route.params.slug
  if (slug) articleStore.fetchArticleDetail(slug)
}

// SEO: 文章加载后更新 Meta 标签
watch(
  () => articleStore.currentArticle,
  (article) => {
    if (article) {
      updateMetaTags(article)
      addStructuredData(article)
    }
  },
  { immediate: true }
)

onMounted(loadArticle)
watch(() => route.params.slug, loadArticle)

onUnmounted(() => {
  resetMetaTags()
})
</script>

<style scoped>
/* 文章详情页样式在全局main.css中定义 */

/* ==== 详情页新增：全宽封面图 ==== */
.post-cover-wrap {
  width: 100%;
  max-height: 360px;
  margin: 0 0 28px;
  border-radius: 14px;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(var(--accent-color-rgb), 0.12), rgba(var(--accent-color-rgb), 0.04));
  box-shadow: 0 8px 24px -8px rgba(0, 0, 0, 0.1);
}
.post-cover-img {
  width: 100%;
  height: 100%;
  max-height: 360px;
  object-fit: cover;
  display: block;
}

/* ==== 详情页：分类+标签 同一行横排，尺寸对齐主页 category-pill ==== */
.post-pills-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin: 0 0 18px;
  line-height: 1;
}
.post-category-pill,
.post-tag-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
  padding: 0 18px;
  border-radius: var(--tag-border-radius);
  font-size: 0.95rem;
  font-weight: 700;
  letter-spacing: 0.01em;
  transition: opacity 0.15s ease, filter 0.15s ease, background 0.15s ease, transform 0.15s ease;
  user-select: none;
  text-decoration: none;
}
/* 分类：主色实底 + 白字（突出「归类」） */
.post-category-pill {
  background: var(--accent-color);
  color: #fff;
}
/* 标签：主色浅背景 + 主色深字（低调，和分类区分） */
.post-tag-pill {
  background: rgba(var(--accent-color-rgb), 0.1);
  color: var(--accent-color);
}
.post-category-pill:hover {
  filter: brightness(1.08);
}
.post-tag-pill:hover {
  background: rgba(var(--accent-color-rgb), 0.18);
}

.post-meta {
  display: flex;
  align-items: center;
  gap: 20px;
  margin: 12px 0 24px;
  color: var(--card-text-color-tertiary);
  font-size: 0.9rem;
}

.meta-item {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.meta-item svg {
  color: var(--card-text-color-tertiary);
  flex-shrink: 0;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-state p {
  color: var(--card-text-color-secondary);
  margin-bottom: 16px;
}

.back-home-link {
  color: var(--accent-color);
  text-decoration: none;
}

.back-home-link:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .post-cover-wrap {
    max-height: 240px;
    margin-bottom: 22px;
    border-radius: 12px;
  }
  .post-cover-img {
    max-height: 240px;
  }
  .post-pills-row {
    gap: 8px;
    margin-bottom: 14px;
  }
  .post-category-pill,
  .post-tag-pill {
    min-height: 32px;
    padding: 0 14px;
    font-size: 0.85rem;
    font-weight: 600;
  }
}

@media (max-width: 480px) {
  .post-cover-wrap {
    max-height: 200px;
    border-radius: 10px;
  }
  .post-cover-img {
    max-height: 200px;
  }
  .post-pills-row {
    gap: 6px;
    margin-bottom: 12px;
  }
  .post-category-pill,
  .post-tag-pill {
    min-height: 30px;
    padding: 0 12px;
    font-size: 0.8rem;
    font-weight: 600;
  }
}
</style>
