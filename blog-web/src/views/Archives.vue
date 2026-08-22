<template>
  <div class="page-view">
    <h1 class="page-heading">文章归档</h1>

    <!-- 加载状态 -->
    <Loading v-if="articleStore.loading" text="加载归档中..." />

    <!-- 错误状态 -->
    <div v-else-if="articleStore.error" class="error-container">
      <p class="error-message">{{ articleStore.error }}</p>
      <button class="retry-btn" @click="retryFetch">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/></svg>
        重试
      </button>
    </div>

    <!-- 归档内容 -->
    <template v-else>
      <div v-for="group in articleStore.archives" :key="group.year" class="year-group">
        <h2 class="year-heading">{{ group.year }}</h2>
        <ul class="archive-list">
          <li v-for="article in group.articles" :key="article.id">
            <router-link :to="'/post/' + article.slug">
              <span class="archive-title">{{ article.title }}</span>
              <span class="archive-date">{{ formatDate(article.created_at) }}</span>
            </router-link>
          </li>
        </ul>
      </div>

      <div v-if="articleStore.archives.length === 0" class="empty">暂无归档文章</div>
    </template>
  </div>
</template>

<script setup>
import { onMounted, nextTick } from 'vue'
import { useRoute, onBeforeRouteLeave } from 'vue-router'
import { useArticleStore } from '../stores/article'
import { saveListState, takeListState } from '../composables/useListState'
import Loading from '../components/common/Loading.vue'
import { formatDate } from '../utils/date'

const articleStore = useArticleStore()
const route = useRoute()

const retryFetch = () => {
  articleStore.fetchArchives()
}

// 离开归档页（去往文章详情）时保存滚动位置
onBeforeRouteLeave(() => {
  saveListState(route.fullPath, { page: 1 })
})

onMounted(async () => {
  retryFetch()
  // 从文章详情返回时恢复滚动位置（归档页无分页）
  const state = takeListState(route.fullPath)
  if (state && state.scrollY > 0) {
    // 等归档列表 DOM 渲染完成后再恢复
    await nextTick()
    await new Promise(requestAnimationFrame)
    window.scrollTo({ top: state.scrollY, behavior: 'instant' })
  }
})
</script>

<style scoped>
/* 归档页面样式在全局main.css中定义 */
</style>
