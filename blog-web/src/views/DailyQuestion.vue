<template>
  <div class="dq-page">
    <div class="dq-main">
      <div class="dq-content">
        <template v-for="(item, index) in allQuestions" :key="item.id">
        <div
          :id="'dq-' + item.date"
          :data-question-id="item.id"
          class="dq-question-section"
        >
          <h1 class="dq-title">{{ item.question }}</h1>
          <div class="dq-meta">
            <span class="dq-date">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
              {{ item.date }}
            </span>
          </div>
          <div class="dq-answer post-content" v-html="renderMarkdown(item.answer)"></div>
        </div>
        <div class="dq-separator" v-if="index < allQuestions.length - 1"></div>
        </template>
      </div>

      <aside class="dq-sidebar">
        <h3 class="dq-sidebar-title">所有每日一问</h3>
        <div class="dq-sidebar-list">
          <div
            v-for="item in allQuestions"
            :key="item.id"
            class="dq-sidebar-item"
            :class="{ active: item.date === currentDate }"
            @click="scrollToDate(item.date)"
          >
            <span class="dq-sidebar-date">{{ item.date }}</span>
            <span class="dq-sidebar-question">{{ item.question }}</span>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { marked } from 'marked'
import { recordContentView } from '../api/view'
import { createPageViewRecorder } from '../utils/viewRecorder'

// 配置 marked
marked.setOptions({
  breaks: true,
  gfm: true
})

const allQuestions = ref([])
const currentDate = ref('')
const recordPageView = createPageViewRecorder(recordContentView)
let questionObserver = null

// 渲染 markdown
const renderMarkdown = (text) => {
  if (!text) return ''
  return marked(text)
}

const loadAllQuestions = async () => {
  try {
    // 先加载所有问题的完整数据
    const { getAllPublishedQuestions } = await import('../api/daily')
    const briefRes = await getAllPublishedQuestions()
    const briefList = briefRes.data || []

    allQuestions.value = briefList

    // 设置当前日期为第一个问题
    if (allQuestions.value.length > 0) {
      currentDate.value = allQuestions.value[0].date
    }
  } catch (e) {
    console.error(e)
  }
}

const applyRecordedView = async (questionID) => {
  try {
    const result = await recordPageView('daily_question', questionID)
    if (!result?.data) return
    const question = allQuestions.value.find(item => item.id === questionID)
    if (question) question.view_count = result.data.view_count
  } catch (error) {
    console.warn('记录每日一问浏览量失败:', error)
  }
}

const observeQuestions = () => {
  const sections = document.querySelectorAll('.dq-question-section[data-question-id]')
  if (typeof IntersectionObserver === 'undefined') {
    const firstID = Number(sections[0]?.dataset.questionId || 0)
    if (firstID) applyRecordedView(firstID)
    return
  }
  questionObserver = new IntersectionObserver(entries => {
    entries.forEach(entry => {
      if (!entry.isIntersecting) return
      const questionID = Number(entry.target.dataset.questionId || 0)
      if (questionID) applyRecordedView(questionID)
    })
  }, { threshold: 0.5 })
  sections.forEach(section => questionObserver.observe(section))
}

const scrollToDate = (date) => {
  const el = document.getElementById('dq-' + date)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    currentDate.value = date
  }
}

// 监听滚动，更新当前日期
const handleScroll = () => {
  const sections = allQuestions.value.map(q => ({
    date: q.date,
    el: document.getElementById('dq-' + q.date)
  })).filter(s => s.el)

  const scrollTop = window.scrollY + 100
  for (let i = sections.length - 1; i >= 0; i--) {
    if (sections[i].el.offsetTop <= scrollTop) {
      currentDate.value = sections[i].date
      break
    }
  }
}

onMounted(async () => {
  await loadAllQuestions()
  await nextTick()
  observeQuestions()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  questionObserver?.disconnect()
})
</script>

<style scoped>
.dq-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.dq-main {
  display: flex;
  gap: 24px;
}

.dq-content {
  flex: 1;
  min-width: 0;
}

.dq-sidebar {
  width: 280px;
  flex-shrink: 0;
  background: var(--card-background);
  border-radius: 12px;
  padding: 16px;
  height: fit-content;
  position: sticky;
  top: 80px;
  box-shadow: var(--shadow-l1);
}

.dq-sidebar-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color, #eee);
}

.dq-sidebar-list {
  max-height: 600px;
  overflow-y: auto;
}

.dq-sidebar-item {
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  margin-bottom: 4px;
}

.dq-sidebar-item:hover {
  background: var(--hover-bg, #f0f0f0);
}

.dq-sidebar-item.active {
  background: var(--accent-color, #007bff);
  color: white;
}

.dq-sidebar-date {
  display: block;
  font-size: 12px;
  opacity: 0.7;
  margin-bottom: 4px;
}

.dq-sidebar-question {
  display: block;
  font-size: 13px;
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dq-question-section {
  scroll-margin-top: 80px;
}

.dq-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 12px;
  line-height: 1.4;
}

.dq-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.dq-date {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: var(--card-text-color-tertiary, #8c8c8c);
}

.dq-date svg {
  flex-shrink: 0;
}

.dq-answer {
  background: var(--card-background, #fdfdfb);
  border-radius: 12px;
  padding: 24px;
  line-height: 1.8;
  font-size: 16px;
  box-shadow: var(--shadow-l1);
}

.dq-separator {
  height: 4px;
  margin: 48px 0;
  background: var(--accent-color);
  border-radius: 2px;
}

/* Markdown 内容样式 */
.dq-answer :deep(h1),
.dq-answer :deep(h2),
.dq-answer :deep(h3),
.dq-answer :deep(h4),
.dq-answer :deep(h5),
.dq-answer :deep(h6) {
  margin-top: 24px;
  margin-bottom: 12px;
  font-weight: 700;
  color: var(--text-primary, #333);
}

.dq-answer :deep(h1) { font-size: 1.5em; }
.dq-answer :deep(h2) { font-size: 1.3em; }
.dq-answer :deep(h3) { font-size: 1.1em; }

.dq-answer :deep(p) {
  margin-bottom: 12px;
}

.dq-answer :deep(ul),
.dq-answer :deep(ol) {
  margin-bottom: 12px;
  padding-left: 24px;
}

.dq-answer :deep(li) {
  margin-bottom: 6px;
}

.dq-answer :deep(code) {
  background: rgba(0, 0, 0, 0.06);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 0.9em;
}

.dq-answer :deep(pre) {
  background: #282c34;
  color: #abb2bf;
  padding: 16px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 16px;
}

.dq-answer :deep(pre code) {
  background: none;
  padding: 0;
  color: inherit;
}

.dq-answer :deep(blockquote) {
  border-left: 4px solid var(--accent-color, #1B365D);
  margin: 16px 0;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.03);
  color: var(--text-secondary, #666);
}

.dq-answer :deep(a) {
  color: var(--accent-color, #1B365D);
  text-decoration: none;
}

.dq-answer :deep(a:hover) {
  text-decoration: underline;
}

.dq-answer :deep(strong) {
  font-weight: 700;
}

@media (max-width: 768px) {
  .dq-main {
    flex-direction: column;
  }

  .dq-sidebar {
    width: 100%;
    position: static;
  }

  .dq-sidebar-list {
    max-height: 200px;
  }

  .dq-title {
    font-size: 22px;
  }
}
</style>
