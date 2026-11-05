package model

import "time"

// Notification 通知模型
type Notification struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"not null;index" json:"user_id"`
	Type      string      `gorm:"size:20" json:"type"`
	Title     string      `gorm:"size:100" json:"title"`
	Content   string      `gorm:"type:text" json:"content"`
	IsRead    bool        `gorm:"default:false" json:"is_read"`
	Data      interface{} `gorm:"-" json:"data,omitempty"`
	CreatedAt time.Time   `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Notification) TableName() string {
	return "notifications"
}
