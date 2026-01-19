package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// PaymentService handles payment processing
type PaymentService struct {
	BaseService
}

// NewPaymentService creates a new payment service
func NewPaymentService(db *gorm.DB, cfg *config.Config) *PaymentService {
	return &PaymentService{
		BaseService: NewBaseService(db, cfg, "payment"),
	}
}

// ProcessPaymentRequest represents payment processing request
type ProcessPaymentRequest struct {
	UserID         string                 `json:"user_id" binding:"required"`
	OrganizationID *string                `json:"organization_id"`
	Amount         float64                `json:"amount" binding:"required"`
	Currency       string                 `json:"currency"`
	Provider       string                 `json:"provider" binding:"required"` // stripe, paypal, upi
	PaymentMethod  string                 `json:"payment_method"`
	Description    string                 `json:"description"`
	Metadata       map[string]interface{} `json:"metadata"`
}

// CreatePaymentIntentRequest represents payment intent creation
type CreatePaymentIntentRequest struct {
	UserID         string  `json:"user_id" binding:"required"`
	OrganizationID *string `json:"organization_id"`
	Amount         float64 `json:"amount" binding:"required"`
	Currency       string  `json:"currency"`
	Provider       string  `json:"provider" binding:"required"`
	Description    string  `json:"description"`
}

// ProcessPayment processes a payment
func (s *PaymentService) ProcessPayment(req *ProcessPaymentRequest) (*models.Payment, error) {
	payment := &models.Payment{
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         "pending",
		Provider:       req.Provider,
		ProviderID:     "",
		Description:    req.Description,
		Metadata:       models.MapToJSON(req.Metadata),
	}

	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	// Process payment based on provider
	switch req.Provider {
	case "stripe":
		return s.processStripePayment(payment, req)
	case "paypal":
		return s.processPayPalPayment(payment, req)
	case "upi":
		return s.processUPIPayment(payment, req)
	default:
		return nil, errors.New("unsupported payment provider")
	}
}

// processStripePayment processes a Stripe payment
func (s *PaymentService) processStripePayment(payment *models.Payment, req *ProcessPaymentRequest) (*models.Payment, error) {
	// TODO: Implement actual Stripe payment processing
	// For now, simulate payment processing
	time.Sleep(2 * time.Second)

	payment.Status = "completed"
	payment.ProviderID = "pi_" + payment.ID[:8]

	if err := s.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

// processPayPalPayment processes a PayPal payment
func (s *PaymentService) processPayPalPayment(payment *models.Payment, req *ProcessPaymentRequest) (*models.Payment, error) {
	// TODO: Implement actual PayPal payment processing
	time.Sleep(2 * time.Second)

	payment.Status = "completed"
	payment.ProviderID = "pay_" + payment.ID[:8]

	if err := s.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

// processUPIPayment processes a UPI payment
func (s *PaymentService) processUPIPayment(payment *models.Payment, req *ProcessPaymentRequest) (*models.Payment, error) {
	// TODO: Implement actual UPI payment processing
	time.Sleep(2 * time.Second)

	payment.Status = "completed"
	payment.ProviderID = "upi_" + payment.ID[:8]

	if err := s.db.Save(payment).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

// CreatePaymentIntent creates a payment intent
func (s *PaymentService) CreatePaymentIntent(req *CreatePaymentIntentRequest) (map[string]interface{}, error) {
	// TODO: Implement actual payment intent creation
	// For now, return a placeholder
	intent := map[string]interface{}{
		"id":            "pi_" + uuid.New().String()[:8],
		"amount":        req.Amount,
		"currency":      req.Currency,
		"status":        "requires_payment_method",
		"client_secret": "pi_" + uuid.New().String() + "_secret_" + uuid.New().String()[:8],
	}

	return intent, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(id string) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.Preload("User").Preload("Organization").First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// ListPayments retrieves payments with filtering and pagination
func (s *PaymentService) ListPayments(userID string, orgID *string, status string, page, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := s.db.Model(&models.Payment{}).Preload("User").Preload("Organization")

	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	} else {
		query = query.Where("user_id = ?", userID)
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
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// RefundPayment refunds a payment
func (s *PaymentService) RefundPayment(id string, userID string, reason string) (*models.Payment, error) {
	var payment models.Payment
	if err := s.db.First(&payment, id).Error; err != nil {
		return nil, err
	}

	// Check if user owns this payment
	if payment.UserID != userID {
		return nil, errors.New("unauthorized to refund this payment")
	}

	// Check if payment is completed
	if payment.Status != "completed" {
		return nil, errors.New("can only refund completed payments")
	}

	// Create refund
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
		return nil, err
	}

	// Process refund
	go s.processRefund(refund)

	return refund, nil
}

// processRefund processes a refund asynchronously
func (s *PaymentService) processRefund(refund *models.Payment) {
	// TODO: Implement actual refund processing
	time.Sleep(1 * time.Second)

	refund.Status = "completed"
	refund.ProviderID = "ref_" + refund.ID[:8]

	if err := s.db.Save(refund).Error; err != nil {
		fmt.Printf("Error saving refund: %v\n", err)
	}
}

// GetPaymentMethods retrieves available payment methods
func (s *PaymentService) GetPaymentMethods() ([]map[string]interface{}, error) {
	methods := []map[string]interface{}{
		{
			"id":          "stripe",
			"name":        "Stripe",
			"description": "Credit/Debit Cards via Stripe",
			"enabled":     true,
			"currencies":  []string{"USD", "EUR", "GBP", "CAD", "AUD"},
		},
		{
			"id":          "paypal",
			"name":        "PayPal",
			"description": "PayPal Express Checkout",
			"enabled":     true,
			"currencies":  []string{"USD", "EUR", "GBP", "CAD", "AUD"},
		},
		{
			"id":          "upi",
			"name":        "UPI",
			"description": "Unified Payments Interface (India)",
			"enabled":     true,
			"currencies":  []string{"INR"},
		},
	}

	return methods, nil
}

// GetPaymentStats retrieves payment statistics
func (s *PaymentService) GetPaymentStats(userID string, orgID *string) (map[string]interface{}, error) {
	var totalPayments int64
	var totalAmount float64
	var successfulPayments int64
	var failedPayments int64
	var pendingPayments int64

	query := s.db.Model(&models.Payment{})
	if orgID != nil {
		query = query.Where("organization_id = ?", orgID)
	} else {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&totalPayments)
	query.Where("status = ?", "completed").Count(&successfulPayments)
	query.Where("status = ?", "failed").Count(&failedPayments)
	query.Where("status = ?", "pending").Count(&pendingPayments)
	query.Where("status = ?", "completed").Select("SUM(amount)").Scan(&totalAmount)

	stats := map[string]interface{}{
		"total_payments":      totalPayments,
		"successful_payments": successfulPayments,
		"failed_payments":     failedPayments,
		"pending_payments":    pendingPayments,
		"total_amount":        totalAmount,
		"success_rate":        float64(successfulPayments) / float64(totalPayments) * 100,
	}

	return stats, nil
}

// ValidatePayment validates a payment
func (s *PaymentService) ValidatePayment(paymentID string) (bool, error) {
	var payment models.Payment
	if err := s.db.First(&payment, paymentID).Error; err != nil {
		return false, err
	}

	return payment.Status == "completed", nil
}
