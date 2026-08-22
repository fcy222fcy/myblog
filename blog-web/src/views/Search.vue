<template>
  <div class="page-view" id="view-search">
    <!-- 统计卡片 -->
    <div class="stats-grid" v-if="keyword && !loading && !errorMsg">
      <div class="stat-card">
        <div class="stat-label">搜索关键词</div>
        <div class="stat-value search-keyword-stat">"{{ keyword }}"</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">相关文章</div>
        <div class="stat-value">{{ total }}</div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <div class="card-title">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
          文章列表
        </div>
      </div>

      <div class="article-cards-container">
        <!-- 加载状态 -->
        <div v-if="loading" class="search-loading">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="animate-spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg>
          搜索中...
        </div>

        <!-- 错误状态 -->
        <div v-else-if="errorMsg" class="error-inline">
          <p class="error-message">{{ errorMsg }}</p>
          <button class="btn btn-sm btn-primary" @click="fetchResults">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/></svg>
            重试
          </button>
        </div>

        <!-- 真实列表 -->
        <template v-else>
          <router-link
            v-for="article in articles"
            :key="article.id"
            :to="'/post/' + article.slug"
            class="article-list-item"
          >
            <div class="article-list-body">
              <h3 class="article-list-title">
                <template v-for="(seg, idx) in highlightTokens(article.title, keyword)" :key="'t-'+article.id+'-'+idx">
                  <mark v-if="seg.matched" class="search-hl">{{ seg.text }}</mark>
                  <span v-else>{{ seg.text }}</span>
                </template>
              </h3>
              <p class="article-list-summary" v-if="displaySummary(article)">
                <template v-for="(seg, idx) in highlightTokens(displaySummary(article), keyword)" :key="'s-'+article.id+'-'+idx">
                  <mark v-if="seg.matched" class="search-hl">{{ seg.text }}</mark>
                  <span v-else>{{ seg.text }}</span>
                </template>
              </p>
              <p v-else class="article-list-summary article-list-summary--empty">暂无摘要</p>
            </div>
          </router-link>

          <div v-if="articles.length === 0 && keyword" class="list-empty">
            <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line><line x1="8" y1="8" x2="14" y2="14"></line><line x1="14" y1="8" x2="8" y2="14"></line></svg>
            <p class="list-empty-title">没有找到与「{{ keyword }}」相关的文章</p>
            <p class="list-empty-hint">换个关键词试试？</p>
          </div>
        </template>
      </div>

      <!-- 分页 -->
      <div class="card-body" v-if="totalPage > 1 && !loading && !errorMsg && articles.length > 0">
        <div class="pagination">
          <div class="pagination-info">
            共 <span>{{ total }}</span> 篇相关文章
          </div>
          <div class="pagination-buttons">
            <button
              class="pagination-btn"
              :disabled="currentPage <= 1"
              @click="goToPage(currentPage - 1)"
            >
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m15 18-6-6 6-6"></path></svg>
              上一页
            </button>
            <button class="pagination-btn active">{{ currentPage }}</button>
            <button
              class="pagination-btn"
              :disabled="currentPage >= totalPage"
              @click="goToPage(currentPage + 1)"
            >
              下一页
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m9 18 6-6-6-6"></path></svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router'
import { searchArticles } from '../api/article'
import { saveListState, takeListState } from '../composables/useListState'

const route = useRoute()
const router = useRouter()

const keyword = ref('')
const articles = ref([])
const loading = ref(false)
const errorMsg = ref('')
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10
const totalPage = ref(0)

const escapeRegExp = (s) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

/**
 * 安全的关键词高亮：把文本切分成「普通段 / 命中段」数组，避免 v-html 导致 XSS
 * @param {string} text - 原始纯文本
 * @param {string} kw   - 搜索关键词（可能包含空白、特殊字符）
 * @returns {Array<{text:string, matched:boolean}>}
 */
const highlightTokens = (text, kw) => {
  const src = String(text ?? '')
  const key = String(kw ?? '').trim()
  if (!src || !key) return [{ text: src, matched: false }]
  const keys = key.split(/\s+/).filter(Boolean)
  const pattern = new RegExp(keys.map(escapeRegExp).join('|'), 'ig')
  const out = []
  let lastEnd = 0
  let m
  while ((m = pattern.exec(src)) !== null) {
    if (m.index > lastEnd) out.push({ text: src.slice(lastEnd, m.index), matched: false })
    out.push({ text: m[0], matched: true })
    lastEnd = m.index + m[0].length
    if (m.index === pattern.lastIndex) pattern.lastIndex++ // 防零宽死循环
  }
  if (lastEnd < src.length) out.push({ text: src.slice(lastEnd), matched: false })
  if (out.length === 0) out.push({ text: src, matched: false })
  return out
}

// 搜索摘要展示顺序：优先正文命中片段(SearchSnippet) → 文章原 Summary
const displaySummary = (article) => {
  if (!article) return ''
  if (article.search_snippet && article.search_snippet.trim()) return article.search_snippet
  return article.summary || ''
}

const fetchResults = async () => {
  if (!keyword.value.trim()) return

  loading.value = true
  errorMsg.value = ''
  try {
    const res = await searchArticles({
      keyword: keyword.value,
      page: currentPage.value,
      page_size: pageSize
    })
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
    totalPage.value = res.data?.total_page || 0
  } catch (err) {
    errorMsg.value = '搜索失败，请稍后重试'
    console.error('搜索失败:', err)
  } finally {
    loading.value = false
  }
}

const goToPage = (page) => {
  if (page < 1 || page > totalPage.value) return
  currentPage.value = page
  fetchResults()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 离开搜索页（去往文章详情）时保存快照，返回时恢复页码与滚动位置
onBeforeRouteLeave(() => {
  saveListState(route.fullPath, { page: currentPage.value })
})

// 关键词变化属于「新搜索」，主动消费掉旧快照（避免重搜同词时跳到旧位置）
watch(
  () => route.query.q,
  (newQ) => {
    if (newQ) {
      keyword.value = newQ
      currentPage.value = 1
      takeListState(route.fullPath)
      window.scrollTo({ top: 0 })
      fetchResults()
    }
  }
)

onMounted(async () => {
  const q = route.query.q
  if (!q) return
  keyword.value = q
  // 从文章详情返回时优先消费快照，决定是否恢复页码与滚动位置
  const state = takeListState(route.fullPath)
  if (state && state.page > 1) {
    currentPage.value = state.page
  }
  await fetchResults()
  if (state && state.scrollY > 0) {
    // 等搜索结果 DOM 渲染完成后再恢复滚动，避免内容尚未铺开
    await nextTick()
    await new Promise(requestAnimationFrame)
    window.scrollTo({ top: state.scrollY, behavior: 'instant' })
  }
})
</script>

<style scoped>
/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}
.stat-card {
  padding: 18px 20px;
  background: var(--card-background);
  border: 1px solid var(--card-border);
  border-radius: var(--card-border-radius);
  box-shadow: var(--shadow-l1);
  transition: all 0.2s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-l2);
  border-color: rgba(var(--accent-color-rgb), 0.3);
}
.stat-label {
  font-size: 0.82rem;
  color: var(--card-text-color-secondary);
  margin-bottom: 6px;
  font-weight: 500;
}
.stat-value {
  font-size: 1.6rem;
  font-weight: 800;
  color: var(--card-text-color-main);
  line-height: 1.2;
}
.search-keyword-stat {
  color: var(--accent-color);
  font-size: 1.2rem;
  word-break: break-all;
  line-height: 1.3;
}

/* 大卡片容器 */
.card {
  background: var(--card-background);
  border: 1px solid var(--card-border);
  border-radius: var(--card-border-radius);
  box-shadow: var(--shadow-l1);
  overflow: hidden;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 24px;
  border-bottom: 1px solid var(--card-separator-color);
  flex-wrap: wrap;
}
.card-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--card-text-color-main);
  margin: 0;
}
.card-title svg { color: var(--accent-color); }

/* 通用按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 600;
  transition: all 0.2s ease;
}
.btn-sm {
  padding: 6px 12px;
  font-size: 0.82rem;
}
.btn-primary {
  background: var(--accent-color);
  color: var(--accent-color-text);
}
.btn-primary:hover {
  background: var(--accent-color-darker);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(var(--accent-color-rgb), 0.25);
}

/* 列表容器 */
.article-cards-container {
  padding: 0;
  width: 100%;
}
.search-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 48px 24px;
  color: var(--card-text-color-tertiary);
  font-size: 0.95rem;
}
.animate-spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.error-inline {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  padding: 48px 24px;
}
.error-inline .error-message {
  color: var(--danger-color);
  margin: 0;
  font-weight: 500;
}

/* 简洁搜索结果项（无封面无分类） */
.article-cards-container {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 4px 0;
}
.article-list-item {
  display: block;
  text-decoration: none;
  color: inherit;
  padding: 20px 24px;
  border-bottom: 1px solid var(--card-separator-color);
  transition: background 0.2s ease, padding 0.2s ease;
}
.article-list-item:last-child {
  border-bottom: none;
}
.article-list-item:hover {
  background: rgba(var(--accent-color-rgb), 0.05);
}
.article-list-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}
.article-list-title {
  font-size: 1.08rem;
  font-weight: 700;
  line-height: 1.55;
  color: var(--card-text-color-main);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-word;
  transition: color 0.15s ease;
}
.article-list-item:hover .article-list-title {
  color: var(--accent-color);
}
.article-list-summary {
  font-size: 0.9rem;
  line-height: 1.75;
  color: var(--card-text-color-secondary);
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-word;
}
.article-list-summary--empty {
  color: var(--card-text-color-tertiary);
  opacity: 0.55;
  font-style: italic;
}

/* 搜索关键词高亮（安全，纯文本分片后 mark 标签） */
.search-hl {
  background: linear-gradient(180deg, transparent 60%, rgba(255, 212, 59, 0.55) 40%);
  color: inherit;
  font-weight: 700;
  padding: 0 2px;
  border-radius: 2px;
}
[data-scheme="dark"] .search-hl {
  background: linear-gradient(180deg, transparent 58%, rgba(255, 193, 7, 0.42) 42%);
  color: #fff7cc;
}

/* 空状态 */
.list-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 56px 24px;
  color: var(--card-text-color-tertiary);
  text-align: center;
}
.list-empty svg {
  opacity: 0.5;
  margin-bottom: 14px;
}
.list-empty-title {
  margin: 0 0 4px;
  font-size: 0.95rem;
  color: var(--card-text-color-secondary);
  font-weight: 600;
}
.list-empty-hint {
  margin: 0;
  font-size: 0.85rem;
  opacity: 0.7;
}

/* card-body 分页区 */
.card-body {
  padding: 14px 24px;
  border-top: 1px solid var(--card-separator-color);
}
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}
.pagination-info {
  color: var(--card-text-color-secondary);
  font-size: 0.9rem;
  font-weight: 500;
}
.pagination-info span {
  color: var(--accent-color);
  font-weight: 700;
}
.pagination-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}
.pagination-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 7px 14px;
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  background: var(--card-background);
  color: var(--card-text-color-secondary);
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.18s ease;
}
.pagination-btn:hover:not(:disabled) {
  border-color: var(--accent-color);
  color: var(--accent-color);
  background: rgba(var(--accent-color-rgb), 0.06);
  transform: translateY(-1px);
}
.pagination-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.pagination-btn.active {
  background: var(--accent-color);
  border-color: var(--accent-color);
  color: var(--accent-color-text);
}
.pagination-btn.active:hover {
  transform: none;
}

@media (max-width: 768px) {
  .article-cards-container {
    gap: 16px;
  }
  .card-header {
    padding: 16px 20px;
  }
  .pagination {
    justify-content: center;
  }
}
@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 10px;
  }
  .stat-card {
    padding: 14px 16px;
  }
  .stat-value {
    font-size: 1.3rem;
  }
  .search-keyword-stat {
    font-size: 1rem;
  }
}
</style>
