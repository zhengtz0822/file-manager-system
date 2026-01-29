package model

import "time"

// UploadChunk 分片上传记录
type UploadChunk struct {
	UploadID       string    `json:"upload_id" gorm:"primaryKey;size:36"`
	FileName       string    `json:"file_name" gorm:"size:255;not null"`
	FileSize       int64     `json:"file_size" gorm:"not null"`
	ChunkSize      int       `json:"chunk_size" gorm:"not null"`
	TotalChunks    int       `json:"total_chunks" gorm:"not null"`
	UploadedChunks int       `json:"uploaded_chunks" gorm:"default:0"`
	StoragePath    string    `json:"storage_path" gorm:"size:500"`
	Status         int       `json:"status" gorm:"default:1;comment:1-上传中,2-已完成,0-已取消"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	ExpiredAt      time.Time `json:"expired_at"`
}

// TableName 指定表名
func (UploadChunk) TableName() string {
	return "upload_chunks"
}

// UploadStatus 上传状态
const (
	UploadStatusUploading = 1 // 上传中
	UploadStatusCompleted = 2 // 已完成
	UploadStatusCancelled = 0 // 已取消
)
