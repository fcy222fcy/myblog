import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getArticleList, getArticleDetail, getArchives } from '../api/article'
import { handleError } from '../utils/errorHandler'
import { recordContentView } from '../api/view'

export const useArticleStore = defineStore('article', () => {
  const articles = ref([])
  const currentArticle = ref(null)
  const archives = ref([])
  const total = ref(0)
  const loading = ref(false)
  const error = ref(null)

  const fetchArticles = async (params = {}) => {
    loading.value = true
    error.value = null
    try {
      const res = await getArticleList(params)
      articles.value = res.data.list || []
      total.value = res.data.total || 0
    } catch (err) {
      error.value = handleError(err, { showMessage: false })
      console.error('获取文章列表失败:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchArticleDetail = async (slug) => {
    loading.value = true
    error.value = null
    try {
      const res = await getArticleDetail(slug)
      currentArticle.value = res.data
      try {
        const viewRes = await recordContentView('article', res.data.id)
        if (currentArticle.value?.id === res.data.id && viewRes.data) {
          currentArticle.value.view_count = viewRes.data.view_count
        }
      } catch (viewError) {
        console.warn('记录文章浏览量失败:', viewError)
      }
    } catch (err) {
      error.value = handleError(err, { showMessage: false })
      console.error('获取文章详情失败:', err)
    } finally {
      loading.value = false
    }
  }

  const fetchArchives = async () => {
    loading.value = true
    error.value = null
    try {
      const res = await getArchives()
      archives.value = res.data || []
    } catch (err) {
      error.value = handleError(err, { showMessage: false })
      console.error('获取文章归档失败:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    articles,
    currentArticle,
    archives,
    total,
    loading,
    error,
    fetchArticles,
    fetchArticleDetail,
    fetchArchives
  }
})
