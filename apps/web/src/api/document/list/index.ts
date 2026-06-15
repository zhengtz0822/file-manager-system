/**
 * 文档管理模块 - API 方法
 */
import { apiClient } from '@/lib/api-client'
import type {
  Document,
  DocumentListResponse,
  GetDocumentListParams,
} from './types'

// 导出类型，方便外部使用
export type * from './types'

/**
 * 获取文档列表
 * @param params 查询参数
 * @returns 文档列表
 */
export async function getDocumentList(
  params: GetDocumentListParams
): Promise<DocumentListResponse> {
  return apiClient.get<DocumentListResponse>('/documents', {
    params,
  }) as unknown as Promise<DocumentListResponse>
}

/**
 * 获取文档详情
 * @param id 文档ID
 * @returns 文档详情
 */
export async function getDocument(id: string): Promise<Document> {
  return apiClient.get<Document>(`/documents/${id}`) as unknown as Promise<Document>
}

/**
 * 删除文档
 * @param id 文档ID
 */
export async function deleteDocument(id: string): Promise<void> {
  await apiClient.delete(`/documents/${id}`)
}
