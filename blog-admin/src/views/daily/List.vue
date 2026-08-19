<template>
  <div>
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-label">问题总数</div>
        <div class="stat-value">{{ total }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">当前页已发布</div>
        <div class="stat-value">{{ questions.filter(q => q.status === 1).length }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">当前页浏览量</div>
        <div class="stat-value">{{ questions.reduce((s, q) => s + (q.view_count || 0), 0) }}</div>
      </div>
      <div class="stat-card">
        <div class="stat-label">当前页点赞数</div>
        <div class="stat-value">{{ questions.reduce((s, q) => s + (q.like_count || 0), 0) }}</div>
      </div>
    </div>

    <div class="card" style="margin-bottom: 20px;">
      <div class="card-header">
        <div class="card-title">问题管理</div>
        <div class="filter-group">
          <div class="search-box">
            <span class="search-box-icon">⌕</span>
            <input type="text" v-model="keyword" placeholder="搜索问题..." @input="onKeywordInput">
          </div>
          <div class="custom-select-wrapper" v-click-outside="closeStatusDropdown">
            <div class="custom-select" @click="toggleStatusDropdown">
              <span class="select-value">{{ selectedStatusLabel }}</span>
              <span class="select-arrow">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
              </span>
            </div>
            <div class="custom-dropdown" v-if="isStatusOpen">
              <div class="dropdown-item" :class="{ active: statusFilter === '' }" @click="selectStatus('')">全部状态</div>
              <div class="dropdown-item" :class="{ active: statusFilter === '1' }" @click="selectStatus('1')">已发布</div>
              <div class="dropdown-item" :class="{ active: statusFilter === '0' }" @click="selectStatus('0')">禁用</div>
              <div class="dropdown-item" :class="{ active: statusFilter === '2' }" @click="selectStatus('2')">定时发布</div>
            </div>
          </div>
          <button class="btn btn-primary" @click="showModal = true; resetForm()">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
            <span>新建问题</span>
          </button>
        </div>
      </div>
      <div class="card-body">
        <table class="table">
          <thead>
            <tr>
              <th>问题</th>
              <th>显示日期</th>
              <th>状态</th>
              <th>浏览/点赞</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="q in filteredQuestions" :key="q.id">
              <td style="max-width: 300px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ q.question }}</td>
              <td>{{ q.date }}</td>
              <td>
                <span class="status-badge" :class="q.status === 1 ? 'status-published' : (q.status === 2 ? 'status-scheduled' : 'status-draft')">
                  {{ statusText(q.status) }}
                </span>
              </td>
              <td>{{ q.view_count || 0 }} / {{ q.like_count || 0 }}</td>
              <td>
                <div style="display: flex; gap: 6px;">
                  <button class="action-btn btn-edit btn-sm" @click="editQuestion(q)">编辑</button>
                  <button v-if="q.status !== 2" class="action-btn btn-sm" :class="q.status === 1 ? 'btn-hide' : 'btn-edit'" @click="toggleStatus(q)">
                    {{ q.status === 1 ? '禁用' : '发布' }}
                  </button>
                  <button class="action-btn btn-delete btn-sm" @click="handleDelete(q.id)">删除</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredQuestions.length === 0">
              <td colspan="5" style="text-align: center; color: var(--card-text-color-tertiary); padding: 40px;">暂无问题</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="pagination" v-if="total > pageSize">
      <button class="btn btn-secondary btn-sm" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
      <span>第 {{ page }} / {{ totalPages }} 页，共 {{ total }} 条</span>
      <button class="btn btn-secondary btn-sm" :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
    </div>

    <div class="modal-overlay" :class="{ active: showModal }" @click.self="showModal = false">
      <div class="modal" style="max-width: 900px;">
        <div class="modal-header">
          <h3 class="modal-title">{{ editingId ? '编辑问题' : '新建问题' }}</h3>
          <button class="modal-close" @click="showModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">日期</label>
            <input type="date" class="form-input" v-model="form.date">
          </div>
          <div class="form-group">
            <label class="form-label">状态</label>
            <select class="form-input" v-model.number="form.status">
              <option :value="1">已发布</option>
              <option :value="0">禁用</option>
              <option :value="2">定时发布</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">问题 <span class="required">*</span></label>
            <textarea class="form-textarea" v-model="form.question" placeholder="输入问题..." rows="3"></textarea>
          </div>
          <div class="form-group" style="margin-bottom: 0;">
            <label class="form-label">答案（支持 Markdown）</label>
            <MdEditor
              v-model="form.answer"
              :theme="editorTheme"
              class="md-editor"
              :preview="true"
              previewTheme="github"
              @onUploadImg="onUploadImg"
              style="height: 400px;"
            />
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn btn-primary" @click="handleSave">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { getDailyQuestionList, createDailyQuestion, updateDailyQuestion, deleteDailyQuestion, updateDailyQuestionStatus } from '../../api/daily'
import { uploadFile, MEDIA_CATEGORIES } from '../../api/media'

const questions = ref([])
const page = ref(1)
const pageSize = 20
const total = ref(0)
const keyword = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const editingId = ref(null)
const form = ref({ question: '', answer: '', date: '', status: 1 })
const editorTheme = ref(localStorage.getItem('scheme') === 'dark' ? 'dark' : 'light')

// 自定义下拉菜单逻辑
const isStatusOpen = ref(false)

const selectedStatusLabel = computed(() => {
  const map = { '': '全部状态', '1': '已发布', '0': '禁用', '2': '定时发布' }
  return map[statusFilter.value] || '全部状态'
})
const statusText = (status) => ({ 0: '禁用', 1: '已发布', 2: '定时发布' }[status] || '未知')

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
  loadQuestions()
}

const filteredQuestions = computed(() => questions.value)

const resetForm = () => { editingId.value = null; form.value = { question: '', answer: '', date: '', status: 1 } }

const editQuestion = (q) => {
  editingId.value = q.id
  form.value = { question: q.question, answer: q.answer, date: q.date, status: q.status }
  showModal.value = true
}

const loadQuestions = async () => {
  try {
    const params = { page: page.value, page_size: pageSize }
    if (statusFilter.value !== '') params.status = Number(statusFilter.value)
    if (keyword.value.trim()) params.keyword = keyword.value.trim()
    const res = await getDailyQuestionList(params)
    questions.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) { console.error(e) }
}

const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize)))
const changePage = (nextPage) => {
  if (nextPage < 1 || nextPage > totalPages.value) return
  page.value = nextPage
  loadQuestions()
}

const onKeywordInput = () => {
  page.value = 1
  loadQuestions()
}

// 图片上传处理
const onUploadImg = async (files, callback) => {
  const res = []
  for (const file of files) {
    try {
      const uploadRes = await uploadFile(file, MEDIA_CATEGORIES.DAILY)
      if (uploadRes.code === 0) {
        res.push(uploadRes.data.url)
      } else {
        ElMessage.error('图片上传失败')
      }
    } catch (e) {
      console.error('上传失败:', e)
      ElMessage.error('图片上传失败')
    }
  }
  callback(res)
}

const handleSave = async () => {
  if (!form.value.question) { ElMessage.warning('请输入问题'); return }
  try {
    if (editingId.value) { await updateDailyQuestion(editingId.value, form.value) }
    else { await createDailyQuestion(form.value) }
    ElMessage.success('保存成功')
    showModal.value = false
    loadQuestions()
  } catch (e) { console.error(e) }
}

const toggleStatus = async (q) => {
  try {
    await updateDailyQuestionStatus(q.id, q.status === 1 ? 0 : 1)
    ElMessage.success('操作成功')
    loadQuestions()
  } catch (e) { console.error(e) }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个问题吗？', '确认删除', { confirmButtonText: '删除', cancelButtonText: '取消', type: 'warning' })
    await deleteDailyQuestion(id)
    ElMessage.success('删除成功')
    loadQuestions()
  } catch (e) { if (e !== 'cancel') console.error(e) }
}

onMounted(loadQuestions)

// 自定义指令：点击外部关闭下拉菜单
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
  flex-wrap: wrap;
  align-items: center;
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
  min-width: 200px;
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

.btn-hide { background: rgba(245, 158, 11, 0.08); color: var(--warning-color); }
.btn-hide:hover { background: var(--warning-color); color: white; }
.md-editor { border-radius: 8px; overflow: hidden; }
</style>
