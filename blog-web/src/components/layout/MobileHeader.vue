<template>
  <header class="mobile-header">
    <router-link class="mobile-brand" to="/">
      <img :src="siteProfile.avatar" :alt="siteProfile.nickname">
      <span>{{ siteProfile.nickname }}</span>
    </router-link>
    <nav class="mobile-actions">
      <button class="icon-button" type="button" @click="appStore.toggleScheme" title="切换主题" aria-label="切换主题">
        {{ appStore.scheme === 'dark' ? '☀' : '☾' }}
      </button>
      <button class="icon-button" type="button" @click="goSearch" title="搜索" aria-label="搜索">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
      </button>
      <button
        ref="menuToggleRef"
        class="icon-button mobile-menu-toggle"
        type="button"
        :title="menuOpen ? '关闭菜单' : '打开菜单'"
        :aria-label="menuOpen ? '关闭导航菜单' : '打开导航菜单'"
        :aria-expanded="menuOpen"
        aria-controls="mobile-navigation-drawer"
        @click="$emit('toggle-menu')"
      >
        <svg v-if="!menuOpen" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="4" y1="7" x2="20" y2="7"></line><line x1="4" y1="12" x2="20" y2="12"></line><line x1="4" y1="17" x2="20" y2="17"></line></svg>
        <svg v-else xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="6" y1="6" x2="18" y2="18"></line><line x1="18" y1="6" x2="6" y2="18"></line></svg>
      </button>
    </nav>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'
defineProps({
  menuOpen: Boolean,
  siteProfile: {
    type: Object,
    required: true
  }
})
defineEmits(['toggle-menu'])
const appStore = useAppStore()
const router = useRouter()
const menuToggleRef = ref(null)
const goSearch = () => {
  router.push({ name: 'Search' })
}

defineExpose({
  focusMenuToggle: () => menuToggleRef.value?.focus()
})
</script>

<style scoped>
/* 移动端头部样式在全局main.css中定义 */
</style>
