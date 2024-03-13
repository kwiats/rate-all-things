package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string `gorm:"unique"`
	Password []byte
	Reviews  []Review `gorm:"foreignKey:UserID"`
}
