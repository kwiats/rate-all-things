package model

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Username       string `gorm:"uniqueIndex"`
	Email          string `gorm:"uniqueIndex"`
	Password       string
	ProfilePicture string

	Chats        []*Chat    `gorm:"many2many:chat_users;"`
	MessagesRead []*Message `gorm:"many2many:message_reads;"`
}
