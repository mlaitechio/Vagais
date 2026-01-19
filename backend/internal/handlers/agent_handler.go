package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// AgentHandler handles agent-related requests
type AgentHandler struct {
	*BaseHandler
	agentService *services.AgentService
}

// NewAgentHandler creates a new agent handler
func NewAgentHandler(db *gorm.DB, cfg *config.Config) *AgentHandler {
	return &AgentHandler{
		BaseHandler:  NewBaseHandler(db, cfg),
		agentService: services.AgentServiceInstance,
	}
}

// CreateAgent creates a new agent
func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req services.CreateAgentRequest
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

	agent, err := h.agentService.CreateAgent(&req, userID, orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, agent)
}

// GetAgent gets an agent by ID
func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentID := c.Param("id")
	agent, err := h.agentService.GetAgent(agentID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Agent not found")
		return
	}

	h.sendSuccess(c, agent)
}

// UpdateAgent updates an agent
func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	agentID := c.Param("id")

	var req services.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	agent, err := h.agentService.UpdateAgent(agentID, &req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, agent)
}

// DeleteAgent deletes an agent
func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	agentID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.agentService.DeleteAgent(agentID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Agent deleted successfully"})
}

// ListAgents lists agents with pagination and filters
func (h *AgentHandler) ListAgents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")
	search := c.Query("search")
	isPublicStr := c.Query("is_public")

	var isPublic *bool
	if isPublicStr != "" {
		if parsed, err := strconv.ParseBool(isPublicStr); err == nil {
			isPublic = &parsed
		}
	}

	agents, total, err := h.agentService.ListAgents(page, limit, category, search, isPublic)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"agents": agents,
		"total":  total,
		"page":   page,
		"limit":  limit,
	})
}

// EnableAgent enables an agent for a user
func (h *AgentHandler) EnableAgent(c *gin.Context) {
	agentID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.agentService.EnableAgent(agentID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Agent enabled successfully"})
}

// DisableAgent disables an agent for a user
func (h *AgentHandler) DisableAgent(c *gin.Context) {
	agentID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.agentService.DisableAgent(agentID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Agent disabled successfully"})
}

// ExecuteAgent executes an agent
func (h *AgentHandler) ExecuteAgent(c *gin.Context) {
	agentID := c.Param("id")

	var req struct {
		Input map[string]interface{} `json:"input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	execution, err := h.agentService.ExecuteAgent(agentID, userID, req.Input)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, execution)
}

// GetAgentCategories gets all agent categories
func (h *AgentHandler) GetAgentCategories(c *gin.Context) {
	categories, err := h.agentService.GetAgentCategories()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, categories)
}

// GetAgentStats gets agent statistics
func (h *AgentHandler) GetAgentStats(c *gin.Context) {
	agentID := c.Param("id")

	stats, err := h.agentService.GetAgentStats(agentID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}
