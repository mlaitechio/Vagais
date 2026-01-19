package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// IntegrationService handles external integrations and webhooks
type IntegrationService struct {
	BaseService
}

// NewIntegrationService creates a new integration service
func NewIntegrationService(db *gorm.DB, cfg *config.Config) *IntegrationService {
	return &IntegrationService{
		BaseService: NewBaseService(db, cfg, "integration"),
	}
}

// WebhookRequest represents webhook request
type WebhookRequest struct {
	URL      string                 `json:"url" binding:"required"`
	Events   []string               `json:"events" binding:"required"`
	Secret   string                 `json:"secret"`
	IsActive bool                   `json:"is_active"`
	Metadata map[string]interface{} `json:"metadata"`
}

// LLMProvider represents LLM provider configuration
type LLMProvider struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"` // openai, anthropic, local, etc.
	Config    map[string]interface{} `json:"config"`
	IsActive  bool                   `json:"is_active"`
	RateLimit int                    `json:"rate_limit"`
	MaxTokens int                    `json:"max_tokens"`
}

// CreateWebhook creates a new webhook
func (s *IntegrationService) CreateWebhook(req *WebhookRequest, userID string, orgID *string) (*models.Webhook, error) {
	webhook := &models.Webhook{
		URL:            req.URL,
		Events:         req.Events,
		Secret:         req.Secret,
		IsActive:       req.IsActive,
		Headers:        models.MapToJSON(req.Metadata), // Using Headers field instead of Metadata
		UserID:         userID,
		OrganizationID: orgID,
	}

	if err := s.db.Create(webhook).Error; err != nil {
		return nil, err
	}

	return webhook, nil
}

// GetWebhook retrieves a webhook by ID
func (s *IntegrationService) GetWebhook(id string) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := s.db.First(&webhook, id).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

// ListWebhooks retrieves webhooks for a user/organization
func (s *IntegrationService) ListWebhooks(userID string, orgID *string) ([]models.Webhook, error) {
	var webhooks []models.Webhook
	query := s.db.Where("user_id = ?", userID)

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	}

	if err := query.Find(&webhooks).Error; err != nil {
		return nil, err
	}

	return webhooks, nil
}

// UpdateWebhook updates a webhook
func (s *IntegrationService) UpdateWebhook(id string, req *WebhookRequest, userID string) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := s.db.First(&webhook, id).Error; err != nil {
		return nil, err
	}

	// Check if user owns this webhook
	if webhook.UserID != userID {
		return nil, errors.New("unauthorized to update this webhook")
	}

	updates := make(map[string]interface{})
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Events != nil {
		updates["events"] = req.Events
	}
	if req.Secret != "" {
		updates["secret"] = req.Secret
	}
	updates["is_active"] = req.IsActive
	if req.Metadata != nil {
		updates["headers"] = models.MapToJSON(req.Metadata)
	}

	if err := s.db.Model(&webhook).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &webhook, nil
}

// DeleteWebhook deletes a webhook
func (s *IntegrationService) DeleteWebhook(id string, userID string) error {
	var webhook models.Webhook
	if err := s.db.First(&webhook, id).Error; err != nil {
		return err
	}

	// Check if user owns this webhook
	if webhook.UserID != userID {
		return errors.New("unauthorized to delete this webhook")
	}

	return s.db.Delete(&webhook).Error
}

// SendWebhook sends a webhook notification
func (s *IntegrationService) SendWebhook(webhook *models.Webhook, event string, payload map[string]interface{}) error {
	// Check if webhook is active and subscribed to this event
	if !webhook.IsActive {
		return errors.New("webhook is inactive")
	}

	eventSubscribed := false
	for _, subscribedEvent := range webhook.Events {
		if subscribedEvent == event || subscribedEvent == "*" {
			eventSubscribed = true
			break
		}
	}

	if !eventSubscribed {
		return errors.New("webhook not subscribed to this event")
	}

	// Prepare webhook payload
	webhookPayload := map[string]interface{}{
		"event":     event,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	}

	// Add signature if secret is provided
	if webhook.Secret != "" {
		// TODO: Implement signature generation
		webhookPayload["signature"] = "sha256=..."
	}

	// Send HTTP request
	jsonData, err := json.Marshal(webhookPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "AGAI-Webhook/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Log webhook delivery
	s.logWebhookDelivery(webhook.ID, event, resp.StatusCode)

	return nil
}

// logWebhookDelivery logs webhook delivery attempt
func (s *IntegrationService) logWebhookDelivery(webhookID string, event string, statusCode int) {
	// TODO: Implement webhook delivery logging
	fmt.Printf("Webhook %s event %s delivered with status %d\n", webhookID, event, statusCode)
}

// GetLLMProviders retrieves available LLM providers
func (s *IntegrationService) GetLLMProviders() ([]LLMProvider, error) {
	var dbProviders []models.LLMProvider
	if err := s.db.Where("is_active = ?", true).Find(&dbProviders).Error; err != nil {
		return nil, err
	}

	providers := make([]LLMProvider, len(dbProviders))
	for i, dbProvider := range dbProviders {
		var config map[string]interface{}
		if err := json.Unmarshal(dbProvider.Config, &config); err != nil {
			config = make(map[string]interface{})
		}

		providers[i] = LLMProvider{
			Name:      dbProvider.Name,
			Type:      dbProvider.Type,
			IsActive:  dbProvider.IsActive,
			RateLimit: dbProvider.RateLimit,
			MaxTokens: dbProvider.MaxTokens,
			Config:    config,
		}
	}

	return providers, nil
}

// TestLLMConnection tests connection to an LLM provider
func (s *IntegrationService) TestLLMConnection(provider LLMProvider) error {
	// TODO: Implement actual connection testing
	switch provider.Type {
	case "openai":
		return s.testOpenAIConnection(provider)
	case "anthropic":
		return s.testAnthropicConnection(provider)
	case "local":
		return s.testLocalConnection(provider)
	default:
		return errors.New("unsupported provider type")
	}
}

// testOpenAIConnection tests OpenAI connection
func (s *IntegrationService) testOpenAIConnection(provider LLMProvider) error {
	// TODO: Implement OpenAI connection test
	return nil
}

// testAnthropicConnection tests Anthropic connection
func (s *IntegrationService) testAnthropicConnection(provider LLMProvider) error {
	// TODO: Implement Anthropic connection test
	return nil
}

// testLocalConnection tests local LLM connection
func (s *IntegrationService) testLocalConnection(provider LLMProvider) error {
	// TODO: Implement local LLM connection test
	return nil
}

// GetIntegrationStats retrieves integration statistics
func (s *IntegrationService) GetIntegrationStats(userID string) (map[string]interface{}, error) {
	var webhookCount int64
	var activeWebhooks int64

	s.db.Model(&models.Webhook{}).Where("user_id = ?", userID).Count(&webhookCount)
	s.db.Model(&models.Webhook{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&activeWebhooks)

	stats := map[string]interface{}{
		"total_webhooks":  webhookCount,
		"active_webhooks": activeWebhooks,
		"llm_providers":   3, // TODO: Get from actual data
	}

	return stats, nil
}
