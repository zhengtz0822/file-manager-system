package model

import "time"

// Application 应用模型
type Application struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AppName        string    `json:"app_name" gorm:"uniqueIndex;size:100;not null"`
	AppIdentifier  string    `json:"app_identifier" gorm:"uniqueIndex;size:16;not null;comment:应用标识（英文数字，最长16位，用于文件存储路径）"`
	AppAccount     string    `json:"app_account" gorm:"uniqueIndex;size:50;not null"`
	AppSecret      string    `json:"app_secret" gorm:"size:100;not null"`
	Status         int       `json:"status" gorm:"default:1;comment:状态:1-启用,0-禁用"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Application) TableName() string {
	return "applications"
}
