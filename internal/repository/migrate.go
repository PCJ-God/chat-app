package repository

import (
	"chat-app/internal/model"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Message{},
		&model.Room{},
		&model.RoomMember{},
		&model.UserSession{},
		&model.Notification{},
	)
}
