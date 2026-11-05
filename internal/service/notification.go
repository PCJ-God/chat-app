package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
)

type NotificationService interface {
	SendNotification(userID uint, notification *model.Notification) error
	GetNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error)
	MarkAsRead(notificationID uint, userID uint) error
	PushOfflineMessages(userID uint) error
}

type NotificationServiceImpl struct {
	notificationRepo repository.NotificationRepository
	messageRepo      repository.MessageRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository, messageRepo repository.MessageRepository) NotificationService {
	return &NotificationServiceImpl{
		notificationRepo: notificationRepo,
		messageRepo:      messageRepo,
	}
}

func (s *NotificationServiceImpl) SendNotification(userID uint, notification *model.Notification) error {
	if notification.UserID == 0 {
		notification.UserID = userID
	}
	return s.notificationRepo.Create(notification)
}

func (s *NotificationServiceImpl) GetNotifications(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	return s.notificationRepo.GetByUserID(userID, page, pageSize)
}

func (s *NotificationServiceImpl) MarkAsRead(notificationID uint, userID uint) error {
	return s.notificationRepo.MarkAsRead(notificationID, userID)
}

func (s *NotificationServiceImpl) PushOfflineMessages(userID uint) error {
	// 获取离线消息
	messages, err := s.messageRepo.GetUnreadMessages(userID)
	if err != nil {
		return err
	}

	// 发送推送通知
	for _, msg := range messages {
		notification := &model.Notification{
			Type:    "message",
			UserID:  userID,
			Title:   "新消息",
			Content: msg.Content,
			Data:    msg,
		}
		if err := s.SendNotification(userID, notification); err != nil {
			return err
		}
	}

	return nil
}
