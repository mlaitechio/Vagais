package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/handlers"
	"github.com/mlaitechio/vagais/internal/middleware"
)

// SetupRoutes sets up all the application routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	userHandler := handlers.NewUserHandler(db, cfg)
	agentHandler := handlers.NewAgentHandler(db, cfg)
	marketplaceHandler := handlers.NewMarketplaceHandler(db, cfg)
	runtimeHandler := handlers.NewRuntimeHandler(db, cfg)
	integrationHandler := handlers.NewIntegrationHandler(db, cfg)
	analyticsHandler := handlers.NewAnalyticsHandler(db, cfg)
	adminHandler := handlers.NewAdminHandler(db, cfg)
	chatHandler := handlers.NewChatHandler(db, cfg)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "healthy"})
		})

		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(), authHandler.Logout)
			auth.POST("/validate", authHandler.ValidateToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// User routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/stats", userHandler.GetUserStats)
			users.GET("/:id", middleware.RoleMiddleware("admin"), userHandler.GetUser)
			users.GET("", middleware.RoleMiddleware("admin"), userHandler.ListUsers)
			users.PUT("/:id/deactivate", middleware.RoleMiddleware("admin"), userHandler.DeactivateUser)
			users.PUT("/:id/activate", middleware.RoleMiddleware("admin"), userHandler.ActivateUser)
			users.PUT("/:id/role", middleware.RoleMiddleware("admin"), userHandler.UpdateUserRole)
		}

		// Organization routes
		orgs := v1.Group("/organizations")
		orgs.Use(middleware.AuthMiddleware())
		{
			orgs.GET("/:id", userHandler.GetOrganization)
			orgs.POST("", userHandler.CreateOrganization)
			orgs.PUT("/:id", userHandler.UpdateOrganization)
			orgs.GET("", userHandler.ListOrganizations)
			orgs.GET("/:id/users", userHandler.GetOrganizationUsers)
			orgs.POST("/:id/users", userHandler.InviteUserToOrganization)
		}

		// Agent routes
		agents := v1.Group("/agents")
		agents.Use(middleware.AuthMiddleware())
		{
			agents.POST("", middleware.RoleMiddleware("admin"), agentHandler.CreateAgent)
			agents.GET("/:id", agentHandler.GetAgent)
			agents.PUT("/:id", middleware.RoleMiddleware("admin"), agentHandler.UpdateAgent)
			agents.DELETE("/:id", middleware.RoleMiddleware("admin"), agentHandler.DeleteAgent)
			agents.GET("", agentHandler.ListAgents)
			agents.POST("/:id/enable", agentHandler.EnableAgent)
			agents.POST("/:id/disable", agentHandler.DisableAgent)
			agents.POST("/:id/execute", agentHandler.ExecuteAgent)
			agents.GET("/categories", agentHandler.GetAgentCategories)
			agents.GET("/:id/stats", agentHandler.GetAgentStats)
		}

		// Public marketplace routes (no auth required)
		public := v1.Group("/public")
		{
			public.GET("/agents", marketplaceHandler.ListAllAgents)
			public.GET("/agents/search", marketplaceHandler.SearchAgentsPublic)
			public.GET("/agents/categories", marketplaceHandler.GetAgentCategories)
			public.GET("/agents/featured", marketplaceHandler.GetFeaturedAgents)
			public.GET("/agents/trending", marketplaceHandler.GetTrendingAgents)
		}

		// Marketplace routes
		marketplace := v1.Group("/marketplace")
		{
			marketplace.GET("/search", marketplaceHandler.SearchAgents)
			marketplace.GET("/featured", marketplaceHandler.GetFeaturedAgents)
			marketplace.GET("/trending", marketplaceHandler.GetTrendingAgents)
			marketplace.GET("/categories", marketplaceHandler.GetAgentCategories)
			marketplace.GET("/stats", marketplaceHandler.GetMarketplaceStats)
			marketplace.GET("/agents", marketplaceHandler.ListMarketplaceAgents)
			marketplace.GET("/agents/:id", marketplaceHandler.GetMarketplaceAgent)
			marketplace.POST("/agents/:id/try", marketplaceHandler.TryMarketplaceAgent)
			marketplace.POST("/agents/:id/purchase", marketplaceHandler.PurchaseMarketplaceAgent)
			marketplace.GET("/agents/:id/reviews", marketplaceHandler.GetAgentReviews)
			marketplace.POST("/agents/:id/reviews", marketplaceHandler.CreateAgentReview)
		}

		// Review routes
		reviews := v1.Group("/reviews")
		reviews.Use(middleware.AuthMiddleware())
		{
			reviews.POST("", marketplaceHandler.CreateReview)
			reviews.GET("/:id", marketplaceHandler.GetReview)
			reviews.PUT("/:id", marketplaceHandler.UpdateReview)
			reviews.DELETE("/:id", marketplaceHandler.DeleteReview)
			reviews.GET("/agent/:agent_id", marketplaceHandler.ListReviews)
			reviews.POST("/:id/helpful", marketplaceHandler.MarkReviewHelpful)
		}

		// Runtime routes
		runtime := v1.Group("/runtime")
		runtime.Use(middleware.AuthMiddleware())
		{
			runtime.POST("/execute", runtimeHandler.ExecuteAgent)
			runtime.GET("/executions/:id", runtimeHandler.GetExecution)
			runtime.GET("/executions", runtimeHandler.ListExecutions)
			runtime.POST("/executions/:id/cancel", runtimeHandler.CancelExecution)
			runtime.GET("/executions/stats", runtimeHandler.GetExecutionStats)
			runtime.GET("/agents/:agent_id/executions/stats", runtimeHandler.GetAgentExecutionStats)
			runtime.GET("/executions/active", runtimeHandler.GetActiveExecutions)
		}

		// Integration routes
		integrations := v1.Group("/integrations")
		integrations.Use(middleware.AuthMiddleware())
		{
			integrations.POST("/webhooks", integrationHandler.CreateWebhook)
			integrations.GET("/webhooks/:id", integrationHandler.GetWebhook)
			integrations.GET("/webhooks", integrationHandler.ListWebhooks)
			integrations.PUT("/webhooks/:id", integrationHandler.UpdateWebhook)
			integrations.DELETE("/webhooks/:id", integrationHandler.DeleteWebhook)
			integrations.GET("/llm-providers", integrationHandler.GetLLMProviders)
			integrations.POST("/llm-providers/test", integrationHandler.TestLLMConnection)
			integrations.GET("/stats", integrationHandler.GetIntegrationStats)
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		analytics.Use(middleware.AuthMiddleware())
		{
			analytics.POST("/track", analyticsHandler.TrackEvent)
			analytics.GET("/usage", analyticsHandler.GetUsageStats)
			analytics.GET("/agents/:agent_id/usage", analyticsHandler.GetAgentUsageStats)
			analytics.GET("/user-behavior", analyticsHandler.GetUserBehaviorAnalytics)
			analytics.GET("/marketplace-trends", analyticsHandler.GetMarketplaceTrends)
			analytics.GET("/revenue", analyticsHandler.GetRevenueAnalytics)
			analytics.GET("/developer-metrics", analyticsHandler.GetDeveloperMetrics)
			analytics.POST("/reports", analyticsHandler.GenerateCustomReport)
		}

		// Billing routes
		//billing := v1.Group("/billing")
		//billing.Use(middleware.AuthMiddleware())
		//{
		//	billing.POST("/subscriptions", billingHandler.CreateSubscription)
		//	billing.GET("/subscriptions/:id", billingHandler.GetSubscription)
		//	billing.GET("/subscriptions", billingHandler.ListSubscriptions)
		//	billing.POST("/subscriptions/:id/cancel", billingHandler.CancelSubscription)
		//	billing.POST("/subscriptions/:id/reactivate", billingHandler.ReactivateSubscription)
		//	billing.GET("/plans", billingHandler.GetAvailablePlans)
		//	billing.GET("/validate", billingHandler.ValidateSubscription)
		//}

		//// Payment routes
		//payments := v1.Group("/payments")
		//payments.Use(middleware.AuthMiddleware())
		//{
		//	payments.POST("/process", paymentHandler.ProcessPayment)
		//	payments.POST("/intent", paymentHandler.CreatePaymentIntent)
		//	payments.GET("/:id", paymentHandler.GetPayment)
		//	payments.GET("", paymentHandler.ListPayments)
		//	payments.POST("/:id/refund", paymentHandler.RefundPayment)
		//	payments.GET("/methods", paymentHandler.GetPaymentMethods)
		//	payments.GET("/stats", paymentHandler.GetPaymentStats)
		//	payments.GET("/:id/validate", paymentHandler.ValidatePayment)
		//}

		//// License routes
		//licenses := v1.Group("/licenses")
		//licenses.Use(middleware.AuthMiddleware())
		//{
		//	licenses.POST("", middleware.RoleMiddleware("admin"), licenseHandler.CreateLicense)
		//	licenses.GET("/:id", licenseHandler.GetLicense)
		//	licenses.GET("/key/:key", licenseHandler.GetLicenseByKey)
		//	licenses.POST("/validate", licenseHandler.ValidateLicense)
		//	licenses.PUT("/:id", middleware.RoleMiddleware("admin"), licenseHandler.UpdateLicense)
		//	licenses.DELETE("/:id", middleware.RoleMiddleware("admin"), licenseHandler.RevokeLicense)
		//	licenses.GET("", middleware.RoleMiddleware("admin"), licenseHandler.ListLicenses)
		//	licenses.GET("/:id/usage", licenseHandler.CheckLicenseUsage)
		//	licenses.POST("/:id/offline", middleware.RoleMiddleware("admin"), licenseHandler.GenerateOfflineLicense)
		//	licenses.POST("/offline/validate", licenseHandler.ValidateOfflineLicense)
		//	licenses.GET("/stats", middleware.RoleMiddleware("admin"), licenseHandler.GetLicenseStats)
		//}

		//// Notification routes
		//notifications := v1.Group("/notifications")
		//notifications.Use(middleware.AuthMiddleware())
		//{
		//	notifications.POST("", notificationHandler.SendNotification)
		//	notifications.GET("/:id", notificationHandler.GetNotification)
		//	notifications.GET("", notificationHandler.ListNotifications)
		//	notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		//	notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		//	notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		//	notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
		//	notifications.POST("/bulk", middleware.RoleMiddleware("admin"), notificationHandler.SendBulkNotification)
		//	notifications.GET("/stats", notificationHandler.GetNotificationStats)
		//}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
		{
			admin.GET("/stats", adminHandler.GetSystemStats)
			admin.GET("/health", adminHandler.GetSystemHealth)
			admin.GET("/logs", adminHandler.GetSystemLogs)
			admin.PUT("/config", adminHandler.UpdateSystemConfig)
			admin.GET("/domains/blocked", adminHandler.GetBlockedDomains)
			admin.POST("/domains/blocked", adminHandler.AddBlockedDomain)
			admin.DELETE("/domains/blocked/:domain", adminHandler.RemoveBlockedDomain)
			admin.GET("/metrics", adminHandler.GetSystemMetrics)
			admin.GET("/audit-logs", adminHandler.GetAuditLogs)
			admin.GET("/backup", adminHandler.GetSystemBackup)
			admin.POST("/backup", adminHandler.CreateSystemBackup)
			admin.GET("/updates", adminHandler.GetSystemUpdates)
			admin.POST("/update", adminHandler.UpdateSystem)
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.GET("/organizations", adminHandler.GetAllOrganizations)
		}
	}

	// WebSocket routes for real-time features
	ws := r.Group("/ws")
	ws.Use(middleware.AuthMiddleware())
	{
		ws.GET("/chat/:agentId", chatHandler.HandleChatWebSocket)
		ws.GET("/executions/:id", func(c *gin.Context) {
			// WebSocket handler for real-time execution updates
			c.JSON(200, gin.H{"message": "WebSocket endpoint"})
		})
	}

	// Admin panel routes (separate from API)
	adminPanel := r.Group("/admin-panel")
	adminPanel.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("admin"))
	{
		adminPanel.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin panel endpoint"})
		})
	}
}
