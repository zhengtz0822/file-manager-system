package handler

import (
	"net/http"
	"strings"

	"file-manager-service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Logout 用户注销
// @Summary 注销登录
// @Description 撤销当前 Token
// @Tags 认证
// @Security Bearer
// @Success 200 {object} Response
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	// 获取 Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "未提供认证信息",
		})
		return
	}

	// 解析 Bearer Token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "认证格式错误",
		})
		return
	}

	tokenString := parts[1]

	// 撤销 token
	if err := jwt.RevokeToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "注销失败",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "注销成功",
	})
}
