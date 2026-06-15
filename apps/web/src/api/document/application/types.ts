/**
 * 应用管理模块 - 类型定义
 */

// 应用实体
export interface Application {
  id: number
  app_name: string
  app_identifier: string
  app_account: string
  app_secret: string
  status: number
  created_at: string
}

// 创建应用请求
export interface CreateApplicationRequest {
  app_name: string
  app_identifier: string
}

// 应用列表响应
export interface ApplicationListResponse {
  applications: Application[]
  total: number
  page: number
  page_size: number
}
