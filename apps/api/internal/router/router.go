package router

import (
	"file-manager-service/internal/config"
	"file-manager-service/internal/handler"
	"file-manager-service/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			auth.POST("/app-token", authHandler.GetAppToken) // 应用获取访问令牌
			auth.POST("/logout", middleware.AuthMiddleware(), handler.Logout)
		}

		// 文档路由（需要认证）
		docHandler := handler.NewDocumentHandler()
		authMiddleware := middleware.AuthMiddleware()

		documents := v1.Group("/documents")
		documents.Use(authMiddleware)
		{
			// 单文件上传（适用于小文件）
			documents.POST("/upload", docHandler.UploadSingleFile)

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

		// 应用管理路由（需要认证）
		appHandler := handler.NewApplicationHandler()
		applications := v1.Group("/applications")
		applications.Use(authMiddleware)
		{
			applications.POST("", appHandler.CreateApplication)
			applications.GET("", appHandler.ListApplications)
			applications.GET("/options", appHandler.GetApplicationOptions) // 获取应用选项（不含敏感信息）
			applications.GET("/:id", appHandler.GetApplication)
			applications.PUT("/:id/status", appHandler.UpdateApplicationStatus)
			applications.DELETE("/:id", appHandler.DeleteApplication)
		}
	}

	// Swagger 文档（根据配置启用）
	if config.GlobalConfig.Swagger.Enabled {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
