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
			Email:          "admin@mlaidigital.com",
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
			Email:          "developer@mlaidigital.com",
			Username:       "developer",
			FirstName:      "MLAI",
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
			Email:          "researcher@mlaidigital.com",
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
			Email:          "founder@mlaidigital.com",
			Username:       "founder",
			FirstName:      "MLAI",
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

func seedAgents(db *gorm.DB) error {
	fmt.Println("Seeding agents...")

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}
	if len(users) < 2 {
		return fmt.Errorf("not enough users found")
	}

	commonConfig := models.MapToJSON(map[string]interface{}{
		"max_tokens":  2000,
		"temperature": 0.7,
	})
	commonTagsSales, _ := json.Marshal([]string{"Sales", "AI"})

	commonTagsInsurance, _ := json.Marshal([]string{"Insurance", "AI", "RAG"})
	commonTagsFinance, _ := json.Marshal([]string{"Finance", "AI", "OCR"})
	commonScreens, _ := json.Marshal([]string{
		"https://agai.studio/agents/default/screenshot1.png",
	})

	agents := []models.Agent{

		// =======================
		// ABHI – TOF (Insurance)
		// =======================

		{
			Name:        "General Health Information Agent",
			Description: "Provides general health awareness, preventive care information, and health insurance related regulatory explanations without offering medical diagnosis or treatment advice.",
			Slug:        "general-health-information-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true, PricingModel: "free",
			Icon:        "https://agai.studio/agents/default/icon.png",
			Screenshots: commonScreens,
			HowItWorks:  "Analyze scraped general health information data to deliver clear, non-diagnostic insights.",
		},
		{
			Name:        "Wellness Programs & Health Calculator Agent",
			Description: "Handles ABHI company-level queries, wellness programs, Activ DayZ, Activ Age, Healthy Heart Score, FAQs, and health calculator URLs.",
			Slug:        "abhi-wellness-programs-health-calculator-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			Screenshots: commonScreens,
			HowItWorks:  "Process scraped ABHI website data to answer wellness and calculator-related queries.",
		},
		{
			Name:        "Qualification & Data Collection Agent",
			Description: "Collects user profile and health-related inputs required to assess eligibility for ABHI plans.",
			Slug:        "qualification-data-collection-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Ask seven fixed qualification questions and store responses.",
		},
		{
			Name:        "Plan Recommendation Agent",
			Description: "Recommends suitable ABHI health insurance plans based on PED and age.",
			Slug:        "abhi-plan-recommendation-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Evaluate qualification answers and recommend best-matching ABHI plans.",
		},
		{
			Name:        "Plan Specialist (Product Details Agent)",
			Description: "Provides detailed product information including coverage, exclusions, and waiting periods.",
			Slug:        "abhi-plan-specialist-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Analyze ABHI product PDFs and present structured plan details.",
		},
		{
			Name:        "Enrollment & Contact Collection Agent",
			Description: "Collects user name, email, and mobile number for enrollment and follow-up.",
			Slug:        "enrollment-contact-collection-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Prompt and securely capture contact details.",
		},

		// =======================
		// Agentic AI – KYC (Finance)
		// =======================

		{
			Name:        "Aadhaar Document Processing Agent",
			Description: "Processes Aadhaar card documents and extracts structured identity information.",
			Slug:        "aadhaar-document-processing-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Extract Aadhaar identity fields using OCR and AI.",
		},
		{
			Name:        "Passport Document Processing Agent",
			Description: "Extracts structured passport details such as name, number, and address.",
			Slug:        "passport-document-processing-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
		},
		{
			Name:        "Driving License Document Processing Agent",
			Description: "Processes driving license images and extracts identity information.",
			Slug:        "driving-license-document-processing-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
		},
		{
			Name:        "PAN Card Document Processing Agent",
			Description: "Extracts structured PAN card identity details.",
			Slug:        "pan-card-document-processing-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
		},

		// =======================
		// iFinance (GCP)
		// =======================

		{
			Name:        "iFinance Agent",
			Description: "Multimodal RAG agent built on 17,000+ documents using OCR, hybrid search, and metadata grounding.",
			Slug:        "ifinance-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Answer complex finance queries over large OCR-processed datasets using multimodal RAG.",
		},
		// =======================
		// Additional Enterprise Agents
		// =======================
		// =======================
		// ABCD Doc Comparison – Insurance
		// =======================

		{
			Name:        "ActivFit Plus & Preferred Plan Information Agent",
			Description: "Retrieves and answers user queries using semantic search over ActivFit Plus and Preferred healthcare plan documents.",
			Slug:        "activfit-plus-preferred-plan-information-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Uses semantic search to retrieve relevant information from ActivFit Plus and Preferred plan documents.",
		},

		{
			Name:        "ActivFit Policy Wording Agent",
			Description: "Provides exact policy wording answers from ActivFit healthcare documents.",
			Slug:        "activfit-policy-wording-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Retrieves precise policy wording responses using semantic search over ActivFit policy documents.",
		},

		{
			Name:        "ActivHealth Policy Wording Agent",
			Description: "Answers user questions using official ActivHealth policy wording documents.",
			Slug:        "activhealth-policy-wording-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Searches ActivHealth policy wording documents to deliver accurate and compliant responses.",
		},

		{
			Name:        "ActivHealth Product Benefit Table Agent",
			Description: "Retrieves and summarizes benefit details from ActivHealth product benefit tables.",
			Slug:        "activhealth-product-benefit-table-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Uses semantic search to retrieve benefit-level details from ActivHealth tables.",
		},

		{
			Name:        "ActivOne NXT Information Agent",
			Description: "Responds to user queries using ActivOne NXT healthcare plan documents.",
			Slug:        "activone-nxt-information-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Retrieves plan-level information from ActivOne NXT documents using semantic search.",
		},

		{
			Name:        "Super Health Top-Up Plus Benefit Table Agent",
			Description: "Answers benefit-related queries using Super Health Top-Up Plus benefit tables.",
			Slug:        "super-health-top-up-plus-benefit-table-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Retrieves benefit details from Super Health Top-Up Plus tables using semantic search.",
		},

		{
			Name:        "Super Health Top-Up Plus Policy Wording Agent",
			Description: "Provides precise answers from Super Health Top-Up Plus policy wording documents.",
			Slug:        "super-health-top-up-plus-policy-wording-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Uses semantic search to return exact wording from Super Health Top-Up Plus policy documents.",
		},

		{
			Name:        "ActivOne Max Information Agent",
			Description: "Provides accurate plan-related information from ActivOne Max healthcare documents.",
			Slug:        "activone-max-information-agent",
			Category:    "Insurance",
			Tags:        commonTagsInsurance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Retrieves answers from ActivOne Max plan documents using semantic search.",
		},
		// =======================
		// ABHFL Sales Pitch – Sales
		// =======================

		{
			Name:        "Product Details Agent",
			Description: "Central product knowledge agent covering eligibility, policy rules, risk norms, and internal documentation.",
			Slug:        "product-details-agent",
			Category:    "Sales",
			Tags:        commonTagsSales,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Acts as the single source of truth for product policies and sales enablement.",
		},

		{
			Name:        "Calculation Agent",
			Description: "Executes deterministic financial calculations for loan eligibility and repayment scenarios.",
			Slug:        "calculation-agent",
			Category:    "Sales",
			Tags:        commonTagsSales,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Performs EMI, BTS, FOIR, LTV, pension eligibility, and part-payment calculations.",
		},

		{
			Name:        "Finverse Guide Agent",
			Description: "Guides users through end-to-end Finverse workflows and operational processes.",
			Slug:        "finverse-guide-agent",
			Category:    "Sales",
			Tags:        commonTagsSales,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Provides step-by-step guidance across sourcing, underwriting, disbursement, and Salesforce workflows.",
		},

		{
			Name:        "Sales Pitch Generation Agent",
			Description: "Generates structured sales pitches, objection handling, and competitive positioning content.",
			Slug:        "sales-pitch-generation-agent",
			Category:    "Sales",
			Tags:        commonTagsSales,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Creates sales narratives, mitigant strategies, and objection-handling points.",
		},
		// =======================
		// NDC – Finance Automation
		// =======================

		{
			Name:        "Fields Extraction from Website Agent",
			Description: "Automates extraction of application fields and documents from web portals.",
			Slug:        "fields-extraction-from-website-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Uses RPA to log in, navigate portals, extract fields, and download documents.",
		},

		{
			Name:        "Fields Extraction from Document Agent",
			Description: "Extracts required fields from documents using OCR and document intelligence.",
			Slug:        "fields-extraction-from-document-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Uses Azure Document Intelligence to extract structured data from documents.",
		},

		{
			Name:        "Fields Comparison Agent",
			Description: "Compares and verifies fields extracted from website and documents.",
			Slug:        "fields-comparison-agent",
			Category:    "Finance",
			Tags:        commonTagsFinance,
			Config:      commonConfig,
			LLMProvider: "openai", LLMModel: "gpt-4",
			EmbeddingProvider: "openai", EmbeddingModel: "text-embedding-ada-002",
			CreatorID: users[1].ID, OrganizationID: users[1].OrganizationID,
			IsPublic: true, IsEnabled: true,
			HowItWorks: "Compares RPA and OCR outputs using rule-based logic and AI reasoning to generate verification status.",
		},

		{
			Name:              "Document Copilot",
			Description:       "AI-powered document intelligence agent for contract analysis, summarization, clause extraction, and contract creation workflows.",
			Slug:              "document-copilot",
			Category:          "AI Ops",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			PricingModel:      "free",
			HowItWorks:        "Uses OCR and LLMs to analyze documents, extract clauses, summarize contracts, detect risks, and generate new contracts using templates.",
		},

		{
			Name:              "Trade Finance Copilot",
			Description:       "Automates KYC checks and end-to-end trade finance workflows including document verification and compliance validation.",
			Slug:              "trade-finance-copilot",
			Category:          "Finance",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Automates KYC verification, trade document validation, discrepancy detection, and regulatory compliance using AI-driven workflows.",
		},

		{
			Name:              "Call Center Analytics",
			Description:       "Agentic AI framework for analyzing call center conversations, IVR flows, sentiment, compliance, and agent performance.",
			Slug:              "call-center-analytics",
			Category:          "AI Ops",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Analyzes voice transcripts and IVR flows to generate insights on sentiment, intent, compliance breaches, and agent effectiveness.",
		},

		{
			Name:              "Sales Intelligent Advisor",
			Description:       "AI-powered wealth advisory and sales assistant that provides personalized recommendations and sales guidance.",
			Slug:              "sales-intelligent-advisor",
			Category:          "Sales",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Combines customer data, portfolio insights, and sales rules to deliver intelligent investment advice and sales strategies.",
		},

		{
			Name:              "Website Search BOT",
			Description:       "GenAI-based website chatbot designed for customer support, upselling, and cross-selling use cases.",
			Slug:              "website-search-bot",
			Category:          "AI Ops",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Uses RAG over website content to answer queries, recommend products, and drive upsell and cross-sell opportunities.",
		},

		{
			Name:              "Loan Underwriting Copilot",
			Description:       "Automates underwriting decisions for loans and insurance using alternative data including social media signals.",
			Slug:              "loan-underwriting-copilot",
			Category:          "Finance",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Evaluates borrower risk using traditional financial data combined with social media and behavioral signals.",
		},

		{
			Name:              "SharePoint Agents",
			Description:       "Intelligent agents for SharePoint data exploration and Power BI reporting using natural language.",
			Slug:              "sharepoint-agents",
			Category:          "AI Ops",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Allows users to query SharePoint documents and generate Power BI reports using conversational AI.",
		},

		{
			Name:              "Investment Research Tool for Stocks",
			Description:       "LLM-based investment research agent for stock analysis, insights, and recommendations.",
			Slug:              "investment-research-stocks",
			Category:          "Finance",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Analyzes financial statements, news, and market data to generate stock insights and investment recommendations.",
		},

		{
			Name:              "AgenticAI for SOC/NOC",
			Description:       "Agentic AI framework to support SOC and NOC operations including alerts, triage, and remediation.",
			Slug:              "agentic-ai-soc-noc",
			Category:          "Security",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Uses multiple autonomous agents to analyze alerts, correlate incidents, recommend remediation, and reduce MTTR.",
		},

		{
			Name:              "Fraud Analytics",
			Description:       "Real-time fraud detection platform combining AI/LLM models with traditional machine learning.",
			Slug:              "fraud-analytics",
			Category:          "Finance",
			Config:            commonConfig,
			LLMProvider:       "openai",
			LLMModel:          "gpt-4",
			EmbeddingProvider: "openai",
			EmbeddingModel:    "text-embedding-ada-002",
			CreatorID:         users[1].ID,
			OrganizationID:    users[1].OrganizationID,
			IsPublic:          true,
			IsEnabled:         true,
			HowItWorks:        "Monitors real-time transactions to detect fraud patterns using AI, LLM reasoning, and ML anomaly detection.",
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
