package model

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name         string        `gorm:"type:varchar(100);not null;unique"`
	CustomFields []CustomField `gorm:"many2many:category_custom_fields;joinForeignKey:CategoryID;JoinReferences:CustomFieldID"`
}

type CustomField struct {
	gorm.Model
	Type       string     `gorm:"type:varchar(100);not null"`
	Categories []Category `gorm:"many2many:category_custom_fields;joinForeignKey:CustomFieldID;JoinReferences:CategoryID"`
}

type CategoryCustomField struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	CategoryID    uint
	CustomFieldID uint
	Title         string
}
