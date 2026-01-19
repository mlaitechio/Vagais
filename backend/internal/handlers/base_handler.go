package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(db *gorm.DB, cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		db:  db,
		cfg: cfg,
	}
}

// sendSuccess sends a successful response
func (h *BaseHandler) sendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// sendError sends an error response
func (h *BaseHandler) sendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// sendCreated sends a created response
func (h *BaseHandler) sendCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}

// getCurrentUserID gets the current user ID from context
func (h *BaseHandler) getCurrentUserID(c *gin.Context) (string, bool) {
	user, exists := c.Get("user")
	if !exists {
		return "", false
	}
	userObj, ok := user.(*models.User)
	if !ok {
		return "", false
	}
	return userObj.ID, true
}

// getCurrentUser gets the current user from context
func (h *BaseHandler) getCurrentUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}
