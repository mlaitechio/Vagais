package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// PaymentHandler handles payment-related requests
type PaymentHandler struct {
	*BaseHandler
	paymentService *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(db *gorm.DB, cfg *config.Config) *PaymentHandler {
	return &PaymentHandler{
		BaseHandler:    NewBaseHandler(db, cfg),
		paymentService: services.PaymentServiceInstance,
	}
}

// ProcessPayment processes a payment
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req services.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment, err := h.paymentService.ProcessPayment(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, payment)
}

// CreatePaymentIntent creates a payment intent
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req services.CreatePaymentIntentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	intent, err := h.paymentService.CreatePaymentIntent(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, intent)
}

// GetPayment gets a payment by ID
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetPayment(paymentID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Payment not found")
		return
	}

	h.sendSuccess(c, payment)
}

// ListPayments lists payments for a user/organization
func (h *PaymentHandler) ListPayments(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	payments, total, err := h.paymentService.ListPayments(userID, &orgID, status, page, limit)
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
func (h *PaymentHandler) RefundPayment(c *gin.Context) {
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

	refund, err := h.paymentService.RefundPayment(paymentID, userID, req.Reason)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, refund)
}

// GetPaymentMethods gets available payment methods
func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.paymentService.GetPaymentMethods()
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, methods)
}

// GetPaymentStats gets payment statistics
func (h *PaymentHandler) GetPaymentStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orgID := c.Query("org_id")

	stats, err := h.paymentService.GetPaymentStats(userID, &orgID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}

// ValidatePayment validates a payment
func (h *PaymentHandler) ValidatePayment(c *gin.Context) {
	paymentID := c.Param("id")

	valid, err := h.paymentService.ValidatePayment(paymentID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"valid": valid})
}
