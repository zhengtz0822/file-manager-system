package repository

import (
	"file-manager-service/internal/model"
)

type ApplicationRepository struct{}

func NewApplicationRepository() *ApplicationRepository {
	return &ApplicationRepository{}
}

// Create 创建应用
func (r *ApplicationRepository) Create(app *model.Application) error {
	return model.DB.Create(app).Error
}

// GetByID 根据ID获取应用
func (r *ApplicationRepository) GetByID(id int) (*model.Application, error) {
	var app model.Application
	err := model.DB.Where("id = ?", id).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByAccount 根据应用账号获取应用
func (r *ApplicationRepository) GetByAccount(account string) (*model.Application, error) {
	var app model.Application
	err := model.DB.Where("app_account = ?", account).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetByIdentifier 根据应用标识获取应用
func (r *ApplicationRepository) GetByIdentifier(identifier string) (*model.Application, error) {
	var app model.Application
	err := model.DB.Where("app_identifier = ?", identifier).First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List 获取应用列表
func (r *ApplicationRepository) List(offset, limit int) ([]model.Application, int64, error) {
	var apps []model.Application
	var total int64

	err := model.DB.Model(&model.Application{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = model.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&apps).Error
	return apps, total, err
}

// UpdateStatus 更新应用状态
func (r *ApplicationRepository) UpdateStatus(id int, status int) error {
	return model.DB.Model(&model.Application{}).Where("id = ?", id).Update("status", status).Error
}

// Delete 删除应用
func (r *ApplicationRepository) Delete(id int) error {
	return model.DB.Delete(&model.Application{}, id).Error
}
