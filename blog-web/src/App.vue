<template>
  <div class="app" :data-scheme="scheme">
    <MobileHeader ref="mobileHeaderRef" :menu-open="menuOpen" :site-profile="siteProfile" @toggle-menu="toggleMenu" />
    <div class="site-shell" :class="{ 'menu-open': menuOpen }">
      <button
        class="mobile-menu-backdrop"
        type="button"
        tabindex="-1"
        aria-hidden="true"
        aria-label="关闭导航菜单"
        @click="closeMenu"
      ></button>
      <AppSidebar ref="sidebarRef" :menu-open="menuOpen" :site-profile="siteProfile" @close-menu="menuOpen = false" />
      <main class="content-area" :inert="menuOpen">
        <router-view />
      </main>
    </div>
    <BackToTop />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from './stores/app'
import MobileHeader from './components/layout/MobileHeader.vue'
import AppSidebar from './components/layout/AppSidebar.vue'
import BackToTop from './components/common/BackToTop.vue'
import { normalizeSiteProfile } from './utils/siteProfile'

const appStore = useAppStore()
const route = useRoute()
const scheme = ref(appStore.scheme)
const menuOpen = ref(false)
const siteProfile = ref(normalizeSiteProfile())
const mobileHeaderRef = ref(null)
const sidebarRef = ref(null)

const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
}

const closeMenu = () => {
  menuOpen.value = false
}

const handleViewportChange = () => {
  if (window.innerWidth > 860) closeMenu()
}

const handleKeydown = (event) => {
  if (event.key === 'Escape') closeMenu()
}

const loadSiteProfile = async () => {
  try {
    const response = await fetch('/api/v1/user/info')
    const result = await response.json()
    if (response.ok && result.code === 0) {
      siteProfile.value = normalizeSiteProfile(result.data)
    }
  } catch {
    siteProfile.value = normalizeSiteProfile()
  }
}

watch(menuOpen, async (open, wasOpen) => {
  document.body.classList.toggle('mobile-menu-open', open)
  await nextTick()
  if (open) {
    sidebarRef.value?.focusInitial()
  } else if (wasOpen && window.innerWidth <= 860) {
    mobileHeaderRef.value?.focusMenuToggle()
  }
})

watch(() => route.fullPath, closeMenu)

onMounted(() => {
  loadSiteProfile()
  window.addEventListener('resize', handleViewportChange)
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleViewportChange)
  document.removeEventListener('keydown', handleKeydown)
  document.body.classList.remove('mobile-menu-open')
})
</script>

<style>
@import './assets/styles/main.css';
</style>
