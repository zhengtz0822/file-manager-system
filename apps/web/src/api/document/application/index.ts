/**
 * 应用管理模块 - API 方法
 */
import { apiClient } from '@/lib/api-client'
import type {
  Application,
  CreateApplicationRequest,
  ApplicationListResponse,
} from './types'

// 导出类型，方便外部使用
export type * from './types'

/**
 * 创建应用
 * @param data 应用信息
 * @returns 创建的应用信息（包含生成的账号和密钥）
 */
export async function createApplication(
  data: CreateApplicationRequest
): Promise<Application> {
  return apiClient.post<Application>('/applications', data) as unknown as Promise<Application>
}

/**
 * 获取应用列表
 * @param params 分页参数
 * @returns 应用列表
 */
export async function getApplicationList(params: {
  page?: number
  page_size?: number
}): Promise<ApplicationListResponse> {
  return apiClient.get<ApplicationListResponse>('/applications', {
    params,
  }) as unknown as Promise<ApplicationListResponse>
}

/**
 * 获取应用详情
 * @param id 应用ID
 * @returns 应用详情
 */
export async function getApplication(id: number): Promise<Application> {
  return apiClient.get<Application>(`/applications/${id}`) as unknown as Promise<Application>
}

/**
 * 更新应用状态
 * @param id 应用ID
 * @param status 状态（1-启用，0-禁用）
 */
export async function updateApplicationStatus(
  id: number,
  status: number
): Promise<void> {
  await apiClient.put(`/applications/${id}/status`, { status })
}

/**
 * 删除应用
 * @param id 应用ID
 */
export async function deleteApplication(id: number): Promise<void> {
  await apiClient.delete(`/applications/${id}`)
}

