package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/mlaitechio/vagais/internal/config"
	"github.com/mlaitechio/vagais/internal/models"
)

// NotificationService handles notifications
type NotificationService struct {
	BaseService
}

// NewNotificationService creates a new notification service
func NewNotificationService(db *gorm.DB, cfg *config.Config) *NotificationService {
	return &NotificationService{
		BaseService: NewBaseService(db, cfg, "notification"),
	}
}

// CreateNotificationRequest represents notification creation request
type CreateNotificationRequest struct {
	UserID   string                 `json:"user_id" binding:"required"`
	Type     string                 `json:"type" binding:"required"` // email, sms, in_app
	Title    string                 `json:"title" binding:"required"`
	Message  string                 `json:"message" binding:"required"`
	Priority string                 `json:"priority"` // low, normal, high, urgent
	Metadata map[string]interface{} `json:"metadata"`
	IsRead   bool                   `json:"is_read"`
}

// SendNotification sends a notification
func (s *NotificationService) SendNotification(req *CreateNotificationRequest) (*models.Notification, error) {
	notification := &models.Notification{
		UserID:   req.UserID,
		Type:     req.Type,
		Title:    req.Title,
		Message:  req.Message,
		Priority: req.Priority,
		Metadata: models.MapToJSON(req.Metadata),
		Status:   "unread",
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	// Send notification based on type
	switch req.Type {
	case "email":
		go s.sendEmailNotification(notification)
	case "sms":
		go s.sendSMSNotification(notification)
	case "in_app":
		// In-app notifications are already created in database
		notification.Status = "sent"
		s.db.Save(notification)
	}

	return notification, nil
}

// sendEmailNotification sends an email notification
func (s *NotificationService) sendEmailNotification(notification *models.Notification) {
	// TODO: Implement actual email sending
	// For now, we'll just mark it as sent
	time.Sleep(1 * time.Second) // Simulate email sending

	notification.Status = "sent"
	s.db.Save(notification)
}

// sendSMSNotification sends an SMS notification
func (s *NotificationService) sendSMSNotification(notification *models.Notification) {
	// TODO: Implement actual SMS sending
	// For now, we'll just mark it as sent
	time.Sleep(500 * time.Millisecond) // Simulate SMS sending

	notification.Status = "sent"
	s.db.Save(notification)
}

// GetNotification retrieves a notification by ID
func (s *NotificationService) GetNotification(id string) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.Preload("User").First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

// ListNotifications retrieves notifications for a user
func (s *NotificationService) ListNotifications(userID string, page, limit int, isRead *bool) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if isRead != nil {
		if *isRead {
			query = query.Where("status = ?", "read")
		} else {
			query = query.Where("status = ?", "unread")
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(id string, userID string) error {
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return err
	}

	// Check if user owns this notification
	if notification.UserID != userID {
		return errors.New("unauthorized to mark this notification as read")
	}

	return s.db.Model(&notification).Update("status", "read").Error
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(userID string) error {
	return s.db.Model(&models.Notification{}).Where("user_id = ? AND status = ?", userID, "unread").Update("status", "read").Error
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(id string, userID string) error {
	var notification models.Notification
	if err := s.db.First(&notification, id).Error; err != nil {
		return err
	}

	// Check if user owns this notification
	if notification.UserID != userID {
		return errors.New("unauthorized to delete this notification")
	}

	return s.db.Delete(&notification).Error
}

// GetUnreadCount retrieves unread notification count for a user
func (s *NotificationService) GetUnreadCount(userID string) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Notification{}).Where("user_id = ? AND status = ?", userID, "unread").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SendBulkNotification sends notifications to multiple users
func (s *NotificationService) SendBulkNotification(userIDs []string, req *CreateNotificationRequest) error {
	for _, userID := range userIDs {
		bulkReq := *req
		bulkReq.UserID = userID
		if _, err := s.SendNotification(&bulkReq); err != nil {
			// Log error but continue with other users
			fmt.Printf("Failed to send notification to user %s: %v\n", userID, err)
		}
	}
	return nil
}

// GetNotificationStats retrieves notification statistics
func (s *NotificationService) GetNotificationStats(userID string) (map[string]interface{}, error) {
	var totalNotifications int64
	var unreadNotifications int64
	var emailNotifications int64
	var smsNotifications int64
	var inAppNotifications int64

	s.db.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&totalNotifications)
	s.db.Model(&models.Notification{}).Where("user_id = ? AND status = ?", userID, "unread").Count(&unreadNotifications)
	s.db.Model(&models.Notification{}).Where("user_id = ? AND type = ?", userID, "email").Count(&emailNotifications)
	s.db.Model(&models.Notification{}).Where("user_id = ? AND type = ?", userID, "sms").Count(&smsNotifications)
	s.db.Model(&models.Notification{}).Where("user_id = ? AND type = ?", userID, "in_app").Count(&inAppNotifications)

	stats := map[string]interface{}{
		"total_notifications":  totalNotifications,
		"unread_notifications": unreadNotifications,
		"email_notifications":  emailNotifications,
		"sms_notifications":    smsNotifications,
		"in_app_notifications": inAppNotifications,
	}

	return stats, nil
}

// CleanupOldNotifications removes old notifications
func (s *NotificationService) CleanupOldNotifications(daysOld int) error {
	cutoff := time.Now().AddDate(0, 0, -daysOld)
	return s.db.Where("created_at < ? AND status = ?", cutoff, "read").Delete(&models.Notification{}).Error
}
