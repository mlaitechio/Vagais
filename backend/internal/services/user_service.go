package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// UserService handles user operations
type UserService struct {
	BaseService
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, cfg *config.Config) *UserService {
	return &UserService{
		BaseService: NewBaseService(db, cfg, "user"),
	}
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Avatar      string                 `json:"avatar"`
	Preferences map[string]interface{} `json:"preferences"`
}

// CreateOrganizationRequest represents organization creation request
type CreateOrganizationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

// UpdateOrganizationRequest represents organization update request
type UpdateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
	IsActive    *bool  `json:"is_active"`
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Organization").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Organization").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id string, req *UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Preferences != nil {
		updates["preferences"] = models.MapToJSON(req.Preferences)
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(page, limit int, orgID *string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{}).Preload("Organization")

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// DeactivateUser deactivates a user
func (s *UserService) DeactivateUser(id string, adminID string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" {
		return errors.New("unauthorized to deactivate users")
	}

	return s.db.Model(&user).Update("is_active", false).Error
}

// ActivateUser activates a user
func (s *UserService) ActivateUser(id string, adminID string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" {
		return errors.New("unauthorized to activate users")
	}

	return s.db.Model(&user).Update("is_active", true).Error
}

// UpdateUserRole updates a user's role
func (s *UserService) UpdateUserRole(id string, role string, adminID string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" {
		return errors.New("unauthorized to update user roles")
	}

	validRoles := []string{"admin", "staff", "maintainer", "user"}
	valid := false
	for _, r := range validRoles {
		if r == role {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("invalid role")
	}

	return s.db.Model(&user).Update("role", role).Error
}

// GetOrganization retrieves an organization by ID
func (s *UserService) GetOrganization(id string) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.Preload("Users").First(&org, id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

// CreateOrganization creates a new organization
func (s *UserService) CreateOrganization(req *CreateOrganizationRequest, creatorID string) (*models.Organization, error) {
	org := &models.Organization{
		Name:        req.Name,
		Slug:        s.generateSlug(req.Name),
		Description: req.Description,
		Website:     req.Website,
		Logo:        req.Logo,
		IsActive:    true,
		Plan:        "free",
	}

	if err := s.db.Create(org).Error; err != nil {
		return nil, err
	}

	// Add creator as admin
	if err := s.db.Model(&models.User{}).Where("id = ?", creatorID).Update("organization_id", org.ID).Error; err != nil {
		return nil, err
	}

	return org, nil
}

// UpdateOrganization updates an organization
func (s *UserService) UpdateOrganization(id string, req *UpdateOrganizationRequest, userID string) (*models.Organization, error) {
	var org models.Organization
	if err := s.db.First(&org, id).Error; err != nil {
		return nil, err
	}

	// Check if user is admin of this organization
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	if user.OrganizationID == nil || *user.OrganizationID != id || user.Role != "admin" {
		return nil, errors.New("unauthorized to update this organization")
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
		updates["slug"] = s.generateSlug(req.Name)
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := s.db.Model(&org).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &org, nil
}

// ListOrganizations retrieves organizations with pagination
func (s *UserService) ListOrganizations(page, limit int) ([]models.Organization, int64, error) {
	var orgs []models.Organization
	var total int64

	query := s.db.Model(&models.Organization{})

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// AddUserToOrganization adds a user to an organization
func (s *UserService) AddUserToOrganization(userID, orgID string, role string, adminID string) error {
	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" || (admin.OrganizationID != nil && *admin.OrganizationID != orgID) {
		return errors.New("unauthorized to add users to organization")
	}

	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"organization_id": orgID,
		"role":            role,
	}).Error
}

// RemoveUserFromOrganization removes a user from an organization
func (s *UserService) RemoveUserFromOrganization(userID string, adminID string) error {
	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" {
		return errors.New("unauthorized to remove users from organization")
	}

	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"organization_id": nil,
		"role":            "user",
	}).Error
}

// GetUserStats retrieves user statistics
func (s *UserService) GetUserStats(userID string) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	var agentCount int64
	var executionCount int64
	var reviewCount int64

	s.db.Model(&models.Agent{}).Where("creator_id = ?", userID).Count(&agentCount)
	s.db.Model(&models.Execution{}).Where("user_id = ?", userID).Count(&executionCount)
	s.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&reviewCount)

	stats := map[string]interface{}{
		"agents_created":    agentCount,
		"executions_total":  executionCount,
		"reviews_written":   reviewCount,
		"credits_remaining": user.Credits,
		"last_login":        user.LastLoginAt,
	}

	return stats, nil
}

// generateSlug generates a URL-friendly slug
func (s *UserService) generateSlug(name string) string {
	// Simple slug generation - in production, use a proper slug library
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}

// GetOrganizationUsers retrieves users in an organization
func (s *UserService) GetOrganizationUsers(orgID string, page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("organization_id = ?", orgID).Preload("Organization")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// InviteUserRequest represents user invitation request
type InviteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

// InviteUserToOrganization invites a user to an organization
func (s *UserService) InviteUserToOrganization(req *InviteUserRequest, orgID string, adminID string) error {
	// Check if admin has permission
	var admin models.User
	if err := s.db.First(&admin, adminID).Error; err != nil {
		return err
	}

	if admin.Role != "admin" || (admin.OrganizationID != nil && *admin.OrganizationID != orgID) {
		return errors.New("unauthorized to invite users to organization")
	}

	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// User exists, add to organization
		return s.db.Model(&existingUser).Updates(map[string]interface{}{
			"organization_id": orgID,
			"role":            req.Role,
		}).Error
	}

	// User doesn't exist, create invitation (in a real app, you'd send an email)
	// For now, we'll just return success
	return nil
}
