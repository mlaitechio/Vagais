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

			Name: "General Health Information Agent,ABHI Wellness Programs & Health Calculator Agent,Qualification& Data Collection Agent,ABHI Plan Recommendation Agent,ABHI Plan Specialist (Product Details Agent),Enrollment & Contact Collection Agent,",

			Description: "Provides general health awareness, preventive care information, and health insurance related regulatory explanations without offering medical diagnosis or treatment advice.,Handles ABHI company-level queries, wellness programs, Activ DayZ, Activ Age, Healthy Heart Score, ABHI-related FAQs, and health calculator URLS for wellness assessment,Collects user details and health-related inputs required to assess eligibility and suitability for personalized ABHI health insurance plan recommendations,Analyzes collected qualification details and recommends suitable ABHI health insurance plans based on PED and Age,Provides detailed information on ABHI health insurance products including coverage, benefits, exclusions, and waiting periods,Collects user contact details for enrollment or follow-up purposes, including name, email address, and mobile number",

			Slug: "General-Health-Information-Agent,ABHI-Wellness-Programs & Health-Calculator-Agent,Qualification & Data-Collection-Agent,ABHI-Plan-Recommendation-Agent,ABHI-Plan-Specialist (Product-DetailsAgent),Enrollment & Contact -Collection-Agent",

			Version: "1.0.0",

			Status: "published",

			Type: "custom",

			Category: "Insurance",

			Tags: tags1,

			Config: models.MapToJSON(map[string]interface{}{

				"max_tokens": 2000,

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

			LLMProvider: "openai",

			LLMModel: "gpt-4",

			EmbeddingProvider: "openai",

			EmbeddingModel: "text-embedding-ada-002",

			CreatorID: users[1].ID, // developer

			OrganizationID: users[1].OrganizationID,

			IsPublic: true,

			IsEnabled: true,

			Price: 0.0,

			Currency: "USD",

			PricingModel: "free",

			Rating: 4.8,

			ReviewCount: 24,

			UsageCount: 3200,

			Downloads: 156,

			Icon: "https://agai.studio/agents/insurance-insight-copilot/icon.png",

			Screenshots: screens1,

			Documentation: "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",

			Repository: "https://github.com/agai-studio/insurance-insight-copilot",

			VideoURL: "https://youtu.be/upowcf0JB0U",

			HowItWorks: "Analyze scraped general health information data to deliver clear, non-diagnostic insights on health awareness, preventive care, and insurance regulations in a user-friendly format.,Process scraped ABHI website data to provide accurate information on ABHI wellness programs, health scores, Activ initiatives, FAQs, and direct users to relevant health calculator tools for wellness assessment.,Ask a fixed set of seven qualification questions to collect user profile and health-related information needed to assess eligibility and readiness for personalized ABHI health insurance plan recommendations.,Use responses from the seven qualification questions to evaluate age and pre-existing disease (PED) status, then recommend the most suitable ABHI health insurance plans accordingly.,Analyze scraped ABHI product PDF data to deliver clear, structured details on plan coverage, benefits, exclusions, and waiting periods to help users understand each health insurance product.,Prompt users to provide their name, email address, and mobile number, then securely capture these contact details for enrollment and follow-up communication purposes.",
		},
		{
			Name:        "Aadhaar Document Processing Agent,Passport Document Processing Agent,Driving License Document Processing Agent,PAN Card Document Processing Agent,",
			Description: "Processes Aadhaar card documents and extracts structured information for downstream use.,Processes passport images to extract structured details including name, passport number, and address in JSON format.,Processes driving license images to extract structured details such as name, license number, and address.,Processes PAN card images to extract structured details including name, PAN number, and address in JSON format.",
			Slug:        "Aadhaar-Document-Processing-Agent,Passport-Document-Processing-Agent,Driving-License-Document-Processing-Agent,PAN-Card-Document-Processing-Agent",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Banking",
			Tags:        tags2,
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
			Screenshots:       screens2,
			Documentation:     "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",
			Repository:        "https://github.com/agai-studio/insurance-insight-copilot",
			VideoURL:          "https://youtu.be/upowcf0JB0U",
			HowItWorks:        "Process Aadhaar card documents to extract structured identity information for downstream workflows.,Extract structured identity details such as name, passport number, and address from passport images.,Process driving license images to extract structured identity data including name, license number, and address.,Extract structured identity information such as name, PAN number, and address from PAN card images.",
		},
		{
			Name:        "ActivFit Plus & Preferred Plan Information Agent,ActivFit Policy Wording Agent,ActivHealth Policy Wording Agent,ActivHealth Product Benefit Table Agent,ActivOne NXT Information Agent,Super Health Top-Up Plus Benefit Table Agent,Super Health Top-Up Plus Policy Wording Agent,ActivOne Max Information Agent,",
			Description: "Retrieves and answers user queries using semantic search over ActivFit Plus and Preferred healthcare plan data.,Retrieves exact policy wording answers from ActivFit healthcare documents using semantic search.,Retrieves accurate answers from ActivHealth policy wording documents using semantic search.,Retrieves product benefit details from ActivHealth benefit tables using semantic search.,Retrieves accurate answers from ActivOne NXT healthcare plan documents using semantic search.,Retrieves benefit details from Super Health Top-Up Plus tables using semantic search.,Retrieves accurate answers from Super Health Top-Up Plus policy wording documents using semantic search.,Retrieves answers from ActivOne Max healthcare plan documents using semantic search.",
			Slug:        "ActivFit-Plus & Preferred-Plan-Information-Agent,ActivFit-Policy-Wording-Agent,ActivHealth-Policy-Wording-Agent,ActivHealth-Product-Benefit-Table- Agent,ActivOne-NXT-Information-Agent,Super-Health-Top-Up-Plus-Benefit-Table-Agent,Super-Health-Top-Up-Plus-Policy-Wording-Agent,ActivOne-Max-Information-Agent",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Insurance",
			Tags:        tags3,
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
			Screenshots:       screens3,
			Documentation:     "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",
			Repository:        "https://github.com/agai-studio/insurance-insight-copilot",
			VideoURL:          "https://youtu.be/upowcf0JB0U",
			HowItWorks:        "Answer user queries by retrieving relevant information from ActivFit Plus and Preferred plan documents using semantic search.,Provide exact policy wording responses by semantically searching ActivFit policy wording documents.,Answer user questions using official ActivHealth policy wording documents via semantic search,Retrieve and summarize benefit details from ActivHealth product benefit tables based on user queries.,Retrieve and respond to user queries using ActivOne NXT healthcare plan documents.,Answer benefit-related queries by retrieving data from Super Health Top-Up Plus benefit tables.,Provide precise answers from Super Health Top-Up Plus policy wording documents using semantic search.,Provide accurate plan-related information by searching ActivOne Max healthcare documents.,",
		},
		{
			Name:        "Product Details Agent,Calculation Agent,Finverse Guide Agent,Sales Pitch Generation Agent,",
			Description: "Responsible for answering all product-related queries such as eligibility programs, product features, policy rules, deviations, LTV/FOIR grids, risk norms, and internal product documentation. This agent dynamically injects product policies into the system context and acts as the knowledge backbone of FinWise.,Guides users through the complete Finverse and operational workflow including sourcing apps, dedupe, legal, technical, underwriting, disbursement, Salesforce processes, and troubleshooting. This agent acts as a digital process handbook for sales and operations teams,Focused on business enablement. Generates structured sales pitches, product positioning narratives, objection-handling points, mitigant recommendations, and competitive talking points based on internal sales strategy and program rules,",
			Slug:        "Product-Details-Agent,Calculation-Agent,Finverse-Guide-Agent,Sales-Pitch-Generation-Agent",
			Version:     "1.0.0",
			Status:      "published",
			Type:        "custom",
			Category:    "Sales",
			Tags:        tags4,
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
			Screenshots:       screens4,
			Documentation:     "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",
			Repository:        "https://github.com/agai-studio/insurance-insight-copilot",
			VideoURL:          "https://youtu.be/upowcf0JB0U",
			HowItWorks:        "Serve as the central product knowledge agent, resolving eligibility, policy, risk, and product rule queries.,Execute deterministic financial calculations for eligibility, repayment models, and business metrics.,Provide step-by-step guidance across Finverse workflows and operational processes.,Enable sales teams with structured pitches, objections handling, and competitive positioning.",
		},
		{

			Name: "Fields Extraction from Website Agent,Fields Extraction from Document Agent,FIelds Comparison Agent,iFinance Agent",

			Description: "Automates login and navigation of the Finverse web portal to extract required application fields and download relevant documents. Extracted data is stored in structured JSON format and documents are saved for downstream processing.,Processes downloaded documents using Azure Document Intelligence to perform OCR and extract required fields and then Converts unstructured document content into structured JSON data for validation.,Compares and verifies fields extracted from the website and documents, Used Agno Agent–based AI validation and rule-based matching logic. Generates final verification status (Verified / Not Applicable / Not Verified) along with structured reasoning.,iFinance agent is built on a corpus of 17,000+ diverse documents containing tables, graphs, and images, extracted using OCR techniques. It leverages a multimodal RAG architecture to answer image-based queries, employs hybrid search for more accurate retrieval, and uses metadata to ensure precise output localization. The chatbot is developed using the LLaMA model for OCR processing, the Gemma model for LLM intelligence, and Stella embeddings to generate vector representations across the entire document set.",

			Slug: "Fields-Extraction-from-Website-Agent,Fields-Extraction-from-Document-Agent,FIelds-Comparison-Agent,iFinance-Agent",

			Version: "1.0.0",

			Status: "published",

			Type: "custom",

			Category: "Finance",

			Tags: tags1,

			Config: models.MapToJSON(map[string]interface{}{

				"max_tokens": 2000,

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

			LLMProvider: "openai",

			LLMModel: "gpt-4",

			EmbeddingProvider: "openai",

			EmbeddingModel: "text-embedding-ada-002",

			CreatorID: users[1].ID, // developer

			OrganizationID: users[1].OrganizationID,

			IsPublic: true,

			IsEnabled: true,

			Price: 0.0,

			Currency: "USD",

			PricingModel: "free",

			Rating: 4.8,

			ReviewCount: 24,

			UsageCount: 3200,

			Downloads: 156,

			Icon: "https://agai.studio/agents/insurance-insight-copilot/icon.png",

			Screenshots: screens1,

			Documentation: "Comprehensive insurance analytics with Azure OpenAI and Cosmos DB integration",

			Repository: "https://github.com/agai-studio/insurance-insight-copilot",

			VideoURL: "https://youtu.be/upowcf0JB0U",

			HowItWorks: "Automates the extraction of fields and download the required documents from the website,Extracts required fields from Documents,Compares and Verify the fields extracted by RPA and OCR,iFinance Agent is a robust RAG-powered solution that efficiently responds to user queries by leveraging large and complex datasets.",
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
