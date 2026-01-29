// Shared TypeScript types for File Manager System
// This file is reserved for future use

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

export interface PaginationParams {
  page: number;
  page_size: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  page_size: number;
}
