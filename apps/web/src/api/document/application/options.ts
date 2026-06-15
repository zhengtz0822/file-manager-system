/**
 * 应用管理模块 - 应用选项 API
 * 用于下拉选择等场景，不包含敏感信息
 */
import { apiClient } from '@/lib/api-client'

/**
 * 应用选项（轻量级，不包含敏感信息）
 */
export interface ApplicationOption {
  id: number
  app_name: string
  app_identifier: string
  app_account: string
}

/**
 * 获取应用选项列表
 * @returns 应用选项列表（不含密钥等敏感信息）
 */
export async function getApplicationOptions(): Promise<ApplicationOption[]> {
  return apiClient.get<ApplicationOption[]>('/applications/options') as unknown as Promise<ApplicationOption[]>
}
