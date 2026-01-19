package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/database"
	"github.com/mlaitechio/vagais/internal/models"
	"github.com/mlaitechio/vagais/internal/services"
	"gorm.io/gorm"
)

func main() {
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

	// Initialize services
	redisClient, _ := database.InitializeRedis(cfg)
	services.InitializeServices(db, redisClient, cfg)

	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [migrate|seed|reset]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "migrate":
		if err := runMigrations(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ Migrations completed successfully")
	case "seed":
		if err := runSeeds(db); err != nil {
			log.Fatalf("Seeding failed: %v", err)
		}
		fmt.Println("✅ Seeding completed successfully")
	case "reset":
		if err := resetDatabase(db); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
		fmt.Println("✅ Database reset completed successfully")
	default:
		fmt.Println("Unknown command. Use: migrate, seed, or reset")
		os.Exit(1)
	}
}

// runMigrations performs database migrations
func runMigrations(db *gorm.DB) error {
	fmt.Println("Running database migrations...")
	defer func() {
		fmt.Println("Finished migration step.")
	}()
	if err := db.AutoMigrate(
		&models.Organization{},
		&models.User{},
		&models.Agent{},
		&models.Review{},
		&models.Execution{},
		&models.License{},
		&models.Payment{},
		&models.Subscription{},
		&models.Analytics{},
		&models.Webhook{},
		&models.Notification{},
		&models.LLMProvider{},
		&models.BillingPlan{},
	); err != nil {
		fmt.Printf("[ERROR] Migration failed: %v\n", err)
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	fmt.Println("✅ All tables migrated successfully")
	return nil
}

// runSeeds seeds the database with initial data
func runSeeds(db *gorm.DB) error {
	fmt.Println("Seeding database with initial data...")
	if err := seedOrganizations(db); err != nil {
		fmt.Printf("[ERROR] Organization seed failed: %v\n", err)
		return err
	}
	if err := seedUsers(db); err != nil {
		fmt.Printf("[ERROR] User seed failed: %v\n", err)
		return err
	}
	if err := seedAgents(db); err != nil {
		fmt.Printf("[ERROR] Agent seed failed: %v\n", err)
		return err
	}
	if err := seedLicenses(db); err != nil {
		fmt.Printf("[ERROR] License seed failed: %v\n", err)
		return err
	}
	if err := seedLLMProviders(db); err != nil {
		fmt.Printf("[ERROR] LLM Provider seed failed: %v\n", err)
		return err
	}
	if err := seedBillingPlans(db); err != nil {
		fmt.Printf("[ERROR] Billing Plan seed failed: %v\n", err)
		return err
	}
	if err := seedAnalytics(db); err != nil {
		fmt.Printf("[ERROR] Analytics seed failed: %v\n", err)
		return err
	}
	fmt.Println("✅ Database seeded successfully")
	return nil
}

// resetDatabase resets the database
func resetDatabase(db *gorm.DB) error {
	fmt.Println("Resetting database...")
	fmt.Println("Dropping all tables...")
	if err := db.Migrator().DropTable(
		&models.Notification{},
		&models.Webhook{},
		&models.Analytics{},
		&models.Subscription{},
		&models.Payment{},
		&models.License{},
		&models.Execution{},
		&models.Review{},
		&models.Agent{},
		&models.User{},
		&models.Organization{},
		&models.LLMProvider{},
		&models.BillingPlan{},
	); err != nil {
		fmt.Printf("[ERROR] Drop tables failed: %v\n", err)
		return fmt.Errorf("failed to drop tables: %v", err)
	}
	fmt.Println("Tables dropped. Running migrations...")
	if err := runMigrations(db); err != nil {
		fmt.Printf("[ERROR] Migration after drop failed: %v\n", err)
		return err
	}
	fmt.Println("Migrations done. Running seeds...")
	if err := runSeeds(db); err != nil {
		fmt.Printf("[ERROR] Seeding after migration failed: %v\n", err)
		return err
	}
	fmt.Println("Database reset and seeded.")
	return nil
}

// seedOrganizations seeds organizations
func seedOrganizations(db *gorm.DB) error {
	fmt.Println("Seeding organizations...")

	organizations := []models.Organization{
		{
			Name:        "AGAI Studio",
			Slug:        "agai-studio",
			Description: "Official AGAI development studio",
			Website:     "https://agai.studio",
			Logo:        "https://agai.studio/logo.png",
			IsActive:    true,
			Plan:        "enterprise",
		},
		{
			Name:        "AI Research Lab",
			Slug:        "ai-research-lab",
			Description: "Research organization focused on AI development",
			Website:     "https://airesearchlab.com",
			Logo:        "https://airesearchlab.com/logo.png",
			IsActive:    true,
			Plan:        "pro",
		},
		{
			Name:        "Startup Inc",
			Slug:        "startup-inc",
			Description: "Innovative startup building AI solutions",
			Website:     "https://startupinc.ai",
			Logo:        "https://startupinc.ai/logo.png",
			IsActive:    true,
			Plan:        "basic",
		},
	}

	for _, org := range organizations {
		if err := db.Create(&org).Error; err != nil {
			return fmt.Errorf("failed to seed organization %s: %v", org.Name, err)
		}
	}

	fmt.Printf("✅ Created %d organizations\n", len(organizations))
	return nil
}

// seedUsers seeds users
func seedUsers(db *gorm.DB) error {
	fmt.Println("Seeding users...")

	// Get organization IDs
	var orgs []models.Organization
	if err := db.Find(&orgs).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %v", err)
	}

	if len(orgs) == 0 {
		return fmt.Errorf("no organizations found for seeding users")
	}

	users := []models.User{
		{
			Email:          "admin@agai.studio",
			Username:       "admin",
			FirstName:      "Admin",
			LastName:       "User",
			PasswordHash:   "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:           "admin",
			IsActive:       true,
			EmailVerified:  true,
			Avatar:         "https://agai.studio/avatars/admin.png",
			OrganizationID: &orgs[0].ID, // AGAI Studio
			Credits:        10000,
			Preferences:    models.MapToJSON(nil),
		},
		{
			Email:          "developer@agai.studio",
			Username:       "developer",
			FirstName:      "John",
			LastName:       "Developer",
			PasswordHash:   "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:           "developer",
			IsActive:       true,
			EmailVerified:  true,
			Avatar:         "https://agai.studio/avatars/developer.png",
			OrganizationID: &orgs[0].ID, // AGAI Studio
			Credits:        5000,
			Preferences:    models.MapToJSON(nil),
		},
		{
			Email:          "researcher@airesearchlab.com",
			Username:       "researcher",
			FirstName:      "Dr. Sarah",
			LastName:       "Researcher",
			PasswordHash:   "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:           "researcher",
			IsActive:       true,
			EmailVerified:  true,
			Avatar:         "https://airesearchlab.com/avatars/researcher.png",
			OrganizationID: &orgs[1].ID, // AI Research Lab
			Credits:        3000,
			Preferences:    models.MapToJSON(nil),
		},
		{
			Email:          "founder@startupinc.ai",
			Username:       "founder",
			FirstName:      "Mike",
			LastName:       "Founder",
			PasswordHash:   "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Role:           "founder",
			IsActive:       true,
			EmailVerified:  true,
			Avatar:         "https://startupinc.ai/avatars/founder.png",
			OrganizationID: &orgs[2].ID, // Startup Inc
			Credits:        1000,
			Preferences:    models.MapToJSON(nil),
		},
	}

	for _, user := range users {
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to seed user %s: %v", user.Username, err)
		}
	}

	fmt.Printf("✅ Created %d users\n", len(users))
	return nil
}

// seedAgents seeds agents
func seedAgents(db *gorm.DB) error {
	fmt.Println("Seeding agents...")

	// Get user IDs
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found for seeding agents")
	}

	tags1, _ := json.Marshal([]string{"support", "customer-service", "chatbot"})
	screens1, _ := json.Marshal([]string{"https://agai.studio/agents/customer-support-bot/screenshot1.png"})
	tags2, _ := json.Marshal([]string{"data-analysis", "visualization", "statistics"})
	screens2, _ := json.Marshal([]string{"https://airesearchlab.com/agents/data-analysis-assistant/screenshot1.png"})
	tags3, _ := json.Marshal([]string{"content-writing", "blogging", "marketing"})
	screens3, _ := json.Marshal([]string{"https://startupinc.ai/agents/content-writer-pro/screenshot1.png"})
	tags4, _ := json.Marshal([]string{"code-review", "programming", "development"})
	screens4, _ := json.Marshal([]string{"https://agai.studio/agents/code-review-assistant/screenshot1.png"})
	tags5, _ := json.Marshal([]string{"email", "communication", "automation"})
	screens5, _ := json.Marshal([]string{"https://agai.studio/agents/email-assistant/screenshot1.png"})
	tags6, _ := json.Marshal([]string{"research", "academic", "writing"})
	screens6, _ := json.Marshal([]string{"https://airesearchlab.com/agents/research-assistant/screenshot1.png"})
	tags7, _ := json.Marshal([]string{"sales", "lead-generation", "crm"})
	screens7, _ := json.Marshal([]string{"https://startupinc.ai/agents/sales-assistant/screenshot1.png"})
	tags8, _ := json.Marshal([]string{"translation", "multilingual", "communication"})
	screens8, _ := json.Marshal([]string{"https://agai.studio/agents/translation-assistant/screenshot1.png"})

	agents := []models.Agent{
		{
			Name:              "Customer Support Bot",
			Description:       "AI-powered customer support agent that can handle common inquiries and provide helpful responses",
			Slug:              "customer-support-bot",
			Version:           "1.0.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Customer Service",
			Tags:              tags1,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 1000, "temperature": 0.7}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID, // developer
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             0.0,
			Currency:          "USD",
			PricingModel:      "free",
			LicenseRequired:   false,
			Rating:            4.5,
			ReviewCount:       12,
			UsageCount:        1500,
			Downloads:         45,
			Icon:              "https://agai.studio/agents/customer-support-bot/icon.png",
			Screenshots:       screens1,
			Documentation:     "Comprehensive documentation for the Customer Support Bot",
			Repository:        "https://github.com/agai-studio/customer-support-bot",
		},
		{
			Name:              "Data Analysis Assistant",
			Description:       "Advanced data analysis agent that can process, analyze, and visualize complex datasets",
			Slug:              "data-analysis-assistant",
			Version:           "2.1.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Data Science",
			Tags:              tags2,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 2000, "temperature": 0.3}),
			LLMProvider:       "anthropic",
			LLMModel:          "claude-3-sonnet",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[2].ID, // researcher
			OrganizationID:    users[2].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             29.99,
			Currency:          "USD",
			PricingModel:      "one-time",
			LicenseRequired:   true,
			Rating:            4.8,
			ReviewCount:       8,
			UsageCount:        750,
			Downloads:         23,
			Icon:              "https://airesearchlab.com/agents/data-analysis-assistant/icon.png",
			Screenshots:       screens2,
			Documentation:     "Advanced data analysis capabilities with statistical modeling",
			Repository:        "https://github.com/ai-research-lab/data-analysis-assistant",
		},
		{
			Name:              "Content Writer Pro",
			Description:       "Professional content writing assistant that creates high-quality articles, blogs, and marketing copy",
			Slug:              "content-writer-pro",
			Version:           "1.5.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Content Creation",
			Tags:              tags3,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 1500, "temperature": 0.8}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[3].ID, // founder
			OrganizationID:    users[3].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             19.99,
			Currency:          "USD",
			PricingModel:      "subscription",
			LicenseRequired:   false,
			Rating:            4.2,
			ReviewCount:       15,
			UsageCount:        2200,
			Downloads:         67,
			Icon:              "https://startupinc.ai/agents/content-writer-pro/icon.png",
			Screenshots:       screens3,
			Documentation:     "Professional content writing with SEO optimization",
			Repository:        "https://github.com/startup-inc/content-writer-pro",
		},
		{
			Name:              "Code Review Assistant",
			Description:       "AI-powered code review agent that analyzes code quality, security, and best practices",
			Slug:              "code-review-assistant",
			Version:           "1.2.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Development",
			Tags:              tags4,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 3000, "temperature": 0.2}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID, // developer
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             39.99,
			Currency:          "USD",
			PricingModel:      "one-time",
			LicenseRequired:   true,
			Rating:            4.7,
			ReviewCount:       6,
			UsageCount:        420,
			Downloads:         18,
			Icon:              "https://agai.studio/agents/code-review-assistant/icon.png",
			Screenshots:       screens4,
			Documentation:     "Advanced code analysis with security scanning",
			Repository:        "https://github.com/agai-studio/code-review-assistant",
		},
		{
			Name:              "Email Assistant",
			Description:       "Smart email management agent that drafts, categorizes, and responds to emails automatically",
			Slug:              "email-assistant",
			Version:           "1.1.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Communication",
			Tags:              tags5,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 800, "temperature": 0.6}),
			LLMProvider:       "anthropic",
			LLMModel:          "claude-3-haiku",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID, // developer
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             0.0,
			Currency:          "USD",
			PricingModel:      "free",
			LicenseRequired:   false,
			Rating:            4.3,
			ReviewCount:       9,
			UsageCount:        1800,
			Downloads:         52,
			Icon:              "https://agai.studio/agents/email-assistant/icon.png",
			Screenshots:       screens5,
			Documentation:     "Email automation and smart categorization",
			Repository:        "https://github.com/agai-studio/email-assistant",
		},
		{
			Name:              "Research Assistant",
			Description:       "Academic research agent that helps with literature review, citation management, and paper writing",
			Slug:              "research-assistant",
			Version:           "2.0.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Research",
			Tags:              tags6,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 2500, "temperature": 0.4}),
			LLMProvider:       "anthropic",
			LLMModel:          "claude-3-sonnet",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[2].ID, // researcher
			OrganizationID:    users[2].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             49.99,
			Currency:          "USD",
			PricingModel:      "subscription",
			LicenseRequired:   true,
			Rating:            4.9,
			ReviewCount:       4,
			UsageCount:        320,
			Downloads:         12,
			Icon:              "https://airesearchlab.com/agents/research-assistant/icon.png",
			Screenshots:       screens6,
			Documentation:     "Academic research with citation management",
			Repository:        "https://github.com/ai-research-lab/research-assistant",
		},
		{
			Name:              "Sales Assistant",
			Description:       "Sales automation agent that generates leads, qualifies prospects, and manages customer relationships",
			Slug:              "sales-assistant",
			Version:           "1.3.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Sales",
			Tags:              tags7,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 1200, "temperature": 0.7}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[3].ID, // founder
			OrganizationID:    users[3].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             24.99,
			Currency:          "USD",
			PricingModel:      "subscription",
			LicenseRequired:   false,
			Rating:            4.1,
			ReviewCount:       11,
			UsageCount:        950,
			Downloads:         38,
			Icon:              "https://startupinc.ai/agents/sales-assistant/icon.png",
			Screenshots:       screens7,
			Documentation:     "Sales automation and lead generation",
			Repository:        "https://github.com/startup-inc/sales-assistant",
		},
		{
			Name:              "Translation Assistant",
			Description:       "Multilingual translation agent supporting 50+ languages with context-aware translations",
			Slug:              "translation-assistant",
			Version:           "1.0.0",
			Status:            "published",
			Type:              "custom",
			Category:          "Translation",
			Tags:              tags8,
			Config:            models.MapToJSON(map[string]interface{}{"max_tokens": 1000, "temperature": 0.5}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID, // developer
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             0.0,
			Currency:          "USD",
			PricingModel:      "free",
			LicenseRequired:   false,
			Rating:            4.6,
			ReviewCount:       7,
			UsageCount:        2800,
			Downloads:         89,
			Icon:              "https://agai.studio/agents/translation-assistant/icon.png",
			Screenshots:       screens8,
			Documentation:     "Multilingual translation with cultural context",
			Repository:        "https://github.com/agai-studio/translation-assistant",
		},
	}

	for _, agent := range agents {
		if err := db.Create(&agent).Error; err != nil {
			return fmt.Errorf("failed to seed agent %s: %v", agent.Name, err)
		}
	}

	fmt.Printf("✅ Created %d agents\n", len(agents))
	return nil
}

// seedLicenses seeds licenses
func seedLicenses(db *gorm.DB) error {
	fmt.Println("Seeding licenses...")

	// Get organization IDs
	var orgs []models.Organization
	if err := db.Find(&orgs).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %v", err)
	}

	if len(orgs) == 0 {
		return fmt.Errorf("no organizations found for seeding licenses")
	}

	features1, _ := json.Marshal([]string{"unlimited_agents", "advanced_analytics", "priority_support", "custom_integrations"})
	features2, _ := json.Marshal([]string{"advanced_analytics", "priority_support"})
	features3, _ := json.Marshal([]string{"basic_analytics"})

	licenses := []models.License{
		{
			Key:            "AGAI-ENT-2024-001",
			Type:           "enterprise",
			Status:         "active",
			OrganizationID: &orgs[0].ID, // AGAI Studio
			IssuedAt:       time.Now(),
			Features:       features1,
			MaxUsers:       100,
			MaxAgents:      50,
			IsValid:        true,
		},
		{
			Key:            "AGAI-PRO-2024-002",
			Type:           "pro",
			Status:         "active",
			OrganizationID: &orgs[1].ID, // AI Research Lab
			IssuedAt:       time.Now(),
			Features:       features2,
			MaxUsers:       25,
			MaxAgents:      20,
			IsValid:        true,
		},
		{
			Key:            "AGAI-BASIC-2024-003",
			Type:           "basic",
			Status:         "active",
			OrganizationID: &orgs[2].ID, // Startup Inc
			IssuedAt:       time.Now(),
			Features:       features3,
			MaxUsers:       10,
			MaxAgents:      5,
			IsValid:        true,
		},
	}

	for _, license := range licenses {
		if err := db.Create(&license).Error; err != nil {
			return fmt.Errorf("failed to seed license %s: %v", license.Key, err)
		}
	}

	fmt.Printf("✅ Created %d licenses\n", len(licenses))
	return nil
}

// seedLLMProviders seeds LLM providers
func seedLLMProviders(db *gorm.DB) error {
	fmt.Println("Seeding LLM providers...")

	config1, _ := json.Marshal(map[string]interface{}{
		"api_key":  "",
		"base_url": "https://api.openai.com/v1",
	})
	config2, _ := json.Marshal(map[string]interface{}{
		"api_key":  "",
		"base_url": "https://api.anthropic.com",
	})
	config3, _ := json.Marshal(map[string]interface{}{
		"model_path": "",
		"endpoint":   "http://localhost:8000",
	})

	providers := []models.LLMProvider{
		{
			Name:      "OpenAI",
			Type:      "openai",
			IsActive:  true,
			RateLimit: 1000,
			MaxTokens: 4096,
			Config:    config1,
		},
		{
			Name:      "Anthropic",
			Type:      "anthropic",
			IsActive:  true,
			RateLimit: 500,
			MaxTokens: 100000,
			Config:    config2,
		},
		{
			Name:      "Local",
			Type:      "local",
			IsActive:  false,
			RateLimit: 100,
			MaxTokens: 2048,
			Config:    config3,
		},
	}

	for _, provider := range providers {
		if err := db.Create(&provider).Error; err != nil {
			return fmt.Errorf("failed to seed LLM provider %s: %v", provider.Name, err)
		}
	}

	fmt.Printf("✅ Created %d LLM providers\n", len(providers))
	return nil
}

// seedBillingPlans seeds billing plans
func seedBillingPlans(db *gorm.DB) error {
	fmt.Println("Seeding billing plans...")

	features1, _ := json.Marshal([]string{"5 agents", "100 executions/month", "Basic support"})
	features2, _ := json.Marshal([]string{"Unlimited agents", "1000 executions/month", "Priority support", "Analytics"})
	features3, _ := json.Marshal([]string{"Unlimited agents", "Unlimited executions", "24/7 support", "Custom integrations", "SLA"})

	plans := []models.BillingPlan{
		{
			Name:          "Free",
			Slug:          "free",
			Price:         0.0,
			Currency:      "USD",
			Interval:      "monthly",
			Features:      features1,
			MaxAgents:     5,
			MaxExecutions: 100,
			IsActive:      true,
			Description:   "Basic plan for getting started",
			SortOrder:     1,
		},
		{
			Name:          "Pro",
			Slug:          "pro",
			Price:         29.99,
			Currency:      "USD",
			Interval:      "monthly",
			Features:      features2,
			MaxAgents:     -1, // Unlimited
			MaxExecutions: 1000,
			IsActive:      true,
			Description:   "Professional plan for growing teams",
			SortOrder:     2,
		},
		{
			Name:          "Enterprise",
			Slug:          "enterprise",
			Price:         99.99,
			Currency:      "USD",
			Interval:      "monthly",
			Features:      features3,
			MaxAgents:     -1, // Unlimited
			MaxExecutions: -1, // Unlimited
			IsActive:      true,
			Description:   "Enterprise plan for large organizations",
			SortOrder:     3,
		},
	}

	for _, plan := range plans {
		if err := db.Create(&plan).Error; err != nil {
			return fmt.Errorf("failed to seed billing plan %s: %v", plan.Name, err)
		}
	}

	fmt.Printf("✅ Created %d billing plans\n", len(plans))
	return nil
}

// seedAnalytics seeds analytics data
func seedAnalytics(db *gorm.DB) error {
	fmt.Println("Seeding analytics data...")

	// Get organization IDs
	var orgs []models.Organization
	if err := db.Find(&orgs).Error; err != nil {
		return fmt.Errorf("failed to get organizations: %v", err)
	}

	if len(orgs) == 0 {
		return fmt.Errorf("no organizations found for seeding analytics")
	}

	analytics := []models.Analytics{
		{
			OrganizationID: &orgs[0].ID, // AGAI Studio
			Type:           "execution",
			Metric:         "total_executions",
			Value:          1250.0,
			Date:           time.Now().AddDate(0, 0, -1),
			Metadata:       models.MapToJSON(map[string]interface{}{"success_rate": 0.95, "avg_duration": 2.3}),
		},
		{
			OrganizationID: &orgs[0].ID,
			Type:           "revenue",
			Metric:         "monthly_revenue",
			Value:          4500.0,
			Date:           time.Now().AddDate(0, -1, 0),
			Metadata:       models.MapToJSON(map[string]interface{}{"subscriptions": 15, "one_time_sales": 3}),
		},
		{
			OrganizationID: &orgs[1].ID, // AI Research Lab
			Type:           "execution",
			Metric:         "total_executions",
			Value:          850.0,
			Date:           time.Now().AddDate(0, 0, -1),
			Metadata:       models.MapToJSON(map[string]interface{}{"success_rate": 0.92, "avg_duration": 4.1}),
		},
		{
			OrganizationID: &orgs[1].ID,
			Type:           "revenue",
			Metric:         "monthly_revenue",
			Value:          2800.0,
			Date:           time.Now().AddDate(0, -1, 0),
			Metadata:       models.MapToJSON(map[string]interface{}{"subscriptions": 8, "one_time_sales": 2}),
		},
		{
			OrganizationID: &orgs[2].ID, // Startup Inc
			Type:           "execution",
			Metric:         "total_executions",
			Value:          420.0,
			Date:           time.Now().AddDate(0, 0, -1),
			Metadata:       models.MapToJSON(map[string]interface{}{"success_rate": 0.88, "avg_duration": 1.8}),
		},
		{
			OrganizationID: &orgs[2].ID,
			Type:           "revenue",
			Metric:         "monthly_revenue",
			Value:          1200.0,
			Date:           time.Now().AddDate(0, -1, 0),
			Metadata:       models.MapToJSON(map[string]interface{}{"subscriptions": 4, "one_time_sales": 1}),
		},
		// Add some historical data for trends
		{
			OrganizationID: &orgs[0].ID,
			Type:           "execution",
			Metric:         "total_executions",
			Value:          1100.0,
			Date:           time.Now().AddDate(0, 0, -7),
			Metadata:       models.MapToJSON(map[string]interface{}{"success_rate": 0.93, "avg_duration": 2.1}),
		},
		{
			OrganizationID: &orgs[0].ID,
			Type:           "revenue",
			Metric:         "monthly_revenue",
			Value:          4200.0,
			Date:           time.Now().AddDate(0, -2, 0),
			Metadata:       models.MapToJSON(map[string]interface{}{"subscriptions": 14, "one_time_sales": 2}),
		},
	}

	for _, analytics := range analytics {
		if err := db.Create(&analytics).Error; err != nil {
			return fmt.Errorf("failed to seed analytics %s: %v", analytics.Metric, err)
		}
	}

	fmt.Printf("✅ Created %d analytics records\n", len(analytics))
	return nil
}
