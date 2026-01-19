package services

import (
	"errors"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// MarketplaceService handles marketplace operations
type MarketplaceService struct {
	BaseService
}

// NewMarketplaceService creates a new marketplace service
func NewMarketplaceService(db *gorm.DB, cfg *config.Config) *MarketplaceService {
	return &MarketplaceService{
		BaseService: NewBaseService(db, cfg, "marketplace"),
	}
}

// CreateReviewRequest represents review creation request
type CreateReviewRequest struct {
	AgentID string `json:"agent_id" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// UpdateReviewRequest represents review update request
type UpdateReviewRequest struct {
	Rating   *int   `json:"rating"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Response string `json:"response"` // Developer response
}

// CreateReview creates a new review
func (s *MarketplaceService) CreateReview(req *CreateReviewRequest, userID string) (*models.Review, error) {
	// Check if user has already reviewed this agent
	var existingReview models.Review
	if err := s.db.Where("agent_id = ? AND user_id = ?", req.AgentID, userID).First(&existingReview).Error; err == nil {
		return nil, errors.New("user has already reviewed this agent")
	}

	review := &models.Review{
		AgentID:    req.AgentID,
		UserID:     userID,
		Rating:     req.Rating,
		Title:      req.Title,
		Content:    req.Content,
		IsVerified: false,
		IsHelpful:  0,
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	// Update agent rating
	s.updateAgentRating(req.AgentID)

	return review, nil
}

// GetReview retrieves a review by ID
func (s *MarketplaceService) GetReview(id string) (*models.Review, error) {
	var review models.Review
	if err := s.db.Preload("Agent").Preload("User").First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

// UpdateReview updates a review
func (s *MarketplaceService) UpdateReview(id string, req *UpdateReviewRequest, userID string) (*models.Review, error) {
	var review models.Review
	if err := s.db.First(&review, id).Error; err != nil {
		return nil, err
	}

	// Check if user owns this review or is admin
	if review.UserID != userID {
		return nil, errors.New("unauthorized to update this review")
	}

	updates := make(map[string]interface{})
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Response != "" {
		updates["response"] = req.Response
	}

	if err := s.db.Model(&review).Updates(updates).Error; err != nil {
		return nil, err
	}

	// Update agent rating if rating changed
	if req.Rating != nil {
		s.updateAgentRating(review.AgentID)
	}

	return &review, nil
}

// DeleteReview deletes a review
func (s *MarketplaceService) DeleteReview(id string, userID string) error {
	var review models.Review
	if err := s.db.First(&review, id).Error; err != nil {
		return err
	}

	// Check if user owns this review or is admin
	if review.UserID != userID {
		return errors.New("unauthorized to delete this review")
	}

	if err := s.db.Delete(&review).Error; err != nil {
		return err
	}

	// Update agent rating
	s.updateAgentRating(review.AgentID)

	return nil
}

// ListReviews retrieves reviews for an agent
func (s *MarketplaceService) ListReviews(agentID string, page, limit int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Preload("User")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

// SearchAgents searches for agents with various filters
func (s *MarketplaceService) SearchAgents(query string, category string, minRating float64, maxPrice float64, page, limit int) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	dbQuery := s.db.Model(&models.Agent{}).Where("is_public = ? AND is_enabled = ?", true, true).Preload("Creator").Preload("Organization")

	// Apply search filters
	if query != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ? OR tags::text ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}
	if minRating > 0 {
		dbQuery = dbQuery.Where("rating >= ?", minRating)
	}
	if maxPrice > 0 {
		dbQuery = dbQuery.Where("price <= ?", maxPrice)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Order("rating DESC, usage_count DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// GetFeaturedAgents retrieves featured agents
func (s *MarketplaceService) GetFeaturedAgents(limit int) ([]models.Agent, error) {
	var agents []models.Agent
	if err := s.db.Where("is_public = ? AND is_enabled = ?", true, true).
		Preload("Creator").Preload("Organization").
		Order("rating DESC, usage_count DESC").
		Limit(limit).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// GetTrendingAgents retrieves trending agents
func (s *MarketplaceService) GetTrendingAgents(limit int) ([]models.Agent, error) {
	var agents []models.Agent
	if err := s.db.Where("is_public = ? AND is_enabled = ?", true, true).
		Preload("Creator").Preload("Organization").
		Order("usage_count DESC, rating DESC").
		Limit(limit).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

// GetAgentCategories retrieves all categories with agent counts
func (s *MarketplaceService) GetAgentCategories() (map[string]int64, error) {
	var results []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	if err := s.db.Model(&models.Agent{}).
		Select("category, COUNT(*) as count").
		Where("is_public = ? AND is_enabled = ?", true, true).
		Group("category").
		Find(&results).Error; err != nil {
		return nil, err
	}

	categories := make(map[string]int64)
	for _, result := range results {
		categories[result.Category] = result.Count
	}

	return categories, nil
}

// MarkReviewHelpful marks a review as helpful
func (s *MarketplaceService) MarkReviewHelpful(reviewID string, userID string) error {
	var review models.Review
	if err := s.db.First(&review, reviewID).Error; err != nil {
		return err
	}

	// Increment helpful count
	return s.db.Model(&review).Update("is_helpful", review.IsHelpful+1).Error
}

// VerifyReview verifies a review (admin only)
func (s *MarketplaceService) VerifyReview(reviewID string) error {
	return s.db.Model(&models.Review{}).Where("id = ?", reviewID).Update("is_verified", true).Error
}

// updateAgentRating updates the average rating for an agent
func (s *MarketplaceService) updateAgentRating(agentID string) {
	var avgRating float64
	var reviewCount int64

	// Calculate average rating
	s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Select("AVG(rating)").Scan(&avgRating)

	// Count reviews
	s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Count(&reviewCount)

	// Update agent
	s.db.Model(&models.Agent{}).Where("id = ?", agentID).Updates(map[string]interface{}{
		"rating":       avgRating,
		"review_count": reviewCount,
	})
}

// ListMarketplaceAgents lists marketplace agents with filters
func (s *MarketplaceService) ListMarketplaceAgents(page, limit int, search, category, pricing string, rating float64) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	dbQuery := s.db.Model(&models.Agent{}).Where("is_public = ? AND is_enabled = ?", true, true).Preload("Creator").Preload("Organization")

	// Apply search filter
	if search != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ? OR tags::text ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Apply category filter
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	// Apply pricing filter
	if pricing != "" {
		switch pricing {
		case "free":
			dbQuery = dbQuery.Where("price = ?", 0)
		case "paid":
			dbQuery = dbQuery.Where("price > ?", 0)
		}
	}

	// Apply rating filter
	if rating > 0 {
		dbQuery = dbQuery.Where("rating >= ?", rating)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Order("rating DESC, usage_count DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// GetMarketplaceAgent gets marketplace agent details
func (s *MarketplaceService) GetMarketplaceAgent(agentID string) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Where("id = ? AND is_public = ? AND is_enabled = ?", agentID, true, true).
		Preload("Creator").Preload("Organization").Preload("Reviews").First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// TryMarketplaceAgent tries a marketplace agent demo
func (s *MarketplaceService) TryMarketplaceAgent(agentID string, userID string, input map[string]interface{}) (map[string]interface{}, error) {
	// Check if agent exists and is public
	var agent models.Agent
	if err := s.db.Where("id = ? AND is_public = ? AND is_enabled = ?", agentID, true, true).First(&agent).Error; err != nil {
		return nil, errors.New("agent not found or not available")
	}

	// Check if user has free trial remaining (10 free messages per agent)
	var executionCount int64
	s.db.Model(&models.Execution{}).Where("agent_id = ? AND user_id = ?", agentID, userID).Count(&executionCount)

	if executionCount >= 10 {
		return nil, errors.New("free trial limit reached. Please purchase credits to continue")
	}

	// Create execution record
	execution := &models.Execution{
		AgentID: agentID,
		UserID:  userID,
		Status:  "completed",
		Input:   models.MapToJSON(input),
		Output:  models.MapToJSON(map[string]interface{}{"message": "Demo response from " + agent.Name}),
	}

	if err := s.db.Create(execution).Error; err != nil {
		return nil, err
	}

	// Update agent usage count
	s.db.Model(&agent).Update("usage_count", agent.UsageCount+1)

	return map[string]interface{}{
		"execution_id":     execution.ID,
		"output":           map[string]interface{}{"message": "Demo response from " + agent.Name},
		"remaining_trials": 10 - executionCount - 1,
	}, nil
}

// PurchaseMarketplaceAgent purchases a marketplace agent
func (s *MarketplaceService) PurchaseMarketplaceAgent(agentID string, userID string, pricingTier string, paymentMethodID string, organizationID string) (map[string]interface{}, error) {
	// Check if agent exists and is public
	var agent models.Agent
	if err := s.db.Where("id = ? AND is_public = ? AND is_enabled = ?", agentID, true, true).First(&agent).Error; err != nil {
		return nil, errors.New("agent not found or not available")
	}

	// For now, just return success - in a real implementation, you would integrate with payment providers
	return map[string]interface{}{
		"purchase_id":  "purchase_" + agentID + "_" + userID,
		"agent_id":     agentID,
		"user_id":      userID,
		"pricing_tier": pricingTier,
		"status":       "completed",
		"message":      "Agent purchased successfully",
	}, nil
}

// GetAgentReviews gets reviews for a marketplace agent
func (s *MarketplaceService) GetAgentReviews(agentID string, page, limit int) ([]models.Review, map[string]interface{}, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Preload("User")

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&reviews).Error; err != nil {
		return nil, nil, err
	}

	// Calculate summary statistics
	var avgRating float64
	var ratingCount int64
	s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Select("AVG(rating)").Scan(&avgRating)
	s.db.Model(&models.Review{}).Where("agent_id = ?", agentID).Count(&ratingCount)

	summary := map[string]interface{}{
		"average_rating": avgRating,
		"total_reviews":  ratingCount,
		"rating_breakdown": map[string]int64{
			"5_star": 0,
			"4_star": 0,
			"3_star": 0,
			"2_star": 0,
			"1_star": 0,
		},
	}

	return reviews, summary, nil
}

// CreateAgentReview creates a review for a marketplace agent
func (s *MarketplaceService) CreateAgentReview(agentID string, userID string, rating int, title string, comment string) (*models.Review, error) {
	// Check if user has already reviewed this agent
	var existingReview models.Review
	if err := s.db.Where("agent_id = ? AND user_id = ?", agentID, userID).First(&existingReview).Error; err == nil {
		return nil, errors.New("user has already reviewed this agent")
	}

	review := &models.Review{
		AgentID:    agentID,
		UserID:     userID,
		Rating:     rating,
		Title:      title,
		Content:    comment,
		IsVerified: false,
		IsHelpful:  0,
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	// Update agent rating
	s.updateAgentRating(agentID)

	return review, nil
}

// GetMarketplaceStats retrieves marketplace statistics
func (s *MarketplaceService) GetMarketplaceStats() (map[string]interface{}, error) {
	var totalAgents int64
	var totalUsers int64
	var totalReviews int64
	var totalExecutions int64

	s.db.Model(&models.Agent{}).Where("is_public = ?", true).Count(&totalAgents)
	s.db.Model(&models.User{}).Count(&totalUsers)
	s.db.Model(&models.Review{}).Count(&totalReviews)
	s.db.Model(&models.Execution{}).Count(&totalExecutions)

	stats := map[string]interface{}{
		"total_agents":     totalAgents,
		"total_users":      totalUsers,
		"total_reviews":    totalReviews,
		"total_executions": totalExecutions,
	}

	return stats, nil
}

// ListAllAgents retrieves all public agents with pagination and filtering
func (s *MarketplaceService) ListAllAgents(page, limit int, category, sortBy, sortOrder string) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	dbQuery := s.db.Model(&models.Agent{}).Where("is_public = ? AND is_enabled = ?", true, true).Preload("Creator").Preload("Organization")

	// Apply category filter
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderClause := "rating DESC, usage_count DESC"
	if sortBy != "" {
		switch sortBy {
		case "rating":
			orderClause = "rating " + sortOrder + ", usage_count DESC"
		case "usage_count":
			orderClause = "usage_count " + sortOrder + ", rating DESC"
		case "created_at":
			orderClause = "created_at " + sortOrder + ", rating DESC"
		case "price":
			orderClause = "price " + sortOrder + ", rating DESC"
		case "downloads":
			orderClause = "downloads " + sortOrder + ", rating DESC"
		}
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Order(orderClause).Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// SearchAgentsPublic searches for agents publicly with various filters
func (s *MarketplaceService) SearchAgentsPublic(query string, category string, minRating float64, maxPrice float64, page, limit int, sortBy, sortOrder string) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	dbQuery := s.db.Model(&models.Agent{}).Where("is_public = ? AND is_enabled = ?", true, true).Preload("Creator").Preload("Organization")

	// Apply search filters
	if query != "" {
		dbQuery = dbQuery.Where("name ILIKE ? OR description ILIKE ? OR tags::text ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%")
	}
	if category != "" {
		dbQuery = dbQuery.Where("category = ?", category)
	}
	if minRating > 0 {
		dbQuery = dbQuery.Where("rating >= ?", minRating)
	}
	if maxPrice > 0 {
		dbQuery = dbQuery.Where("price <= ?", maxPrice)
	}

	// Get total count
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderClause := "rating DESC, usage_count DESC"
	if sortBy != "" {
		switch sortBy {
		case "rating":
			orderClause = "rating " + sortOrder + ", usage_count DESC"
		case "usage_count":
			orderClause = "usage_count " + sortOrder + ", rating DESC"
		case "created_at":
			orderClause = "created_at " + sortOrder + ", rating DESC"
		case "price":
			orderClause = "price " + sortOrder + ", rating DESC"
		case "downloads":
			orderClause = "downloads " + sortOrder + ", rating DESC"
		}
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := dbQuery.Offset(offset).Limit(limit).Order(orderClause).Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}
