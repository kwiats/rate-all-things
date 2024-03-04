package service

import (
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
)

type CustomFieldService struct {
	repository repository.ICustomFieldRepository
}

type ICustomFieldService interface {
	CreateCustomField(model.CustomFieldDTO) (bool, error)
	GetCustomField(uint) (model.CustomFieldDTO, error)
	GetCategories() ([]model.CustomFieldDTO, error)
	DeleteCustomField(uint, bool) (bool, error)
	UpdateCustomField(uint, model.CustomFieldDTO) (bool, error)
}

func NewCustomFieldService(repository repository.ICustomFieldRepository) *CustomFieldService {
	return &CustomFieldService{repository: repository}
}

func (service *CustomFieldService) CreateCustomField(customFieldDTO model.CustomFieldDTO) (bool, error) {

	customField := model.CustomField{
		Type:            customFieldDTO.Type,
		DefaultSettings: customFieldDTO.DefaultSettings,
	}

	if _, err := service.repository.CreateCustomField(&customField); err != nil {
		return false, err
	}
	return true, nil
}

func (service *CustomFieldService) GetCustomField(id uint) (model.CustomFieldDTO, error) {
	customField, err := service.repository.GetCustomFieldByID(id)
	if err != nil {
		return model.CustomFieldDTO{}, err
	}

	customFieldDTO := model.CustomFieldDTO{
		//ID:   customField.ID,
		Type: customField.Type,
	}

	return customFieldDTO, nil
}

func (service *CustomFieldService) GetCategories() ([]model.CustomFieldDTO, error) {
	categories, err := service.repository.GetAllCustomFields()
	if err != nil {
		return nil, err
	}

	var categoriesDTO []model.CustomFieldDTO
	for _, customField := range categories {
		categoriesDTO = append(categoriesDTO, model.CustomFieldDTO{
			Type: customField.Type,
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

func (service *CustomFieldService) UpdateCustomField(id uint, customFieldDTO model.CustomFieldDTO) (bool, error) {

	customField := model.CustomField{
		//Name: customFieldDTO.Name,
		Type: customFieldDTO.Type,
	}

	if _, err := service.repository.UpdateCustomField(id, &customField); err != nil {
		return false, err
	}

	return true, nil
}
