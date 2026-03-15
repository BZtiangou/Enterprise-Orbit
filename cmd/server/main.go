package main

import (
	"fmt"
	"log"

	"enterprise-orbit/configs"
	"enterprise-orbit/internal/api/routes"
	"enterprise-orbit/internal/pkg/database"
	"enterprise-orbit/internal/pkg/logger"
	"enterprise-orbit/pkg/auth"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	log.Printf("Database config: host=%s, port=%d, dbname=%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	auth.LoadJWTSecret(cfg.JWT.Secret)

	// 初始化日志
	if err := logger.Init(cfg); err != nil {
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

	// 初始化RBAC数据
	if err := database.InitRBACData(db); err != nil {
		log.Fatal("Cannot initialize RBAC data:", err)
	}

	// 初始化Gin
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// 注入数据库连接
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r, db)

	// 启动服务器
	log.Printf("Starting enterprise-orbit server on :%d", cfg.Server.Port)
	if err := r.Run(":" + fmt.Sprintf("%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
