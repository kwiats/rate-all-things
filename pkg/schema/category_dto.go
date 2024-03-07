package schema

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/datatypes"
)

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryDTO struct {
	Name         string                 `json:"name"`
	CustomFields []CreateCustomFieldDTO `json:"customFields"`
}

type UpdateCategoryDTO struct {
	Name         string                 `json:"name"`
	CustomFields []CreateCustomFieldDTO `json:"customFields"`
}

type CustomFieldDTO struct {
	Id              uint           `json:"id"`
	Type            string         `json:"type"`
	DefaultSettings datatypes.JSON `json:"defaultSettings"`
}
type CreateCustomFieldDTO struct {
	ID            uint           `json:"id"`
	CustomFieldId uint           `json:"customFieldId"`
	Title         string         `json:"Title"`
	Settings      datatypes.JSON `json:"settings"`
}

type CategoryOutputDTO struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	CustomFields []CustomFieldOutputDTO `json:"customFields"`
}

type CustomFieldOutputDTO struct {
	ID            uint           `json:"id"`
	CustomFieldId uint           `json:"customFieldId"`
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Settings      datatypes.JSON `json:"settings"`
}

func MapCCFToModel(categoryId uint, dtoFields []CreateCustomFieldDTO) []*model.CategoryCustomField {
	var categoryCustomFields []*model.CategoryCustomField

	for _, field := range dtoFields {
		categoryCustomField := model.CategoryCustomField{
			CategoryID:    categoryId,
			CustomFieldID: field.CustomFieldId,
			Title:         field.Title,
			Settings:      field.Settings,
		}
		categoryCustomFields = append(categoryCustomFields, &categoryCustomField)
	}

	return categoryCustomFields
}

//func MapCustomFieldsDTOToModel(customFieldsDTO []CustomFieldDTO) []CustomField {
//	var customFields []CustomField
//	for _, cfDTO := range customFieldsDTO {
//		customFields = append(customFields, CustomField{
//			: cfDTO.CustomFieldId,
//		})
//	}
//	return customFields
//}

//func MapCustomFieldsModelToDTO(CustomFields []CustomField) []CustomFieldDTO {
//	var customFieldsDTO []CustomFieldDTO
//	for _, cf := range CustomFields {
//		customFieldsDTO = append(customFieldsDTO, CustomFieldDTO{
//			Type: cf.Type,
//		})
//	}
//	return customFieldsDTO
//}
