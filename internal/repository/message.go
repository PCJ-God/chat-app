package repository

import (
	"chat-app/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(msg *model.Message) error
	GetMessages(userID uint, targetID *uint, roomID *uint, page, pageSize int) ([]model.Message, int64, error)
	GetUnreadMessages(userID uint) ([]model.Message, error)
	MarkAsRead(messageID, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *messageRepository) GetMessages(userID uint, targetID *uint, roomID *uint, page, pageSize int) ([]model.Message, int64, error) {
	var messages []model.Message
	var total int64

	query := r.db.Model(&model.Message{})

	switch {
	case roomID != nil:
		query = query.Where("room_id = ?", *roomID)
	case targetID != nil:
		query = query.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			userID, *targetID, *targetID, userID)
	default:
		query = query.Where("sender_id = ? OR receiver_id = ?", userID, userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	err := query.Preload("Sender").Preload("Receiver").Preload("Room").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error
	return messages, total, err
}

func (r *messageRepository) GetUnreadMessages(userID uint) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("receiver_id = ? AND status <> ?", userID, "read").
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *messageRepository) MarkAsRead(messageID, userID uint) error {
	return r.db.Model(&model.Message{}).
		Where("id = ? AND receiver_id = ?", messageID, userID).
		Update("status", "read").Error
}

func (r *messageRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("receiver_id = ? AND status <> ?", userID, "read").
		Count(&count).Error
	return count, err
}
