/**
 * 文档管理模块 - 类型定义
 */

/**
 * 文档信息
 */
export interface Document {
  id: string
  file_name: string
  storage_path: string
  file_size: number
  file_type: string
  file_extension: string
  md5_hash: string
  upload_id: string
  uploaded_by: 'user' | 'app'
  user_id?: number
  app_id?: number
  status: number
  created_at: string
  updated_at: string
}

/**
 * 文档列表响应
 */
export interface DocumentListResponse {
  documents: Document[]
  total: number
  page: number
  page_size: number
}

/**
 * 获取文档列表请求参数
 */
export interface GetDocumentListParams {
  page: number
  page_size: number
  keyword?: string
  app_identifier?: string
}
