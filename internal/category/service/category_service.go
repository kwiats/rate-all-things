package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kwiats/rate-all-things/pkg/model"
	"github.com/kwiats/rate-all-things/pkg/schema"
	"gorm.io/gorm"
)

type CategoryService struct {
	repository ICategoryRepository
}

type ICategoryRepository interface {
	CreateCategory(*model.Category) (*model.Category, error)
	CreateCategoryWithCustomFields(*model.Category, []*model.CategoryCustomField) (*model.Category, error)
	GetCategoryByID(uint) (*model.Category, error)
	GetCategoryCustomField(id uint) ([]*model.CategoryCustomField, error)
	GetAllCategories() ([]*model.Category, error)
	UpdateCategory(uint, *model.Category, []*model.CategoryCustomField) (*model.Category, error)
	DeleteCategoryByID(uint, bool) error
}

func NewCategoryService(repository ICategoryRepository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (service *CategoryService) CreateCategory(createCategoryDTO schema.CreateCategoryDTO) (bool, error) {
	category := model.Category{
		Name: createCategoryDTO.Name,
	}
	categoryCustomFields := schema.MapCCFToModel(0, createCategoryDTO.CustomFields)
	_, err := service.repository.CreateCategoryWithCustomFields(&category, categoryCustomFields)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (service *CategoryService) GetCategory(id uint) (*schema.CategoryDTO, error) {
	category, err := service.repository.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &schema.CategoryDTO{}, fmt.Errorf("category with ID %d not found", id)
		}
		return &schema.CategoryDTO{}, err
	}

	categoryCustomFields, err := service.repository.GetCategoryCustomField(category.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &schema.CategoryDTO{}, fmt.Errorf("category with ID %d not found", id)
		}
		return &schema.CategoryDTO{}, err
	}
	ccfDTOs := make([]schema.CategoryCustomFieldDTO, 0, len(categoryCustomFields))
	channel := make(chan schema.CategoryCustomFieldDTO, len(categoryCustomFields))
	for _, ccf := range categoryCustomFields {
		go func(data *model.CategoryCustomField) {
			channel <- schema.CategoryCustomFieldDTO{
				ID:            data.ID,
				CategoryID:    data.CategoryID,
				CustomFieldID: data.CustomFieldID,
				Title:         data.Title,
				Settings:      data.Settings,
			}
		}(ccf)
	}
	for i := 0; i < len(categoryCustomFields); i++ {
		ccfDTO := <-channel
		ccfDTOs = append(ccfDTOs, ccfDTO)
	}
	close(channel)

	categoryDTO := schema.CategoryDTO{
		ID:           category.ID,
		Name:         category.Name,
		CustomFields: ccfDTOs,
	}

	return &categoryDTO, nil
}

func (service *CategoryService) GetCategories() ([]*schema.CategoryDTO, error) {
	categories, err := service.repository.GetAllCategories()
	if err != nil {
		return nil, err
	}

	categoriesDTO := make([]*schema.CategoryDTO, 0, len(categories))
	channel := make(chan *schema.CategoryDTO, len(categories))

	var wg sync.WaitGroup

	for _, category := range categories {
		wg.Add(1)
		go func(c *model.Category) {
			defer wg.Done()
			channel <- &schema.CategoryDTO{
				ID:   c.ID,
				Name: c.Name,
			}
		}(category)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for dto := range channel {
		categoriesDTO = append(categoriesDTO, dto)
	}

	return categoriesDTO, nil
}

func (service *CategoryService) DeleteCategory(id uint, forceDelete bool) (bool, error) {
	err := service.repository.DeleteCategoryByID(id, forceDelete)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service *CategoryService) UpdateCategory(id uint, updatedCategory schema.UpdateCategoryDTO) (bool, error) {
	category := model.Category{
		Name: updatedCategory.Name,
	}
	var categoryCustomFields []*model.CategoryCustomField
	for _, customField := range updatedCategory.CustomFields {
		categoryCustomFields = append(categoryCustomFields, &model.CategoryCustomField{
			CustomFieldID: customField.CustomFieldID,
			CategoryID:    id,
			Title:         customField.Title,
			Settings:      customField.Settings,
		})
	}

	if _, err := service.repository.UpdateCategory(id, &category, categoryCustomFields); err != nil {
		return false, err
	}

	return true, nil
}
