package model

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryDTO struct {
	Name         string                 `json:"name"`
	CustomFields []CreateCustomFieldDTO `json:"customFields"`
}

type CustomFieldDTO struct {
	Id   uint   `json:"id"`
	Type string `json:"Type"`
}
type CreateCustomFieldDTO struct {
	CustomFieldId uint   `json:"customFieldId"`
	Title         string `json:"Title"`
}

type CategoryOutputDTO struct {
	ID           uint                   `json:"id"`
	Name         string                 `json:"name"`
	CustomFields []CustomFieldOutputDTO `json:"customFields"`
}

type CustomFieldOutputDTO struct {
	CustomFieldId uint   `json:"customFieldId"`
	Type          string `json:"type"`
	Title         string `json:"title"`
}

func MapCCFToModel(categoryId uint, dtoFields []CreateCustomFieldDTO) []CategoryCustomField {
	var categoryCustomFields []CategoryCustomField

	for _, field := range dtoFields {
		categoryCustomField := CategoryCustomField{
			CategoryID:    categoryId,
			CustomFieldID: field.CustomFieldId,
			Title:         field.Title,
		}
		categoryCustomFields = append(categoryCustomFields, categoryCustomField)
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
