package model

type CategoryDTO struct {
	ID           uint             `json:"id"`
	Name         string           `json:"name"`
	CustomFields []CustomFieldDTO `json:"customFields"`
}
type CustomFieldDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func MapCustomFieldsDTOToModel(customFieldsDTO []CustomFieldDTO) []CustomField {
	var customFields []CustomField
	for _, cfDTO := range customFieldsDTO {
		customFields = append(customFields, CustomField{
			Name: cfDTO.Name,
			Type: cfDTO.Type,
		})
	}
	return customFields
}

func MapCustomFieldsModelToDTO(CustomFields []CustomField) []CustomFieldDTO {
	var customFieldsDTO []CustomFieldDTO
	for _, cfDTO := range CustomFields {
		customFieldsDTO = append(customFieldsDTO, CustomFieldDTO{
			Name: cfDTO.Name,
			Type: cfDTO.Type,
		})
	}
	return customFieldsDTO
}
