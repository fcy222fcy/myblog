import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const scheme = ref(localStorage.getItem('scheme') || 'light')
  // 立即同步到 <html> 根元素，确保刷新后全局 CSS 变量（body 背景等）正确跟随主题
  document.documentElement.setAttribute('data-scheme', scheme.value)
  const searchOpen = ref(false)

  const toggleScheme = () => {
    scheme.value = scheme.value === 'light' ? 'dark' : 'light'
    localStorage.setItem('scheme', scheme.value)
    document.documentElement.setAttribute('data-scheme', scheme.value)
  }

  const openSearch = () => {
    searchOpen.value = true
  }

  const closeSearch = () => {
    searchOpen.value = false
  }

  return {
    scheme,
    searchOpen,
    toggleScheme,
    openSearch,
    closeSearch
  }
})
