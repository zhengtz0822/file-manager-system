import { apiClient } from '@/lib/api-client'

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: {
    id: number
    userid: string
    name: string
    username: string
    created_at: string
  }
}

/**
 * 用户登录
 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
  return apiClient.post<LoginResponse>('/auth/login', data) as unknown as Promise<LoginResponse>
}

/**
 * 用户注册
 */
export async function register(data: {
  username: string
  password: string
}): Promise<void> {
  await apiClient.post<void>('/auth/register', data)
}

/**
 * 用户登出
 */
export async function logout(): Promise<void> {
  await apiClient.post<void>('/auth/logout')
}
