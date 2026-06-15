package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"file-manager-service/internal/model"
	"file-manager-service/internal/repository"
	"regexp"
	"strings"
)

type ApplicationService struct {
	appRepo *repository.ApplicationRepository
}

func NewApplicationService() *ApplicationService {
	return &ApplicationService{
		appRepo: repository.NewApplicationRepository(),
	}
}

// CreateApplicationRequest 创建应用请求
type CreateApplicationRequest struct {
	AppName       string `json:"app_name" binding:"required"`
	AppIdentifier string `json:"app_identifier" binding:"required"`
}

// CreateApplicationResponse 创建应用响应
type CreateApplicationResponse struct {
	ID            int    `json:"id"`
	AppName       string `json:"app_name"`
	AppIdentifier string `json:"app_identifier"`
	AppAccount    string `json:"app_account"`
	AppSecret     string `json:"app_secret"`
	Status        int    `json:"status"`
	CreatedAt     string `json:"created_at"`
}

// ListApplicationsResponse 应用列表响应
type ListApplicationsResponse struct {
	Applications []CreateApplicationResponse `json:"applications"`
	Total        int64                       `json:"total"`
	Page         int                         `json:"page"`
	PageSize     int                         `json:"page_size"`
}

// generateAppAccount 生成应用账号
func (s *ApplicationService) generateAppAccount() (string, error) {
	// 生成16位随机字符串
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "APP_" + hex.EncodeToString(bytes), nil
}

// generateAppSecret 生成应用密钥
func (s *ApplicationService) generateAppSecret() (string, error) {
	// 生成64位随机字符串（32字节转64位十六进制）
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Create 创建应用
func (s *ApplicationService) Create(req *CreateApplicationRequest) (*CreateApplicationResponse, error) {
	// 验证应用名称
	if strings.TrimSpace(req.AppName) == "" {
		return nil, errors.New("应用名称不能为空")
	}

	// 验证应用标识
	req.AppIdentifier = strings.TrimSpace(req.AppIdentifier)
	if req.AppIdentifier == "" {
		return nil, errors.New("应用标识不能为空")
	}
	if len(req.AppIdentifier) > 16 {
		return nil, errors.New("应用标识最长16位")
	}
	// 验证只包含英文和数字
	matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", req.AppIdentifier)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, errors.New("应用标识只能包含英文和数字，不能包含符号")
	}

	// 检查应用标识是否已存在
	_, err = s.appRepo.GetByIdentifier(req.AppIdentifier)
	if err == nil {
		return nil, errors.New("应用标识已存在")
	}

	// 生成应用账号
	appAccount, err := s.generateAppAccount()
	if err != nil {
		return nil, err
	}

	// 检查账号是否已存在
	_, err = s.appRepo.GetByAccount(appAccount)
	if err == nil {
		// 账号已存在，重新生成（极低概率）
		return s.Create(req)
	}

	// 生成应用密钥
	appSecret, err := s.generateAppSecret()
	if err != nil {
		return nil, err
	}

	app := &model.Application{
		AppName:       req.AppName,
		AppIdentifier: req.AppIdentifier,
		AppAccount:    appAccount,
		AppSecret:     appSecret,
		Status:        1,
	}

	if err := s.appRepo.Create(app); err != nil {
		return nil, err
	}

	return &CreateApplicationResponse{
		ID:            app.ID,
		AppName:       app.AppName,
		AppIdentifier: app.AppIdentifier,
		AppAccount:    app.AppAccount,
		AppSecret:     app.AppSecret,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// List 获取应用列表
func (s *ApplicationService) List(page, pageSize int) (*ListApplicationsResponse, error) {
	offset := (page - 1) * pageSize
	apps, total, err := s.appRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	var result []CreateApplicationResponse
	for _, app := range apps {
		result = append(result, CreateApplicationResponse{
			ID:            app.ID,
			AppName:       app.AppName,
			AppIdentifier: app.AppIdentifier,
			AppAccount:    app.AppAccount,
			AppSecret:     app.AppSecret,
			Status:        app.Status,
			CreatedAt:     app.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &ListApplicationsResponse{
		Applications: result,
		Total:        total,
		Page:         page,
		PageSize:     pageSize,
	}, nil
}

// GetByID 根据ID获取应用
func (s *ApplicationService) GetByID(id int) (*CreateApplicationResponse, error) {
	app, err := s.appRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &CreateApplicationResponse{
		ID:            app.ID,
		AppName:       app.AppName,
		AppIdentifier: app.AppIdentifier,
		AppAccount:    app.AppAccount,
		AppSecret:     app.AppSecret,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateStatus 更新应用状态
func (s *ApplicationService) UpdateStatus(id int, status int) error {
	if status != 0 && status != 1 {
		return errors.New("状态值必须是0或1")
	}
	return s.appRepo.UpdateStatus(id, status)
}

// Delete 删除应用
func (s *ApplicationService) Delete(id int) error {
	return s.appRepo.Delete(id)
}

// ApplicationOption 应用选项（不含敏感信息）
type ApplicationOption struct {
	ID            int    `json:"id"`
	AppName       string `json:"app_name"`
	AppIdentifier string `json:"app_identifier"`
	AppAccount    string `json:"app_account"`
}

// GetOptions 获取应用选项列表（用于下拉选择等场景，不含敏感信息）
func (s *ApplicationService) GetOptions() ([]ApplicationOption, error) {
	apps, _, err := s.appRepo.List(0, 1000) // 获取所有应用，不分页
	if err != nil {
		return nil, err
	}

	var options []ApplicationOption
	for _, app := range apps {
		options = append(options, ApplicationOption{
			ID:            app.ID,
			AppName:       app.AppName,
			AppIdentifier: app.AppIdentifier,
			AppAccount:    app.AppAccount,
		})
	}

	return options, nil
}
