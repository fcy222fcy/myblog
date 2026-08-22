/**
 * 列表页滚动位置记忆（基于 sessionStorage，单次消费）
 *
 * 使用场景：首页文章列表 → 点击标签/分类进入归档页 → 点击「返回首页」
 * 返回时恢复离开前的滚动位置，保持用户浏览上下文。
 *
 * - 离开首页前调用 saveScrollPosition() 保存当前滚动位置
 * - 返回首页、文章数据渲染完成后调用 takeScrollPosition() 取出并恢复
 * - 记忆为一次性：读取即清除，避免过期位置残留
 * - sessionStorage 不可用时静默降级（不影响页面正常功能）
 */
const STORAGE_KEY = 'blog:scroll-restore'

export function saveScrollPosition() {
  try {
    const y = window.scrollY || document.documentElement.scrollTop || 0
    sessionStorage.setItem(STORAGE_KEY, String(y))
  } catch {
    /* 忽略：sessionStorage 不可用或超出配额 */
  }
}

export function takeScrollPosition() {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    sessionStorage.removeItem(STORAGE_KEY)
    if (raw === null) return null
    const y = Number(raw)
    return Number.isFinite(y) && y > 0 ? y : null
  } catch {
    return null
  }
}

export function hasScrollMemory() {
  try {
    return sessionStorage.getItem(STORAGE_KEY) !== null
  } catch {
    return false
  }
}
