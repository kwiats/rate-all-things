package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReviewItem struct {
	gorm.Model
	CategoryID     uint
	Category       Category
	UserID         uint
	User           User
	ReviewItemData []ReviewItemData `gorm:"foreignKey:ReviewItemID"`
	Reviews        []Review         `gorm:"foreignKey:ReviewItemID"`
}
type Review struct {
	gorm.Model
	ReviewItemID uint
	ReviewItem   ReviewItem
	UserID       uint
	User         User
	Rating       int
	Comment      string
}

type ReviewItemData struct {
	gorm.Model
	ReviewItemID          uint
	CategoryCustomFieldID uint
	CategoryCustomField   CategoryCustomField
	Value                 datatypes.JSON
}


