package routes

import (
    "enterprise-orbit/internal/api/handlers"
    "enterprise-orbit/internal/middleware"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "healthy",
            "service": "enterprise-orbit",
        })
    })
    
    // 认证路由
    authHandler := handlers.NewAuthHandler(db)
    r.POST("/auth/login", authHandler.Login)
    r.POST("/auth/register", authHandler.Register)
    
    // 需要认证的API路由
    api := r.Group("/api")
    api.Use(middleware.JWTAuth())
    {
        // 用户管理
        userHandler := handlers.NewUserHandler(db)
        userRoutes := api.Group("/users")
        {
            userRoutes.GET("", middleware.RBAC("user:read"), userHandler.GetUsers)
            userRoutes.GET("/:id", middleware.RBAC("user:read"), userHandler.GetUser)
            userRoutes.POST("", middleware.RBAC("user:create"), userHandler.CreateUser)
            userRoutes.PUT("/:id", middleware.RBAC("user:update"), userHandler.UpdateUser)
            userRoutes.DELETE("/:id", middleware.RBAC("user:delete"), userHandler.DeleteUser)
        }
        
        // 客户管理
        clientHandler := handlers.NewClientHandler(db)
        clientRoutes := api.Group("/clients")
        {
            clientRoutes.GET("", middleware.RBAC("client:read"), clientHandler.GetClients)
            clientRoutes.GET("/:id", middleware.RBAC("client:read"), clientHandler.GetClient)
            clientRoutes.POST("", middleware.RBAC("client:create"), clientHandler.CreateClient)
            clientRoutes.PUT("/:id", middleware.RBAC("client:update"), clientHandler.UpdateClient)
            clientRoutes.DELETE("/:id", middleware.RBAC("client:delete"), clientHandler.DeleteClient)
        }
    }
}