package category_model

import "gorm.io/gorm"


type Category struct {
	gorm.Model        
	Name       string `gorm:"type:varchar(100);not null;unique"` 
	CustomFields []CustomField `gorm:"many2many:category_custom_fields;"`
}

type CustomField struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100);not null"`
	Type         string `gorm:"type:varchar(100);not null"`
	Categories   []Category `gorm:"many2many:category_custom_fields;"`
}

type CategoryCustomField struct {
	CategoryID   uint `gorm:"primaryKey"`
	CustomFieldID uint `gorm:"primaryKey"`
}