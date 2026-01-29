package main

import (
	"fmt"
	"log"
	"file-manager-service/internal/config"
	"file-manager-service/internal/model"
	"file-manager-service/internal/pkg/jwt"
	"file-manager-service/internal/pkg/storage"
	"file-manager-service/internal/router"
)

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
