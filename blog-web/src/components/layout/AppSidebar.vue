<template>
  <aside
    id="mobile-navigation-drawer"
    ref="drawerRef"
    class="left-sidebar"
    :class="{ open: menuOpen }"
    :role="menuOpen ? 'dialog' : undefined"
    :aria-modal="menuOpen ? 'true' : undefined"
    :aria-label="menuOpen ? '站点导航' : undefined"
    @keydown="handleDrawerKeydown"
  >
    <button ref="closeButtonRef" class="sidebar-mobile-close" type="button" aria-label="关闭导航菜单" @click="$emit('close-menu')">
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="6" y1="6" x2="18" y2="18"></line><line x1="18" y1="6" x2="6" y2="18"></line></svg>
    </button>
    <header>
      <figure class="site-avatar">
        <router-link to="/" @click="$emit('close-menu')">
          <img :src="siteProfile.avatar" class="site-logo" :alt="siteProfile.nickname">
        </router-link>
        <span class="avatar-badge">🤖</span>
      </figure>
      <div class="site-meta">
        <h1 class="site-name"><router-link to="/" @click="$emit('close-menu')">{{ siteProfile.nickname }}</router-link></h1>
        <h2 class="site-description">{{ siteProfile.bio }}</h2>
      </div>
    </header>

    <ol class="menu-social">
      <li><a href="mailto:3181088318@qq.com" title="Email">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect width="20" height="16" x="2" y="4" rx="2"></rect><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"></path></svg>
      </a></li>
      <li><a href="https://github.com/fcy222fcy" title="GitHub" target="_blank">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36-.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"></path><path d="M9 18c-4.51 2-5-2-7-2"></path></svg>
      </a></li>
      <li><a href="https://space.bilibili.com/1156362409?" title="Bilibili" target="_blank">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M17.813 4.653h.854c1.51.054 2.769.578 3.773 1.574 1.004.995 1.524 2.249 1.56 3.76v7.36c-.036 1.51-.556 2.769-1.56 3.773s-2.262 1.524-3.773 1.56H5.333c-1.51-.036-2.769-.556-3.773-1.56S.036 18.858 0 17.347v-7.36c.036-1.511.556-2.765 1.56-3.76 1.004-.996 2.262-1.52 3.773-1.574h.774l-1.174-1.12a1.234 1.234 0 0 1-.373-.906c0-.356.124-.658.373-.907l.027-.027c.267-.249.573-.373.92-.373.347 0 .653.124.92.373L9.653 4.44c.071.071.134.142.187.213h4.267a.836.836 0 0 1 .16-.213l2.853-2.747c.267-.249.573-.373.92-.373.347 0 .662.124.929.373.249.249.373.551.373.907 0 .355-.124.657-.373.906zM5.333 7.24c-.746.018-1.373.276-1.88.773-.506.498-.769 1.13-.786 1.894v7.52c.017.764.28 1.395.786 1.893.507.498 1.134.756 1.88.773h13.334c.746-.017 1.373-.275 1.88-.773.506-.498.769-1.129.786-1.893v-7.52c-.017-.765-.28-1.396-.786-1.894-.507-.497-1.134-.755-1.88-.773zM8 11.107c.373 0 .684.124.933.373.25.249.383.569.4.96v1.173c-.017.391-.15.711-.4.96-.249.25-.56.374-.933.374s-.684-.125-.933-.374c-.25-.249-.383-.569-.4-.96V12.44c.017-.391.15-.711.4-.96.249-.249.56-.373.933-.373zm8 0c.373 0 .684.124.933.373.25.249.383.569.4.96v1.173c-.017.391-.15.711-.4.96-.249.25-.56.374-.933.374s-.684-.125-.933-.374c-.25-.249-.383-.569-.4-.96V12.44c.017-.391.15-.711.4-.96.249-.249.56-.373.933-.373z"/></svg>
      </a></li>
      <li><a href="#" title="WeChat">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 512 512" fill="currentColor">
          <path d="M408.67,298.53a21,21,0,1,1,20.9-21,20.85,20.85,0,0,1-20.9,21m-102.17,0a21,21,0,1,1,20.9-21,20.84,20.84,0,0,1-20.9,21M458.59,417.39C491.1,394.08,512,359.13,512,319.51c0-71.08-68.5-129.35-154.41-129.35S203.17,248.43,203.17,319.51s68.5,129.34,154.42,129.34c17.41,0,34.83-2.33,49.92-7,2.49-.86,3.48-1.17,4.64-1.17a16.67,16.67,0,0,1,8.13,2.34L454,462.83a11.62,11.62,0,0,0,3.48,1.17,5,5,0,0,0,4.65-4.66,14.27,14.27,0,0,0-.77-3.86c-.41-1.46-5-16-7.36-25.27a18.94,18.94,0,0,1-.33-3.47,11.4,11.4,0,0,1,5-9.35"/>
          <path d="M246.13,178.51a24.47,24.47,0,0,1,0-48.94c12.77,0,24.38,11.65,24.38,24.47,1.16,12.82-10.45,24.47-24.38,24.47m-123.06,0A24.47,24.47,0,1,1,147.45,154a24.57,24.57,0,0,1-24.38,24.47M184.6,48C82.43,48,0,116.75,0,203c0,46.61,24.38,88.56,63.85,116.53C67.34,321.84,68,327,68,329a11.38,11.38,0,0,1-.66,4.49C63.85,345.14,59.4,364,59.21,365s-1.16,3.5-1.16,4.66a5.49,5.49,0,0,0,5.8,5.83,7.15,7.15,0,0,0,3.49-1.17L108,351c3.49-2.33,5.81-2.33,9.29-2.33a16.33,16.33,0,0,1,5.81,1.16c18.57,5.83,39.47,8.16,60.37,8.16h10.45a133.24,133.24,0,0,1-5.81-38.45c0-78.08,75.47-141,168.35-141h10.45C354.1,105.1,277.48,48,184.6,48"/>
        </svg>
      </a></li>
      <li><a href="/api/v1/feed.xml" title="RSS 订阅" target="_blank">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 11a9 9 0 0 1 9 9"></path><path d="M4 4a16 16 0 0 1 16 16"></path><circle cx="5" cy="19" r="1"></circle></svg>
      </a></li>
    </ol>

    <form class="sidebar-search" @submit.prevent="handleSearchSubmit" novalidate>
      <svg class="sidebar-search-icon" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line></svg>
      <input
        v-model="searchKeyword"
        type="search"
        placeholder="搜索文章，回车搜索..."
        class="sidebar-search-input"
        @keydown.enter.prevent="handleSearchSubmit"
      >
      <button type="submit" class="sidebar-search-btn" title="搜索">
        <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="m9 18 6-6-6-6"/></svg>
      </button>
    </form>

    <nav class="menu" id="main-menu">
      <li><router-link to="/" @click="$emit('close-menu')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="m3 9 9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>
        <span>主页</span>
      </router-link></li>
      <li><router-link to="/archives" @click="$emit('close-menu')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"></path></svg>
        <span>文章</span>
      </router-link></li>

      <li><router-link to="/about" @click="$emit('close-menu')">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
        <span>关于我</span>
      </router-link></li>
    </nav>

    <div class="sidebar-bottom">
      <button class="sidebar-bottom-item theme-toggle" type="button" @click="appStore.toggleScheme">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9Z"></path></svg>
        {{ appStore.scheme === 'dark' ? '亮色模式' : '暗色模式' }}
      </button>
      <div class="sidebar-bottom-item site-runtime" :title="'建站日期：' + siteStartDate">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"></circle>
          <polyline points="12 6 12 12 16 14"></polyline>
        </svg>
        <div class="site-runtime-text">
          <div class="site-runtime-label">运行时长</div>
          <div class="site-runtime-value">{{ runtime.days }}天 {{ runtime.hours }}:{{ runtime.minutes }}:{{ runtime.seconds }}</div>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'

const props = defineProps({
  menuOpen: Boolean,
  siteProfile: {
    type: Object,
    required: true
  }
})
const emit = defineEmits(['close-menu'])
const appStore = useAppStore()
const router = useRouter()
const drawerRef = ref(null)
const closeButtonRef = ref(null)

const searchKeyword = ref('')

const focusInitial = async () => {
  await nextTick()
  closeButtonRef.value?.focus()
}

const handleDrawerKeydown = (event) => {
  if (!props.menuOpen || event.key !== 'Tab' || !drawerRef.value) return

  const focusable = [...drawerRef.value.querySelectorAll('a[href], button:not([disabled]), input:not([disabled])')]
    .filter((element) => element.offsetParent !== null)
  if (!focusable.length) return

  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

defineExpose({ focusInitial })

const handleSearchSubmit = () => {
  const kw = searchKeyword.value.trim()
  if (!kw) return
  emit('close-menu')
  router.push({ name: 'Search', query: { q: kw } })
}

const siteStartDate = '2026-06-20 00:00:00'

const runtime = reactive({
  days: 0,
  hours: '00',
  minutes: '00',
  seconds: '00'
})

let runtimeTimer = null

const pad = (n) => String(n).padStart(2, '0')

const calculateRuntime = () => {
  const start = new Date(siteStartDate.replace(/-/g, '/')).getTime()
  const now = Date.now()
  const diff = Math.max(0, now - start)

  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  runtime.days = days
  runtime.hours = pad(hours)
  runtime.minutes = pad(minutes)
  runtime.seconds = pad(seconds)
}

onMounted(() => {
  calculateRuntime()
  runtimeTimer = setInterval(calculateRuntime, 1000)
})

onUnmounted(() => {
  if (runtimeTimer) {
    clearInterval(runtimeTimer)
    runtimeTimer = null
  }
})
</script>
