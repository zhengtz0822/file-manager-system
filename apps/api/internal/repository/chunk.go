package repository

import (
	"file-manager-service/internal/model"
	"time"
)

// ChunkRepository 分片上传数据访问层
type ChunkRepository struct{}

// NewChunkRepository 创建分片仓库
func NewChunkRepository() *ChunkRepository {
	return &ChunkRepository{}
}

// Create 创建分片上传记录
func (r *ChunkRepository) Create(chunk *model.UploadChunk) error {
	return model.DB.Create(chunk).Error
}

// FindByID 根据 UploadID 查找分片上传记录
func (r *ChunkRepository) FindByID(uploadID string) (*model.UploadChunk, error) {
	var chunk model.UploadChunk
	err := model.DB.Where("upload_id = ?", uploadID).First(&chunk).Error
	if err != nil {
		return nil, err
	}
	return &chunk, nil
}

// UpdateUploadedChunks 更新已上传分片数
func (r *ChunkRepository) UpdateUploadedChunks(uploadID string, count int) error {
	return model.DB.Model(&model.UploadChunk{}).
		Where("upload_id = ?", uploadID).
		Update("uploaded_chunks", count).Error
}

// UpdateStatus 更新上传状态
func (r *ChunkRepository) UpdateStatus(uploadID string, status int) error {
	return model.DB.Model(&model.UploadChunk{}).
		Where("upload_id = ?", uploadID).
		Update("status", status).Error
}

// UpdateStoragePath 更新存储路径
func (r *ChunkRepository) UpdateStoragePath(uploadID string, path string) error {
	return model.DB.Model(&model.UploadChunk{}).
		Where("upload_id = ?", uploadID).
		Update("storage_path", path).Error
}

// Delete 删除分片上传记录
func (r *ChunkRepository) Delete(uploadID string) error {
	return model.DB.Where("upload_id = ?", uploadID).Delete(&model.UploadChunk{}).Error
}

// CleanExpiredChunks 清理过期的分片记录
func (r *ChunkRepository) CleanExpiredChunks() error {
	return model.DB.Where("expired_at < ? AND status != ?", time.Now(), model.UploadStatusCompleted).
		Delete(&model.UploadChunk{}).Error
}
