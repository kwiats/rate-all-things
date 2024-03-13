package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CustomField struct {
	gorm.Model
	Type                 string                // "textField", "fileField"
	DefaultSettings      datatypes.JSON        // '{"defaultValue": "", "maxLength": 255}'
	CategoryCustomFields []CategoryCustomField `gorm:"foreignKey:CustomFieldID"`
}
