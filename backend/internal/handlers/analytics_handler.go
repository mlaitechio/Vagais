package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// AnalyticsHandler handles analytics-related requests
type AnalyticsHandler struct {
	*BaseHandler
	analyticsService *services.AnalyticsService
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(db *gorm.DB, cfg *config.Config) *AnalyticsHandler {
	return &AnalyticsHandler{
		BaseHandler:      NewBaseHandler(db, cfg),
		analyticsService: services.AnalyticsServiceInstance,
	}
}

// TrackEvent tracks an analytics event
func (h *AnalyticsHandler) TrackEvent(c *gin.Context) {
	var req struct {
		EventType string                 `json:"event_type" binding:"required"`
		Metric    string                 `json:"metric" binding:"required"`
		Value     float64                `json:"value"`
		AgentID   *string                `json:"agent_id"`
		OrgID     *string                `json:"org_id"`
		Metadata  map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, _ := h.getCurrentUserID(c)

	if err := h.analyticsService.TrackEvent(req.EventType, req.Metric, req.Value, &userID, req.AgentID, req.OrgID, req.Metadata); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Event tracked successfully"})
}

// GetUsageStats gets usage statistics
func (h *AnalyticsHandler) GetUsageStats(c *gin.Context) {
	userID := c.Query("user_id")
	orgID := c.Query("org_id")
	timeRange := c.DefaultQuery("time_range", "7d")

	stats, err := h.analyticsService.GetUsageStats(&userID, &orgID, timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetAgentUsageStats gets agent usage statistics
func (h *AnalyticsHandler) GetAgentUsageStats(c *gin.Context) {
	agentIDStr := c.Param("agent_id")
	agentID, err := uuid.Parse(agentIDStr)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid agent ID")
		return
	}

	timeRange := c.DefaultQuery("time_range", "7d")

	stats, err := h.analyticsService.GetAgentUsageStats(agentID, timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetUserBehaviorAnalytics gets user behavior analytics
func (h *AnalyticsHandler) GetUserBehaviorAnalytics(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	timeRange := c.DefaultQuery("time_range", "7d")

	analytics, err := h.analyticsService.GetUserBehaviorAnalytics(userUUID, timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, analytics)
}

// GetMarketplaceTrends gets marketplace trends
func (h *AnalyticsHandler) GetMarketplaceTrends(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "7d")

	trends, err := h.analyticsService.GetMarketplaceTrends(timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, trends)
}

// GetRevenueAnalytics gets revenue analytics
func (h *AnalyticsHandler) GetRevenueAnalytics(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "7d")

	analytics, err := h.analyticsService.GetRevenueAnalytics(timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, analytics)
}

// GetDeveloperMetrics gets developer metrics
func (h *AnalyticsHandler) GetDeveloperMetrics(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	timeRange := c.DefaultQuery("time_range", "7d")

	metrics, err := h.analyticsService.GetDeveloperMetrics(userUUID, timeRange)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, metrics)
}

// GenerateCustomReport generates a custom report
func (h *AnalyticsHandler) GenerateCustomReport(c *gin.Context) {
	var req struct {
		ReportType string                 `json:"report_type" binding:"required"`
		Filters    map[string]interface{} `json:"filters"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	report, err := h.analyticsService.GenerateCustomReport(req.ReportType, req.Filters)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, report)
}
