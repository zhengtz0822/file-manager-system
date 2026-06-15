package handler

import (
	"net/http"

	"file-manager-service/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录系统，获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息" SchemaExample({"username":"admin","password":"admin123"})
// @Success 200 {object} Response{data=service.LoginResponse} "登录成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "用户名或密码错误"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "登录失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "登录成功",
		Data:    resp,
	})
}

// Register 用户注册
// @Summary 用户注册
// @Description 注册新用户账号
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息" SchemaExample({"username":"newuser","password":"password123"})
// @Success 200 {object} Response "注册成功"
// @Failure 400 {object} Response "参数错误或用户已存在"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "注册失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "注册成功",
	})
}

// GetAppToken 应用获取访问令牌
// @Summary 应用获取访问令牌
// @Description 外部应用使用应用账号和密钥获取访问令牌，用于调用需要认证的 API
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body service.AppTokenRequest true "应用认证信息" SchemaExample({"app_account":"APP_abc123","app_secret":"64位十六进制密钥"})
// @Success 200 {object} Response{data=service.AppTokenResponse} "获取令牌成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "应用认证失败（应用不存在、密钥错误或应用已禁用）"
// @Router /auth/app-token [post]
func (h *AuthHandler) GetAppToken(c *gin.Context) {
	var req service.AppTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.authService.GetAppToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "应用认证失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取令牌成功",
		Data:    resp,
	})
}

// Response 通用响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
