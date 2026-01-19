package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// BillingService handles billing and payment processing
type BillingService struct {
	BaseService
}

// NewBillingService creates a new billing service
func NewBillingService(db *gorm.DB, cfg *config.Config) *BillingService {
	return &BillingService{
		BaseService: NewBaseService(db, cfg, "billing"),
	}
}

// CreateSubscriptionRequest represents subscription creation request
type CreateSubscriptionRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	OrganizationID *string `json:"organization_id"`
	Plan           string  `json:"plan" binding:"required"`
	Provider       string  `json:"provider" binding:"required"` // stripe, paypal, upi
	Amount         float64 `json:"amount" binding:"required"`
	Currency       string  `json:"currency"`
	Interval       string  `json:"interval"` // monthly, yearly
}

// CreatePaymentRequest represents payment creation request
type CreatePaymentRequest struct {
	UserID         string                 `json:"user_id" binding:"required"`
	OrganizationID *string                `json:"organization_id"`
	Amount         float64                `json:"amount" binding:"required"`
	Currency       string                 `json:"currency"`
	Provider       string                 `json:"provider" binding:"required"`
	Description    string                 `json:"description"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// CreateSubscription creates a new subscription
func (s *BillingService) CreateSubscription(req *CreateSubscriptionRequest) (*models.Subscription, error) {
	subscription := &models.Subscription{
		UserID:             req.UserID,
		OrganizationID:     req.OrganizationID,
		Plan:               req.Plan,
		Status:             "active",
		Provider:           req.Provider,
		ProviderID:         "", // TODO: Get from payment provider
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0), // Default to 1 month
		CancelAtPeriodEnd:  false,
		Amount:             req.Amount,
		Currency:           req.Currency,
	}

	if err := s.db.Create(subscription).Error; err != nil {
		return nil, err
	}

	return subscription, nil
}

// GetSubscription retrieves a subscription by ID
func (s *BillingService) GetSubscription(id string) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := s.db.Preload("User").Preload("Organization").First(&subscription, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

// ListSubscriptions retrieves subscriptions for a user/organization
func (s *BillingService) ListSubscriptions(userID string, orgID *string) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	query := s.db.Model(&models.Subscription{}).Preload("User").Preload("Organization")

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&subscriptions).Error; err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// CancelSubscription cancels a subscription
func (s *BillingService) CancelSubscription(id string, userID string) error {
	var subscription models.Subscription
	if err := s.db.First(&subscription, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user owns this subscription
	if subscription.UserID != userID {
		return errors.New("unauthorized to cancel this subscription")
	}

	return s.db.Model(&subscription).Update("status", "cancelled").Error
}

// ReactivateSubscription reactivates a cancelled subscription
func (s *BillingService) ReactivateSubscription(id string, userID string) error {
	var subscription models.Subscription
	if err := s.db.First(&subscription, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user owns this subscription
	if subscription.UserID != userID {
		return errors.New("unauthorized to reactivate this subscription")
	}

	return s.db.Model(&subscription).Update("status", "active").Error
}

// CreatePayment creates a new payment
func (s *BillingService) CreatePayment(req *CreatePaymentRequest) (*models.Payment, error) {
	payment := &models.Payment{
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "pending",
		Provider:       req.Provider,
		ProviderID:     "", // TODO: Get from payment provider
		Description:    req.Description,
		Metadata:       models.MapToJSON(req.Metadata),
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	// Process payment
	go s.processPayment(payment)

	return payment, nil
}

// processPayment processes a payment asynchronously
func (s *BillingService) processPayment(payment *models.Payment) {
	// TODO: Implement actual payment processing
	// For now, we'll simulate payment processing
	time.Sleep(2 * time.Second)

	// Simulate successful payment
	payment.Status = "completed"
	payment.ProviderID = "pay_" + payment.ID[:8]

	if err := s.db.Save(payment).Error; err != nil {
		// Log error but don't fail
		fmt.Printf("Error saving payment: %v\n", err)
	}
}

// GetPayment retrieves a payment by ID
func (s *BillingService) GetPayment(id string) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.Preload("User").Preload("Organization").First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// ListPayments retrieves payments for a user/organization
func (s *BillingService) ListPayments(userID string, orgID *string, page, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := s.db.Model(&models.Payment{}).Preload("User").Preload("Organization")

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// RefundPayment refunds a payment
func (s *BillingService) RefundPayment(id string, userID string, reason string) error {
	var payment models.Payment
	if err := s.db.First(&payment, "id = ?", id).Error; err != nil {
		return err
	}

	// Check if user owns this payment
	if payment.UserID != userID {
		return errors.New("unauthorized to refund this payment")
	}

	// Check if payment is completed
	if payment.Status != "completed" {
		return errors.New("can only refund completed payments")
	}

	// Create refund payment
	refund := &models.Payment{
		UserID:         payment.UserID,
		OrganizationID: payment.OrganizationID,
		Amount:         -payment.Amount, // Negative amount for refund
		Currency:       payment.Currency,
		Status:         "pending",
		Provider:       payment.Provider,
		Description:    "Refund: " + payment.Description,
		Metadata:       models.MapToJSON(map[string]interface{}{"refund_reason": reason, "original_payment_id": payment.ID}),
	}

	if err := s.db.Create(refund).Error; err != nil {
		return err
	}

	// Process refund
	go s.processRefund(refund)

	return nil
}

// processRefund processes a refund asynchronously
func (s *BillingService) processRefund(refund *models.Payment) {
	// TODO: Implement actual refund processing
	time.Sleep(1 * time.Second)

	refund.Status = "completed"
	refund.ProviderID = "ref_" + refund.ID[:8]

	if err := s.db.Save(refund).Error; err != nil {
		fmt.Printf("Error saving refund: %v\n", err)
	}
}

// GetBillingStats retrieves billing statistics
func (s *BillingService) GetBillingStats(userID string, orgID *string) (map[string]interface{}, error) {
	var totalPayments int64
	var totalAmount float64
	var activeSubscriptions int64
	var cancelledSubscriptions int64

	// Get payment stats
	paymentQuery := s.db.Model(&models.Payment{})
	if orgID != nil {
		paymentQuery = paymentQuery.Where("organization_id = ?", orgID)
	} else {
		paymentQuery = paymentQuery.Where("user_id = ?", userID)
	}
	paymentQuery.Where("status = ?", "completed").Count(&totalPayments)
	paymentQuery.Where("status = ?", "completed").Select("SUM(amount)").Scan(&totalAmount)

	// Get subscription stats
	subQuery := s.db.Model(&models.Subscription{})
	if orgID != nil {
		subQuery = subQuery.Where("organization_id = ?", orgID)
	} else {
		subQuery = subQuery.Where("user_id = ?", userID)
	}
	subQuery.Where("status = ?", "active").Count(&activeSubscriptions)
	subQuery.Where("status = ?", "cancelled").Count(&cancelledSubscriptions)

	stats := map[string]interface{}{
		"total_payments":          totalPayments,
		"total_amount":            totalAmount,
		"active_subscriptions":    activeSubscriptions,
		"cancelled_subscriptions": cancelledSubscriptions,
	}

	return stats, nil
}

// GetAvailablePlans retrieves available subscription plans
func (s *BillingService) GetAvailablePlans() ([]map[string]interface{}, error) {
	var dbPlans []models.BillingPlan
	if err := s.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&dbPlans).Error; err != nil {
		return nil, err
	}

	plans := make([]map[string]interface{}, len(dbPlans))
	for i, dbPlan := range dbPlans {
		var features []string
		if err := json.Unmarshal(dbPlan.Features, &features); err != nil {
			features = []string{}
		}

		plans[i] = map[string]interface{}{
			"id":             dbPlan.Slug,
			"name":           dbPlan.Name,
			"price":          dbPlan.Price,
			"currency":       dbPlan.Currency,
			"interval":       dbPlan.Interval,
			"features":       features,
			"max_agents":     dbPlan.MaxAgents,
			"max_executions": dbPlan.MaxExecutions,
			"description":    dbPlan.Description,
		}
	}

	return plans, nil
}

// ValidateSubscription validates if a user has an active subscription
func (s *BillingService) ValidateSubscription(userID string, orgID *string) (bool, string, error) {
	var subscription models.Subscription
	query := s.db.Model(&models.Subscription{})

	if orgID != nil {
		query = query.Where("organization_id = ? AND status = ?", orgID, "active")
	} else {
		query = query.Where("user_id = ? AND status = ?", userID, "active")
	}

	if err := query.First(&subscription).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, "no_subscription", nil
		}
		return false, "", err
	}

	// Check if subscription is expired
	if subscription.CurrentPeriodEnd.Before(time.Now()) {
		return false, "expired", nil
	}

	return true, subscription.Plan, nil
}
