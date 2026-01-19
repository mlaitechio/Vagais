package services

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// AnalyticsService handles analytics and business intelligence
type AnalyticsService struct {
	BaseService
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *gorm.DB, cfg *config.Config) *AnalyticsService {
	return &AnalyticsService{
		BaseService: NewBaseService(db, cfg, "analytics"),
	}
}

// TrackEvent tracks an analytics event
func (s *AnalyticsService) TrackEvent(eventType, metric string, value float64, userID, agentID, orgID *string, metadata map[string]interface{}) error {
	analytics := &models.Analytics{
		Type:           eventType,
		Metric:         metric,
		Value:          value,
		UserID:         userID,
		AgentID:        agentID,
		OrganizationID: orgID,
		Metadata:       models.MapToJSON(metadata),
		Date:           time.Now(),
	}
	return s.db.Create(analytics).Error
}

// GetUsageStats retrieves usage statistics
func (s *AnalyticsService) GetUsageStats(userID *string, orgID *string, timeRange string) (map[string]interface{}, error) {
	var totalExecutions int64
	var totalAgents int64
	var totalUsers int64
	var totalRevenue float64

	query := s.db.Model(&models.Execution{})

	// Apply filters
	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	}

	// Apply time range filter
	switch timeRange {
	case "today":
		query = query.Where("created_at >= ?", time.Now().Truncate(24*time.Hour))
	case "week":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", time.Now().AddDate(-1, 0, 0))
	}

	query.Count(&totalExecutions)

	// Get agent count
	agentQuery := s.db.Model(&models.Agent{})
	if userID != nil {
		agentQuery = agentQuery.Where("creator_id = ?", userID)
	}
	if orgID != nil {
		agentQuery = agentQuery.Where("organization_id = ?", orgID)
	}
	agentQuery.Count(&totalAgents)

	// Get user count
	userQuery := s.db.Model(&models.User{})
	if orgID != nil {
		userQuery = userQuery.Where("organization_id = ?", orgID)
	}
	userQuery.Count(&totalUsers)

	// Get revenue (if payment service is available)
	if s.cfg.Payment.StripeSecretKey != "" {
		// TODO: Implement actual revenue calculation
		totalRevenue = 0.0
	}

	stats := map[string]interface{}{
		"total_executions": totalExecutions,
		"total_agents":     totalAgents,
		"total_users":      totalUsers,
		"total_revenue":    totalRevenue,
	}

	return stats, nil
}

// GetAgentUsageStats retrieves agent usage statistics
func (s *AnalyticsService) GetAgentUsageStats(agentID uuid.UUID, timeRange string) (map[string]interface{}, error) {
	var totalExecutions int64
	var successfulExecutions int64
	var failedExecutions int64
	var avgDuration float64
	var uniqueUsers int64

	query := s.db.Model(&models.Execution{}).Where("agent_id = ?", agentID.String())

	// Apply time range filter
	switch timeRange {
	case "today":
		query = query.Where("created_at >= ?", time.Now().Truncate(24*time.Hour))
	case "week":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", time.Now().AddDate(-1, 0, 0))
	}

	// Get counts
	query.Count(&totalExecutions)
	query.Where("status = ?", "completed").Count(&successfulExecutions)
	query.Where("status = ?", "failed").Count(&failedExecutions)

	// Get average duration
	query.Where("status = ?", "completed").Select("AVG(duration)").Scan(&avgDuration)

	// Get unique users
	query.Distinct("user_id").Count(&uniqueUsers)

	stats := map[string]interface{}{
		"total_executions":      totalExecutions,
		"successful_executions": successfulExecutions,
		"failed_executions":     failedExecutions,
		"success_rate":          float64(successfulExecutions) / float64(totalExecutions) * 100,
		"avg_duration_ms":       avgDuration,
		"unique_users":          uniqueUsers,
	}

	return stats, nil
}

// GetUserBehaviorAnalytics retrieves user behavior analytics
func (s *AnalyticsService) GetUserBehaviorAnalytics(userID uuid.UUID, timeRange string) (map[string]interface{}, error) {
	var totalExecutions int64
	var favoriteAgents []string
	var peakUsageHour int
	var avgSessionDuration float64

	query := s.db.Model(&models.Execution{}).Where("user_id = ?", userID.String())

	// Apply time range filter
	switch timeRange {
	case "today":
		query = query.Where("created_at >= ?", time.Now().Truncate(24*time.Hour))
	case "week":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", time.Now().AddDate(-1, 0, 0))
	}

	query.Count(&totalExecutions)

	// Get favorite agents (most used)
	var agentUsage []struct {
		AgentName string `json:"agent_name"`
		Count     int64  `json:"count"`
	}

	query.Select("agents.name as agent_name, COUNT(*) as count").
		Joins("JOIN agents ON executions.agent_id = agents.id").
		Group("agents.name").
		Order("count DESC").
		Limit(5).
		Find(&agentUsage)

	for _, usage := range agentUsage {
		favoriteAgents = append(favoriteAgents, usage.AgentName)
	}

	// TODO: Implement peak usage hour calculation
	peakUsageHour = 14 // Placeholder

	// TODO: Implement average session duration calculation
	avgSessionDuration = 1800.0 // 30 minutes in seconds

	stats := map[string]interface{}{
		"total_executions":       totalExecutions,
		"favorite_agents":        favoriteAgents,
		"peak_usage_hour":        peakUsageHour,
		"avg_session_duration_s": avgSessionDuration,
	}

	return stats, nil
}

// GetMarketplaceTrends retrieves marketplace trends
func (s *AnalyticsService) GetMarketplaceTrends(timeRange string) (map[string]interface{}, error) {
	var totalAgents int64
	var totalUsers int64
	var totalReviews int64
	var topCategories []string

	// Apply time range filter
	var timeFilter time.Time
	switch timeRange {
	case "today":
		timeFilter = time.Now().Truncate(24 * time.Hour)
	case "week":
		timeFilter = time.Now().AddDate(0, 0, -7)
	case "month":
		timeFilter = time.Now().AddDate(0, -1, 0)
	case "year":
		timeFilter = time.Now().AddDate(-1, 0, 0)
	default:
		timeFilter = time.Now().AddDate(0, -1, 0) // Default to last month
	}

	// Get counts
	s.db.Model(&models.Agent{}).Where("created_at >= ?", timeFilter).Count(&totalAgents)
	s.db.Model(&models.User{}).Where("created_at >= ?", timeFilter).Count(&totalUsers)
	s.db.Model(&models.Review{}).Where("created_at >= ?", timeFilter).Count(&totalReviews)

	// Get top categories
	var categoryStats []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	s.db.Model(&models.Agent{}).
		Select("category, COUNT(*) as count").
		Where("created_at >= ?", timeFilter).
		Group("category").
		Order("count DESC").
		Limit(5).
		Find(&categoryStats)

	for _, stat := range categoryStats {
		topCategories = append(topCategories, stat.Category)
	}

	trends := map[string]interface{}{
		"new_agents":     totalAgents,
		"new_users":      totalUsers,
		"new_reviews":    totalReviews,
		"top_categories": topCategories,
		"growth_rate":    s.calculateGrowthRate(timeRange),
	}

	return trends, nil
}

// calculateGrowthRate calculates growth rate for the given time range
func (s *AnalyticsService) calculateGrowthRate(timeRange string) float64 {
	var currentCount, previousCount int64
	var timeFilter, previousTimeFilter time.Time

	// Calculate time filters
	switch timeRange {
	case "today":
		timeFilter = time.Now().Truncate(24 * time.Hour)
		previousTimeFilter = timeFilter.AddDate(0, 0, -1)
	case "week":
		timeFilter = time.Now().AddDate(0, 0, -7)
		previousTimeFilter = timeFilter.AddDate(0, 0, -7)
	case "month":
		timeFilter = time.Now().AddDate(0, -1, 0)
		previousTimeFilter = timeFilter.AddDate(0, -1, 0)
	case "year":
		timeFilter = time.Now().AddDate(-1, 0, 0)
		previousTimeFilter = timeFilter.AddDate(-1, 0, 0)
	default:
		timeFilter = time.Now().AddDate(0, -1, 0)
		previousTimeFilter = timeFilter.AddDate(0, -1, 0)
	}

	// Get current period count
	s.db.Model(&models.Agent{}).Where("created_at >= ?", timeFilter).Count(&currentCount)

	// Get previous period count
	s.db.Model(&models.Agent{}).Where("created_at >= ? AND created_at < ?", previousTimeFilter, timeFilter).Count(&previousCount)

	// Calculate growth rate
	if previousCount == 0 {
		if currentCount > 0 {
			return 100.0 // 100% growth if there were no agents before
		}
		return 0.0
	}

	growthRate := float64(currentCount-previousCount) / float64(previousCount) * 100
	return growthRate
}

// GetRevenueAnalytics retrieves revenue analytics
func (s *AnalyticsService) GetRevenueAnalytics(timeRange string) (map[string]interface{}, error) {
	var totalRevenue float64
	var subscriptionCount, oneTimeSalesCount int64
	var timeFilter time.Time

	// Calculate time filter
	switch timeRange {
	case "today":
		timeFilter = time.Now().Truncate(24 * time.Hour)
	case "week":
		timeFilter = time.Now().AddDate(0, 0, -7)
	case "month":
		timeFilter = time.Now().AddDate(0, -1, 0)
	case "year":
		timeFilter = time.Now().AddDate(-1, 0, 0)
	default:
		timeFilter = time.Now().AddDate(0, -1, 0)
	}

	// Get total revenue from payments
	s.db.Model(&models.Payment{}).
		Where("created_at >= ? AND status = ?", timeFilter, "completed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue)

	// Get subscription count
	s.db.Model(&models.Subscription{}).
		Where("created_at >= ? AND status = ?", timeFilter, "active").
		Count(&subscriptionCount)

	// Get one-time sales count (payments that are not subscriptions)
	s.db.Model(&models.Payment{}).
		Where("created_at >= ? AND status = ?", timeFilter, "completed").
		Count(&oneTimeSalesCount)

	// Get payment methods distribution
	var paymentMethods []struct {
		Provider string  `json:"provider"`
		Amount   float64 `json:"amount"`
	}

	s.db.Model(&models.Payment{}).
		Select("provider, SUM(amount) as amount").
		Where("created_at >= ? AND status = ?", timeFilter, "completed").
		Group("provider").
		Find(&paymentMethods)

	paymentMethodsMap := make(map[string]float64)
	for _, pm := range paymentMethods {
		paymentMethodsMap[pm.Provider] = pm.Amount
	}

	// Get top products (agents with highest revenue)
	var topProducts []struct {
		Name   string  `json:"name"`
		Revenue float64 `json:"revenue"`
	}

	s.db.Model(&models.Agent{}).
		Select("name, price * downloads as revenue").
		Where("created_at >= ? AND is_public = ?", timeFilter, true).
		Order("revenue DESC").
		Limit(5).
		Find(&topProducts)

	topProductNames := make([]string, len(topProducts))
	for i, product := range topProducts {
		topProductNames[i] = product.Name
	}

	revenue := map[string]interface{}{
		"total_revenue":   totalRevenue,
		"revenue_growth":  s.calculateRevenueGrowth(timeRange),
		"top_products":    topProductNames,
		"payment_methods": paymentMethodsMap,
		"subscriptions":   subscriptionCount,
		"one_time_sales":  oneTimeSalesCount,
	}

	return revenue, nil
}

// calculateRevenueGrowth calculates revenue growth rate for the given time range
func (s *AnalyticsService) calculateRevenueGrowth(timeRange string) float64 {
	var currentRevenue, previousRevenue float64
	var timeFilter, previousTimeFilter time.Time

	// Calculate time filters
	switch timeRange {
	case "today":
		timeFilter = time.Now().Truncate(24 * time.Hour)
		previousTimeFilter = timeFilter.AddDate(0, 0, -1)
	case "week":
		timeFilter = time.Now().AddDate(0, 0, -7)
		previousTimeFilter = timeFilter.AddDate(0, 0, -7)
	case "month":
		timeFilter = time.Now().AddDate(0, -1, 0)
		previousTimeFilter = timeFilter.AddDate(0, -1, 0)
	case "year":
		timeFilter = time.Now().AddDate(-1, 0, 0)
		previousTimeFilter = timeFilter.AddDate(-1, 0, 0)
	default:
		timeFilter = time.Now().AddDate(0, -1, 0)
		previousTimeFilter = timeFilter.AddDate(0, -1, 0)
	}

	// Get current period revenue
	s.db.Model(&models.Payment{}).
		Where("created_at >= ? AND status = ?", timeFilter, "completed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&currentRevenue)

	// Get previous period revenue
	s.db.Model(&models.Payment{}).
		Where("created_at >= ? AND created_at < ? AND status = ?", previousTimeFilter, timeFilter, "completed").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&previousRevenue)

	// Calculate growth rate
	if previousRevenue == 0 {
		if currentRevenue > 0 {
			return 100.0 // 100% growth if there was no revenue before
		}
		return 0.0
	}

	growthRate := (currentRevenue - previousRevenue) / previousRevenue * 100
	return growthRate
}

// GetDeveloperMetrics retrieves developer performance metrics
func (s *AnalyticsService) GetDeveloperMetrics(userID uuid.UUID, timeRange string) (map[string]interface{}, error) {
	var agentsCreated int64
	var totalExecutions int64
	var avgRating float64
	var totalReviews int64

	query := s.db.Model(&models.Agent{}).Where("creator_id = ?", userID.String())

	// Apply time range filter
	switch timeRange {
	case "today":
		query = query.Where("created_at >= ?", time.Now().Truncate(24*time.Hour))
	case "week":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, 0, -7))
	case "month":
		query = query.Where("created_at >= ?", time.Now().AddDate(0, -1, 0))
	case "year":
		query = query.Where("created_at >= ?", time.Now().AddDate(-1, 0, 0))
	}

	query.Count(&agentsCreated)

	// Get total executions for user's agents
	s.db.Model(&models.Execution{}).
		Joins("JOIN agents ON executions.agent_id = agents.id").
		Where("agents.creator_id = ?", userID.String()).
		Count(&totalExecutions)

	// Get average rating
	s.db.Model(&models.Agent{}).Where("creator_id = ?", userID.String()).Select("AVG(rating)").Scan(&avgRating)

	// Get total reviews
	s.db.Model(&models.Review{}).
		Joins("JOIN agents ON reviews.agent_id = agents.id").
		Where("agents.creator_id = ?", userID.String()).
		Count(&totalReviews)

	metrics := map[string]interface{}{
		"agents_created":    agentsCreated,
		"total_executions":  totalExecutions,
		"avg_rating":        avgRating,
		"total_reviews":     totalReviews,
		"success_rate":      95.5, // Placeholder
		"user_satisfaction": 4.2,  // Placeholder
	}

	return metrics, nil
}

// GenerateCustomReport generates a custom analytics report
func (s *AnalyticsService) GenerateCustomReport(reportType string, filters map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Implement custom report generation
	// This is a placeholder implementation

	report := map[string]interface{}{
		"report_type":  reportType,
		"filters":      filters,
		"data":         []interface{}{},
		"summary":      map[string]interface{}{},
		"generated_at": time.Now(),
	}

	return report, nil
}
