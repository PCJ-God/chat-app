package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"chat-app/pkg/websocket"
)

type MessageService interface {
	SendMessage(msg *model.Message) error
	GetMessages(userID uint, targetID *uint, roomID *uint, page, pageSize int) ([]model.Message, int64, error)
	MarkAsRead(messageID uint, userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
	hub         *websocket.Hub
}

func NewMessageService(messageRepo repository.MessageRepository, hub *websocket.Hub) MessageService {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
		hub:         hub,
	}
}

func (s *MessageServiceImpl) SendMessage(msg *model.Message) error {
	if msg.Type == "" {
		msg.Type = "text"
	}
	if msg.Status == "" {
		msg.Status = "sent"
	}

	// 保存消息到数据库
	if err := s.messageRepo.Create(msg); err != nil {
		return err
	}

	// 推送消息到WebSocket
	s.broadcastMessage(msg)

	return nil
}

func (s *MessageServiceImpl) broadcastMessage(msg *model.Message) {
	messageData := map[string]interface{}{
		"type":       "message",
		"id":         msg.ID,
		"sender_id":  msg.SenderID,
		"content":    msg.Content,
		"created_at": msg.CreatedAt,
	}

	// 私聊消息
	if msg.ReceiverID != nil {
		messageData["receiver_id"] = *msg.ReceiverID
		s.hub.SendToUser(*msg.ReceiverID, messageData)
	}

	// 群聊消息
	if msg.RoomID != nil {
		messageData["room_id"] = *msg.RoomID
		s.hub.BroadcastToRoom(*msg.RoomID, messageData)
	}
}

func (s *MessageServiceImpl) GetMessages(userID uint, targetID *uint, roomID *uint, page, pageSize int) ([]model.Message, int64, error) {
	return s.messageRepo.GetMessages(userID, targetID, roomID, page, pageSize)
}

func (s *MessageServiceImpl) MarkAsRead(messageID uint, userID uint) error {
	return s.messageRepo.MarkAsRead(messageID, userID)
}

func (s *MessageServiceImpl) GetUnreadCount(userID uint) (int64, error) {
	return s.messageRepo.GetUnreadCount(userID)
}
