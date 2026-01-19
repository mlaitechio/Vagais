package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	*BaseHandler
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(db, cfg),
		authService: services.AuthServiceInstance,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.sendSuccess(c, response)
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req services.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	h.sendCreated(c, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.sendSuccess(c, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.authService.Logout(userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	h.sendSuccess(c, gin.H{"message": "Logged out successfully"})
}

// ValidateToken validates a JWT token
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		h.sendError(c, http.StatusUnauthorized, "No token provided")
		return
	}

	// Remove "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		h.sendError(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	h.sendSuccess(c, gin.H{
		"valid":      true,
		"user_id":    claims.UserID,
		"expires_at": claims.ExpiresAt,
	})
}

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid email address")
		return
	}

	if err := h.authService.ForgotPassword(req.Email); err != nil {
		h.sendError(c, http.StatusInternalServerError, "Failed to send reset email")
		return
	}

	h.sendSuccess(c, gin.H{"message": "Password reset email sent"})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.Password); err != nil {
		h.sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Password reset successfully"})
}
