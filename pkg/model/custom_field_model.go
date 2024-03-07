package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CustomField struct {
	gorm.Model
	Type            string         `gorm:"type:varchar(100);not null"`
	DefaultSettings datatypes.JSON `gorm:"type:json"`
	Categories      []Category     `gorm:"many2many:category_custom_fields;joinForeignKey:CustomFieldID;JoinReferences:CategoryID"`
}
