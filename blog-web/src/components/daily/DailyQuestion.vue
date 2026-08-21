<template>
  <article class="article-card dq-daily-card">
    <div class="card-link">
      <div class="dq-daily-header">
        <span class="category-pill category-life">每日一问</span>
        <div class="dq-daily-nav">
          <button class="dq-daily-nav-btn" @click.stop="prevDay" :disabled="!hasPrev">◀ 前一天</button>
          <div class="dq-daily-date-wrapper">
            <button class="dq-daily-date-btn" @click.stop="toggleCalendar">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"></rect><line x1="16" x2="16" y1="2" y2="6"></line><line x1="8" x2="8" y1="2" y2="6"></line><line x1="3" x2="21" y1="10" y2="10"></line></svg>
              <span>{{ question?.date || '加载中...' }}</span>
            </button>
            <!-- 日历选择器 -->
            <div v-if="showCalendar" class="dq-calendar-popup" @click.stop>
              <div class="dq-calendar-header">
                <button class="dq-calendar-nav" @click="prevMonth">◀</button>
                <span class="dq-calendar-title">{{ calendarYear }}年{{ calendarMonth + 1 }}月</span>
                <button class="dq-calendar-nav" @click="nextMonth">▶</button>
              </div>
              <div class="dq-calendar-weekdays">
                <span v-for="day in ['日', '一', '二', '三', '四', '五', '六']" :key="day">{{ day }}</span>
              </div>
              <div class="dq-calendar-days">
                <button
                  v-for="(day, index) in calendarDays"
                  :key="index"
                  class="dq-calendar-day"
                  :class="{
                    'other-month': day.otherMonth,
                    'is-today': day.isToday,
                    'is-selected': day.date === question?.date,
                    'has-question': day.hasQuestion
                  }"
                  :disabled="day.otherMonth"
                  @click="selectDate(day.date)"
                >
                  {{ day.day }}
                </button>
              </div>
              <button class="dq-calendar-close" @click="showCalendar = false">关闭</button>
            </div>
          </div>
          <button class="dq-daily-nav-btn" @click.stop="nextDay" :disabled="!hasNext">后一天 ▶</button>
        </div>
      </div>

      <h2 class="dq-daily-title">{{ question?.question || '暂无问题' }}</h2>

      <div class="dq-daily-answer-wrapper" v-if="question" :class="{ expanded: answerVisible }">
        <p class="dq-daily-preview">{{ question.answer }}</p>
        <div v-if="!answerVisible" class="dq-daily-mask" @click.stop="answerVisible = true">
          <span class="dq-daily-mask-text">点击查看答案</span>
        </div>
        <div v-if="answerVisible" class="dq-daily-hide-btn" @click.stop="answerVisible = false">
          <span class="dq-daily-hide-text">点击收起答案</span>
        </div>
      </div>

      <router-link :to="`/daily/${question?.date}`" class="dq-view-btn" v-if="question">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"></path><circle cx="12" cy="12" r="3"></circle></svg>
        <span>查看所有每日一问</span>
      </router-link>
    </div>
  </article>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { getLatestQuestion, getQuestionByDate, getPreviousQuestion, getNextQuestion } from '../../api/daily'

const question = ref(null)
const answerVisible = ref(false)
const hasPrev = ref(false)
const hasNext = ref(false)
const showCalendar = ref(false)
const calendarYear = ref(new Date().getFullYear())
const calendarMonth = ref(new Date().getMonth())

// 有每日一问的日期集合（用于日历标亮）
const questionDates = ref(new Set())

// 拉取所有已发布问题的日期
const loadQuestionDates = async () => {
  try {
    const res = await getAllPublishedQuestions()
    const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
    questionDates.value = new Set(list.map((q) => q.date).filter(Boolean))
  } catch (e) {
    console.error('获取每日一问日期失败', e)
  }
}

// 计算日历天数
const calendarDays = computed(() => {
  const year = calendarYear.value
  const month = calendarMonth.value
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDayOfWeek = firstDay.getDay()
  const daysInMonth = lastDay.getDate()

  const days = []
  const today = new Date()
  const todayStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`

  // 填充上个月的日期
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = startDayOfWeek - 1; i >= 0; i--) {
    const day = prevMonthLastDay - i
    const m = month === 0 ? 12 : month
    const y = month === 0 ? year - 1 : year
    const dateStr = `${y}-${String(m).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    days.push({
      day,
      date: dateStr,
      otherMonth: true,
      isToday: false,
      hasQuestion: questionDates.value.has(dateStr)
    })
  }

  // 填充当月日期
  for (let day = 1; day <= daysInMonth; day++) {
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    days.push({
      day,
      date: dateStr,
      otherMonth: false,
      isToday: dateStr === todayStr,
      hasQuestion: questionDates.value.has(dateStr)
    })
  }

  // 填充下个月的日期
  const remaining = 42 - days.length
  for (let day = 1; day <= remaining; day++) {
    const m = month + 2 > 12 ? 1 : month + 2
    const y = month + 2 > 12 ? year + 1 : year
    const dateStr = `${y}-${String(m).padStart(2, '0')}-${String(day).padStart(2, '0')}`
    days.push({
      day,
      date: dateStr,
      otherMonth: true,
      isToday: false,
      hasQuestion: questionDates.value.has(dateStr)
    })
  }

  return days
})

const loadQuestion = async (date) => {
  try {
    const res = date ? await getQuestionByDate(date).catch(() => getLatestQuestion()) : await getLatestQuestion()
    if (res.data) {
      question.value = res.data
      answerVisible.value = false
      // 更新日历月份到问题所在月份
      if (res.data.date) {
        const [y, m] = res.data.date.split('-').map(Number)
        calendarYear.value = y
        calendarMonth.value = m - 1
        await checkNavigation(res.data.date)
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const checkNavigation = async (date) => {
  try {
    await getPreviousQuestion(date)
    hasPrev.value = true
  } catch {
    hasPrev.value = false
  }
  try {
    await getNextQuestion(date)
    hasNext.value = true
  } catch {
    hasNext.value = false
  }
}

const toggleCalendar = () => {
  showCalendar.value = !showCalendar.value
}

const prevMonth = () => {
  if (calendarMonth.value === 0) {
    calendarMonth.value = 11
    calendarYear.value--
  } else {
    calendarMonth.value--
  }
}

const nextMonth = () => {
  if (calendarMonth.value === 11) {
    calendarMonth.value = 0
    calendarYear.value++
  } else {
    calendarMonth.value++
  }
}

const selectDate = async (date) => {
  try {
    const res = await getQuestionByDate(date)
    if (res.data) {
      question.value = res.data
      answerVisible.value = false
      await checkNavigation(res.data.date)
    }
    showCalendar.value = false
  } catch {
    alert('该日期没有问题')
  }
}

const prevDay = async () => {
  if (question.value?.date) {
    try {
      const res = await getPreviousQuestion(question.value.date)
      if (res.data) {
        question.value = res.data
        answerVisible.value = false
        const [y, m] = res.data.date.split('-').map(Number)
        calendarYear.value = y
        calendarMonth.value = m - 1
        await checkNavigation(res.data.date)
      }
    } catch (e) {
      console.error('没有前一天的问题了')
    }
  }
}

const nextDay = async () => {
  if (question.value?.date) {
    try {
      const res = await getNextQuestion(question.value.date)
      if (res.data) {
        question.value = res.data
        answerVisible.value = false
        const [y, m] = res.data.date.split('-').map(Number)
        calendarYear.value = y
        calendarMonth.value = m - 1
        await checkNavigation(res.data.date)
      }
    } catch (e) {
      console.error('没有后一天的问题了')
    }
  }
}

// 点击外部关闭日历
const handleClickOutside = (e) => {
  if (showCalendar.value && !e.target.closest('.dq-calendar-popup') && !e.target.closest('.dq-daily-date-btn')) {
    showCalendar.value = false
  }
}

onMounted(() => {
  loadQuestion()
  loadQuestionDates()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.dq-daily-card {
  border-left: 4px solid var(--accent-color);
  display: flex;
  flex-direction: column;
}

.dq-daily-card :deep(.card-link) {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 24px 24px 18px;
}

.dq-daily-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  flex-wrap: wrap;
  gap: 8px;
}

.dq-daily-nav {
  display: flex;
  align-items: center;
  gap: 8px;
}

.dq-daily-nav-btn {
  padding: 6px 12px;
  background: var(--card-background);
  border: 1px solid var(--card-separator-color);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--card-text-color-secondary);
}

.dq-daily-nav-btn:hover:not(:disabled) {
  background: var(--accent-color);
  color: var(--accent-color-text);
  border-color: var(--accent-color);
}

.dq-daily-nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.dq-daily-date-wrapper {
  position: relative;
}

.dq-daily-date-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: var(--card-background);
  border: 1px solid var(--card-separator-color);
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--card-text-color-secondary);
}

.dq-daily-date-btn:hover {
  background: var(--accent-color);
  color: var(--accent-color-text);
  border-color: var(--accent-color);
}

.dq-calendar-popup {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: var(--card-background);
  border: 1px solid var(--card-separator-color);
  border-radius: 8px;
  padding: 12px;
  box-shadow: var(--shadow-l2);
  z-index: 100;
  min-width: 280px;
}

.dq-calendar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.dq-calendar-nav {
  padding: 4px 8px;
  background: none;
  border: none;
  cursor: pointer;
  color: var(--card-text-color-secondary);
  font-size: 12px;
}

.dq-calendar-nav:hover {
  color: var(--accent-color);
}

.dq-calendar-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--card-text-color-main);
}

.dq-calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}

.dq-calendar-weekdays span {
  text-align: center;
  font-size: 11px;
  color: var(--card-text-color-tertiary);
  padding: 4px 0;
}

.dq-calendar-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.dq-calendar-day {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  color: var(--card-text-color-main);
  transition: all 0.15s;
}

.dq-calendar-day:hover:not(:disabled) {
  background: rgba(var(--accent-color-rgb), 0.1);
  color: var(--accent-color);
}

.dq-calendar-day.other-month {
  color: var(--card-text-color-tertiary);
  opacity: 0.4;
  cursor: not-allowed;
}

.dq-calendar-day.is-today {
  font-weight: 700;
  color: var(--accent-color);
}

.dq-calendar-day.is-selected {
  background: var(--accent-color);
  color: var(--accent-color-text);
}

/* 有每日一问的日期：圆点标记 */
.dq-calendar-day.has-question {
  position: relative;
}
.dq-calendar-day.has-question::after {
  content: '';
  position: absolute;
  bottom: 3px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent-color);
}
.dq-calendar-day.has-question.is-selected::after {
  background: var(--accent-color-text);
}

.dq-calendar-close {
  width: 100%;
  margin-top: 8px;
  padding: 6px;
  background: var(--card-separator-color);
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  color: var(--card-text-color-secondary);
}

.dq-calendar-close:hover {
  background: var(--accent-color);
  color: var(--accent-color-text);
}

.dq-daily-title {
  font-size: 1.35rem;
  font-weight: 700;
  margin-bottom: 8px;
  line-height: 1.4;
  color: var(--card-text-color-main);
}

.dq-daily-answer-wrapper {
  position: relative;
  margin-bottom: 12px;
  max-height: 104px;
  overflow: hidden;
  transition: max-height 0.35s ease;
}

.dq-daily-answer-wrapper.expanded {
  max-height: 2000px;
  overflow: visible;
}

.dq-daily-preview {
  color: var(--card-text-color-secondary);
  line-height: 1.6;
  margin-bottom: 0;
}

.dq-daily-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    180deg,
    transparent 0%,
    transparent 18%,
    color-mix(in srgb, var(--card-background) 55%, transparent) 38%,
    var(--card-background) 62%,
    var(--card-background) 100%
  );
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding-bottom: 10px;
  cursor: pointer;
  z-index: 2;
}

.dq-daily-mask::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 56px;
  background: var(--card-background);
  z-index: 0;
}

.dq-daily-mask-text {
  position: relative;
  z-index: 1;
  font-size: 13px;
  color: var(--accent-color);
  font-weight: 600;
  background: var(--card-background);
  padding: 6px 18px;
  border-radius: 20px;
  border: 1px solid var(--accent-color);
  transition: all 0.2s;
}

.dq-daily-mask-text:hover {
  background: var(--accent-color);
  color: var(--accent-color-text);
}

.dq-daily-hide-btn {
  text-align: center;
  margin-top: 8px;
}

.dq-daily-hide-text {
  font-size: 13px;
  color: var(--accent-color);
  font-weight: 500;
  cursor: pointer;
}

.dq-view-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(var(--accent-color-rgb), 0.08);
  color: var(--accent-color);
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s;
  margin-top: auto;
  align-self: flex-start;
}

.dq-view-btn:hover {
  background: var(--accent-color);
  color: var(--accent-color-text);
}

.category-life {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>
