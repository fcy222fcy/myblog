// utils/categoryStyle.js
// 文章分类胶囊配色：保留 6 个老分类的固定配色；新增分类按名称做稳定 hash
// 自动生成色相，避免出现「白底白字」的不可见胶囊。
//
// 使用方式：
//   import { categoryStyle } from '@/utils/categoryStyle'
//   <router-link :style="categoryStyle(article.category)" ... />

// 与 main.css 中 .category-build 等保持一致；新增老分类时两边同步。
const KNOWN_CATEGORY_COLORS = {
  build: '#193c68', // 搭建网站
  dev: '#6d1b85', // 软件开发
  life: '#117361', // 生活记录
  books: '#8B5A3C', // 读书笔记
  tools: '#4A6E8A', // 工具推荐
  study: '#B8860B' // 学习笔记
}

const hashHue = (str) => {
  let h = 0
  for (let i = 0; i < str.length; i++) h = (h * 31 + str.charCodeAt(i)) % 360
  return h
}

export const categoryStyle = (category) => {
  if (!category) return {}
  const slug = (category.slug || '').toString()
  const hex = KNOWN_CATEGORY_COLORS[slug]
  if (hex) return { background: hex }

  // 新分类：基于名称 hash 生成稳定的中等明度色，与老分类观感统一；
  // 白字已由 .category-pill 基础样式保证。
  const seed = (category.name || slug || 'default').toString()
  const h = hashHue(seed)
  return { background: `hsl(${h}, 55%, 38%)` }
}