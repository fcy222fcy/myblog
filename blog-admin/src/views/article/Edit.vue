<template>
  <div class="article-edit-container">
    <!-- 返回按钮 -->
    <button class="back-btn" @click="$router.push('/articles')">
      <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
      <span>返回列表</span>
    </button>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="0"
      class="article-form"
    >
      <!-- 顶部标题栏 -->
      <div class="article-title-bar">
        <el-form-item prop="title" class="title-form-item">
          <input type="text" class="title-input" v-model="form.title" placeholder="输入文章标题...">
        </el-form-item>
      </div>

      <!-- 编辑器主体 -->
      <div class="editor-main" ref="editorContainerRef">
        <MdEditor
          ref="mdEditorRef"
          v-model="form.content"
          :theme="editorTheme"
          class="md-editor"
          :preview="true"
          previewTheme="github"
          @onUploadImg="onUploadImg"
          @onDrop="onEditorDrop"
        />
      </div>

      <!-- 底部操作栏 -->
      <div class="article-actions">
        <div class="action-left">
          <span class="word-count">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>
            {{ wordCount }} 字
          </span>
          <span v-if="autoSaveStatus" class="auto-save-status" :class="autoSaveStatus">
            <svg v-if="autoSaveStatus === 'saving'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle></svg>
            <svg v-else-if="autoSaveStatus === 'saved'" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
            {{ autoSaveText }}
          </span>
        </div>
        <div class="action-right">
          <el-button @click="resetForm">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="1 4 1 10 7 10"></polyline><path d="M3.51 15a9 9 0 1 0 2.13-9.36L1 10"></path></svg>
            <span>重置</span>
          </el-button>
          <el-button type="default" @click="handleSaveDraft">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"></path><polyline points="17 21 17 13 7 13 7 21"></polyline><polyline points="7 3 7 8 15 8"></polyline></svg>
            <span>保存草稿</span>
          </el-button>
          <el-button type="primary" @click="handlePublish">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="22" y1="2" x2="11" y2="13"></line><polygon points="22 2 15 22 11 13 2 9 22 2"></polygon></svg>
            <span>发布文章</span>
          </el-button>
        </div>
      </div>
    </el-form>

    <!-- 发布设置弹窗 -->
    <div class="publish-modal-overlay" :class="{ active: showPublishModal }" @click.self="showPublishModal = false">
      <div class="publish-modal">
        <div class="publish-modal-header">
          <h3>发布设置</h3>
          <button class="publish-modal-close" @click="showPublishModal = false">×</button>
        </div>
        <div class="publish-modal-body">
          <!-- 分类 -->
          <div class="publish-field">
            <label class="publish-label">分类 <span class="required">*</span></label>
            <CustomSelect
              v-model="form.category_id"
              :options="categoryOptions"
              placeholder="请选择分类"
            />
          </div>

          <!-- 标签：多选 + 搜索 + 全选/清空 按钮，复选框勾选一次选多个 -->
          <div class="publish-field">
            <label class="publish-label">标签</label>
            <CustomSelect
              v-model="form.tag_ids"
              :multiple="true"
              :searchable="true"
              :show-select-all="true"
              :options="tagOptions"
              placeholder="选择文章标签（可多选）"
              search-placeholder="搜索标签名..."
            />
            <div class="publish-hint">已选 {{ form.tag_ids?.length || 0 }} 个标签（支持批量勾选）</div>
          </div>

          <!-- 封面图 -->
          <div class="publish-field">
            <label class="publish-label">封面图</label>
            <ImageUpload v-model="form.cover" :category="MEDIA_CATEGORIES.ARTICLE" />
          </div>

          <!-- URL Slug -->
          <div class="publish-field">
            <label class="publish-label">URL Slug</label>
            <input
              type="text"
              class="publish-select"
              v-model="form.slug"
              placeholder="留空则自动根据标题生成"
            >
            <div v-if="form.slug" class="slug-preview">
              文章链接：<code>{{ blogUrl }}/posts/{{ form.slug }}</code>
            </div>
            <div v-else-if="form.title" class="slug-preview slug-auto">
              将自动生成：<code>{{ blogUrl }}/posts/{{ previewSlug }}</code>
            </div>
          </div>

          <!-- 定时发布 -->
          <div class="publish-field">
            <label class="publish-label">定时发布</label>
            <input
              type="datetime-local"
              class="publish-select"
              v-model="form.scheduled_at"
            >
            <div class="publish-hint">选择未来时间，文章将在指定时间自动发布；留空则立即发布</div>
          </div>

          <!-- 置顶 -->
          <div class="publish-field">
            <div class="publish-switch-row">
              <div class="publish-switch-text">
                <label class="publish-label">置顶文章</label>
                <div class="publish-hint">开启后文章将在前台列表顶部优先展示</div>
              </div>
              <el-switch v-model="form.is_top" />
            </div>
          </div>

          <!-- 摘要 -->
          <div class="publish-field">
            <div class="publish-label-row">
              <label class="publish-label">摘要</label>
              <span class="publish-char-count" :class="{ 'exceeded': form.summary.length > 200 }">
                {{ form.summary.length }} / 200
              </span>
            </div>
            <textarea
              v-model="form.summary"
              class="publish-textarea"
              placeholder="请输入文章摘要（选填，最多200字）"
              rows="3"
              maxlength="200"
            ></textarea>
          </div>

          <!-- SEO 设置 -->
          <div class="publish-field">
            <button type="button" class="seo-toggle" @click="showSeoSettings = !showSeoSettings">
              <span>SEO 设置</span>
              <svg :class="{ rotated: showSeoSettings }" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
            </button>
            <div v-show="showSeoSettings" class="seo-settings">
              <div class="seo-field">
                <label class="seo-label">SEO 标题</label>
                <input type="text" class="publish-select" v-model="form.seo_title" placeholder="留空则使用文章标题" maxlength="200">
                <span class="seo-char-count">{{ form.seo_title.length }} / 200</span>
              </div>
              <div class="seo-field">
                <label class="seo-label">SEO 描述</label>
                <textarea class="publish-textarea" v-model="form.seo_description" placeholder="留空则使用文章摘要" rows="2" maxlength="500"></textarea>
                <span class="seo-char-count">{{ form.seo_description.length }} / 500</span>
              </div>
              <div class="seo-field">
                <label class="seo-label">SEO 关键词</label>
                <input type="text" class="publish-select" v-model="form.seo_keywords" placeholder="多个关键词用英文逗号分隔，如：Go,Vue,博客" maxlength="300">
                <span class="seo-char-count">{{ form.seo_keywords.length }} / 300</span>
              </div>
            </div>
          </div>
        </div>
        <div class="publish-modal-footer">
          <button class="publish-btn-secondary" @click="showPublishModal = false">取消</button>
          <button class="publish-btn-primary" @click="confirmPublish">确认发布</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { getArticleDetail, createArticle, updateArticle } from '../../api/article'
import { getCategoryList } from '../../api/category'
import { getTagList } from '../../api/tag'
import { uploadFile, MEDIA_CATEGORIES } from '../../api/media'
import ImageUpload from '../../components/common/ImageUpload.vue'
import CustomSelect from '../../components/common/CustomSelect.vue'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const editorTheme = ref(localStorage.getItem('scheme') === 'dark' ? 'dark' : 'light')
const formRef = ref(null)
const showPublishModal = ref(false)
const showSeoSettings = ref(false)
const mdEditorRef = ref(null)
const editorContainerRef = ref(null)

// 博客首页 URL（用于 Slug 预览）
const blogUrl = computed(() => {
  return window.location.origin
})

// 预览自动生成的 slug
const previewSlug = computed(() => {
  if (!form.value.title) return ''
  let slug = form.value.title.toLowerCase()
  // 保留字母、数字、空格、连字符
  slug = slug.replace(/[^a-z0-9\s-]/g, '')
  slug = slug.trim().replace(/\s+/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '')
  if (!slug) slug = 'post-' + Date.now().toString(36)
  if (slug.length > 200) slug = slug.slice(0, 200)
  return slug
})

const form = ref({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: '',
  tag_ids: [],
  status: 'draft',
  slug: '',
  scheduled_at: null,
  is_top: false,
  seo_title: '',
  seo_description: '',
  seo_keywords: ''
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择文章分类', trigger: 'change' }
  ],
  summary: [
    { max: 200, message: '摘要不能超过 200 个字符', trigger: 'blur' }
  ]
}

const categories = ref([])
const tags = ref([])

// CustomSelect 用的 options
const categoryOptions = computed(() => {
  return (categories.value || []).map(c => ({ value: c.id, label: c.name }))
})
const tagOptions = computed(() => {
  return (tags.value || []).map(t => ({ value: t.id, label: t.name }))
})

const wordCount = computed(() => {
  const content = form.value.content || ''
  // 移除 Markdown 语法字符，统计实际内容字数
  const plainText = content
    .replace(/#{1,6}\s/g, '')
    .replace(/\*{1,3}(.*?)\*{1,3}/g, '$1')
    .replace(/`(.*?)`/g, '$1')
    .replace(/```[\s\S]*?```/g, '')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, '')
    .replace(/[-*+]\s/g, '')
    .replace(/\d+\.\s/g, '')
    .replace(/>\s/g, '')
    .replace(/\|/g, '')
    .replace(/---+/g, '')
    .trim()
  return plainText.length
})

// 将 ISO 时间格式转换为 datetime-local input 所需格式
const formatForDateTimeInput = (isoStr) => {
  if (!isoStr) return null
  const d = new Date(isoStr)
  if (isNaN(d.getTime())) return null
  return d.toISOString().slice(0, 16)
}

// 批量上传图片并返回URL数组
const uploadImages = async (files) => {
  const urls = []
  for (const file of files) {
    try {
      const uploadRes = await uploadFile(file, MEDIA_CATEGORIES.ARTICLE)
      if (uploadRes.code === 0) {
        urls.push(uploadRes.data.url)
      } else {
        ElMessage.error(`图片 ${file.name || '文件'} 上传失败：${uploadRes.message || '未知错误'}`)
      }
    } catch (e) {
      console.error('上传失败:', e)
      ElMessage.error(`图片 ${file.name || '文件'} 上传失败`)
    }
  }
  return urls
}

// 在编辑器光标位置插入Markdown图片语法
const insertImagesAtCursor = (urls) => {
  if (!urls || urls.length === 0) return
  const insertText = urls.map(url => `![image](${url})`).join('\n')

  // 使用 md-editor-v3 暴露的 insert 方法在光标处插入
  if (mdEditorRef.value && typeof mdEditorRef.value.insert === 'function') {
    mdEditorRef.value.insert((selectedText) => {
      return {
        targetValue: selectedText
          ? `${selectedText}\n${insertText}\n`
          : insertText,
        select: false,
        deviation: insertText.length
      }
    })
  } else {
    // fallback：直接追加到末尾
    form.value.content += (form.value.content ? '\n' : '') + insertText
  }
}

// 图片上传处理（工具栏按钮触发）
const onUploadImg = async (files, callback) => {
  const res = await uploadImages(files)
  callback(res)
}

// 拖拽图片到编辑器
const onEditorDrop = (e) => {
  // 检查是否有文件
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return

  const imageFiles = Array.from(files).filter(file => /^image\/.*/.test(file.type))
  if (imageFiles.length === 0) return

  e.preventDefault()
  e.stopPropagation()

  // 上传并在光标处插入
  uploadImages(imageFiles).then(urls => {
    if (urls.length > 0) {
      insertImagesAtCursor(urls)
    }
  })
}

// 处理粘贴事件（支持从剪贴板粘贴图片）
const handleEditorPaste = async (e) => {
  // 检查焦点是否在编辑器内或编辑器容器中
  const activeEl = document.activeElement
  const isFocusInEditor =
    editorContainerRef.value?.contains(activeEl) ||
    (e.target && (
      e.target === editorContainerRef.value ||
      editorContainerRef.value?.contains(e.target)
    ))

  // 如果焦点不在编辑器内，跳过处理
  if (!isFocusInEditor) return

  const clipboardData = e.clipboardData || window.clipboardData
  if (!clipboardData) return

  // 检查剪贴板中是否有文件（图片）
  const items = clipboardData.items
  if (!items) return

  const imageFiles = []
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    if (item.type && item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        imageFiles.push(file)
      }
    }
  }

  // 如果有图片文件，处理上传
  if (imageFiles.length > 0) {
    e.preventDefault()
    e.stopPropagation()
    const urls = await uploadImages(imageFiles)
    if (urls.length > 0) {
      insertImagesAtCursor(urls)
    }
  }
}

const loadData = async () => {
  try {
    const [catRes, tagRes] = await Promise.all([getCategoryList(), getTagList()])
    categories.value = catRes.data || []
    tags.value = tagRes.data || []
  } catch (e) { console.error(e) }

  // 新文章时尝试恢复 localStorage 草稿
  if (!isEdit.value) {
    await restoreLocalDraft()
    return
  }

  if (isEdit.value) {
    try {
      const res = await getArticleDetail(route.params.id)
      if (res.data) {
        form.value = {
          ...form.value,
          ...res.data,
          category_id: res.data.category_id || '',
          tag_ids: res.data.tag_ids || [],
          status: res.data.status || 'draft',
          slug: res.data.slug || '',
          scheduled_at: res.data.scheduled_at ? formatForDateTimeInput(res.data.scheduled_at) : null,
          is_top: !!res.data.is_top,
          seo_title: res.data.seo_title || '',
          seo_description: res.data.seo_description || '',
          seo_keywords: res.data.seo_keywords || ''
        }
      }
    } catch (e) { console.error(e) }
  }
}

// 自动保存草稿
const autoSaveStatus = ref('') // '' | 'saving' | 'saved' | 'error'
const autoSaveText = computed(() => {
  switch (autoSaveStatus.value) {
    case 'saving': return '正在保存...'
    case 'saved': return '草稿已自动保存'
    case 'error': return '自动保存失败'
    default: return ''
  }
})
let autoSaveTimer = null
let lastAutoSavedTitle = ''
let lastAutoSavedContent = ''

// localStorage 草稿相关
const DRAFT_KEY = 'article_draft_new'

const saveLocalDraft = () => {
  if (!form.value.title && !form.value.content) return
  const draft = { ...form.value, savedAt: new Date().toISOString() }
  localStorage.setItem(DRAFT_KEY, JSON.stringify(draft))
}

const loadLocalDraft = () => {
  try {
    const raw = localStorage.getItem(DRAFT_KEY)
    if (!raw) return null
    return JSON.parse(raw)
  } catch {
    return null
  }
}

const clearLocalDraft = () => {
  localStorage.removeItem(DRAFT_KEY)
}

const restoreLocalDraft = async () => {
  if (isEdit.value) return
  const draft = loadLocalDraft()
  if (!draft) return

  try {
    await ElMessageBox.confirm(
      `检测到未保存的草稿（${draft.title || '无标题'}），是否恢复？`,
      '恢复草稿',
      {
        confirmButtonText: '恢复',
        cancelButtonText: '丢弃',
        type: 'info'
      }
    )
    form.value = {
      ...form.value,
      title: draft.title || '',
      content: draft.content || '',
      summary: draft.summary || '',
      cover: draft.cover || '',
      category_id: draft.category_id || '',
      tag_ids: draft.tag_ids || [],
      status: draft.status || 'draft',
      slug: draft.slug || '',
      scheduled_at: draft.scheduled_at || null,
      is_top: draft.is_top || false,
      seo_title: draft.seo_title || '',
      seo_description: draft.seo_description || '',
      seo_keywords: draft.seo_keywords || ''
    }
    ElMessage.success('草稿已恢复')
  } catch {
    clearLocalDraft()
  }
}

const autoSaveDraft = async () => {
  // 新文章：保存到 localStorage
  if (!isEdit.value) {
    if (!form.value.title && !form.value.content) return
    saveLocalDraft()
    autoSaveStatus.value = 'saved'
    setTimeout(() => { autoSaveStatus.value = '' }, 3000)
    return
  }

  // 编辑模式：保存到服务端
  if (form.value.title === lastAutoSavedTitle && form.value.content === lastAutoSavedContent) return
  if (!form.value.title || !form.value.title.trim()) return

  autoSaveStatus.value = 'saving'
  try {
    await updateArticle(route.params.id, { ...form.value, status: 'draft' })
    lastAutoSavedTitle = form.value.title
    lastAutoSavedContent = form.value.content
    autoSaveStatus.value = 'saved'
    setTimeout(() => { autoSaveStatus.value = '' }, 3000)
  } catch (e) {
    console.error('自动保存失败:', e)
    autoSaveStatus.value = 'error'
    setTimeout(() => { autoSaveStatus.value = '' }, 3000)
  }
}

const debouncedAutoSave = () => {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(autoSaveDraft, 3000)
}

// 监听标题和内容变化触发自动保存
watch(
  () => [form.value.title, form.value.content],
  () => { debouncedAutoSave() }
)

onBeforeUnmount(() => {
  if (autoSaveTimer) clearTimeout(autoSaveTimer)
  // 移除粘贴事件监听器
  if (editorContainerRef.value) {
    editorContainerRef.value.removeEventListener('paste', handleEditorPaste, true)
  }
  document.removeEventListener('paste', handleEditorPaste, true)
})

const handleSaveDraft = async () => {
  form.value.status = 'draft'
  await handleSave(true)
}

const handlePublish = async () => {
  // 验证标题
  if (!form.value.title || !form.value.title.trim()) {
    ElMessage.warning('请先输入文章标题')
    return
  }

  // 打开发布设置弹窗
  showPublishModal.value = true
}

const confirmPublish = async () => {
  // 验证分类
  if (!form.value.category_id) {
    ElMessage.warning('请选择文章分类')
    return
  }

  // 如果设置了定时发布时间，状态设为 scheduled，否则为 published
  form.value.status = form.value.scheduled_at ? 'scheduled' : 'published'
  showPublishModal.value = false
  await handleSave()
}

const handleSave = async (isDraft = false) => {
  // 草稿保存只验证标题
  if (!form.value.title || !form.value.title.trim()) {
    ElMessage.warning('请输入文章标题')
    return
  }

  try {
    // 转换 scheduled_at 为 RFC3339 格式（Go time.Time 要求）
    const payload = { ...form.value }
    if (payload.scheduled_at && typeof payload.scheduled_at === 'string') {
      const dt = new Date(payload.scheduled_at)
      if (!isNaN(dt.getTime())) {
        payload.scheduled_at = dt.toISOString()
      }
    }
    if (isEdit.value) {
      await updateArticle(route.params.id, payload)
    } else {
      await createArticle(payload)
    }
    // 清除 localStorage 草稿
    clearLocalDraft()
    // 重置自动保存跟踪
    lastAutoSavedTitle = form.value.title
    lastAutoSavedContent = form.value.content
    if (isDraft) {
      autoSaveStatus.value = 'saved'
      setTimeout(() => { autoSaveStatus.value = '' }, 3000)
      ElMessage.success('草稿保存成功')
      return
    }
    ElMessage.success('保存成功')
    router.push('/articles')
  } catch (e) {
    console.error(e)
    ElMessage.error('保存失败，请重试')
  }
}

// 重置表单
const resetForm = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要重置表单吗？所有未保存的更改将丢失。',
      '确认重置',
      {
        confirmButtonText: '确定重置',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // 重置表单验证状态
    if (formRef.value) {
      formRef.value.resetFields()
    }

    // 重置表单数据
    form.value = {
      title: '',
      content: '',
      summary: '',
      cover: '',
      category_id: '',
      tag_ids: [],
      status: 'draft',
      slug: '',
      scheduled_at: null,
      is_top: false,
      seo_title: '',
      seo_description: '',
      seo_keywords: ''
    }

    // 清除 localStorage 草稿
    clearLocalDraft()

    // 重置自动保存跟踪
    lastAutoSavedTitle = ''
    lastAutoSavedContent = ''
    autoSaveStatus.value = ''

    // 如果是编辑模式，重新加载数据
    if (isEdit.value) {
      await loadData()
    }

    ElMessage.success('表单已重置')
  } catch (error) {
    // 用户取消操作
    if (error !== 'cancel') {
      console.error('重置失败:', error)
    }
  }
}

// 初始化编辑器事件监听
const initEditorEventListeners = () => {
  // 使用 nextTick 确保 DOM 已渲染
  setTimeout(() => {
    // 优先绑定到编辑器容器
    if (editorContainerRef.value) {
      editorContainerRef.value.addEventListener('paste', handleEditorPaste, true)
    }
    // 同时在 document 级别也绑定一个（捕获阶段），防止事件被编辑器内部拦截
    document.addEventListener('paste', handleEditorPaste, true)
  }, 100)
}

onMounted(() => {
  loadData()
  initEditorEventListeners()
})
</script>

<style scoped>
.article-edit-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding-bottom: 100px;
  /* 突破外层 content-area 的 32px padding，让编辑器更宽敞 */
  margin: -32px;
  padding: 24px 32px 100px;
}

/* 表单样式 */
.article-form {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.article-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.article-form :deep(.el-form-item__error) {
  padding-top: 4px;
  font-size: 12px;
}

/* 返回按钮 */
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: transparent;
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  color: var(--card-text-color-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  width: fit-content;
}

.back-btn:hover {
  background: var(--card-background);
  color: var(--card-text-color-main);
  border-color: var(--card-text-color-secondary);
}

/* 标题栏 - 参考CSDN扁平化设计，去掉卡片阴影包裹 */
.article-title-bar {
  background: var(--card-background);
  border-radius: 12px;
  padding: 24px 32px;
  border: 1px solid var(--card-separator-color);
}

.title-form-item {
  width: 100%;
}

.title-form-item :deep(.el-form-item__content) {
  line-height: normal;
}

.title-input {
  width: 100%;
  padding: 10px 0;
  border: none;
  font-size: 28px;
  font-weight: 700;
  background: transparent;
  color: var(--card-text-color-main);
  line-height: 1.4;
}

.title-input::placeholder {
  color: var(--card-text-color-tertiary);
  font-weight: 500;
}

.title-input:focus {
  outline: none;
}

/* 编辑器主体 - CSDN风格，扁平简洁 */
.editor-main {
  background: var(--card-background);
  border-radius: 12px;
  border: 1px solid var(--card-separator-color);
  overflow: hidden;
}

.md-editor {
  min-height: 650px;
}

/* 让编辑器内部区域自适应宽度，移除硬编码的固定尺寸 */
.md-editor :deep(.cm-scroller) {
  width: 100% !important;
  min-height: 550px;
  max-height: calc(100vh - 300px);
  overflow: auto;
}

/* 让 md-editor 工具栏和编辑器内部更宽敞 */
.md-editor :deep(.md-editor-toolbar) {
  padding: 10px 16px;
}

.md-editor :deep(.md-editor-content) {
  /* 编辑区和预览区的padding */
  padding-left: 0;
  padding-right: 0;
}

.md-editor :deep(.md-editor-preview) {
  /* 预览区左右留白 */
  padding: 24px 32px !important;
}

.md-editor :deep(.md-editor-textarea-wrapper) {
  /* 编辑输入区左右留白 - 参考CSDN */
  padding: 16px 24px !important;
  border-right: 1px solid var(--card-separator-color);
  box-sizing: border-box;
}

/* 底部操作栏 */
.article-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 32px;
  background: var(--card-background);
  border-top: 1px solid var(--card-separator-color);
  box-shadow: 0 -2px 12px var(--card-separator-color);
  position: fixed;
  bottom: 0;
  left: var(--sidebar-width);
  right: 0;
  z-index: 100;
}

.action-left,
.action-right {
  display: flex;
  gap: 14px;
  align-items: center;
}

/* 字数统计 */
.word-count {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--card-text-color-secondary);
  padding: 6px 14px;
  background: var(--body-background);
  border-radius: 6px;
}

.word-count svg {
  opacity: 0.6;
}

/* 自动保存状态 */
.auto-save-status {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  padding: 4px 12px;
  border-radius: 12px;
  animation: fadeIn 0.3s ease;
}

.auto-save-status.saving {
  color: var(--card-text-color-secondary);
  background: var(--body-background);
}

.auto-save-status.saved {
  color: var(--success-color);
  background: rgba(var(--success-color-rgb), 0.08);
}

.auto-save-status.error {
  color: var(--danger-color);
  background: rgba(var(--danger-color-rgb), 0.08);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 发布设置弹窗 */
.publish-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: none;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(4px);
}

.publish-modal-overlay.active {
  display: flex;
}

.publish-modal {
  background: var(--card-background);
  border-radius: 12px;
  width: 90%;
  max-width: 560px;
  min-height: 60vh;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.publish-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--card-separator-color);
}

.publish-modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: var(--card-text-color-main);
  margin: 0;
}

.publish-modal-close {
  background: none;
  font-size: 24px;
  color: var(--card-text-color-tertiary);
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s;
}

.publish-modal-close:hover {
  background: var(--body-background);
  color: var(--card-text-color-main);
}

.publish-modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.publish-field {
  margin-bottom: 24px;
}

.publish-field:last-child {
  margin-bottom: 0;
}

.publish-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: var(--card-text-color-main);
  margin-bottom: 10px;
}

.publish-label .required {
  color: var(--danger-color);
}

.publish-label-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.publish-select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--card-background);
  color: var(--card-text-color-main);
  cursor: pointer;
  transition: all 0.2s;
}

.publish-select:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px rgba(var(--accent-color-rgb), 0.1);
}

.publish-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  font-size: 14px;
  background: var(--card-background);
  color: var(--card-text-color-main);
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
  transition: all 0.2s;
}

.publish-textarea:focus {
  outline: none;
  border-color: var(--accent-color);
  box-shadow: 0 0 0 3px rgba(var(--accent-color-rgb), 0.1);
}

.publish-char-count {
  font-size: 12px;
  color: var(--card-text-color-tertiary);
}

.publish-char-count.exceeded {
  color: var(--danger-color);
  font-weight: 600;
}

.publish-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--card-separator-color);
}

.publish-btn-secondary {
  padding: 10px 20px;
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  background: var(--card-background);
  color: var(--card-text-color-main);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.publish-btn-secondary:hover {
  border-color: var(--card-text-color-tertiary);
}

.publish-btn-primary {
  padding: 10px 24px;
  border: none;
  border-radius: 8px;
  background: var(--accent-color);
  color: var(--accent-color-text);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.publish-btn-primary:hover {
  background: var(--accent-color-darker);
}

/* Slug 预览 */
.slug-preview {
  margin-top: 8px;
  font-size: 12px;
  color: var(--card-text-color-secondary);
}
.slug-preview code {
  background: var(--body-background);
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: var(--accent-color);
}
.slug-preview.slug-auto code {
  color: var(--card-text-color-tertiary);
}

/* 定时发布提示 */
.publish-hint {
  margin-top: 6px;
  font-size: 12px;
  color: var(--card-text-color-tertiary);
}

/* 置顶开关（与 Element Plus 默认蓝不同，改用项目强调色） */
.publish-switch-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.publish-switch-text {
  min-width: 0;
}

.publish-switch-text .publish-label {
  margin-bottom: 4px;
}

.publish-switch-text .publish-hint {
  margin-top: 0;
}

.publish-switch-row :deep(.el-switch.is-checked .el-switch__core) {
  background-color: var(--accent-color);
  border-color: var(--accent-color);
}

/* SEO 折叠按钮 */
.seo-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 10px 0;
  background: none;
  border: none;
  font-size: 14px;
  font-weight: 600;
  color: var(--card-text-color-main);
  cursor: pointer;
}
.seo-toggle svg {
  transition: transform 0.2s ease;
}
.seo-toggle svg.rotated {
  transform: rotate(180deg);
}

/* SEO 设置区域 */
.seo-settings {
  padding: 12px 0 0;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.seo-field {
  position: relative;
}
.seo-label {
  display: block;
  font-size: 13px;
  color: var(--card-text-color-secondary);
  margin-bottom: 6px;
}
.seo-char-count {
  position: absolute;
  right: 8px;
  bottom: 8px;
  font-size: 11px;
  color: var(--card-text-color-tertiary);
}

/* 响应式 */
@media (max-width: 768px) {
  .publish-modal {
    width: 95%;
    max-height: 90vh;
  }

  .publish-modal-body {
    padding: 20px;
  }
}
</style>
