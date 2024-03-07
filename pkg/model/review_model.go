package model

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Rate       int `gorm:"not null"`
	CategoryID uint
	Category   Category
}
