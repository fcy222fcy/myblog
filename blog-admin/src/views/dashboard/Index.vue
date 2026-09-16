<template>
  <div>
    <!-- 统计接口失败提示：避免失败时静默显示 0，让人误以为博客是空的 -->
    <div v-if="statsError" class="dash-alert">
      <span>统计数据加载失败，下方数字不可信</span>
      <button class="dash-alert-retry" @click="loadAll">重新加载</button>
    </div>

    <!-- 统计卡片（加载态 + 真实态） -->
    <SkeletonLoader v-if="loading" type="stats" :count="4" />
    <div v-else class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">文章总数</div>
        <div class="stat-value">{{ formatCount(stats.article_count) }}</div>
        <div class="stat-sub">已发布 {{ formatCount(stats.published_count) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">今日浏览</div>
        <div class="stat-value">{{ formatCount(stats.today_views) }}</div>
        <div class="stat-sub">累计 {{ formatCount(stats.total_views) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">评论总数</div>
        <div class="stat-value">{{ formatCount(stats.comment_count) }}</div>
        <div class="stat-sub" :class="{ 'is-alert': pendingCount > 0 }">
          待审核 {{ formatCount(stats.pending_count) }}
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-label">待发布</div>
        <div class="stat-value">{{ formatCount(pendingPublishCount) }}</div>
        <div class="stat-sub">
          草稿 {{ formatCount(stats.draft_count) }} · 定时 {{ formatCount(stats.scheduled_count) }}
        </div>
      </div>
    </div>

    <!-- 最近文章 -->
    <div class="dash-panel">
      <div class="dash-panel-head">
        <div class="dash-panel-title">最近文章</div>
        <div class="dash-panel-more" @click="$router.push('/articles')">全部文章</div>
      </div>
      <SkeletonLoader v-if="loading" type="table" :count="5" />
      <div v-else-if="recentArticles.length" class="recent-list">
        <div
          v-for="item in recentArticles"
          :key="item.id"
          class="recent-item"
          @click="$router.push(`/articles/edit/${item.id}`)"
        >
          <div class="recent-main">
            <div class="recent-title">{{ item.title }}</div>
            <div class="recent-meta">{{ item.created_at }} · {{ formatCount(item.view_count) }} 次浏览</div>
          </div>
          <div class="recent-arrow">→</div>
        </div>
      </div>
      <div v-else class="dash-panel-empty">
        {{ recentError ? '最近文章加载失败' : '还没有文章，去写一篇吧' }}
      </div>
    </div>

    <div class="dashboard-grid">
      <div class="dashboard-card" @click="$router.push('/articles')">
        <div class="dashboard-card-icon">📝</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">文章管理</div>
          <div class="dashboard-card-desc">管理博客文章，发布、编辑、删除</div>
          <div class="dashboard-card-meta">
            共 {{ formatCount(stats.article_count) }} 篇 · 已发布 {{ formatCount(stats.published_count) }}
          </div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>

      <div class="dashboard-card" @click="$router.push('/categories')">
        <div class="dashboard-card-icon">📂</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">分类管理</div>
          <div class="dashboard-card-desc">管理文章分类，建立内容结构</div>
          <div class="dashboard-card-meta">
            {{ categoryCount === null ? '分类的增删改查' : `共 ${categoryCount} 个分类` }}
          </div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>

      <div class="dashboard-card" @click="$router.push('/tags')">
        <div class="dashboard-card-icon">🏷️</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">标签管理</div>
          <div class="dashboard-card-desc">管理文章标签，方便内容检索</div>
          <div class="dashboard-card-meta">
            {{ tagCount === null ? '标签的增删改查' : `共 ${tagCount} 个标签` }}
          </div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>

      <div class="dashboard-card" @click="$router.push('/comments')">
        <div class="dashboard-card-icon">💬</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">评论管理</div>
          <div class="dashboard-card-desc">审核和管理读者评论</div>
          <div class="dashboard-card-meta">
            待审核 {{ formatCount(stats.pending_count) }} · 共 {{ formatCount(stats.comment_count) }} 条
          </div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>

      <div class="dashboard-card" @click="$router.push('/daily-question')">
        <div class="dashboard-card-icon">💡</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">每日一问</div>
          <div class="dashboard-card-desc">管理每日问答内容</div>
          <div class="dashboard-card-meta">
            {{ dailyCount === null ? '问题与回答管理' : `共 ${dailyCount} 条` }}
          </div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>

      <div class="dashboard-card" @click="$router.push('/about')">
        <div class="dashboard-card-icon">👤</div>
        <div class="dashboard-card-info">
          <div class="dashboard-card-title">关于我</div>
          <div class="dashboard-card-desc">编辑个人简介和社交信息</div>
          <div class="dashboard-card-meta">头像、简介、社交链接</div>
        </div>
        <div class="dashboard-card-arrow">→</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import request from '../../api/request'
import { getCategoryList } from '../../api/category'
import { getTagList } from '../../api/tag'
import { getDailyQuestionList } from '../../api/daily'
import SkeletonLoader from '../../components/common/SkeletonLoader.vue'

const stats = ref({})
const recentArticles = ref([])
const categoryCount = ref(null)
const tagCount = ref(null)
const dailyCount = ref(null)

const loading = ref(true)
const statsError = ref(false)
const recentError = ref(false)

// 未取到数据时显示破折号，而不是会被误读成「真的是 0」的数字
const formatCount = (value) => {
  if (value === undefined || value === null) return '—'
  const num = Number(value)
  if (Number.isNaN(num)) return '—'
  return num.toLocaleString('zh-CN')
}

const pendingCount = computed(() => Number(stats.value.pending_count || 0))

const pendingPublishCount = computed(() => {
  if (stats.value.draft_count === undefined && stats.value.scheduled_count === undefined) return undefined
  return Number(stats.value.draft_count || 0) + Number(stats.value.scheduled_count || 0)
})

const loadStats = async () => {
  try {
    const res = await request.get('/admin/dashboard/stats')
    stats.value = res.data || {}
  } catch (e) {
    statsError.value = true
  }
}

const loadRecentArticles = async () => {
  try {
    const res = await request.get('/admin/dashboard/recent-articles', { params: { limit: 5 } })
    recentArticles.value = res.data || []
  } catch (e) {
    recentError.value = true
  }
}

// 分类 / 标签 / 每日一问没有专用的计数接口，用各自的列表接口取长度或 total
const loadCounts = async () => {
  const [categories, tags, daily] = await Promise.allSettled([
    getCategoryList(),
    getTagList(),
    getDailyQuestionList({ page: 1, page_size: 1 })
  ])
  if (categories.status === 'fulfilled') categoryCount.value = categories.value.data?.length ?? null
  if (tags.status === 'fulfilled') tagCount.value = tags.value.data?.length ?? null
  if (daily.status === 'fulfilled') dailyCount.value = daily.value.data?.total ?? null
}

const loadAll = async () => {
  loading.value = true
  statsError.value = false
  recentError.value = false
  await Promise.all([loadStats(), loadRecentArticles(), loadCounts()])
  loading.value = false
}

onMounted(loadAll)
</script>
