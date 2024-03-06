package model

import (
	categorymodel "github.com/kwiats/rate-all-things/internal/category/model"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	Rate       int `gorm:"not null"`
	CategoryID uint
	Category   categorymodel.Category
}
