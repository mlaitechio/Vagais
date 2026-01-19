package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// RuntimeService handles agent execution and runtime management
type RuntimeService struct {
	BaseService
}

// NewRuntimeService creates a new runtime service
func NewRuntimeService(db *gorm.DB, cfg *config.Config) *RuntimeService {
	return &RuntimeService{
		BaseService: NewBaseService(db, cfg, "runtime"),
	}
}

// ExecuteAgentRequest represents agent execution request
type ExecuteAgentRequest struct {
	AgentID   string                 `json:"agent_id" binding:"required"`
	Input     map[string]interface{} `json:"input" binding:"required"`
	SessionID string                 `json:"session_id,omitempty"`
}

// ExecuteAgent executes an agent with the given input
func (s *RuntimeService) ExecuteAgent(req *ExecuteAgentRequest, userID string, orgID *string) (*models.Execution, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", req.AgentID).Error; err != nil {
		return nil, errors.New("agent not found")
	}
	if !agent.IsEnabled {
		return nil, errors.New("agent is not enabled")
	}
	if !agent.IsPublic && agent.CreatorID != userID {
		return nil, errors.New("unauthorized to execute this agent")
	}
	execution := &models.Execution{
		AgentID:        req.AgentID,
		UserID:         userID,
		OrganizationID: orgID,
		Status:         "running",
		Input:          models.MapToJSON(req.Input),
		Output:         models.MapToJSON(map[string]interface{}{}),
		SessionID:      req.SessionID,
	}
	if err := s.db.Create(execution).Error; err != nil {
		return nil, err
	}
	go s.executeAgentAsync(execution, &agent)
	return execution, nil
}

// executeAgentAsync executes the agent asynchronously
func (s *RuntimeService) executeAgentAsync(execution *models.Execution, agent *models.Agent) {
	startTime := time.Now()

	// TODO: Implement actual agent execution logic
	// This is a placeholder for the actual execution
	time.Sleep(2 * time.Second) // Simulate processing time

	// Update execution with result
	output := map[string]interface{}{
		"result":    fmt.Sprintf("Agent '%s' executed successfully", agent.Name),
		"data":      execution.Input,
		"timestamp": time.Now().Unix(),
	}

	execution.Status = "completed"
	execution.Output = models.MapToJSON(output)
	execution.Duration = int64(time.Since(startTime).Milliseconds())

	if err := s.db.Save(execution).Error; err != nil {
		// Log error but don't fail the execution
		fmt.Printf("Error saving execution: %v\n", err)
	}

	// Update agent usage count
	s.db.Model(agent).Update("usage_count", agent.UsageCount+1)
}

// GetExecution retrieves an execution by ID
func (s *RuntimeService) GetExecution(id string) (*models.Execution, error) {
	var execution models.Execution
	if err := s.db.Preload("Agent").Preload("User").Preload("Organization").First(&execution, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

// ListExecutions retrieves executions with filtering and pagination
func (s *RuntimeService) ListExecutions(userID string, agentID *string, status string, page, limit int) ([]models.Execution, int64, error) {
	var executions []models.Execution
	var total int64

	query := s.db.Model(&models.Execution{}).Preload("Agent").Preload("User").Preload("Organization")

	// Apply filters
	query = query.Where("user_id = ?", userID)
	if agentID != nil {
		query = query.Where("agent_id = ?", *agentID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// CancelExecution cancels a running execution
func (s *RuntimeService) CancelExecution(id string, userID string) error {
	var execution models.Execution
	if err := s.db.First(&execution, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user owns this execution
	if execution.UserID != userID {
		return errors.New("unauthorized to cancel this execution")
	}

	// Only allow cancellation of running executions
	if execution.Status != "running" {
		return errors.New("can only cancel running executions")
	}

	return s.db.Model(&execution).Update("status", "cancelled").Error
}

// GetExecutionStats retrieves execution statistics
func (s *RuntimeService) GetExecutionStats(userID string, timeRange string) (map[string]interface{}, error) {
	var totalExecutions int64
	var successfulExecutions int64
	var failedExecutions int64
	var avgDuration float64

	query := s.db.Model(&models.Execution{}).Where("user_id = ?", userID)

	// Apply time range filter
	switch timeRange {
	case "today":
		query = query.Where("created_at >= ?", time.Now().Truncate(24*time.Hour))
	case "week":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
	}

	// Get counts
	query.Count(&totalExecutions)
	query.Where("status = ?", "completed").Count(&successfulExecutions)
	query.Where("status = ?", "failed").Count(&failedExecutions)

	// Get average duration
	query.Where("status = ?", "completed").Select("AVG(duration)").Scan(&avgDuration)

	stats := map[string]interface{}{
		"total_executions":      totalExecutions,
		"successful_executions": successfulExecutions,
		"failed_executions":     failedExecutions,
		"success_rate":          float64(successfulExecutions) / float64(totalExecutions) * 100,
		"avg_duration_ms":       avgDuration,
	}

	return stats, nil
}

// GetAgentExecutionStats retrieves execution statistics for a specific agent
func (s *RuntimeService) GetAgentExecutionStats(agentID string) (map[string]interface{}, error) {
	var totalExecutions int64
	var successfulExecutions int64
	var failedExecutions int64
	var avgDuration float64

	query := s.db.Model(&models.Execution{}).Where("agent_id = ?", agentID)

	// Get counts
	query.Count(&totalExecutions)
	query.Where("status = ?", "completed").Count(&successfulExecutions)
	query.Where("status = ?", "failed").Count(&failedExecutions)

	// Get average duration
	query.Where("status = ?", "completed").Select("AVG(duration)").Scan(&avgDuration)

	stats := map[string]interface{}{
		"total_executions":      totalExecutions,
		"successful_executions": successfulExecutions,
		"failed_executions":     failedExecutions,
		"success_rate":          float64(successfulExecutions) / float64(totalExecutions) * 100,
		"avg_duration_ms":       avgDuration,
	}

	return stats, nil
}

// StreamExecution streams execution updates via WebSocket
func (s *RuntimeService) StreamExecution(ctx context.Context, executionID string, updates chan<- map[string]interface{}) {
	defer close(updates)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var execution models.Execution
			if err := s.db.First(&execution, "id = ?", executionID).Error; err != nil {
				updates <- map[string]interface{}{
					"error": "Execution not found",
				}
				return
			}

			update := map[string]interface{}{
				"id":     execution.ID,
				"status": execution.Status,
				"output": execution.Output,
			}

			if execution.Status == "completed" || execution.Status == "failed" {
				updates <- update
				return
			}

			updates <- update
		}
	}
}

// GetActiveExecutions retrieves currently running executions
func (s *RuntimeService) GetActiveExecutions(userID string) ([]models.Execution, error) {
	var executions []models.Execution
	if err := s.db.Where("user_id = ? AND status = ?", userID, "running").
		Preload("Agent").Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// CleanupOldExecutions removes old execution records
func (s *RuntimeService) CleanupOldExecutions(daysOld int) error {
	cutoff := time.Now().AddDate(0, 0, -daysOld)
	return s.db.Where("created_at < ? AND status IN (?)", cutoff, []string{"completed", "failed", "cancelled"}).Delete(&models.Execution{}).Error
}
