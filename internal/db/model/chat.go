package model

import (
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	*gorm.Model

	Users    []User    `gorm:"many2many:chat_users;"`
	Messages []Message `gorm:"foreignKey:ChatID"`
}

type Message struct {
	*gorm.Model
	Content  string `gorm:"type:text"`
	SenderID uint
	ChatID   uint

	Sender       *User         `gorm:"foreignKey:SenderID"`
	Chat         *Chat         `gorm:"foreignKey:ChatID"`
	MessageReads []MessageRead `gorm:"foreignKey:MessageID"`
}

type MessageRead struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	MessageID uint
	ReadAt    time.Time

	User    *User    `gorm:"foreignKey:UserID"`
	Message *Message `gorm:"foreignKey:MessageID"`
}
