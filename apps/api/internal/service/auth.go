package service

import (
	"errors"
	"fmt"

	"file-manager-service/internal/model"
	"file-manager-service/internal/pkg/jwt"
	"file-manager-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	appRepo  *repository.ApplicationRepository
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
		appRepo:  repository.NewApplicationRepository(),
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"` // 改为 interface{} 以支持自定义格式
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成用户 Token
	token, err := jwt.GenerateUserToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	// 构造前端需要的用户数据格式
	userData := map[string]interface{}{
		"id":         user.ID,
		"userid":     fmt.Sprintf("%d", user.ID), // 前端需要的 userid
		"name":       user.Username,               // 前端需要的 name 字段
		"username":   user.Username,
		"created_at": user.CreatedAt,
	}

	return &LoginResponse{
		Token: token,
		User:  userData,
	}, nil
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) error {
	// 检查用户是否已存在
	_, err := s.userRepo.FindByUsername(req.Username)
	if err == nil {
		return errors.New("用户已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(user)
}

// AppTokenRequest 应用获取令牌请求
type AppTokenRequest struct {
	AppAccount string `json:"app_account" binding:"required"`
	AppSecret  string `json:"app_secret" binding:"required"`
}

// AppTokenResponse 应用令牌响应
type AppTokenResponse struct {
	Token string      `json:"token"`
	App   interface{} `json:"app"`
}

// GetAppToken 应用获取访问令牌
func (s *AuthService) GetAppToken(req *AppTokenRequest) (*AppTokenResponse, error) {
	// 查找应用
	app, err := s.appRepo.GetByAccount(req.AppAccount)
	if err != nil {
		return nil, errors.New("应用不存在")
	}

	// 验证密钥
	if app.AppSecret != req.AppSecret {
		return nil, errors.New("应用密钥错误")
	}

	// 检查应用状态
	if app.Status != 1 {
		return nil, errors.New("应用已禁用")
	}

	// 生成应用 Token
	token, err := jwt.GenerateAppToken(app.ID, app.AppName, app.AppIdentifier)
	if err != nil {
		return nil, err
	}

	// 构造响应数据
	appData := map[string]interface{}{
		"id":             app.ID,
		"app_name":       app.AppName,
		"app_identifier": app.AppIdentifier,
		"app_account":    app.AppAccount,
		"created_at":     app.CreatedAt,
	}

	return &AppTokenResponse{
		Token: token,
		App:   appData,
	}, nil
}
