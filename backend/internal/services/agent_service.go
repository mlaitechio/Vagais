package services

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// AgentService handles agent operations
type AgentService struct {
	BaseService
}

// NewAgentService creates a new agent service
func NewAgentService(db *gorm.DB, cfg *config.Config) *AgentService {
	return &AgentService{
		BaseService: NewBaseService(db, cfg, "agent"),
	}
}

// CreateAgentRequest represents agent creation request
type CreateAgentRequest struct {
	Name              string                 `json:"name" binding:"required"`
	Description       string                 `json:"description"`
	Category          string                 `json:"category"`
	Tags              []string               `json:"tags"`
	Config            map[string]interface{} `json:"config"`
	LLMProvider       string                 `json:"llm_provider"`
	LLMModel          string                 `json:"llm_model"`
	EmbeddingProvider string                 `json:"embedding_provider"`
	EmbeddingModel    string                 `json:"embedding_model"`
	IsPublic          bool                   `json:"is_public"`
	Price             float64                `json:"price"`
	PricingModel      string                 `json:"pricing_model"`
}

// UpdateAgentRequest represents agent update request
type UpdateAgentRequest struct {
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	Category          string                 `json:"category"`
	Tags              []string               `json:"tags"`
	Config            map[string]interface{} `json:"config"`
	LLMProvider       string                 `json:"llm_provider"`
	LLMModel          string                 `json:"llm_model"`
	EmbeddingProvider string                 `json:"embedding_provider"`
	EmbeddingModel    string                 `json:"embedding_model"`
	IsPublic          *bool                  `json:"is_public"`
	IsEnabled         *bool                  `json:"is_enabled"`
	Price             *float64               `json:"price"`
	PricingModel      string                 `json:"pricing_model"`
}

// CreateAgent creates a new agent
func (s *AgentService) CreateAgent(req *CreateAgentRequest, creatorID string, orgID *string) (*models.Agent, error) {
	tagsJSON, _ := json.Marshal(req.Tags)
	agent := &models.Agent{
		Name:              req.Name,
		Description:       req.Description,
		Slug:              s.generateSlug(req.Name),
		Category:          req.Category,
		Tags:              tagsJSON,
		Config:            models.MapToJSON(req.Config),
		LLMProvider:       req.LLMProvider,
		LLMModel:          req.LLMModel,
		EmbeddingProvider: req.EmbeddingProvider,
		EmbeddingModel:    req.EmbeddingModel,
		CreatorID:         creatorID,
		OrganizationID:    orgID,
		IsPublic:          req.IsPublic,
		Price:             req.Price,
		PricingModel:      req.PricingModel,
		Status:            "draft",
		Type:              "custom",
		Version:           "1.0.0",
	}

	if err := s.db.Create(agent).Error; err != nil {
		return nil, err
	}

	return agent, nil
}

// GetAgent retrieves an agent by ID
func (s *AgentService) GetAgent(id string) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Preload("Creator").Preload("Organization").First(&agent, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetAgentBySlug retrieves an agent by slug
func (s *AgentService) GetAgentBySlug(slug string) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.Preload("Creator").Preload("Organization").Where("slug = ?", slug).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// UpdateAgent updates an agent
func (s *AgentService) UpdateAgent(id string, req *UpdateAgentRequest) (*models.Agent, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
		updates["slug"] = s.generateSlug(req.Name)
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != nil {
		tagsJSON, _ := json.Marshal(req.Tags)
		updates["tags"] = tagsJSON
	}
	if req.Config != nil {
		updates["config"] = models.MapToJSON(req.Config)
	}
	if req.LLMProvider != "" {
		updates["llm_provider"] = req.LLMProvider
	}
	if req.LLMModel != "" {
		updates["llm_model"] = req.LLMModel
	}
	if req.EmbeddingProvider != "" {
		updates["embedding_provider"] = req.EmbeddingProvider
	}
	if req.EmbeddingModel != "" {
		updates["embedding_model"] = req.EmbeddingModel
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.PricingModel != "" {
		updates["pricing_model"] = req.PricingModel
	}

	if err := s.db.Model(&agent).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &agent, nil
}

// DeleteAgent deletes an agent
func (s *AgentService) DeleteAgent(id string, userID string) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user is the creator or has admin rights
	if agent.CreatorID != userID {
		return errors.New("unauthorized to delete this agent")
	}

	return s.db.Delete(&agent).Error
}

// ListAgents retrieves agents with filtering and pagination
func (s *AgentService) ListAgents(page, limit int, category, search string, isPublic *bool) ([]models.Agent, int64, error) {
	var agents []models.Agent
	var total int64

	query := s.db.Model(&models.Agent{}).Preload("Creator").Preload("Organization")

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// EnableAgent enables an agent for use
func (s *AgentService) EnableAgent(id string, userID string) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user can enable this agent
	if agent.CreatorID != userID {
		return errors.New("unauthorized to enable this agent")
	}

	return s.db.Model(&agent).Update("is_enabled", true).Error
}

// DisableAgent disables an agent
func (s *AgentService) DisableAgent(id string, userID string) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user can disable this agent
	if agent.CreatorID != userID {
		return errors.New("unauthorized to disable this agent")
	}

	return s.db.Model(&agent).Update("is_enabled", false).Error
}

// ExecuteAgent executes an agent with given input
func (s *AgentService) ExecuteAgent(id string, userID string, input map[string]interface{}) (*models.Execution, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Check if agent is enabled
	if !agent.IsEnabled {
		return nil, errors.New("agent is not enabled")
	}

	// Create execution record
	execution := &models.Execution{
		AgentID: id,
		UserID:  userID,
		Status:  "running",
		Input:   models.MapToJSON(input),
		Output:  models.JSON{},
	}

	if err := s.db.Create(execution).Error; err != nil {
		return nil, err
	}

	// TODO: Implement actual agent execution logic
	// For now, we'll simulate execution
	time.Sleep(2 * time.Second)

	// Update execution with result
	output := map[string]interface{}{
		"result": "Agent execution completed successfully",
		"data":   input,
	}

	execution.Status = "completed"
	execution.Output = models.MapToJSON(output)
	execution.Duration = 2000 // 2 seconds in milliseconds

	if err := s.db.Save(execution).Error; err != nil {
		return nil, err
	}

	// Update agent usage count
	s.db.Model(&agent).Update("usage_count", agent.UsageCount+1)

	return execution, nil
}

// GetAgentCategories retrieves all available agent categories
func (s *AgentService) GetAgentCategories() ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.Agent{}).Distinct("category").Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAgentStats retrieves agent statistics
func (s *AgentService) GetAgentStats(agentID string) (map[string]interface{}, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", agentID).Error; err != nil {
		return nil, err
	}

	var executionCount int64
	s.db.Model(&models.Execution{}).Where("agent_id = ?", agentID).Count(&executionCount)

	var avgDuration float64
	s.db.Model(&models.Execution{}).Where("agent_id = ? AND status = ?", agentID, "completed").Select("AVG(duration)").Scan(&avgDuration)

	stats := map[string]interface{}{
		"total_executions": executionCount,
		"avg_duration":     avgDuration,
		"rating":           agent.Rating,
		"review_count":     agent.ReviewCount,
		"usage_count":      agent.UsageCount,
		"downloads":        agent.Downloads,
	}

	return stats, nil
}

// generateSlug generates a URL-friendly slug
func (s *AgentService) generateSlug(name string) string {
	// Simple slug generation - in production, use a proper slug library
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}
