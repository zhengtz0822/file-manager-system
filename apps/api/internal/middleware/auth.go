package middleware

import (
	"net/http"
	"strings"

	"file-manager-service/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证信息",
			})
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "认证格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 检查 token 是否在黑名单
		if jwt.IsRevoked(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token 已被撤销",
			})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token 无效或已过期",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("claims", claims)

		// 根据 Token 类型设置不同的上下文信息
		if claims.Type == jwt.TokenTypeApp {
			// 应用 Token
			c.Set("token_type", "app")
			c.Set("app_id", claims.AppID)
			c.Set("app_name", claims.AppName)
			c.Set("app_identifier", claims.AppIdentifier)
		} else {
			// 用户 Token
			c.Set("token_type", "user")
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
		}

		c.Next()
	}
}
