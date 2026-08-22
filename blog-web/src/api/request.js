import axios from 'axios'
import { handleError } from '@/utils/errorHandler'
import { getOrCreateVisitorID } from '@/utils/visitorIdentity'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 可以在这里添加 token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    const visitorID = getOrCreateVisitorID(globalThis.localStorage, globalThis.crypto)
    if (visitorID) {
      config.headers['X-Visitor-ID'] = visitorID
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response.data
  },
  error => {
    // 使用统一的错误处理，但不自动显示消息
    // 让调用方决定是否显示错误消息
    handleError(error, { showMessage: false })
    return Promise.reject(error)
  }
)

export default request
