import axios from 'axios'
import { useAuthStore } from '@/stores/auth-store'

/**
 * 通用 API 响应类型
 */
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  error?: string
}

// 创建 axios 实例
export const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加认证 token
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().auth.accessToken
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 统一错误处理 + 完全解包
apiClient.interceptors.response.use(
  (response) => {
    // 直接返回 response.data.data
    // 这样使用者就可以直接拿到业务数据
    return response.data.data
  },
  (error) => {
    // 处理 401 未授权错误
    if (error.response?.status === 401) {
      useAuthStore.getState().auth.reset()
      window.location.href = '/sign-in'
    }
    // 抛出业务错误信息
    const businessError = error.response?.data?.message || error.message
    return Promise.reject(new Error(businessError))
  }
)
