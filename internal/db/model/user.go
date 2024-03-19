package model

import "gorm.io/gorm"

type User struct {
	*gorm.Model
	Username       string
	Email          string
	Password       string
	ProfilePicture string
}
