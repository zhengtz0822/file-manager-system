package main

import (
	"fmt"
	"log"

	_ "file-manager-service/docs" // Swagger 文档

	"file-manager-service/internal/config"
	"file-manager-service/internal/model"
	"file-manager-service/internal/pkg/jwt"
	"file-manager-service/internal/pkg/storage"
	"file-manager-service/internal/router"
)

// @title 文件管理系统 API
// @version 1.0
// @description 文件管理系统后端 API，支持大文件分片上传、下载和预览
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@filemanager.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请输入 Bearer {token} 格式的 Token，例如：Bearer eyJhbGc...

func main() {
	// 加载配置
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化存储
	storage.Init(&cfg.Storage)

	// 初始化 JWT
	jwt.Init(&cfg.JWT)

	// 初始化数据库
	if err := model.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 自动迁移数据库表
	if err := model.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
