package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// UserHandler handles user-related requests
type UserHandler struct {
	*BaseHandler
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *gorm.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(db, cfg),
		userService: services.UserServiceInstance,
	}
}

// GetProfile gets the current user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "User not found")
		return
	}

	h.sendSuccess(c, user)
}

// UpdateProfile updates the current user's profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, user)
}

// GetUser gets a user by ID (admin only)
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUser(userID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "User not found")
		return
	}

	h.sendSuccess(c, user)
}

// ListUsers lists users with pagination (admin only)
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	orgID := c.Query("org_id")

	users, total, err := h.userService.ListUsers(page, limit, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// DeactivateUser deactivates a user (admin only)
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID := c.Param("id")

	currentUserID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.userService.DeactivateUser(userID, currentUserID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "User deactivated successfully"})
}

// ActivateUser activates a user (admin only)
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID := c.Param("id")

	currentUserID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.userService.ActivateUser(userID, currentUserID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "User activated successfully"})
}

// UpdateUserRole updates a user's role (admin only)
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	currentUserID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.userService.UpdateUserRole(userID, req.Role, currentUserID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "User role updated successfully"})
}

// GetUserStats gets user statistics
func (h *UserHandler) GetUserStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := h.userService.GetUserStats(userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetOrganization gets an organization by ID
func (h *UserHandler) GetOrganization(c *gin.Context) {
	orgID := c.Param("id")

	org, err := h.userService.GetOrganization(orgID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Organization not found")
		return
	}

	h.sendSuccess(c, org)
}

// CreateOrganization creates a new organization
func (h *UserHandler) CreateOrganization(c *gin.Context) {
	var req services.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	org, err := h.userService.CreateOrganization(&req, userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, org)
}

// UpdateOrganization updates an organization
func (h *UserHandler) UpdateOrganization(c *gin.Context) {
	orgID := c.Param("id")
	var req services.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	org, err := h.userService.UpdateOrganization(orgID, &req, userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, org)
}

// ListOrganizations lists organizations with pagination
func (h *UserHandler) ListOrganizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orgs, total, err := h.userService.ListOrganizations(page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"organizations": orgs,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// GetOrganizationUsers gets users in an organization
func (h *UserHandler) GetOrganizationUsers(c *gin.Context) {
	orgID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userService.GetOrganizationUsers(orgID, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"users": users,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// InviteUserToOrganization invites a user to an organization
func (h *UserHandler) InviteUserToOrganization(c *gin.Context) {
	orgID := c.Param("id")
	var req services.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.userService.InviteUserToOrganization(&req, orgID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "User invited successfully"})
}
