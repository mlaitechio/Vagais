package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
)

var (
	AuthServiceInstance         *AuthService
	UserServiceInstance         *UserService
	AgentServiceInstance        *AgentService
	MarketplaceServiceInstance  *MarketplaceService
	RuntimeServiceInstance      *RuntimeService
	IntegrationServiceInstance  *IntegrationService
	NotificationServiceInstance *NotificationService
	AnalyticsServiceInstance    *AnalyticsService
	BillingServiceInstance      *BillingService
	LicenseServiceInstance      *LicenseService
	PaymentServiceInstance      *PaymentService
)

// InitializeServices initializes all services with graceful fallbacks
func InitializeServices(db *gorm.DB, redisClient *redis.Client, cfg *config.Config) {
	// Initialize core services
	AuthServiceInstance = NewAuthService(db, cfg)
	UserServiceInstance = NewUserService(db, cfg)
	AgentServiceInstance = NewAgentService(db, cfg)
	MarketplaceServiceInstance = NewMarketplaceService(db, cfg)
	RuntimeServiceInstance = NewRuntimeService(db, cfg)
	IntegrationServiceInstance = NewIntegrationService(db, cfg)

	// Initialize optional services with fallbacks
	NotificationServiceInstance = NewNotificationService(db, cfg)
	AnalyticsServiceInstance = NewAnalyticsService(db, cfg)

	// Initialize payment and licensing services with fallbacks
	BillingServiceInstance = NewBillingService(db, cfg)
	LicenseServiceInstance = NewLicenseService(db, cfg)
	PaymentServiceInstance = NewPaymentService(db, cfg)
}

// Service interface for common service operations
type Service interface {
	IsAvailable() bool
	GetStatus() string
}

// BaseService provides common functionality for all services
type BaseService struct {
	db   *gorm.DB
	cfg  *config.Config
	name string
}

// NewBaseService creates a new base service
func NewBaseService(db *gorm.DB, cfg *config.Config, name string) BaseService {
	return BaseService{
		db:   db,
		cfg:  cfg,
		name: name,
	}
}

// IsAvailable checks if the service is available
func (s *BaseService) IsAvailable() bool {
	return s.db != nil
}

// GetStatus returns the service status
func (s *BaseService) GetStatus() string {
	if s.IsAvailable() {
		return "available"
	}
	return "unavailable"
}

// GetDB returns the database instance
func (s *BaseService) GetDB() *gorm.DB {
	return s.db
}

// GetConfig returns the configuration
func (s *BaseService) GetConfig() *config.Config {
	return s.cfg
}

// GetName returns the service name
func (s *BaseService) GetName() string {
	return s.name
}
