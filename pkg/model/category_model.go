package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name         string        `gorm:"type:varchar(100);not null;unique"`
	CustomFields []CustomField `gorm:"many2many:category_custom_fields;joinForeignKey:CategoryID;JoinReferences:CustomFieldID"`
}

type CategoryCustomField struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	CategoryID    uint
	CustomFieldID uint
	Title         string
	Settings      datatypes.JSON `gorm:"type:json"`
}
