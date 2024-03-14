package service

import (
	"sync"

	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
)

type CustomFieldService struct {
	repository ICustomFieldRepository
}

type ICustomFieldRepository interface {
	CreateCustomField(*model.CustomField) (*model.CustomField, error)
	GetCustomFieldByID(uint) (*model.CustomField, error)
	GetAllCustomFields() ([]*model.CustomField, error)
	UpdateCustomField(uint, *model.CustomField) (*model.CustomField, error)
	DeleteCustomFieldByID(uint, bool) error
}



func NewCustomFieldService(repository ICustomFieldRepository) *CustomFieldService{
	return &CustomFieldService{repository: repository}
}

func (service *CustomFieldService) CreateCustomField(customFieldDTO schema.CreateCustomFieldDTO) (bool, error) {

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

func (service *CustomFieldService) GetCustomFields() ([]*schema.CustomFieldDTO, error) {
	customFields, err := service.repository.GetAllCustomFields()
	if err != nil {
		return nil, err
	}

	lenghtCustomFields := len(customFields)
	customFieldsDTO := make([]*schema.CustomFieldDTO, 0, lenghtCustomFields)
	channel := make(chan *schema.CustomFieldDTO, lenghtCustomFields)

	var wg sync.WaitGroup

	for _, customField := range customFields {
		wg.Add(1)
		go func(cf *model.CustomField) {
			defer wg.Done()
			channel <- &schema.CustomFieldDTO{
				ID:              cf.ID,
				Type:            cf.Type,
				DefaultSettings: cf.DefaultSettings,
			}
		}(customField)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for dto := range channel {
		customFieldsDTO = append(customFieldsDTO, dto)
	}

	return customFieldsDTO, nil
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
