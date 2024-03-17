package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name                 string
	CategoryCustomFields []CategoryCustomField `gorm:"foreignKey:CategoryID"`
}
type CategoryCustomField struct {
	gorm.Model
	CategoryID    uint
	CustomFieldID uint
	CustomField   CustomField
	Title         string
	Settings      datatypes.JSON
}
