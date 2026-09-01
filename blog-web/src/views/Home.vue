<template>
  <div class="page-view" id="view-home">
    <!-- 每日一问卡片 -->
    <DailyQuestion />

    <!-- 骨架屏加载 -->
    <ArticleListSkeleton v-if="articleStore.loading" />

    <!-- 错误状态 -->
    <ErrorState v-else-if="articleStore.error" :message="articleStore.error" @retry="retryFetch" />

    <!-- 文章卡片列表 -->
    <template v-else>
      <article v-for="article in articleStore.articles" :key="article.id" class="article-card">
        <router-link :to="'/post/' + article.slug" class="card-link">
          <div class="card-pills-row">
            <router-link
              :to="'/category/' + (article.category?.id || 0)"
              class="category-pill"
              :class="'category-' + (article.category?.slug || 'default')"
              :style="categoryStyle(article.category)"
              :title="'查看「' + (article.category?.name || '未分类') + '」分类下的文章'"
              @click.stop
            >{{ article.category?.name || '未分类' }}</router-link>
            <router-link
              v-for="t in (article.tags || [])"
              :key="t.id"
              :to="'/tag/' + t.id"
              class="tag-pill"
              :style="tagStyle(t)"
              :title="'查看「' + t.name + '」标签下的文章'"
              @click.stop
            >#{{ t.name }}</router-link>
          </div>
          <h2>{{ article.title }}</h2>
          <p class="article-summary">{{ displaySummary(article) }}</p>
          <dl class="article-meta">
            <div>
              <dt>日期</dt>
              <dd>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect><line x1="16" x2="16" y1="2" y2="6"></line><line x1="8" x2="8" y1="2" y2="6"></line><line x1="3" x2="21" y1="10" y2="10"></line></svg>
                {{ formatDate(article.created_at) }}
              </dd>
            </div>
            <div>
              <dt>时长</dt>
              <dd>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                {{ article.reading_time || 3 }} 分钟
              </dd>
            </div>
            <div>
              <dt>浏览量</dt>
              <dd>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>
                {{ article.view_count }}
              </dd>
            </div>
            <div>
              <dt>语言</dt>
              <dd>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"></circle><line x1="2" x2="22" y1="12" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>
                中文
              </dd>
            </div>
          </dl>
        </router-link>
      </article>

      <div v-if="articleStore.articles.length === 0" class="empty">没有更多文章了</div>
    </template>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted } from 'vue'
import { useArticleStore } from '../stores/article'
import { saveScrollPosition, takeScrollPosition } from '../composables/useScrollRestore'
import DailyQuestion from '../components/daily/DailyQuestion.vue'
import ArticleListSkeleton from '../components/common/ArticleListSkeleton.vue'
import ErrorState from '../components/common/ErrorState.vue'
import { formatDate } from '../utils/date'
import { categoryStyle } from '../utils/categoryStyle'

const articleStore = useArticleStore()

const retryFetch = async () => {
  await articleStore.fetchArticles({ page: 1, page_size: 20 })
  restoreScrollPosition()
}

// 返回首页时恢复到离开前的位置（如从标签/分类归档页点「返回首页」返回）
const restoreScrollPosition = async () => {
  const y = takeScrollPosition()
  if (y === null) return
  await nextTick()
  requestAnimationFrame(() => {
    window.scrollTo({ top: y, behavior: 'instant' })
  })
}

// 离开首页前记住当前滚动位置，供返回时恢复
onBeforeUnmount(saveScrollPosition)

// 标签颜色：根据 tag.name 字符串 hash 分配稳定色相（不随主题变化）
const tagStyle = (tag) => {
  const name = (tag && tag.name) || ''
  let h = 0
  for (let i = 0; i < name.length; i++) h = (h * 31 + name.charCodeAt(i)) % 360
  return {
    background: `hsl(${h}, 55%, 92%)`,
    color: `hsl(${h}, 60%, 28%)`,
    borderColor: `hsl(${h}, 50%, 82%)`
  }
}

// 摘要展示：优先后端 summary，否则从 Markdown content 提取纯文本截取
const displaySummary = (article) => {
  if (!article) return ''
  if (article.summary && article.summary.trim()) return article.summary
  const md = article.content || ''
  const text = md
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]+)\]\([^)]*\)/g, '$1')
    .replace(/[#>*_`\-\|]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
  if (!text) return ''
  return text.length > 130 ? text.slice(0, 130) + '…' : text
}

onMounted(() => {
  retryFetch()
})
</script>

<style scoped>
.card-pills-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-bottom: 4px;
  line-height: 1;
}
.tag-pill {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 32px;
  padding: 0 14px;
  border-radius: var(--tag-border-radius);
  font-size: 0.88rem;
  font-weight: 600;
  background: rgba(var(--accent-color-rgb), 0.1);
  color: var(--accent-color);
  text-decoration: none;
  transition: background 0.15s ease, transform 0.15s ease;
  user-select: none;
}
.tag-pill:hover {
  background: rgba(var(--accent-color-rgb), 0.2);
  transform: translateY(-1px);
}
.card-pills-row .category-pill {
  text-decoration: none;
  transition: filter 0.15s ease, transform 0.15s ease;
}
.card-pills-row .category-pill:hover {
  filter: brightness(1.08);
  transform: translateY(-1px);
}
.article-summary {
  display: -webkit-box;
  height: 4.95em;
  overflow: hidden;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}
.empty {
  text-align: center;
  padding: 40px;
  color: var(--card-text-color-secondary);
}
</style>
