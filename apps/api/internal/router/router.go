package router

import (
	"file-manager-service/internal/handler"
	"file-manager-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS 中间件
	r.Use(CORSMiddleware())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// 认证路由（不需要认证）
		authHandler := handler.NewAuthHandler()
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", middleware.AuthMiddleware(), handler.Logout)
		}

		// 文档路由（需要认证）
		docHandler := handler.NewDocumentHandler()
		authMiddleware := middleware.AuthMiddleware()

		documents := v1.Group("/documents")
		documents.Use(authMiddleware)
		{
			// 分片上传
			documents.POST("/chunks/init", docHandler.InitUpload)
			documents.POST("/chunks/upload", docHandler.UploadChunk)
			documents.POST("/chunks/complete", docHandler.CompleteUpload)
			documents.DELETE("/chunks/:upload_id", docHandler.CancelUpload)

			// 文档管理
			documents.GET("", docHandler.ListDocuments)
			documents.GET("/:id", docHandler.GetDocument)
			documents.DELETE("/:id", docHandler.DeleteDocument)

			// 文档访问
			documents.GET("/:id/download", docHandler.DownloadDocument)
			documents.GET("/:id/preview", docHandler.PreviewDocument)
		}

		// 用户信息
		v1.GET("/user/me", authMiddleware, docHandler.GetCurrentUser)
	}

	return r
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
