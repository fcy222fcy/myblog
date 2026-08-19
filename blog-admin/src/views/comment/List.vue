<template>
  <div>
    <!-- 统计卡片（加载态 + 真实态） -->
    <SkeletonLoader v-if="loading" type="stats" :count="4" />
    <div v-else class="comment-stats-grid">
      <div class="stat-card">
        <div class="stat-label">全部评论</div>
      <div class="stat-value">{{ stats.total }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">待审核</div>
        <div class="stat-value text-warning">{{ stats.pending }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">已通过</div>
        <div class="stat-value text-success">{{ stats.approved }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">已拒绝</div>
        <div class="stat-value text-danger">{{ stats.rejected }}</div>
      </div>
    </div>

    <div class="comment-filter-bar">
      <CustomSelect class="status-select" v-model="statusFilter" :options="statusOptions" @change="onFilterChange" />
      <div class="search-box">
        <span class="search-box-icon">⌕</span>
        <input type="text" v-model="keyword" placeholder="搜索评论..." @input="onKeywordInput">
      </div>
    </div>

    <div class="comment-cards">
      <!-- 骨架屏加载态 -->
      <SkeletonLoader v-if="loading" type="comment" :count="4" />
      <!-- 真实列表 -->
      <template v-else>
      <div v-for="comment in filteredComments" :key="comment.id" class="comment-card">
        <div class="comment-card-header">
          <div class="comment-card-user">
            <div class="comment-avatar" :style="getAvatarStyle(comment.status)">
              {{ (comment.nickname || '匿名').charAt(0) }}
            </div>
            <div class="comment-user-info">
              <div class="comment-user-name">{{ comment.nickname || '匿名用户' }}</div>
              <div class="comment-user-email">{{ comment.email || '' }}</div>
            </div>
          </div>
          <span class="status-badge" :class="'status-' + comment.status">{{ statusText(comment.status) }}</span>
        </div>
        <div class="comment-card-body">
          <p class="comment-content">{{ comment.content }}</p>
          <div class="comment-meta">
            <span>评论于
              <a
                v-if="comment.article && comment.article.slug"
                class="text-accent"
                :href="getArticleUrl(comment.article.slug)"
                target="_blank"
                rel="noopener noreferrer"
                @click.stop
              >{{ comment.article.title }}</a>
              <span v-else-if="comment.article && comment.article.title" class="text-accent">{{ comment.article.title }}</span>
              <span v-else class="text-accent">未知文章</span>
            </span>
            <span>·</span>
            <span>{{ formatDate(comment.created_at) }}</span>
          </div>
        </div>
        <div class="comment-card-footer">
          <button v-if="comment.status !== 'approved'" class="action-btn btn-edit btn-sm" @click="handleStatus(comment.id, 'approved')">通过</button>
          <button v-if="comment.status !== 'rejected'" class="action-btn btn-hide btn-sm" @click="handleStatus(comment.id, 'rejected')">拒绝</button>
          <button class="action-btn btn-delete btn-sm" @click="handleDelete(comment.id)">删除</button>
        </div>
      </div>
      </template>
    </div>

    <div v-if="total > pageSize" class="card-body comment-pagination">
      <button class="pagination-btn" :disabled="page <= 1" @click="page--; loadComments()">上一页</button>
      <span>第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页</span>
      <button class="pagination-btn" :disabled="page * pageSize >= total" @click="page++; loadComments()">下一页</button>
    </div>

    <div v-if="filteredComments.length === 0 && !loading" class="card">
      <div class="card-body empty-state-sm">暂无评论</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCommentList, updateCommentStatus, deleteComment } from '../../api/comment'
import SkeletonLoader from '../../components/common/SkeletonLoader.vue'
import CustomSelect from '../../components/common/CustomSelect.vue'

const comments = ref([])
const statusFilter = ref('')
const keyword = ref('')
const loading = ref(true)
const page = ref(1)
const pageSize = 20
const total = ref(0)
const stats = ref({ total: 0, pending: 0, approved: 0, rejected: 0 })
let keywordTimer

const statusOptions = [
  { value: '', label: '全部状态' },
  { value: 'pending', label: '待审核' },
  { value: 'approved', label: '已通过' },
  { value: 'rejected', label: '已拒绝' }
]

const formatDate = (d) => d ? new Date(d).toLocaleDateString('zh-CN') : ''
const statusText = (s) => ({ pending: '待审核', approved: '已通过', rejected: '已拒绝' }[s] || s)

const getArticleUrl = (slug) => {
  const hashPath = `/#/post/${slug}#comment-section`
  if (typeof window !== 'undefined' && window.location) {
    const host = window.location.hostname
    if (host === 'localhost' || host === '127.0.0.1') {
      return `${window.location.protocol}//${host}:5173${hashPath}`
    }
  }
  return hashPath
}

const getAvatarStyle = (status) => {
  const colors = {
    pending: 'background: rgba(var(--accent-color-rgb), 0.1); color: var(--accent-color);',
    approved: 'background: rgba(var(--success-color-rgb), 0.1); color: var(--success-color);',
    rejected: 'background: rgba(var(--danger-color-rgb), 0.1); color: var(--danger-color);'
  }
  return colors[status] || colors.pending
}

const filteredComments = computed(() => {
  return comments.value
})

const loadComments = async () => {
  try {
    const params = { page: page.value, page_size: pageSize }
    if (statusFilter.value) params.status = statusFilter.value
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    const [res, allRes, pendingRes, approvedRes, rejectedRes] = await Promise.all([
      getCommentList(params),
      getCommentList({ page: 1, page_size: 1, keyword: keyword.value.trim() }),
      getCommentList({ page: 1, page_size: 1, status: 'pending', keyword: keyword.value.trim() }),
      getCommentList({ page: 1, page_size: 1, status: 'approved', keyword: keyword.value.trim() }),
      getCommentList({ page: 1, page_size: 1, status: 'rejected', keyword: keyword.value.trim() })
    ])
    comments.value = res.data?.list || []
    total.value = res.data?.total || 0
    stats.value = {
      total: allRes.data?.total || 0,
      pending: pendingRes.data?.total || 0,
      approved: approvedRes.data?.total || 0,
      rejected: rejectedRes.data?.total || 0
    }
  } catch (e) { console.error(e) }
  loading.value = false
}

const onFilterChange = () => { page.value = 1; loadComments() }
const onKeywordInput = () => {
  clearTimeout(keywordTimer)
  keywordTimer = setTimeout(() => { page.value = 1; loadComments() }, 250)
}

const handleStatus = async (id, status) => {
  try {
    await updateCommentStatus(id, status)
    ElMessage.success('操作成功')
    loadComments()
  } catch (e) { console.error(e) }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '确认删除', { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' })
    await deleteComment(id)
    ElMessage.success('删除成功')
    loadComments()
  } catch (e) { if (e !== 'cancel') console.error(e) }
}

onMounted(loadComments)
</script>

<style scoped>
.comment-stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  margin-bottom: 24px;
  width: 100%;
}
.comment-stats-grid .stat-card {
  padding: 18px 20px;
}
.comment-stats-grid .stat-label {
  font-size: 13px;
  margin-bottom: 6px;
}
.comment-stats-grid .stat-value {
  font-size: 24px;
}
@media (max-width: 900px) {
  .comment-stats-grid {
    grid-template-columns: repeat(2, 1fr);
    max-width: 100%;
  }
}
@media (max-width: 520px) {
  .comment-stats-grid {
    grid-template-columns: 1fr;
  }
}

.comment-filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  align-items: stretch;
  justify-content: flex-end;
  flex-wrap: nowrap;
}
.comment-filter-bar .status-select {
  flex: 0 0 auto;
  width: 150px;
}
.comment-filter-bar .status-select :deep(.custom-select) {
  min-height: 40px;
  height: 40px;
  padding-top: 0;
  padding-bottom: 0;
  border-radius: var(--card-border-radius);
}
.comment-filter-bar .search-box {
  flex: 0 0 auto;
  width: 260px;
}
@media (max-width: 768px) {
  .comment-filter-bar {
    flex-wrap: wrap;
    justify-content: stretch;
  }
  .comment-filter-bar .status-select,
  .comment-filter-bar .search-box {
    flex: 1 1 100%;
    width: 100%;
  }
}

.comment-user-info { display: flex; flex-direction: column; }
.comment-user-name { font-size: 14px; font-weight: 600; color: var(--card-text-color-main); }
.comment-user-email { font-size: 12px; color: var(--card-text-color-tertiary); }
.comment-meta { font-size: 13px; color: var(--card-text-color-tertiary); display: flex; gap: 6px; }
.btn-hide { background: rgba(var(--warning-color-rgb), 0.08); color: var(--warning-color); }
.btn-hide:hover { background: var(--warning-color); color: white; }
</style>
