package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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
		&models.Webhook{},
		&models.Notification{},
		&models.LLMProvider{},
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
	if err := seedLLMProviders(db); err != nil {
		fmt.Printf("[ERROR] LLM Provider seed failed: %v\n", err)
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
		&models.Execution{},
		&models.Review{},
		&models.Agent{},
		&models.User{},
		&models.Organization{},
		&models.LLMProvider{},
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
			FirstName:      "MLAI",
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

	// Tags and screenshots for the 4 new agents
	tags1, _ := json.Marshal([]string{"Conversational Insurance Analytics", "Risk & Propensity Modeling", "Personalized Outreach Generator", "RAG + Document Intelligence", "Cosmos DB + APIs"})
	screens1, _ := json.Marshal([]string{"https://agai.studio/agents/insurance-insight-copilot/screenshot1.png"})
	tags2, _ := json.Marshal([]string{"Conversational AI Interface", "Real-Time Loan Calculations", "Centralized Policy & Product Access", "Multi-Language Support", "Integration Connectors"})
	screens2, _ := json.Marshal([]string{"https://agai.studio/agents/sales-buddy/screenshot1.png"})
	tags3, _ := json.Marshal([]string{"Document & Data Ingestion", "AI-Driven Primary Analysis", "Secondary Analysis & Market Intelligence", "AI-Powered Financial Modelling", "Conversational Research Copilot", "Visualization & Insights Dashboard"})
	screens3, _ := json.Marshal([]string{"https://agai.studio/agents/ai-powered-investment-research/screenshot1.png"})
	tags4, _ := json.Marshal([]string{"Agent for Document Quality", "Document Classification", "OCR + AI Understanding (Mistral, GPT)", "Talk to Your Document", "Data Summarization Agent", "Signature Detection & Extraction", "Stamp Detection Agent", "Central Repository Agent"})
	screens4, _ := json.Marshal([]string{"https://agai.studio/agents/document-copilot/screenshot1.png"})

	agents := []models.Agent{
		{
			Name:        "Insurance Insight Copilot",
			Description: "It consolidates policy data, customer data, claims data, and interaction patterns into a single intelligent AI system that empowers insurers with natural-language analytics, automated churn prediction, and personalized customer retention messaging — all powered by Azure OpenAI, Cosmos DB, and cognitive intelligence layers.",
			Slug:        "insurance-insight-copilot",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Customer facing",
			Tags:        tags1,
			Config: models.MapToJSON(map[string]interface{}{
				"max_tokens":  2000,
				"temperature": 0.7,
				"cloud_components": []string{
					"Azure OpenAI Service",
					"Azure Cosmos DB",
					"Azure AI Search",
					"Azure Functions",
					"Azure Blob",
					"Azure Key Vault",
					"Azure Web App",
					"Azure Monitor & Log Analytics",
					"Microsoft Entra ID",
				},
			}),
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
			Rating:            4.8,
			ReviewCount:       24,
			UsageCount:        3200,
			Downloads:         156,
			Icon:              "https://agai.studio/agents/insurance-insight-copilot/icon.png",
			Screenshots:       screens1,
			Documentation:     "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",
			Repository:        "https://github.com/agai-studio/insurance-insight-copilot",
			VideoURL:          "https://youtu.be/upowcf0JB0U",
			HowItWorks:        "The Insurance Insight Copilot leverages Azure OpenAI for natural language understanding, Cosmos DB for scalable data storage, and Azure AI Search for intelligent querying. It processes policy data, customer interactions, and claims history to provide actionable insights, predict churn, and generate personalized customer communications.",
		},
		{
			Name:        "Sales Buddy",
			Description: "It consolidates scattered knowledge especially—product programs, eligibility rules, rate cards, collateral classifications, and compliance policies—into a single conversational AI assistant accessible via web, mobile or other interfaces.",
			Slug:        "sales-buddy",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Sales",
			Tags:        tags2,
			Config: models.MapToJSON(map[string]interface{}{
				"max_tokens":  1500,
				"temperature": 0.6,
				"cloud_components": []string{
					"Azure OpenAI Service",
					"Azure AI Search",
					"Azure Functions",
					"Azure Blob & SQL/Cosmos DB",
					"Azure VM",
					"Microsoft Entra ID (Azure AD)",
					"Azure Monitor & Log Analytics",
				},
			}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[3].ID, // founder
			OrganizationID:    users[3].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             0.0,
			Currency:          "USD",
			PricingModel:      "free",
			Rating:            4.7,
			ReviewCount:       18,
			UsageCount:        2850,
			Downloads:         124,
			Icon:              "https://agai.studio/agents/sales-buddy/icon.png",
			Screenshots:       screens2,
			Documentation:     "Sales assistant with real-time loan calculations and multi-language support",
			Repository:        "https://github.com/agai-studio/sales-buddy",
			VideoURL:          "https://youtu.be/QoU3DNLHadk",
			HowItWorks:        "Sales Buddy uses Azure OpenAI to understand natural language queries about products, eligibility, and pricing. It integrates with Azure AI Search to quickly retrieve relevant information from product catalogs, rate cards, and compliance documents. The system supports real-time calculations and provides accurate, contextual responses across multiple languages.",
		},
		{
			Name:        "AI Powered Investment Research",
			Description: "AI-Powered Investment Research is a Snowflake-native solution that automates financial analysis and research using LLMs, semantic search, and NLP. It helps analysts and financial institutions quickly extract insights, summarize reports, compare companies, and generate high-quality financial intelligence. Built on Snowflake Cortex with a secure multi-tenant SaaS architecture, it boosts accuracy, speed, and productivity.",
			Slug:        "ai-powered-investment-research",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Research",
			Tags:        tags3,
			Config: models.MapToJSON(map[string]interface{}{
				"max_tokens":  3000,
				"temperature": 0.5,
				"cloud_components": []string{
					"Azure App Service / Azure Functions",
					"Azure Blob Storage",
					"Azure Event Grid",
					"Azure Bot Service (Microsoft Teams)",
					"Azure API Management",
					"Microsoft Entra ID (Azure AD)",
					"Azure Key Vault",
					"Azure Monitor & Log Analytics",
					"Azure Virtual Network (Private Endpoints)",
				},
			}),
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[2].ID, // researcher
			OrganizationID:    users[2].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			Price:             0.0,
			Currency:          "USD",
			PricingModel:      "free",
			Rating:            4.9,
			ReviewCount:       32,
			UsageCount:        4100,
			Downloads:         198,
			Icon:              "https://agai.studio/agents/ai-powered-investment-research/icon.png",
			Screenshots:       screens3,
			Documentation:     "Comprehensive investment research platform with AI-driven analysis and financial modeling",
			Repository:        "https://github.com/agai-studio/ai-powered-investment-research",
			VideoURL:          "https://youtu.be/xi8j8sZYVt8z",
			HowItWorks:        "The AI-Powered Investment Research platform ingests financial documents and data into Azure Blob Storage, processes them using Azure OpenAI for analysis and summarization, and leverages Snowflake Cortex for advanced analytics. The conversational copilot enables analysts to query financial data naturally, while the visualization dashboard presents insights in an actionable format. The system is built with enterprise-grade security using Microsoft Entra ID and Azure Key Vault.",
		},
		{
			Name:        "Document Copilot",
			Description: "Document Co-Pilot is an AI-powered document intelligence platform on Microsoft Azure that uses Generative AI and LLMs to understand, validate, and act on enterprise documents. It reduces manual effort, accelerates processing, improves compliance, and turns static documents into intelligent, interactive assets.",
			Slug:        "document-copilot",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Document processing",
			Tags:        tags4,
			Config: models.MapToJSON(map[string]interface{}{
				"max_tokens":  2500,
				"temperature": 0.4,
				"cloud_components": []string{
					"AZURE AI Foundry",
					"Azure OpenAI Service (GPT/Mistral)",
					"Azure Blob Storage / SharePoint",
					"Azure Functions / Logic Apps",
					"Azure API Management",
					"Microsoft Entra ID",
					"Azure Monitor & Log Analytics",
				},
			}),
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
			Rating:            4.8,
			ReviewCount:       28,
			UsageCount:        3600,
			Downloads:         172,
			Icon:              "https://agai.studio/agents/document-copilot/icon.png",
			Screenshots:       screens4,
			Documentation:     "Enterprise document intelligence with AI-powered quality checks, classification, and extraction",
			Repository:        "https://github.com/agai-studio/document-copilot",
			VideoURL:          "https://youtu.be/txVbmkYnUs8",
			HowItWorks:        "Document Copilot processes documents stored in Azure Blob Storage or SharePoint using Azure AI Foundry and Azure OpenAI. It employs multiple specialized agents: a quality agent validates document integrity, a classification agent categorizes documents, OCR and AI understanding extract structured data, a conversational interface allows users to query documents, and detection agents identify signatures and stamps. All processing is orchestrated through Azure Functions and secured with Microsoft Entra ID.",
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
