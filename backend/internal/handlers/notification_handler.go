package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/services"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	*BaseHandler
	notificationService *services.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(db *gorm.DB, cfg *config.Config) *NotificationHandler {
	return &NotificationHandler{
		BaseHandler:         NewBaseHandler(db, cfg),
		notificationService: services.NotificationServiceInstance,
	}
}

// SendNotification sends a notification
func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req services.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	notification, err := h.notificationService.SendNotification(&req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendCreated(c, notification)
}

// GetNotification gets a notification by ID
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	notificationID := c.Param("id")

	notification, err := h.notificationService.GetNotification(notificationID)
	if err != nil {
		h.sendError(c, http.StatusNotFound, "Notification not found")
		return
	}

	h.sendSuccess(c, notification)
}

// ListNotifications lists notifications for a user
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	isReadStr := c.Query("is_read")

	var isRead *bool
	if isReadStr != "" {
		if parsed, err := strconv.ParseBool(isReadStr); err == nil {
			isRead = &parsed
		}
	}

	notifications, total, err := h.notificationService.ListNotifications(userID, page, limit, isRead)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"limit":         limit,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.notificationService.MarkAsRead(notificationID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Notification marked as read"})
}

// MarkAllAsRead marks all notifications as read for a user
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "All notifications marked as read"})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationID := c.Param("id")

	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	if err := h.notificationService.DeleteNotification(notificationID, userID); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Notification deleted successfully"})
}

// GetUnreadCount gets unread notification count for a user
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"unread_count": count})
}

// SendBulkNotification sends notifications to multiple users
func (h *NotificationHandler) SendBulkNotification(c *gin.Context) {
	var req struct {
		UserIDs []string               `json:"user_ids" binding:"required"`
		Type    string                 `json:"type" binding:"required"`
		Title   string                 `json:"title" binding:"required"`
		Message string                 `json:"message" binding:"required"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	notificationReq := &services.CreateNotificationRequest{
		Type:     req.Type,
		Title:    req.Title,
		Message:  req.Message,
		Metadata: req.Data,
	}

	if err := h.notificationService.SendBulkNotification(req.UserIDs, notificationReq); err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, gin.H{"message": "Bulk notifications sent successfully"})
}

// GetNotificationStats gets notification statistics for a user
func (h *NotificationHandler) GetNotificationStats(c *gin.Context) {
	userID, exists := h.getCurrentUserID(c)
	if !exists {
		h.sendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	stats, err := h.notificationService.GetNotificationStats(userID)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.sendSuccess(c, stats)
}
