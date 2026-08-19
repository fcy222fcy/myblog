<template>
  <div>
    <!-- 统计卡片（加载态 + 真实态） -->
    <SkeletonLoader v-if="loading" type="stats" :count="4" />
    <div v-else class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">文章总数</div>
        <div class="stat-value">{{ allTotal }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">已发布</div>
        <div class="stat-value">{{ stats.published || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">草稿</div>
        <div class="stat-value">{{ stats.draft || 0 }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">总浏览量</div>
        <div class="stat-value">{{ stats.totalViews || 0 }}</div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <div class="card-title">文章列表</div>
        <div class="filter-group">
          <div class="search-box">
            <span class="search-box-icon">⌕</span>
            <input type="text" v-model="keyword" placeholder="搜索文章..." @input="onKeywordInput">
          </div>
          <CustomSelect class="category-select" v-model="categoryFilter" :options="categoryOptions" @change="onCategoryChange" />
          <div class="custom-select-wrapper" v-click-outside="closeStatusDropdown">
            <div class="custom-select" @click="toggleStatusDropdown">
              <span class="select-value">{{ selectedStatusLabel }}</span>
              <span class="select-arrow">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
              </span>
            </div>
            <div class="custom-dropdown" v-if="isStatusOpen">
              <div class="dropdown-item" :class="{ active: statusFilter === '' }" @click="selectStatus('')">全部状态</div>
              <div class="dropdown-item" :class="{ active: statusFilter === 'published' }" @click="selectStatus('published')">已发布</div>
              <div class="dropdown-item" :class="{ active: statusFilter === 'draft' }" @click="selectStatus('draft')">草稿</div>
              <div class="dropdown-item" :class="{ active: statusFilter === 'scheduled' }" @click="selectStatus('scheduled')">定时发布</div>
            </div>
          </div>
          <button class="btn btn-primary btn-new-article" @click="$router.push('/articles/edit')">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            <span>新建文章</span>
          </button>
        </div>
      </div>

      <div class="article-cards-container">
        <!-- 骨架屏加载态 -->
        <SkeletonLoader v-if="loading" type="article" :count="5" />
        <!-- 真实列表 -->
        <template v-else>
        <div v-for="article in articles" :key="article.id" class="article-card-item" @click="$router.push('/articles/edit/' + article.id)">
          <!-- 封面缩略图 -->
          <div class="article-card-cover" :class="{ 'no-cover': !article.cover }">
            <img v-if="article.cover" :src="article.cover" alt="封面" class="article-cover-img">
            <div v-else class="article-cover-placeholder" aria-label="暂无封面">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" role="img"><title>暂无封面</title><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline></svg>
            </div>
          </div>
          <!-- 文章信息 -->
          <div class="article-card-info">
            <div class="article-card-header">
              <h3 class="article-card-title">{{ article.title }}</h3>
              <span class="status-badge" :class="'status-' + article.status">
                {{ statusLabel(article.status) }}
              </span>
            </div>
            <p class="article-card-summary">{{ article.summary || '暂无摘要' }}</p>
            <div class="article-card-footer">
              <div class="article-card-meta">
                <span>{{ article.category?.name || '未分类' }}</span>
                <span>·</span>
                <span>{{ formatDate(article.created_at) }}</span>
                <span>·</span>
                <span>{{ article.view_count || 0 }} 次浏览</span>
              </div>
              <div class="article-card-actions" @click.stop>
                <button class="action-btn btn-edit btn-sm" @click="$router.push('/articles/edit/' + article.id)">编辑</button>
                <button class="action-btn btn-delete btn-sm" @click="handleDelete(article.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
        <div v-if="articles.length === 0 && !loading" style="text-align: center; padding: 40px; color: var(--card-text-color-tertiary);">
          暂无文章
        </div>
        </template>
      </div>

      <div class="card-body">
        <div class="pagination">
          <div class="pagination-info">
            共 <span>{{ total }}</span> 篇文章
          </div>
          <div class="pagination-buttons">
            <button class="pagination-btn" :disabled="page <= 1" @click="page--; loadArticles()">上一页</button>
            <button class="pagination-btn active">{{ page }}</button>
            <button class="pagination-btn" :disabled="page * 10 >= total" @click="page++; loadArticles()">下一页</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { getArticleList, deleteArticle } from '../../api/article'
import request from '../../api/request'
import { getCategoryList } from '../../api/category'
import { ElMessage, ElMessageBox } from 'element-plus'
import SkeletonLoader from '../../components/common/SkeletonLoader.vue'
import CustomSelect from '../../components/common/CustomSelect.vue'

const articles = ref([])
const categories = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const allTotal = ref(0)
const keyword = ref('')
const categoryFilter = ref('')
const statusFilter = ref('')
const stats = ref({ published: 0, draft: 0, totalViews: 0 })
let keywordTimer

const categoryOptions = computed(() => {
  return [{ value: '', label: '全部分类' }, ...categories.value.map(c => ({ value: c.id, label: c.name }))]
})

const isStatusOpen = ref(false)

const selectedStatusLabel = computed(() => {
  const map = { '': '全部状态', 'published': '已发布', 'draft': '草稿', 'scheduled': '定时发布' }
  return map[statusFilter.value] || '全部状态'
})

const toggleStatusDropdown = () => {
  isStatusOpen.value = !isStatusOpen.value
}

const closeStatusDropdown = () => {
  isStatusOpen.value = false
}

const selectStatus = (value) => {
  statusFilter.value = value
  isStatusOpen.value = false
  page.value = 1
  loadArticles()
}

const onCategoryChange = () => { page.value = 1; loadArticles() }
const onKeywordInput = () => {
  clearTimeout(keywordTimer)
  keywordTimer = setTimeout(() => { page.value = 1; loadArticles() }, 250)
}

const formatDate = (d) => d ? d.split('T')[0] : ''
const statusLabel = (s) => ({ published: '已发布', draft: '草稿', scheduled: '定时发布' }[s] || s)

const loadStats = async () => {
  try {
    const [allRes, publishedRes, draftRes, scheduledRes, dashboardRes] = await Promise.all([
      getArticleList({ page: 1, page_size: 1 }),
      getArticleList({ page: 1, page_size: 1, status: 'published' }),
      getArticleList({ page: 1, page_size: 1, status: 'draft' }),
      getArticleList({ page: 1, page_size: 1, status: 'scheduled' }),
      request.get('/admin/dashboard/stats')
    ])
    allTotal.value = allRes.data?.total || 0
    stats.value.published = publishedRes.data?.total || 0
    stats.value.draft = draftRes.data?.total || 0
    stats.value.scheduled = scheduledRes.data?.total || 0
    stats.value.totalViews = dashboardRes.data?.total_views || 0
  } catch (e) { console.error(e) }
}

const loadArticles = async () => {
  try {
    const params = { page: page.value, page_size: 10 }
    if (keyword.value) params.keyword = keyword.value
    if (categoryFilter.value) params.category_id = categoryFilter.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await getArticleList(params)
    articles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) { console.error(e) }
}

const loadCategories = async () => {
  try {
    const res = await getCategoryList()
    categories.value = res.data || []
  } catch (e) { console.error(e) }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？此操作无法撤销。', '确认删除', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await deleteArticle(id)
    ElMessage.success('删除成功')
    loadArticles()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

onMounted(async () => {
  loading.value = true
  loadStats()
  loadCategories()
  await loadArticles()
  loading.value = false
})

const vClickOutside = {
  mounted(el, binding) {
    el._clickOutside = (event) => {
      if (!el.contains(event.target)) {
        binding.value()
      }
    }
    document.addEventListener('click', el._clickOutside)
  },
  unmounted(el) {
    document.removeEventListener('click', el._clickOutside)
  }
}
</script>

<style scoped>
.filter-group {
  display: flex;
  gap: 12px;
  flex-wrap: nowrap;
  align-items: center;
}
.filter-group .search-box {
  flex: 0 0 auto;
  width: 260px;
}
.filter-group .category-select {
  flex: 0 0 auto;
  width: 150px;
}
.filter-group .category-select :deep(.custom-select) {
  min-height: 40px;
  height: 40px;
  padding-top: 0;
  padding-bottom: 0;
  border-radius: var(--card-border-radius);
}
.filter-group .custom-select-wrapper {
  position: relative;
  display: inline-flex;
  flex: 0 0 auto;
  width: 130px;
}
.filter-group .custom-select-wrapper .custom-select {
  width: 100%;
}
.filter-group .btn-new-article {
  flex: 0 0 auto;
  height: 40px;
  white-space: nowrap;
}

.custom-select-wrapper {
  position: relative;
  display: inline-flex;
}

.custom-select {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 32px 8px 12px;
  border: 1px solid var(--card-separator-color);
  border-radius: var(--card-border-radius);
  font-size: 13px;
  font-weight: 500;
  color: var(--card-text-color-main);
  background: var(--card-background);
  cursor: pointer;
  transition: all 0.15s ease;
  min-width: 100px;
  height: 40px;
}

.custom-select:hover {
  border-color: var(--accent-color);
}

.select-value {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.select-arrow {
  position: absolute;
  right: 10px;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: var(--card-text-color-tertiary);
  display: flex;
  align-items: center;
}

.custom-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  min-width: 140px;
  background: var(--card-background);
  border: 1px solid var(--card-separator-color);
  border-radius: var(--card-border-radius);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  z-index: 100;
}

.dropdown-item {
  padding: 10px 12px;
  font-size: 13px;
  color: var(--card-text-color-main);
  cursor: pointer;
  transition: all 0.15s ease;
}

.dropdown-item:hover {
  background: rgba(var(--accent-color-rgb), 0.06);
  color: var(--accent-color);
}

.dropdown-item.active {
  background: rgba(var(--accent-color-rgb), 0.1);
  color: var(--accent-color);
  font-weight: 600;
}

@media (max-width: 900px) {
  .filter-group {
    flex-wrap: wrap;
    justify-content: stretch;
  }
  .filter-group .search-box,
  .filter-group .category-select,
  .filter-group .custom-select-wrapper,
  .filter-group .btn-new-article {
    flex: 1 1 100%;
    width: 100%;
  }
}

.article-cards-container {
  padding: 0;
  width: 100%;
}
.article-card-item {
  display: flex;
  gap: 16px;
  padding: 20px 28px;
  border-bottom: 1px solid var(--card-separator-color);
  transition: background 0.15s ease;
  cursor: pointer;
  width: 100%;
}
.article-card-item:hover {
  background: rgba(var(--accent-color-rgb), 0.02);
}
.article-card-item:last-child {
  border-bottom: none;
}

/* 封面缩略图 */
.article-card-cover {
  width: 140px;
  height: 90px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, #667eea20, #764ba220);
}
.article-card-cover.no-cover {
  background: rgba(var(--accent-color-rgb), 0.03);
  border: 1px dashed var(--card-separator-color);
}
.article-cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.article-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--card-text-color-tertiary);
  opacity: 0.5;
}

.article-card-info {
  flex: 1;
  min-width: 0;
}
.article-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}
.article-card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--card-text-color-main);
  cursor: pointer;
  margin: 0;
}
.article-card-title:hover {
  color: var(--accent-color);
}
.article-card-summary {
  font-size: 14px;
  color: var(--card-text-color-secondary);
  margin: 0 0 12px;
  line-height: 1.6;
}
.article-card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.article-card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--card-text-color-tertiary);
}
.article-card-actions {
  display: flex;
  gap: 8px;
}
</style>
