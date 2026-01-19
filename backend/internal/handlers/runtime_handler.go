package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// RuntimeHandler handles runtime-related requests
type RuntimeHandler struct {
	*BaseHandler
	runtimeService *services.RuntimeService
}

// NewRuntimeHandler creates a new runtime handler
func NewRuntimeHandler(db *gorm.DB, cfg *config.Config) *RuntimeHandler {
	return &RuntimeHandler{
		BaseHandler:    NewBaseHandler(db, cfg),
		runtimeService: services.RuntimeServiceInstance,
	}
}

// ExecuteAgent executes an agent
func (h *RuntimeHandler) ExecuteAgent(c *gin.Context) {
	var req services.ExecuteAgentRequest
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

	execution, err := h.runtimeService.ExecuteAgent(&req, userID, orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, execution)
}

// GetExecution gets an execution by ID
func (h *RuntimeHandler) GetExecution(c *gin.Context) {
	executionID := c.Param("id")

	execution, err := h.runtimeService.GetExecution(executionID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Execution not found")
		return
	}

	h.sendSuccess(c, execution)
}

// ListExecutions lists executions with filters
func (h *RuntimeHandler) ListExecutions(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	agentID := c.Query("agent_id")
	status := c.Query("status")

	executions, total, err := h.runtimeService.ListExecutions(userID, &agentID, status, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"executions": executions,
		"total":      total,
		"page":       page,
		"limit":      limit,
	})
}

// CancelExecution cancels a running execution
func (h *RuntimeHandler) CancelExecution(c *gin.Context) {
	executionID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.runtimeService.CancelExecution(executionID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Execution cancelled successfully"})
}

// GetExecutionStats gets execution statistics
func (h *RuntimeHandler) GetExecutionStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	timeRange := c.DefaultQuery("time_range", "7d")

	stats, err := h.runtimeService.GetExecutionStats(userID, timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetAgentExecutionStats gets execution statistics for a specific agent
func (h *RuntimeHandler) GetAgentExecutionStats(c *gin.Context) {
	agentID := c.Param("agent_id")

	stats, err := h.runtimeService.GetAgentExecutionStats(agentID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetActiveExecutions gets active executions for a user
func (h *RuntimeHandler) GetActiveExecutions(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	executions, err := h.runtimeService.GetActiveExecutions(userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, executions)
}
