package main

import (
	"log"
	"Enterprise-Orbit/internal/config"
	"Enterprise-Orbit/pkg/database"
	"Enterprise-Orbit/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// 初始化日志
	if err := logger.InitLogger(cfg); err != nil {
		log.Fatal("Cannot initialize logger:", err)
	}

	// 连接数据库
	db, err := database.NewPostgreSQL(&cfg.Database)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	// 自动迁移
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Cannot migrate database:", err)
	}

	// 启动HTTP服务器
	startServer(cfg, db)
}

func startServer(cfg *config.Config, db *gorm.DB) {
	// 这里将初始化Gin路由和中间件
	log.Printf("Starting %s server on port %d", cfg.Server.Name, cfg.Server.Port)
}