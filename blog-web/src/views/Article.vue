<template>
  <div class="page-view">
    <span class="back-btn" role="button" tabindex="0" @click="goBack" @keydown.enter.prevent="goBack">← 返回列表</span>

    <!-- 骨架屏加载 -->
    <ArticleDetailSkeleton v-if="articleStore.loading" />

    <!-- 错误状态 -->
    <ErrorState v-else-if="articleStore.error" :message="articleStore.error" @retry="loadArticle" />

    <!-- 文章内容 -->
    <template v-else-if="articleStore.currentArticle">
      <div class="article-reading-layout">
        <div class="article-primary-column">
          <article class="post-detail">
        <!-- 全宽封面图：有图才渲染，失败自动隐藏；点击可放大预览 -->
        <div
          v-if="hasCover(articleStore.currentArticle)"
          class="post-cover-wrap"
          role="button"
          tabindex="0"
          :aria-label="'查看封面大图：' + articleStore.currentArticle.title"
          @click="openCoverLightbox"
          @keydown.enter.prevent="openCoverLightbox"
        >
          <img
            :src="resolveCover(articleStore.currentArticle.cover)"
            :alt="articleStore.currentArticle.title"
            class="post-cover-img"
            @error="onCoverError"
          >
          <span class="cover-zoom-hint" aria-hidden="true">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line><line x1="11" y1="8" x2="11" y2="14"></line><line x1="8" y1="11" x2="14" y2="11"></line></svg>
          </span>
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

            <div
              ref="postContentRef"
              class="post-content"
              v-html="renderedContent"
              @click="onContentClick"
              @error.capture="onContentImageError"
            ></div>
          </article>

          <!-- 评论区与正文保持同一列 -->
          <div id="comment-section">
            <CommentSection ref="commentSectionRef" />
          </div>
        </div>

        <aside v-if="toc.length" class="article-toc" aria-label="文章目录">
          <div class="article-toc-sticky">
            <div class="article-toc-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="9" x2="19" y2="9"></line>
                <line x1="5" y1="15" x2="19" y2="15"></line>
                <line x1="11" y1="4" x2="7" y2="20"></line>
                <line x1="17" y1="4" x2="13" y2="20"></line>
              </svg>
            </div>
            <h2>目录</h2>
            <nav class="article-toc-card">
              <ol>
                <li
                  v-for="item in toc"
                  :key="item.id"
                  :class="['toc-level-' + item.level, { active: activeHeadingID === item.id }]"
                >
                  <a :href="'#' + item.id">
                    <span class="toc-number" aria-hidden="true">{{ item.number }}</span>
                    <span>{{ item.text }}</span>
                  </a>
                </li>
              </ol>
            </nav>
          </div>
        </aside>
      </div>

      <!-- 图片放大预览（lightbox）：点击正文图片打开，点遮罩/×/Esc 关闭 -->
      <Teleport to="body">
        <div
          v-if="lightbox.visible"
          class="lightbox-overlay"
          role="dialog"
          aria-modal="true"
          aria-label="图片预览"
          @click.self="closeLightbox"
        >
          <button class="lightbox-close" type="button" aria-label="关闭预览" @click="closeLightbox">✕</button>
          <img :src="lightbox.src" :alt="lightbox.alt" class="lightbox-img" @click="closeLightbox">
          <p v-if="lightbox.alt" class="lightbox-caption">{{ lightbox.alt }}</p>
        </div>
      </Teleport>

    </template>

    <!-- 文章不存在 -->
    <div v-else class="empty-state">
      <p>文章不存在</p>
      <router-link to="/" class="back-home-link">返回首页</router-link>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onMounted, watch, onUnmounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useArticleStore } from '../stores/article'
import CommentSection from '../components/comment/CommentSection.vue'
import ArticleDetailSkeleton from '../components/common/ArticleDetailSkeleton.vue'
import ErrorState from '../components/common/ErrorState.vue'
import { updateMetaTags, addStructuredData, resetMetaTags } from '../utils/seo'
import { formatDate } from '../utils/date'
import { renderArticleMarkdown } from '../utils/articleEnhancements'
import 'highlight.js/styles/github-dark.css'

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

const articleRender = computed(() => renderArticleMarkdown(articleStore.currentArticle?.content || ''))
const toc = computed(() => articleRender.value.toc)

const renderedContent = computed(() => {
  if (!articleStore.currentArticle?.content) return ''
  let html = sanitizeHtml(articleRender.value.html)
  html = html.replace(/<table>/g, '<div class="table-wrapper"><table>')
  html = html.replace(/<\/table>/g, '</table></div>')
  return html
})

// ==== 文章目录：当前章节随阅读位置更新 ====
const postContentRef = ref(null)
const activeHeadingID = ref('')
let headingObserver = null

const observeHeadings = async () => {
  if (headingObserver) headingObserver.disconnect()
  headingObserver = null
  activeHeadingID.value = toc.value[0]?.id || ''

  await nextTick()
  if (!postContentRef.value || !toc.value.length || typeof IntersectionObserver === 'undefined') return

  headingObserver = new IntersectionObserver((entries) => {
    const visible = entries
      .filter(entry => entry.isIntersecting)
      .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top)
    if (visible[0]?.target?.id) activeHeadingID.value = visible[0].target.id
  }, { rootMargin: '-88px 0px -68% 0px', threshold: [0, 1] })

  toc.value.forEach((item) => {
    const heading = document.getElementById(item.id)
    if (heading) headingObserver.observe(heading)
  })
}

// ==== 正文图片点击放大（lightbox）====
const lightbox = reactive({ visible: false, src: '', alt: '' })

const normalizeImgSrc = (src) => {
  if (!src) return ''
  if (/^(?:https?:)?\/\//i.test(src)) return src
  return src.startsWith('/') ? src : '/' + src
}

// 事件委托：点击正文中非链接内的图片时打开预览
const copyCode = async (button) => {
  const code = button.closest('.code-block')?.querySelector('code')
  if (!code) return
  const text = [...code.querySelectorAll('.code-line')].map(line => line.textContent || '').join('\n')

  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text)
    } else {
      const textarea = document.createElement('textarea')
      textarea.value = text
      textarea.style.position = 'fixed'
      textarea.style.opacity = '0'
      document.body.appendChild(textarea)
      textarea.select()
      document.execCommand('copy')
      textarea.remove()
    }
    button.textContent = '已复制'
    button.classList.add('copied')
  } catch {
    button.textContent = '复制失败'
  }

  window.setTimeout(() => {
    if (!button.isConnected) return
    button.textContent = '复制'
    button.classList.remove('copied')
  }, 1800)
}

const onContentClick = (e) => {
  const target = e.target
  if (!(target instanceof HTMLElement)) return

  const copyButton = target.closest('[data-copy-code]')
  if (copyButton) {
    copyCode(copyButton)
    return
  }

  if (target.tagName !== 'IMG') return
  // 图片外层是链接时保留原跳转行为，不劫持
  if (target.closest('a')) return
  const src = normalizeImgSrc(target.getAttribute('src') || '')
  if (!src) return
  lightbox.src = src
  lightbox.alt = target.getAttribute('alt') || ''
  lightbox.visible = true
  e.preventDefault()
}

const onContentImageError = (e) => {
  const image = e.target
  if (!(image instanceof HTMLImageElement)) return
  const frame = image.closest('.article-image-frame')
  if (frame) frame.classList.add('image-failed')
}

const closeLightbox = () => {
  lightbox.visible = false
}

// 封面图点击放大：复用同一 lightbox；封面在卡片里被 object-fit:cover 裁切，
// 放大后 object-fit:contain 显示完整大图
const openCoverLightbox = () => {
  const article = articleStore.currentArticle
  if (!hasCover(article)) return
  lightbox.src = resolveCover(article.cover)
  lightbox.alt = article.title || ''
  lightbox.visible = true
}

const onKeydown = (e) => {
  if (e.key === 'Escape' && lightbox.visible) closeLightbox()
}

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
      observeHeadings()
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadArticle()
  document.addEventListener('keydown', onKeydown)
})
watch(() => route.params.slug, loadArticle)

onUnmounted(() => {
  resetMetaTags()
  document.removeEventListener('keydown', onKeydown)
  if (headingObserver) headingObserver.disconnect()
})
</script>

<style scoped>
/* 文章详情页样式在全局main.css中定义 */

/* ==== 详情页新增：全宽封面图 ==== */
.post-cover-wrap {
  position: relative;
  width: 100%;
  max-height: 360px;
  margin: 0 0 28px;
  border-radius: 14px;
  overflow: hidden;
  background: linear-gradient(135deg, rgba(var(--accent-color-rgb), 0.12), rgba(var(--accent-color-rgb), 0.04));
  box-shadow: 0 8px 24px -8px rgba(0, 0, 0, 0.1);
  cursor: zoom-in;
  transition: box-shadow 0.2s ease;
}
.post-cover-wrap:hover {
  box-shadow: 0 12px 32px -8px rgba(0, 0, 0, 0.18);
}
.post-cover-wrap:focus-visible {
  outline: 2px solid var(--accent-color);
  outline-offset: 3px;
}
.post-cover-img {
  width: 100%;
  height: 100%;
  max-height: 360px;
  object-fit: cover;
  display: block;
  transition: transform 0.3s ease;
}
.post-cover-wrap:hover .post-cover-img {
  transform: scale(1.03);
}
/* 封面「可放大」提示图标，悬停时浮现 */
.cover-zoom-hint {
  position: absolute;
  right: 12px;
  bottom: 12px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  opacity: 0;
  transform: translateY(4px);
  transition: opacity 0.2s ease, transform 0.2s ease;
  pointer-events: none;
}
.post-cover-wrap:hover .cover-zoom-hint,
.post-cover-wrap:focus-visible .cover-zoom-hint {
  opacity: 1;
  transform: translateY(0);
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

.article-reading-layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) clamp(300px, 24vw, 374px);
  align-items: start;
  gap: 30px;
  width: calc(100% + 48px);
}

.article-primary-column {
  min-width: 0;
}

.article-toc {
  min-width: 0;
  align-self: stretch;
}

.article-toc-sticky {
  position: sticky;
  top: 40px;
  max-height: calc(100vh - 80px);
  overflow-y: auto;
  padding: 0 4px 12px 0;
  scrollbar-width: thin;
  scrollbar-color: rgba(var(--accent-color-rgb), 0.25) transparent;
}

.article-toc-icon {
  width: 36px;
  height: 36px;
  color: var(--card-text-color-main);
}

.article-toc-icon svg {
  display: block;
  width: 36px;
  height: 36px;
}

.article-toc h2 {
  margin: 0 0 10px;
  color: var(--card-text-color-main);
  font-size: 1rem;
  font-weight: 700;
  line-height: 1.15;
}

.article-toc-card {
  overflow: hidden;
  padding: 10px 14px 12px;
  border: 1px solid var(--card-separator-color);
  border-radius: 10px;
  background: var(--card-background);
  box-shadow: var(--shadow-l1);
}

.article-toc-card ol {
  margin: 0;
  padding: 0;
  list-style: none;
}

.article-toc li {
  margin: 15px 0;
  line-height: 1.2;
}

.article-toc li:first-child {
  margin-top: 5px;
}

.article-toc li:last-child {
  margin-bottom: 5px;
}

.article-toc .toc-level-3 {
  margin-left: 35px;
}

.article-toc a {
  display: flex;
  align-items: baseline;
  gap: 8px;
  padding: 5px;
  color: var(--accent-color);
  font-size: 1rem;
  text-decoration: none;
  text-wrap: pretty;
  transition: color 0.15s ease, transform 0.15s ease;
}

.article-toc .toc-number {
  flex: 0 0 auto;
  font-variant-numeric: tabular-nums;
}

.article-toc a:hover {
  color: var(--card-text-color-main);
  transform: translateX(2px);
}

.article-toc li.active a {
  color: var(--accent-color);
  font-weight: 700;
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

@media (max-width: 1280px) {
  .article-reading-layout {
    display: block;
    width: 100%;
  }
  .article-toc {
    display: none;
  }
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

/* ==== 正文图片放大预览（lightbox）==== */
.lightbox-overlay {
  position: fixed;
  inset: 0;
  z-index: 3000;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 24px 24px;
  background: rgba(0, 0, 0, 0.88);
  cursor: zoom-out;
  animation: lightboxFadeIn 0.2s ease-out;
}
.lightbox-img {
  max-width: min(92vw, 1400px);
  max-height: 84vh;
  object-fit: contain;
  border-radius: 8px;
  background: #fff;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.6);
}
.lightbox-caption {
  margin-top: 14px;
  max-width: 80vw;
  color: rgba(255, 255, 255, 0.85);
  font-size: 0.92rem;
  text-align: center;
  word-break: break-all;
}
.lightbox-close {
  position: absolute;
  top: 16px;
  right: 20px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.15s ease, transform 0.15s ease;
}
.lightbox-close:hover {
  background: rgba(255, 255, 255, 0.32);
  transform: scale(1.05);
}
@keyframes lightboxFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>
