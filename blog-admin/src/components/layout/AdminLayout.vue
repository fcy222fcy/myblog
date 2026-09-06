<template>
  <div class="app-container" :data-scheme="scheme">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ open: menuOpen }">
      <div class="sidebar-header">
        <div class="sidebar-avatar" @click="handleAvatarClick">
          <span v-if="!userInfo.avatar">{{ userInfo.nickname?.charAt(0) || 'U' }}</span>
          <img v-else :src="userInfo.avatar" alt="头像">
          <div class="sidebar-avatar-edit">更换头像</div>
        </div>
        <input ref="avatarInput" type="file" accept="image/*" style="display: none;" @change="handleAvatarUpload">
        <div class="sidebar-username">{{ userInfo.nickname || '用户' }}</div>
        <div class="sidebar-desc">{{ userInfo.bio || '日常落灰的个人博客' }}</div>
        <button class="sidebar-edit-btn" @click="showProfileDialog = true">编辑资料</button>
        <div class="sidebar-social">
          <a v-for="link in userInfo.socialLinks" :key="link.name" :href="link.url" class="sidebar-social-item" target="_blank" :title="link.name">
            {{ link.icon || '🔗' }}
          </a>
        </div>
      </div>

      <nav class="sidebar-menu">
        <div v-for="item in menuItems" :key="item.path"
             class="menu-item" :class="{ active: currentPath === item.path }"
             @click="navigateTo(item.path)">
          <span class="menu-item-icon" v-html="item.icon"></span>
          <span class="menu-item-text">{{ item.title }}</span>
        </div>
      </nav>

      <div class="sidebar-features">
        <div class="feature-item" @click="toggleScheme">
          <span class="feature-item-icon"><svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path></svg></span>
          <span>{{ scheme === 'dark' ? '亮色模式' : '暗色模式' }}</span>
        </div>
        <div class="feature-item logout-feature" @click="handleLogout">
          <span class="feature-item-icon"><svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path><polyline points="16 17 21 12 16 7"></polyline><line x1="21" y1="12" x2="9" y2="12"></line></svg></span>
          <span>退出登录</span>
        </div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 顶部导航 -->
      <header class="header">
        <div class="header-left">
          <button class="mobile-menu-btn" @click="menuOpen = !menuOpen">☰</button>
          <h1 class="page-title">{{ currentPageTitle }}</h1>
        </div>
        <div class="header-right">
          <button class="header-btn" title="通知">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path><path d="M13.73 21a2 2 0 0 1-3.46 0"></path></svg>
          </button>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="content-area">
        <router-view />
      </main>
    </div>

    <!-- 回到顶部按钮 -->
    <BackToTop />

    <!-- 编辑个人资料弹窗 -->
    <div class="modal-overlay" :class="{ active: showProfileDialog }" @click.self="closeProfileDialog">
      <div class="modal profile-modal">
        <div class="modal-header">
          <h3 class="modal-title">账号设置</h3>
          <button class="modal-close" @click="closeProfileDialog">×</button>
        </div>
        <div class="profile-tabs">
          <div class="profile-tab" :class="{ active: profileTab === 'basic' }" @click="profileTab = 'basic'">基本资料</div>
          <div class="profile-tab" :class="{ active: profileTab === 'password' }" @click="profileTab = 'password'">修改密码</div>
        </div>
        <div class="modal-body">
          <!-- 基本资料 -->
          <div v-if="profileTab === 'basic'">
            <div class="form-group">
              <label class="form-label">昵称</label>
              <input v-model="profileForm.nickname" class="form-input" placeholder="输入昵称">
            </div>
            <div class="form-group">
              <label class="form-label">邮箱 <span class="label-tip">（用于博主身份识别，建议与系统配置一致）</span></label>
              <input v-model="profileForm.email" type="email" class="form-input" placeholder="输入邮箱">
            </div>
            <div class="form-group">
              <label class="form-label">个人简介</label>
              <textarea v-model="profileForm.bio" class="form-textarea" placeholder="输入个人简介..." rows="3"></textarea>
            </div>
          </div>
          <!-- 修改密码 -->
          <div v-else>
            <div class="form-group">
              <label class="form-label">当前密码</label>
              <input v-model="passwordForm.old_password" type="password" class="form-input" placeholder="输入当前密码">
            </div>
            <div class="form-group">
              <label class="form-label">新密码</label>
              <input v-model="passwordForm.new_password" type="password" class="form-input" placeholder="至少 6 位">
            </div>
            <div class="form-group">
              <label class="form-label">确认新密码</label>
              <input v-model="passwordForm.confirm_password" type="password" class="form-input" placeholder="再次输入新密码">
            </div>
            <div v-if="passwordTip" class="password-tip" :class="{ error: passwordTipIsError }">{{ passwordTip }}</div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeProfileDialog">取消</button>
          <button v-if="profileTab === 'basic'" class="btn btn-primary" @click="saveProfile" :disabled="savingProfile">
            {{ savingProfile ? '保存中...' : '保存' }}
          </button>
          <button v-else class="btn btn-primary" @click="savePassword" :disabled="savingPassword">
            {{ savingPassword ? '修改中...' : '修改密码' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getUserInfo, updateUserInfo, changePassword } from '../../api/auth'
import { uploadFile, MEDIA_CATEGORIES } from '../../api/media'
import { ElMessage, ElMessageBox } from 'element-plus'
import BackToTop from '../common/BackToTop.vue'

const route = useRoute()
const router = useRouter()
const menuOpen = ref(false)
const showProfileDialog = ref(false)
const profileTab = ref('basic')
const scheme = ref(localStorage.getItem('scheme') || 'light')
const avatarInput = ref(null)
const savingProfile = ref(false)
const savingPassword = ref(false)
const passwordTip = ref('')
const passwordTipIsError = ref(false)

const userInfo = ref({
  nickname: 'Fu Chengyan',
  bio: '日常落灰的个人博客，分享 Golang、AI 和 NAS 折腾经验',
  avatar: '',
  email: '',
  socialLinks: []
})

const profileForm = ref({
  nickname: '',
  email: '',
  bio: ''
})

const passwordForm = ref({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const closeProfileDialog = () => {
  showProfileDialog.value = false
  profileTab.value = 'basic'
  passwordTip.value = ''
  passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
}

const menuItems = [
  { path: '/dashboard', title: '仪表盘', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path><polyline points="9 22 9 12 15 12 15 22"></polyline></svg>' },
  { path: '/articles', title: '文章', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path><polyline points="14 2 14 8 20 8"></polyline><line x1="16" y1="13" x2="8" y2="13"></line><line x1="16" y1="17" x2="8" y2="17"></line></svg>' },
  { path: '/categories', title: '分类', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"></path></svg>' },
  { path: '/tags', title: '标签', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20.59 13.41l-7.17 7.17a2 2 0 0 1-2.83 0L2 12V2h10l8.59 8.59a2 2 0 0 1 0 2.82z"></path><line x1="7" y1="7" x2="7.01" y2="7"></line></svg>' },
  { path: '/comments', title: '评论', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>' },

  { path: '/daily-question', title: '每日一问', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path><line x1="12" y1="17" x2="12.01" y2="17"></line></svg>' },
  { path: '/about', title: '关于我', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>' },
  { path: '/audit', title: '操作日志', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>' },
]

const currentPath = computed(() => route.path)
const currentPageTitle = computed(() => {
  // 精确匹配
  const exactItem = menuItems.find(i => i.path === route.path)
  if (exactItem) return exactItem.title

  // 前缀匹配（处理子路由，如 /articles/123）
  const prefixItem = menuItems.find(i => route.path.startsWith(i.path + '/') || route.path.startsWith(i.path + '?'))
  if (prefixItem) return prefixItem.title

  return '仪表盘'
})

const navigateTo = (path) => {
  router.push(path)
  menuOpen.value = false
}

const toggleScheme = () => {
  scheme.value = scheme.value === 'light' ? 'dark' : 'light'
  localStorage.setItem('scheme', scheme.value)
}

const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '退出',
      cancelButtonText: '取消',
      type: 'warning'
    })
    localStorage.removeItem('token')
    localStorage.removeItem('comment_user')
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (e) {}
}

const handleAvatarClick = () => {
  avatarInput.value?.click()
}

const handleAvatarUpload = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件')
    return
  }

  // 验证文件大小 (2MB)
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.error('图片大小不能超过 2MB')
    return
  }

  try {
    const result = await uploadFile(file, MEDIA_CATEGORIES.AVATAR)

    if (result.code === 0) {
      // 更新用户头像
      await updateUserInfo({ avatar: result.data.url })
      userInfo.value.avatar = result.data.url
      ElMessage.success('头像更新成功')
    } else {
      ElMessage.error(result.message || '上传失败')
    }
  } catch (error) {
    console.error('头像上传失败:', error)
    ElMessage.error('上传失败，请重试')
  }

  // 清空 input 的值
  event.target.value = ''
}

const saveProfile = async () => {
  savingProfile.value = true
  try {
    const payload = {
      nickname: profileForm.value.nickname,
      bio: profileForm.value.bio
    }
    if (profileForm.value.email) payload.email = profileForm.value.email

    await updateUserInfo(payload)
    userInfo.value.nickname = profileForm.value.nickname
    userInfo.value.bio = profileForm.value.bio
    userInfo.value.email = profileForm.value.email || userInfo.value.email
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error(e?.response?.data?.message || '保存失败')
  } finally {
    savingProfile.value = false
  }
}

const savePassword = async () => {
  passwordTip.value = ''
  passwordTipIsError.value = true

  if (!passwordForm.value.old_password) {
    passwordTip.value = '请输入当前密码'
    return
  }
  if (!passwordForm.value.new_password || passwordForm.value.new_password.length < 6) {
    passwordTip.value = '新密码至少 6 位'
    return
  }
  if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
    passwordTip.value = '两次输入的新密码不一致'
    return
  }

  savingPassword.value = true
  try {
    await changePassword({
      old_password: passwordForm.value.old_password,
      new_password: passwordForm.value.new_password
    })
    passwordTipIsError.value = false
    passwordTip.value = '密码修改成功！'
    passwordForm.value = { old_password: '', new_password: '', confirm_password: '' }
    ElMessage.success('密码修改成功')
  } catch (e) {
    passwordTipIsError.value = true
    passwordTip.value = e?.response?.data?.message || '密码修改失败，请检查当前密码是否正确'
  } finally {
    savingPassword.value = false
  }
}

onMounted(async () => {
  profileForm.value.nickname = userInfo.value.nickname
  profileForm.value.bio = userInfo.value.bio
  profileForm.value.email = userInfo.value.email
  try {
    const res = await getUserInfo()
    if (res.data) {
      userInfo.value = { ...userInfo.value, ...res.data }
      profileForm.value.nickname = res.data.nickname || ''
      profileForm.value.bio = res.data.bio || ''
      profileForm.value.email = res.data.email || ''
    }
  } catch (e) {}
})
</script>

<style scoped>
/* 退出登录按钮 */
.logout-feature {
  color: var(--danger-color) !important;
  font-weight: 500;
  transition: all 0.2s ease;
  border-radius: 6px;
}
.logout-feature:hover {
  background: rgba(239, 68, 68, 0.08) !important;
  color: var(--danger-color) !important;
  opacity: 0.85;
}
.logout-feature:hover .feature-item-icon {
  color: var(--danger-color);
}

/* 资料弹窗 */
.profile-modal {
  width: 460px;
  max-width: 92%;
}

.profile-tabs {
  display: flex;
  border-bottom: 1px solid var(--card-separator-color);
  padding: 0 20px;
  background: var(--card-background-selected);
}

.profile-tab {
  padding: 12px 18px;
  font-size: 0.9rem;
  color: var(--card-text-color-secondary);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: all 0.2s;
  font-weight: 500;
}

.profile-tab:hover {
  color: var(--card-text-color-main);
}

.profile-tab.active {
  color: var(--accent-color);
  border-bottom-color: var(--accent-color);
  font-weight: 600;
}

/* 表单：仅保留 AdminLayout 特有属性，其余复用全局 */
.form-group {
  margin-bottom: 16px;
}

.label-tip {
  font-size: 0.78rem;
  color: var(--card-text-color-tertiary);
  font-weight: normal;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

/* 密码提示 */
.password-tip {
  margin-top: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 0.85rem;
  background: rgba(16, 185, 129, 0.1);
  color: var(--success-color);
  border: 1px solid rgba(16, 185, 129, 0.25);
}

.password-tip.error {
  background: rgba(239, 68, 68, 0.08);
  color: var(--danger-color);
  border-color: rgba(239, 68, 68, 0.25);
}
</style>
