package service

import (
	"errors"
	"time"

	"chat-app/internal/model"
	"chat-app/internal/repository"
)

type RoomService interface {
	CreateRoom(room *model.Room, creatorID uint) error
	JoinRoom(roomID, userID uint) error
	LeaveRoom(roomID, userID uint) error
	GetRoom(roomID uint) (*model.Room, error)
	GetUserRooms(userID uint) ([]model.Room, error)
	SendRoomMessage(roomID uint, senderID uint, content string) error
}

type RoomServiceImpl struct {
	roomRepo       repository.RoomRepository
	roomMemberRepo repository.RoomMemberRepository
	messageService MessageService
}

func NewRoomService(roomRepo repository.RoomRepository, roomMemberRepo repository.RoomMemberRepository, messageService MessageService) RoomService {
	return &RoomServiceImpl{
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
		messageService: messageService,
	}
}

func (s *RoomServiceImpl) CreateRoom(room *model.Room, creatorID uint) error {
	room.CreatedBy = creatorID
	if room.Type == "" {
		room.Type = "public"
	}

	if err := s.roomRepo.Create(room); err != nil {
		return err
	}

	// 创建者自动加入
	member := &model.RoomMember{
		RoomID:   room.ID,
		UserID:   creatorID,
		Role:     "owner",
		JoinedAt: time.Now(),
	}

	return s.roomMemberRepo.Create(member)
}

func (s *RoomServiceImpl) JoinRoom(roomID, userID uint) error {
	if s.roomMemberRepo.IsMember(roomID, userID) {
		return nil
	}

	member := &model.RoomMember{
		RoomID:   roomID,
		UserID:   userID,
		Role:     "member",
		JoinedAt: time.Now(),
	}

	return s.roomMemberRepo.Create(member)
}

func (s *RoomServiceImpl) LeaveRoom(roomID, userID uint) error {
	return s.roomMemberRepo.Delete(roomID, userID)
}

func (s *RoomServiceImpl) GetRoom(roomID uint) (*model.Room, error) {
	return s.roomRepo.GetByID(roomID)
}

func (s *RoomServiceImpl) GetUserRooms(userID uint) ([]model.Room, error) {
	return s.roomRepo.GetUserRooms(userID)
}

func (s *RoomServiceImpl) SendRoomMessage(roomID uint, senderID uint, content string) error {
	// 检查用户是否在房间中
	if !s.roomMemberRepo.IsMember(roomID, senderID) {
		return errors.New("不是房间成员")
	}

	message := &model.Message{
		SenderID: senderID,
		RoomID:   &roomID,
		Content:  content,
		Type:     "text",
		Status:   "sent",
	}

	return s.messageService.SendMessage(message)
}
