import axios from 'axios'
import { getToken, clearToken } from '@/utils/token'

export const api = axios.create({
  baseURL: '/api',
  timeout: 15000,
})

// 请求拦截：自动带上 JWT
api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截：401 时清除 token 并跳转登录
api.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      clearToken()
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  },
)
