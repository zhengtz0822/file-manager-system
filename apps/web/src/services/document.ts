import request from './api';

export interface InitUploadRequest {
  file_name: string;
  file_size: number;
  chunk_size: number;
}

export interface InitUploadResponse {
  upload_id: string;
  total_chunks: number;
  chunk_size: number;
  file_size: number;
}

export interface CompleteUploadRequest {
  upload_id: string;
}

export interface CompleteUploadResponse {
  document_id: string;
  file_name: string;
  file_size: number;
}

export interface Document {
  id: string;
  file_name: string;
  storage_path: string;
  file_size: number;
  file_type: string;
  file_extension: string;
  md5_hash: string;
  status: number;
  created_at: string;
  updated_at: string;
}

export interface DocumentListResponse {
  documents: Document[];
  total: number;
  page: number;
  page_size: number;
}

export interface DocumentListRequest {
  page: number;
  page_size: number;
  keyword?: string;
}

// 初始化上传
export async function initUpload(data: InitUploadRequest): Promise<InitUploadResponse> {
  return request.post('/documents/chunks/init', data);
}

// 上传分片
export async function uploadChunk(
  uploadId: string,
  chunkNumber: number,
  file: File
): Promise<void> {
  const formData = new FormData();
  formData.append('upload_id', uploadId);
  formData.append('chunk_number', String(chunkNumber));
  formData.append('file', file);

  return request.post('/documents/chunks/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

// 完成上传
export async function completeUpload(data: CompleteUploadRequest): Promise<CompleteUploadResponse> {
  return request.post('/documents/chunks/complete', data);
}

// 取消上传
export async function cancelUpload(uploadId: string): Promise<void> {
  return request.delete(`/documents/chunks/${uploadId}`);
}

// 获取文档列表
export async function getDocumentList(params: DocumentListRequest): Promise<DocumentListResponse> {
  return request.get('/documents', { params });
}

// 获取文档详情
export async function getDocument(id: string): Promise<Document> {
  return request.get(`/documents/${id}`);
}

// 删除文档
export async function deleteDocument(id: string): Promise<void> {
  return request.delete(`/documents/${id}`);
}

// 下载文档
export function downloadDocument(id: string): string {
  return `/api/v1/documents/${id}/download`;
}

// 预览文档
export function previewDocument(id: string): string {
  return `/api/v1/documents/${id}/preview`;
}
