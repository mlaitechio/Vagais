package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// LicenseService handles enterprise licensing
type LicenseService struct {
	BaseService
}

// NewLicenseService creates a new license service
func NewLicenseService(db *gorm.DB, cfg *config.Config) *LicenseService {
	return &LicenseService{
		BaseService: NewBaseService(db, cfg, "license"),
	}
}

// CreateLicenseRequest represents license creation request
type CreateLicenseRequest struct {
	Key            string      `json:"key" binding:"required"`
	Type           string      `json:"type"`
	OrganizationID *string     `json:"organization_id"`
	Features       models.JSON `json:"features"`
	MaxUsers       int         `json:"max_users"`
	MaxAgents      int         `json:"max_agents"`
}

// ValidateLicenseRequest represents license validation request
type ValidateLicenseRequest struct {
	LicenseKey string `json:"license_key" binding:"required"`
	Domain     string `json:"domain"`
	IPAddress  string `json:"ip_address"`
}

// CreateLicense creates a new license
func (s *LicenseService) CreateLicense(req *CreateLicenseRequest) (*models.License, error) {
	license := &models.License{
		Key:            req.Key,
		Type:           req.Type,
		Status:         "active",
		OrganizationID: req.OrganizationID,
		IssuedAt:       time.Now(),
		ExpiresAt:      nil, // ExpiresAt is not part of CreateLicenseRequest
		Features:       req.Features,
		MaxUsers:       req.MaxUsers,
		MaxAgents:      req.MaxAgents,
		IsValid:        true,
	}

	if err := s.db.Create(license).Error; err != nil {
		return nil, err
	}

	return license, nil
}

// GetLicense retrieves a license by ID
func (s *LicenseService) GetLicense(id string) (*models.License, error) {
	var license models.License
	if err := s.db.Preload("Organization").First(&license, id).Error; err != nil {
		return nil, err
	}
	return &license, nil
}

// GetLicenseByKey retrieves a license by key
func (s *LicenseService) GetLicenseByKey(key string) (*models.License, error) {
	var license models.License
	if err := s.db.Preload("Organization").Where("key = ?", key).First(&license).Error; err != nil {
		return nil, err
	}
	return &license, nil
}

// ValidateLicense validates a license
func (s *LicenseService) ValidateLicense(req *ValidateLicenseRequest) (*models.License, error) {
	license, err := s.GetLicenseByKey(req.LicenseKey)
	if err != nil {
		return nil, errors.New("invalid license key")
	}

	// Check if license is valid
	if !license.IsValid {
		return nil, errors.New("license is invalid")
	}

	// Check if license is expired
	if license.ExpiresAt != nil && license.ExpiresAt.Before(time.Now()) {
		license.IsValid = false
		s.db.Save(license)
		return nil, errors.New("license has expired")
	}

	// Check domain restrictions (if any)
	if req.Domain != "" {
		// TODO: Implement domain validation
	}

	// Check IP restrictions (if any)
	if req.IPAddress != "" {
		// TODO: Implement IP validation
	}

	return license, nil
}

// UpdateLicense updates a license
func (s *LicenseService) UpdateLicense(id string, updates map[string]interface{}) (*models.License, error) {
	var license models.License
	if err := s.db.First(&license, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&license).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &license, nil
}

// RevokeLicense revokes a license
func (s *LicenseService) RevokeLicense(id string) error {
	var license models.License
	if err := s.db.First(&license, id).Error; err != nil {
		return err
	}

	return s.db.Model(&license).Update("is_valid", false).Error
}

// ListLicenses retrieves licenses with pagination
func (s *LicenseService) ListLicenses(page, limit int, orgID *string) ([]models.License, int64, error) {
	var licenses []models.License
	var total int64

	query := s.db.Model(&models.License{}).Preload("Organization")

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&licenses).Error; err != nil {
		return nil, 0, err
	}

	return licenses, total, nil
}

// CheckLicenseUsage checks if license usage is within limits
func (s *LicenseService) CheckLicenseUsage(licenseID string) (map[string]interface{}, error) {
	var license models.License
	if err := s.db.First(&license, licenseID).Error; err != nil {
		return nil, err
	}

	var userCount int64
	var agentCount int64

	// Count users
	if license.OrganizationID != nil {
		s.db.Model(&models.User{}).Where("organization_id = ?", license.OrganizationID).Count(&userCount)
	}

	// Count agents
	if license.OrganizationID != nil {
		s.db.Model(&models.Agent{}).Where("organization_id = ?", license.OrganizationID).Count(&agentCount)
	}

	usage := map[string]interface{}{
		"users_used":       userCount,
		"users_limit":      license.MaxUsers,
		"agents_used":      agentCount,
		"agents_limit":     license.MaxAgents,
		"users_remaining":  license.MaxUsers - int(userCount),
		"agents_remaining": license.MaxAgents - int(agentCount),
	}

	return usage, nil
}

// GenerateOfflineLicense generates an offline license for air-gapped deployments
func (s *LicenseService) GenerateOfflineLicense(licenseID string, privateKeyPath string) (string, error) {
	var license models.License
	if err := s.db.First(&license, licenseID).Error; err != nil {
		return "", err
	}

	// TODO: Implement offline license generation with digital signature
	// For now, return a placeholder
	offlineLicense := fmt.Sprintf("OFFLINE_LICENSE_%s_%s", license.Key, time.Now().Format("20060102"))

	return offlineLicense, nil
}

// ValidateOfflineLicense validates an offline license
func (s *LicenseService) ValidateOfflineLicense(offlineLicense string, publicKeyPath string) (*models.License, error) {
	// TODO: Implement offline license validation
	// For now, return an error
	return nil, errors.New("offline license validation not implemented")
}

// GetLicenseStats retrieves license statistics
func (s *LicenseService) GetLicenseStats() (map[string]interface{}, error) {
	var totalLicenses int64
	var activeLicenses int64
	var expiredLicenses int64
	var trialLicenses int64

	s.db.Model(&models.License{}).Count(&totalLicenses)
	s.db.Model(&models.License{}).Where("is_valid = ?", true).Count(&activeLicenses)
	s.db.Model(&models.License{}).Where("expires_at < ?", time.Now()).Count(&expiredLicenses)
	s.db.Model(&models.License{}).Where("type = ?", "trial").Count(&trialLicenses)

	stats := map[string]interface{}{
		"total_licenses":   totalLicenses,
		"active_licenses":  activeLicenses,
		"expired_licenses": expiredLicenses,
		"trial_licenses":   trialLicenses,
	}

	return stats, nil
}

// generateLicenseKey generates a unique license key
func (s *LicenseService) generateLicenseKey() string {
	// Generate a unique license key
	key := fmt.Sprintf("AGAI-%s-%s", time.Now().Format("20060102"), uuid.New().String()[:8])
	return key
}

// loadPrivateKey loads a private key from file
func (s *LicenseService) loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	// TODO: Implement private key loading
	return nil, errors.New("private key loading not implemented")
}

// loadPublicKey loads a public key from file
func (s *LicenseService) loadPublicKey(path string) (*rsa.PublicKey, error) {
	// TODO: Implement public key loading
	return nil, errors.New("public key loading not implemented")
}
