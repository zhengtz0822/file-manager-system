package service

import (
	"file-manager-service/internal/model"
	"file-manager-service/internal/pkg/storage"
	"file-manager-service/internal/pkg/uuid"
	"file-manager-service/internal/repository"
	"os"
	"path/filepath"
	"strings"
)

// DocumentService 文档服务
type DocumentService struct {
	docRepo   *repository.DocumentRepository
	chunkRepo *repository.ChunkRepository
}

// NewDocumentService 创建文档服务
func NewDocumentService() *DocumentService {
	return &DocumentService{
		docRepo:   repository.NewDocumentRepository(),
		chunkRepo: repository.NewChunkRepository(),
	}
}

// InitUploadRequest 初始化上传请求
type InitUploadRequest struct {
	FileName  string `json:"file_name" binding:"required"`
	FileSize  int64  `json:"file_size" binding:"required"`
	ChunkSize int    `json:"chunk_size" binding:"required"`
}

// InitUploadResponse 初始化上传响应
type InitUploadResponse struct {
	UploadID     string `json:"upload_id"`
	TotalChunks  int    `json:"total_chunks"`
	ChunkSize    int    `json:"chunk_size"`
	FileSize     int64  `json:"file_size"`
}

// InitUpload 初始化分片上传
func (s *DocumentService) InitUpload(req *InitUploadRequest) (*InitUploadResponse, error) {
	// 检查文件类型
	if !storage.IsAllowedExtension(req.FileName) {
		return nil, ErrInvalidFileType
	}

	// 计算总分片数
	totalChunks := int(req.FileSize) / req.ChunkSize
	if int(req.FileSize)%req.ChunkSize != 0 {
		totalChunks++
	}

	uploadID := uuid.Generate()

	chunk := &model.UploadChunk{
		UploadID:       uploadID,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		ChunkSize:      req.ChunkSize,
		TotalChunks:    totalChunks,
		UploadedChunks: 0,
		Status:         model.UploadStatusUploading,
	}

	if err := s.chunkRepo.Create(chunk); err != nil {
		return nil, err
	}

	return &InitUploadResponse{
		UploadID:    uploadID,
		TotalChunks: totalChunks,
		ChunkSize:   req.ChunkSize,
		FileSize:    req.FileSize,
	}, nil
}

// UploadChunk 上传分片
func (s *DocumentService) UploadChunk(uploadID string, chunkNumber int, filePath string) error {
	chunk, err := s.chunkRepo.FindByID(uploadID)
	if err != nil {
		return ErrUploadNotFound
	}

	if chunk.Status != model.UploadStatusUploading {
		return ErrUploadStatus
	}

	// 更新已上传分片数
	if err := s.chunkRepo.UpdateUploadedChunks(uploadID, chunk.UploadedChunks+1); err != nil {
		return err
	}

	return nil
}

// CompleteUploadRequest 完成上传请求
type CompleteUploadRequest struct {
	UploadID string `json:"upload_id" binding:"required"`
}

// CompleteUploadResponse 完成上传响应
type CompleteUploadResponse struct {
	DocumentID string `json:"document_id"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
}

// CompleteUpload 完成上传
func (s *DocumentService) CompleteUpload(req *CompleteUploadRequest) (*CompleteUploadResponse, error) {
	chunk, err := s.chunkRepo.FindByID(req.UploadID)
	if err != nil {
		return nil, ErrUploadNotFound
	}

	if chunk.UploadedChunks != chunk.TotalChunks {
		return nil, ErrIncompleteChunks
	}

	// 生成文档 ID 和存储路径
	documentID := uuid.Generate()
	storagePath := storage.GenerateStoragePath(documentID, chunk.FileName)

	// 合并分片
	if err := storage.MergeChunks(req.UploadID, chunk.TotalChunks, storagePath); err != nil {
		return nil, err
	}

	// 计算文件 MD5
	md5Hash, _ := storage.GetFileMD5(storagePath)

	// 创建文档记录
	doc := &model.Document{
		ID:            documentID,
		FileName:      chunk.FileName,
		StoragePath:   storagePath,
		FileSize:      chunk.FileSize,
		FileType:      storage.GetMimeType(chunk.FileName),
		FileExtension: strings.TrimPrefix(filepath.Ext(chunk.FileName), "."),
		MD5Hash:       md5Hash,
		UploadID:      req.UploadID,
		Status:        model.DocumentStatusNormal,
	}

	if err := s.docRepo.Create(doc); err != nil {
		return nil, err
	}

	// 更新分片上传状态
	s.chunkRepo.UpdateStatus(req.UploadID, model.UploadStatusCompleted)
	s.chunkRepo.UpdateStoragePath(req.UploadID, storagePath)

	return &CompleteUploadResponse{
		DocumentID: documentID,
		FileName:   doc.FileName,
		FileSize:   doc.FileSize,
	}, nil
}

// GetDocument 获取文档信息
func (s *DocumentService) GetDocument(id string) (*model.Document, error) {
	return s.docRepo.FindByID(id)
}

// DeleteDocument 删除文档
func (s *DocumentService) DeleteDocument(id string) error {
	doc, err := s.docRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 删除物理文件
	if err := os.Remove(doc.StoragePath); err != nil {
		return err
	}

	// 软删除数据库记录
	return s.docRepo.Delete(id)
}

// ListDocuments 文档列表
type ListDocumentsRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// ListDocumentsResponse 文档列表响应
type ListDocumentsResponse struct {
	Documents []*model.Document `json:"documents"`
	Total     int64            `json:"total"`
	Page      int              `json:"page"`
	PageSize  int              `json:"page_size"`
}

// List 获取文档列表
func (s *DocumentService) List(req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	docs, total, err := s.docRepo.List(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, err
	}

	return &ListDocumentsResponse{
		Documents: docs,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}, nil
}

// CancelUpload 取消上传
func (s *DocumentService) CancelUpload(uploadID string) error {
	chunk, err := s.chunkRepo.FindByID(uploadID)
	if err != nil {
		return ErrUploadNotFound
	}

	// 删除已上传的分片文件
	storage.CleanChunks(uploadID, chunk.UploadedChunks)

	// 删除分片记录
	return s.chunkRepo.Delete(uploadID)
}

// 错误定义
var (
	ErrInvalidFileType     = &ServiceError{Code: 400, Message: "不支持的文件类型"}
	ErrUploadNotFound      = &ServiceError{Code: 404, Message: "上传记录不存在"}
	ErrUploadStatus        = &ServiceError{Code: 400, Message: "上传状态不正确"}
	ErrIncompleteChunks    = &ServiceError{Code: 400, Message: "分片不完整"}
	ErrDocumentNotFound    = &ServiceError{Code: 404, Message: "文档不存在"}
)

// ServiceError 服务错误
type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
