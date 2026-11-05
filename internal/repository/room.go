package repository

import (
	"chat-app/internal/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(room *model.Room) error
	GetByID(id uint) (*model.Room, error)
	GetUserRooms(userID uint) ([]model.Room, error)
}

type RoomMemberRepository interface {
	Create(member *model.RoomMember) error
	Delete(roomID, userID uint) error
	IsMember(roomID, userID uint) bool
}

type roomRepository struct {
	db *gorm.DB
}

type roomMemberRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func NewRoomMemberRepository(db *gorm.DB) RoomMemberRepository {
	return &roomMemberRepository{db: db}
}

func (r *roomRepository) Create(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepository) GetByID(id uint) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Members").Preload("Creator").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetUserRooms(userID uint) ([]model.Room, error) {
	var rooms []model.Room
	sub := r.db.Model(&model.RoomMember{}).Select("room_id").Where("user_id = ?", userID)
	if err := r.db.Where("id IN (?)", sub).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomMemberRepository) Create(member *model.RoomMember) error {
	return r.db.Create(member).Error
}

func (r *roomMemberRepository) Delete(roomID, userID uint) error {
	return r.db.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&model.RoomMember{}).Error
}

func (r *roomMemberRepository) IsMember(roomID, userID uint) bool {
	var count int64
	r.db.Model(&model.RoomMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count)
	return count > 0
}
