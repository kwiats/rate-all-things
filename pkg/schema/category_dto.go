package schema

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/datatypes"
)

type CategoryDTO struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	CustomFields []CategoryCustomFieldDTO `json:"customFields,omitempty"`
}

type CreateCategoryDTO struct {
	Name         string                   `json:"name"`
	CustomFields []CreateCategoryFieldDTO `json:"customFields,omitempty"`
}

type UpdateCategoryDTO struct {
	Name         string                   `json:"name"`
	CustomFields []CreateCategoryFieldDTO `json:"customFields,omitempty"`
}

type CustomFieldDTO struct {
	ID              uint           `json:"id"`
	Type            string         `json:"type"`
	DefaultSettings datatypes.JSON `json:"defaultSettings"`
}
type CreateCustomFieldDTO struct {
	Type            string         `json:"type"`
	DefaultSettings datatypes.JSON `json:"defaultSettings"`
}

type CreateCategoryFieldDTO struct {
	CustomFieldID uint           `json:"customFieldId"`
	Title         string         `json:"title"`
	Settings      datatypes.JSON `json:"settings,omitempty"`
}

type CategoryCustomFieldDTO struct {
	ID            uint           `json:"id"`
	CategoryID    uint           `json:"categoryId"`
	CustomFieldID uint           `json:"customFieldId"`
	Title         string         `json:"title"`
	Settings      datatypes.JSON `json:"settings,omitempty"`
}

func MapCCFToModel(categoryId uint, dtoFields []CreateCategoryFieldDTO) []*model.CategoryCustomField {
	var categoryCustomFields []*model.CategoryCustomField

	for _, field := range dtoFields {
		categoryCustomField := model.CategoryCustomField{
			CategoryID:    categoryId,
			CustomFieldID: field.CustomFieldID,
			Title:         field.Title,
			Settings:      field.Settings,
		}
		categoryCustomFields = append(categoryCustomFields, &categoryCustomField)
	}

	return categoryCustomFields
}
