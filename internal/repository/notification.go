package repository

import (
	"chat-app/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(n *model.Notification) error
	GetByUserID(userID uint, page, pageSize int) ([]model.Notification, int64, error)
	MarkAsRead(id, userID uint) error
	CountUnread(userID uint) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *notificationRepository) GetByUserID(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var list []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

func (r *notificationRepository) MarkAsRead(id, userID uint) error {
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}
