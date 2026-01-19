package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// AuthService handles authentication operations
type AuthService struct {
	BaseService
}

// NewAuthService creates a new auth service
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		BaseService: NewBaseService(db, cfg, "auth"),
	}
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=8"`
	FirstName        string `json:"first_name" binding:"required"`
	LastName         string `json:"last_name" binding:"required"`
	OrganizationName string `json:"organization_name"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *models.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Login handles user login
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	var user models.User
	if err := s.db.Preload("Organization").Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	s.db.Model(&user).Update("last_login_at", time.Now())

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Register handles user registration
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "user",
		IsActive:     true,
		Credits:      100, // Default credits
	}

	// If organization name is provided, create organization
	if req.OrganizationName != "" {
		org := &models.Organization{
			Name:        req.OrganizationName,
			Slug:        s.generateSlug(req.OrganizationName),
			Description: "",
			IsActive:    true,
			Plan:        "free",
		}

		if err := s.db.Create(org).Error; err != nil {
			return nil, err
		}

		user.OrganizationID = &org.ID
		user.Role = "admin" // First user in org is admin
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Load organization data
	s.db.Preload("Organization").First(user, user.ID)

	// Generate tokens
	accessToken, refreshToken, expiresAt, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken refreshes an access token
func (s *AuthService) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	var user models.User
	if err := s.db.Preload("Organization").First(&user, "id = ?", claims.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new tokens
	accessToken, newRefreshToken, expiresAt, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:         &user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout handles user logout
func (s *AuthService) Logout(userID string) error {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success
	return nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	return s.validateToken(tokenString)
}

// ForgotPassword handles password reset request
func (s *AuthService) ForgotPassword(email string) error {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	// Generate a reset token
	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour) // Token expires in 1 hour

	// Store the reset token
	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if err := s.db.Create(&resetToken).Error; err != nil {
		return errors.New("failed to create reset token")
	}

	// In a real implementation, you would send an email with the reset link
	// For now, we'll just print the token for testing purposes
	fmt.Printf("Password reset token generated for %s: %s\n", email, token)

	return nil
}

// ResetPassword handles password reset
func (s *AuthService) ResetPassword(token, newPassword string) error {
	var resetToken models.PasswordResetToken
	if err := s.db.Where("token = ? AND used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken).Error; err != nil {
		return errors.New("invalid or expired reset token")
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Update the user's password
	if err := s.db.Model(&models.User{}).Where("id = ?", resetToken.UserID).Update("password_hash", string(hashedPassword)).Error; err != nil {
		return errors.New("failed to update password")
	}

	// Mark the token as used
	if err := s.db.Model(&resetToken).Update("used", true).Error; err != nil {
		fmt.Printf("Failed to mark reset token as used: %v\n", err)
	}

	return nil
}

// generateTokens generates access and refresh tokens
func (s *AuthService) generateTokens(userID string, email, role string) (string, string, time.Time, error) {
	now := time.Now()
	accessExpiresAt := now.Add(15 * time.Minute)
	refreshExpiresAt := now.Add(7 * 24 * time.Hour)

	// Generate access token
	accessClaims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Generate refresh token
	refreshClaims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWT.SecretKey))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessTokenString, refreshTokenString, accessExpiresAt, nil
}

// validateToken validates a JWT token
func (s *AuthService) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateSlug generates a URL-friendly slug
func (s *AuthService) generateSlug(name string) string {
	// Simple slug generation - in production, use a proper slug library
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}

// GetUserFromToken gets user from JWT token
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := s.db.Preload("Organization").First(&user, "id = ?", claims.UserID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// IsAvailable checks if the auth service is available
func (s *AuthService) IsAvailable() bool {
	return s.BaseService.IsAvailable()
}

// GetStatus returns the service status
func (s *AuthService) GetStatus() string {
	return s.BaseService.GetStatus()
}
