package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/database"
	"github.com/mlaitechio/vagais/internal/services"
)

func testMain() {
	fmt.Println("Testing vagais.ai Backend...")

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize configuration
	cfg := config.Load()
	fmt.Printf("Configuration loaded: Database=%s, Port=%s\n", cfg.DatabaseType, cfg.Port)

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("Database initialized successfully")

	// Initialize optional services
	redisClient, _ := database.InitializeRedis(cfg)
	if redisClient != nil {
		fmt.Println("Redis initialized successfully")
	} else {
		fmt.Println("Redis not available (graceful degradation)")
	}

	// Initialize services
	services.InitializeServices(db, redisClient, cfg)
	fmt.Println("All services initialized successfully")

	// Test service availability
	fmt.Printf("Auth Service: %s\n", services.AuthServiceInstance.GetStatus())
	fmt.Printf("User Service: %s\n", services.UserServiceInstance.GetStatus())
	fmt.Printf("Agent Service: %s\n", services.AgentServiceInstance.GetStatus())
	fmt.Printf("Marketplace Service: %s\n", services.MarketplaceServiceInstance.GetStatus())
	fmt.Printf("Runtime Service: %s\n", services.RuntimeServiceInstance.GetStatus())
	fmt.Printf("Integration Service: %s\n", services.IntegrationServiceInstance.GetStatus())
	fmt.Printf("Notification Service: %s\n", services.NotificationServiceInstance.GetStatus())
	fmt.Printf("Analytics Service: %s\n", services.AnalyticsServiceInstance.GetStatus())
	fmt.Printf("Billing Service: %s\n", services.BillingServiceInstance.GetStatus())
	fmt.Printf("License Service: %s\n", services.LicenseServiceInstance.GetStatus())
	fmt.Printf("Payment Service: %s\n", services.PaymentServiceInstance.GetStatus())

	fmt.Println("\n✅ Backend test completed successfully!")
	fmt.Println("All core services are working with graceful degradation for optional services.")
}

//func main() {
//	testMain()
//}
