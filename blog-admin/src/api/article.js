import request from './request'

export const getArticleList = (params) => request.get('/admin/articles', { params })
export const getDraftArticles = (params) => request.get('/admin/articles', { params: { ...params, status: 'draft' } })
export const getPublishedArticles = (params) => request.get('/admin/articles', { params: { ...params, status: 'published' } })
export const getArticleDetail = (id) => request.get(`/admin/articles/${id}`)
export const createArticle = (data) => request.post('/admin/articles', data)
export const updateArticle = (id, data) => request.put(`/admin/articles/${id}`, data)
export const updateArticleStatus = (id, status) => request.put(`/admin/articles/${id}/status`, { status })
export const deleteArticle = (id) => request.delete(`/admin/articles/${id}`)
export const batchDeleteArticles = (ids) => request.post('/admin/articles/batch-delete', { ids })
