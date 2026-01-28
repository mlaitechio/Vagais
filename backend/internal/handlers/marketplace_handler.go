package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// MarketplaceHandler handles marketplace-related requests
type MarketplaceHandler struct {
	*BaseHandler
	marketplaceService *services.MarketplaceService
}

// NewMarketplaceHandler creates a new marketplace handler
func NewMarketplaceHandler(db *gorm.DB, cfg *config.Config) *MarketplaceHandler {
	return &MarketplaceHandler{
		BaseHandler:        NewBaseHandler(db, cfg),
		marketplaceService: services.MarketplaceServiceInstance,
	}
}

// CreateReview creates a new review
func (h *MarketplaceHandler) CreateReview(c *gin.Context) {
	var req services.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	review, err := h.marketplaceService.CreateReview(&req, userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, review)
}

// GetReview gets a review by ID
func (h *MarketplaceHandler) GetReview(c *gin.Context) {
	reviewID := c.Param("id")

	review, err := h.marketplaceService.GetReview(reviewID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Review not found")
		return
	}

	h.sendSuccess(c, review)
}

// UpdateReview updates a review
func (h *MarketplaceHandler) UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	var req services.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	review, err := h.marketplaceService.UpdateReview(reviewID, &req, userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, review)
}

// DeleteReview deletes a review
func (h *MarketplaceHandler) DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.marketplaceService.DeleteReview(reviewID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Review deleted successfully"})
}

// ListReviews lists reviews for an agent
func (h *MarketplaceHandler) ListReviews(c *gin.Context) {
	agentID := c.Param("agent_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, total, err := h.marketplaceService.ListReviews(agentID, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"reviews": reviews,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// SearchAgents searches for agents
func (h *MarketplaceHandler) SearchAgents(c *gin.Context) {
	query := c.Query("query")
	category := c.Query("category")
	minRatingStr := c.Query("min_rating")
	maxPriceStr := c.Query("max_price")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var minRating, maxPrice float64
	if minRatingStr != "" {
		if parsed, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			minRating = parsed
		}
	}
	if maxPriceStr != "" {
		if parsed, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = parsed
		}
	}

	agents, total, err := h.marketplaceService.SearchAgents(query, category, minRating, maxPrice, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"data":        agents,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// GetFeaturedAgents gets featured agents
func (h *MarketplaceHandler) GetFeaturedAgents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	agents, err := h.marketplaceService.GetFeaturedAgents(limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, agents)
}

// GetTrendingAgents gets trending agents
func (h *MarketplaceHandler) GetTrendingAgents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	agents, err := h.marketplaceService.GetTrendingAgents(limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, agents)
}

// GetAgentCategories gets agent categories with counts
func (h *MarketplaceHandler) GetAgentCategories(c *gin.Context) {
	categories, err := h.marketplaceService.GetAgentCategories()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, categories)
}

// MarkReviewHelpful marks a review as helpful
func (h *MarketplaceHandler) MarkReviewHelpful(c *gin.Context) {
	reviewID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.marketplaceService.MarkReviewHelpful(reviewID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Review marked as helpful"})
}

// GetMarketplaceStats gets marketplace statistics
func (h *MarketplaceHandler) GetMarketplaceStats(c *gin.Context) {
	stats, err := h.marketplaceService.GetMarketplaceStats()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// ListAllAgents lists all public agents (no auth required)
func (h *MarketplaceHandler) ListAllAgents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	category := c.Query("category")
	sortBy := c.DefaultQuery("sort_by", "rating")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	agents, total, err := h.marketplaceService.ListAllAgents(page, limit, category, sortBy, sortOrder)
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

// SearchAgentsPublic searches for agents publicly (no auth required)
func (h *MarketplaceHandler) SearchAgentsPublic(c *gin.Context) {
	query := c.Query("q")
	category := c.Query("category")
	minRatingStr := c.Query("min_rating")
	maxPriceStr := c.Query("max_price")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sortBy := c.DefaultQuery("sort_by", "rating")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	var minRating, maxPrice float64
	if minRatingStr != "" {
		if parsed, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			minRating = parsed
		}
	}
	if maxPriceStr != "" {
		if parsed, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = parsed
		}
	}

	agents, total, err := h.marketplaceService.SearchAgentsPublic(query, category, minRating, maxPrice, page, limit, sortBy, sortOrder)
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

// ListMarketplaceAgents lists marketplace agents with filters
func (h *MarketplaceHandler) ListMarketplaceAgents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	category := c.Query("category")
	pricing := c.Query("pricing")
	ratingStr := c.Query("rating")

	var rating float64
	if ratingStr != "" {
		if parsed, err := strconv.ParseFloat(ratingStr, 64); err == nil {
			rating = parsed
		}
	}

	agents, total, err := h.marketplaceService.ListMarketplaceAgents(page, limit, search, category, pricing, rating)
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

// GetMarketplaceAgent gets marketplace agent details
func (h *MarketplaceHandler) GetMarketplaceAgent(c *gin.Context) {
	agentID := c.Param("id")

	agent, err := h.marketplaceService.GetMarketplaceAgent(agentID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Agent not found")
		return
	}

	h.sendSuccess(c, agent)
}

// TryMarketplaceAgent tries a marketplace agent demo
func (h *MarketplaceHandler) TryMarketplaceAgent(c *gin.Context) {
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

	result, err := h.marketplaceService.TryMarketplaceAgent(agentID, userID, req.Input)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, result)
}

// PurchaseMarketplaceAgent purchases a marketplace agent
func (h *MarketplaceHandler) PurchaseMarketplaceAgent(c *gin.Context) {
	agentID := c.Param("id")

	var req struct {
		PricingTier    string `json:"pricing_tier" binding:"required"`
		OrganizationID string `json:"organization_id"`
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

	purchase, err := h.marketplaceService.PurchaseMarketplaceAgent(agentID, userID, req.PricingTier, req.OrganizationID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, purchase)
}

// GetAgentReviews gets reviews for a marketplace agent
func (h *MarketplaceHandler) GetAgentReviews(c *gin.Context) {
	agentID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, summary, err := h.marketplaceService.GetAgentReviews(agentID, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"reviews": reviews,
		"summary": summary,
		"page":    page,
		"limit":   limit,
	})
}

// CreateAgentReview creates a review for a marketplace agent
func (h *MarketplaceHandler) CreateAgentReview(c *gin.Context) {
	agentID := c.Param("id")

	var req struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
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

	review, err := h.marketplaceService.CreateAgentReview(agentID, userID, req.Rating, req.Title, req.Comment)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, review)
}
