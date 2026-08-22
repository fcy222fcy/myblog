/**
 * 列表页状态快照（滚动位置 + 页码），按路由 path 隔离，一次性消费
 *
 * 使用场景：列表页（首页/归档/分类/标签/搜索）→ 点文章进详情 → 返回时
 * 恢复到来源列表页的页码与滚动位置，保持浏览上下文。
 *
 * - 离开列表页时调用 saveListState(path, { page }) 保存当前状态
 * - 列表页重新挂载（popstate 触发）时调用 takeListState(path) 取出并恢复
 * - 读取即清除：避免过期位置残留；同时支持多级返回链路（每个 path 独立存储）
 * - sessionStorage 不可用时静默降级，不影响页面正常功能
 */

const PREFIX = 'blog:list-state:'

function safeSet(key, value) {
  try { sessionStorage.setItem(key, value) } catch { /* ignore */ }
}
function safeGet(key) {
  try { return sessionStorage.getItem(key) } catch { return null }
}
function safeRemove(key) {
  try { sessionStorage.removeItem(key) } catch { /* ignore */ }
}

/**
 * 保存当前列表页状态
 * @param {string} routePath 路由完整路径，如 '/tag/5'、'/category/2'
 * @param {object} [opts]
 * @param {number} [opts.page=1] 当前页码（无分页的列表页可省略）
 */
export function saveListState(routePath, { page = 1 } = {}) {
  if (!routePath) return
  const y = window.scrollY || document.documentElement.scrollTop || 0
  safeSet(PREFIX + routePath, JSON.stringify({ scrollY: y, page }))
}

/**
 * 取出并清除指定列表页的快照（一次性消费）
 * @param {string} routePath
 * @returns {{scrollY: number, page: number} | null}
 */
export function takeListState(routePath) {
  if (!routePath) return null
  const raw = safeGet(PREFIX + routePath)
  safeRemove(PREFIX + routePath)
  if (!raw) return null
  try {
    const s = JSON.parse(raw)
    if (!s || typeof s.scrollY !== 'number') return null
    return {
      scrollY: Math.max(0, s.scrollY | 0),
      page: Number.isFinite(s.page) && s.page > 0 ? (s.page | 0) : 1
    }
  } catch {
    return null
  }
}

/**
 * 主动清除指定列表页的快照（备用，一般无需调用）
 */
export function clearListState(routePath) {
  if (!routePath) return
  safeRemove(PREFIX + routePath)
}
