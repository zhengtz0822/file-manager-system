package model

import "time"

// Document 文档模型
type Document struct {
	ID            string    `json:"id" gorm:"primaryKey;size:36"`
	FileName      string    `json:"file_name" gorm:"size:255;not null"`
	StoragePath   string    `json:"storage_path" gorm:"size:500;not null"`
	FileSize      int64     `json:"file_size" gorm:"not null"`
	FileType      string    `json:"file_type" gorm:"size:100;not null"`
	FileExtension string    `json:"file_extension" gorm:"size:20;not null"`
	MD5Hash       string    `json:"md5_hash" gorm:"size:32"`
	UploadID      string    `json:"upload_id" gorm:"size:36"`
	UploadedBy    string    `json:"uploaded_by" gorm:"size:10;not null;comment:上传者类型:user-用户,app-应用"`
	UserID        *int      `json:"user_id,omitempty" gorm:"comment:用户ID（uploaded_by=user时有值）"`
	AppID         *int      `json:"app_id,omitempty" gorm:"comment:应用ID（uploaded_by=app时有值）"`
	Status        int       `json:"status" gorm:"default:1;comment:1-正常,0-已删除"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Document) TableName() string {
	return "documents"
}

// DocumentStatus 文档状态
const (
	DocumentStatusNormal   = 1 // 正常
	DocumentStatusDeleted = 0 // 已删除
)

// UploadedBy 上传者类型
const (
	UploadedByUser = "user" // 用户上传
	UploadedByApp  = "app"  // 应用上传
)
