package routes

import (
	"enterprise-orbit/internal/api/handlers"
	"enterprise-orbit/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "enterprise-orbit",
		})
	})

	authHandler := handlers.NewAuthHandler(db)
	r.POST("/auth/login", authHandler.Login)
	r.POST("/auth/register", authHandler.Register)

	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	{
		userHandler := handlers.NewUserHandler(db)
		userRoutes := api.Group("/users")
		{
			userRoutes.GET("", middleware.RBAC("user:read"), userHandler.GetUsers)
			userRoutes.GET("/:id", middleware.RBAC("user:read"), userHandler.GetUser)
			userRoutes.POST("", middleware.RBAC("user:create"), userHandler.CreateUser)
			userRoutes.PUT("/:id", middleware.RBAC("user:update"), userHandler.UpdateUser)
			userRoutes.DELETE("/:id", middleware.RBAC("user:delete"), userHandler.DeleteUser)
		}

		clientHandler := handlers.NewClientHandler(db)
		clientRoutes := api.Group("/clients")
		{
			clientRoutes.GET("", middleware.RBAC("client:read"), clientHandler.GetClients)
			clientRoutes.GET("/:id", middleware.RBAC("client:read"), clientHandler.GetClient)
			clientRoutes.POST("", middleware.RBAC("client:create"), clientHandler.CreateClient)
			clientRoutes.PUT("/:id", middleware.RBAC("client:update"), clientHandler.UpdateClient)
			clientRoutes.DELETE("/:id", middleware.RBAC("client:delete"), clientHandler.DeleteClient)
		}

		customerHandler := handlers.NewCustomerHandler(db)
		customerRoutes := api.Group("/customers")
		{
			customerRoutes.GET("", middleware.RBAC("customer:read"), customerHandler.GetCustomers)
			customerRoutes.GET("/:id", middleware.RBAC("customer:read"), customerHandler.GetCustomer)
			customerRoutes.POST("", middleware.RBAC("customer:create"), customerHandler.CreateCustomer)
			customerRoutes.PUT("/:id", middleware.RBAC("customer:update"), customerHandler.UpdateCustomer)
			customerRoutes.DELETE("/:id", middleware.RBAC("customer:delete"), customerHandler.DeleteCustomer)
			customerRoutes.GET("/:id/health", middleware.RBAC("customer:read"), customerHandler.GetCustomerHealth)
			customerRoutes.POST("/:id/calculate-health", middleware.RBAC("customer:update"), customerHandler.CalculateHealthScore)
		}

		contractHandler := handlers.NewContractHandler(db)
		contractRoutes := api.Group("/contracts")
		{
			contractRoutes.GET("", middleware.RBAC("contract:read"), contractHandler.GetContracts)
			contractRoutes.GET("/:id", middleware.RBAC("contract:read"), contractHandler.GetContract)
			contractRoutes.POST("", middleware.RBAC("contract:create"), contractHandler.CreateContract)
			contractRoutes.PUT("/:id", middleware.RBAC("contract:update"), contractHandler.UpdateContract)
			contractRoutes.DELETE("/:id", middleware.RBAC("contract:delete"), contractHandler.DeleteContract)
			contractRoutes.POST("/:id/submit", middleware.RBAC("contract:update"), contractHandler.SubmitForApproval)
			contractRoutes.POST("/:id/approve", middleware.RBAC("contract:approve"), contractHandler.ApproveContract)
			contractRoutes.GET("/expiring", middleware.RBAC("contract:read"), contractHandler.GetExpiringContracts)
			contractRoutes.GET("/:id/performance", middleware.RBAC("contract:read"), contractHandler.GetContractPerformance)
			contractRoutes.POST("/:id/performance", middleware.RBAC("contract:update"), contractHandler.CreatePerformanceRecord)
		}

		interactionHandler := handlers.NewInteractionHandler(db)
		interactionRoutes := api.Group("/interactions")
		{
			interactionRoutes.GET("", middleware.RBAC("interaction:read"), interactionHandler.GetInteractions)
			interactionRoutes.GET("/:id", middleware.RBAC("interaction:read"), interactionHandler.GetInteraction)
			interactionRoutes.POST("", middleware.RBAC("interaction:create"), interactionHandler.CreateInteraction)
			interactionRoutes.PUT("/:id", middleware.RBAC("interaction:update"), interactionHandler.UpdateInteraction)
			interactionRoutes.DELETE("/:id", middleware.RBAC("interaction:delete"), interactionHandler.DeleteInteraction)
			interactionRoutes.GET("/timeline", middleware.RBAC("interaction:read"), interactionHandler.GetInteractionTimeline)
			interactionRoutes.GET("/stats", middleware.RBAC("interaction:read"), interactionHandler.GetInteractionStats)
		}

		dashboardHandler := handlers.NewDashboardHandler(db)
		dashboardRoutes := api.Group("/dashboard")
		{
			dashboardRoutes.GET("/overview", middleware.RBAC("dashboard:read"), dashboardHandler.GetOverview)
			dashboardRoutes.GET("/customer-growth", middleware.RBAC("dashboard:read"), dashboardHandler.GetCustomerGrowthTrend)
			dashboardRoutes.GET("/contract-distribution", middleware.RBAC("dashboard:read"), dashboardHandler.GetContractValueDistribution)
			dashboardRoutes.GET("/interaction-heatmap", middleware.RBAC("dashboard:read"), dashboardHandler.GetInteractionHeatmap)
			dashboardRoutes.GET("/renewal-alerts", middleware.RBAC("dashboard:read"), dashboardHandler.GetRenewalRiskAlerts)
			dashboardRoutes.GET("/top-customers", middleware.RBAC("dashboard:read"), dashboardHandler.GetTopCustomers)
			dashboardRoutes.GET("/recent-activities", middleware.RBAC("dashboard:read"), dashboardHandler.GetRecentActivities)
		}
	}
}
