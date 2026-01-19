package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// BillingHandler handles billing-related requests
type BillingHandler struct {
	*BaseHandler
	billingService *services.BillingService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(db *gorm.DB, cfg *config.Config) *BillingHandler {
	return &BillingHandler{
		BaseHandler:    NewBaseHandler(db, cfg),
		billingService: services.BillingServiceInstance,
	}
}

// CreateSubscription creates a new subscription
func (h *BillingHandler) CreateSubscription(c *gin.Context) {
	var req services.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	subscription, err := h.billingService.CreateSubscription(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, subscription)
}

// GetSubscription gets a subscription by ID
func (h *BillingHandler) GetSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")

	subscription, err := h.billingService.GetSubscription(subscriptionID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Subscription not found")
		return
	}

	h.sendSuccess(c, subscription)
}

// ListSubscriptions lists subscriptions for a user/organization
func (h *BillingHandler) ListSubscriptions(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")

	subscriptions, err := h.billingService.ListSubscriptions(userID, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, subscriptions)
}

// CancelSubscription cancels a subscription
func (h *BillingHandler) CancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.billingService.CancelSubscription(subscriptionID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Subscription cancelled successfully"})
}

// ReactivateSubscription reactivates a cancelled subscription
func (h *BillingHandler) ReactivateSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.billingService.ReactivateSubscription(subscriptionID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Subscription reactivated successfully"})
}

// CreatePayment creates a new payment
func (h *BillingHandler) CreatePayment(c *gin.Context) {
	var req services.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment, err := h.billingService.CreatePayment(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, payment)
}

// GetPayment gets a payment by ID
func (h *BillingHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.billingService.GetPayment(paymentID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Payment not found")
		return
	}

	h.sendSuccess(c, payment)
}

// ListPayments lists payments for a user/organization
func (h *BillingHandler) ListPayments(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, total, err := h.billingService.ListPayments(userID, &orgID, page, limit)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"payments": payments,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// RefundPayment refunds a payment
func (h *BillingHandler) RefundPayment(c *gin.Context) {
	paymentID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
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

	if err := h.billingService.RefundPayment(paymentID, userID, req.Reason); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Payment refunded successfully"})
}

// GetBillingStats gets billing statistics
func (h *BillingHandler) GetBillingStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")

	stats, err := h.billingService.GetBillingStats(userID, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// GetAvailablePlans gets available subscription plans
func (h *BillingHandler) GetAvailablePlans(c *gin.Context) {
	plans, err := h.billingService.GetAvailablePlans()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, plans)
}

// ValidateSubscription validates a subscription
func (h *BillingHandler) ValidateSubscription(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")

	valid, status, err := h.billingService.ValidateSubscription(userID, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"valid":  valid,
		"status": status,
	})
}
