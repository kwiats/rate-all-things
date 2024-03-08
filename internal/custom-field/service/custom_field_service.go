package service

import (
	"github.com/kwiats/rate-all-things/internal/custom-field/repository"
	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
)

type CustomFieldService struct {
	repository repository.ICustomFieldRepository
}

type ICustomFieldService interface {
	CreateCustomField(schema.CustomFieldDTO) (bool, error)
	GetCustomField(uint) (schema.CustomFieldDTO, error)
	GetCustomFields() ([]schema.CustomFieldDTO, error)
	DeleteCustomField(uint, bool) (bool, error)
	UpdateCustomField(uint, schema.CustomFieldDTO) (bool, error)
}

func NewCustomFieldService(repository repository.ICustomFieldRepository) *CustomFieldService {
	return &CustomFieldService{repository: repository}
}

func (service *CustomFieldService) CreateCustomField(customFieldDTO schema.CustomFieldDTO) (bool, error) {

	customField := model.CustomField{
		Type:            customFieldDTO.Type,
		DefaultSettings: customFieldDTO.DefaultSettings,
	}

	if _, err := service.repository.CreateCustomField(&customField); err != nil {
		return false, err
	}
	return true, nil
}

func (service *CustomFieldService) GetCustomField(id uint) (schema.CustomFieldDTO, error) {
	customField, err := service.repository.GetCustomFieldByID(id)
	if err != nil {
		return schema.CustomFieldDTO{}, err
	}

	customFieldDTO := schema.CustomFieldDTO{
		Type:            customField.Type,
		DefaultSettings: customField.DefaultSettings,
	}

	return customFieldDTO, nil
}

func (service *CustomFieldService) GetCustomFields() ([]schema.CustomFieldDTO, error) {
	categories, err := service.repository.GetAllCustomFields()
	if err != nil {
		return nil, err
	}

	var categoriesDTO []schema.CustomFieldDTO
	for _, customField := range categories {
		categoriesDTO = append(categoriesDTO, schema.CustomFieldDTO{
			Id:              customField.ID,
			Type:            customField.Type,
			DefaultSettings: customField.DefaultSettings,
		})
	}

	return categoriesDTO, nil
}

func (service *CustomFieldService) DeleteCustomField(id uint, forceDelete bool) (bool, error) {
	err := service.repository.DeleteCustomFieldByID(id, forceDelete)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *CustomFieldService) UpdateCustomField(id uint, customFieldDTO schema.CustomFieldDTO) (bool, error) {

	customField := model.CustomField{
		Type:            customFieldDTO.Type,
		DefaultSettings: customFieldDTO.DefaultSettings,
	}

	if _, err := service.repository.UpdateCustomField(id, &customField); err != nil {
		return false, err
	}

	return true, nil
}
