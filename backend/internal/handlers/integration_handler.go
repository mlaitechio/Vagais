package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// IntegrationHandler handles integration-related requests
type IntegrationHandler struct {
	*BaseHandler
	integrationService *services.IntegrationService
}

// NewIntegrationHandler creates a new integration handler
func NewIntegrationHandler(db *gorm.DB, cfg *config.Config) *IntegrationHandler {
	return &IntegrationHandler{
		BaseHandler:        NewBaseHandler(db, cfg),
		integrationService: services.IntegrationServiceInstance,
	}
}

// CreateWebhook creates a new webhook
func (h *IntegrationHandler) CreateWebhook(c *gin.Context) {
	var req services.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Organization ID should be passed as a separate parameter or from user context
	// For now, we'll get it from the user's organization
	var orgID *string
	user, exists := h.getCurrentUser(c)
	if exists && user.OrganizationID != nil {
		orgID = user.OrganizationID
	}

	webhook, err := h.integrationService.CreateWebhook(&req, userID, orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, webhook)
}

// GetWebhook gets a webhook by ID
func (h *IntegrationHandler) GetWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	webhook, err := h.integrationService.GetWebhook(webhookID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Webhook not found")
		return
	}

	h.sendSuccess(c, webhook)
}

// ListWebhooks lists webhooks for a user/organization
func (h *IntegrationHandler) ListWebhooks(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")

	webhooks, err := h.integrationService.ListWebhooks(userID, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, webhooks)
}

// UpdateWebhook updates a webhook
func (h *IntegrationHandler) UpdateWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	var req services.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	webhook, err := h.integrationService.UpdateWebhook(webhookID, &req, userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, webhook)
}

// DeleteWebhook deletes a webhook
func (h *IntegrationHandler) DeleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.integrationService.DeleteWebhook(webhookID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Webhook deleted successfully"})
}

// GetLLMProviders gets available LLM providers
func (h *IntegrationHandler) GetLLMProviders(c *gin.Context) {
	providers, err := h.integrationService.GetLLMProviders()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, providers)
}

// TestLLMConnection tests connection to an LLM provider
func (h *IntegrationHandler) TestLLMConnection(c *gin.Context) {
	var req services.LLMProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.integrationService.TestLLMConnection(req); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Connection test successful"})
}

// GetIntegrationStats gets integration statistics
func (h *IntegrationHandler) GetIntegrationStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := h.integrationService.GetIntegrationStats(userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}
