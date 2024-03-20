package model

import (
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

	Sender *User  `gorm:"foreignKey:SenderID"`
	Chat   *Chat  `gorm:"foreignKey:ChatID"`
	ReadBy []User `gorm:"many2many:message_reads;"`
}
