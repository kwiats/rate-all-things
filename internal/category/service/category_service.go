package service

import (
	"errors"
	"fmt"
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
	"gorm.io/gorm"
)

type CategoryService struct {
	repository repository.ICategoryRepository
}

type ICategoryService interface {
	CreateCategory(model.CategoryDTO) (bool, error)
	GetCategory(uint) (model.CategoryDTO, error)
	GetCategories() ([]model.CategoryDTO, error)
	DeleteCategory(uint, bool) (bool, error)
	UpdateCategory(uint, model.CategoryDTO) (bool, error)
}

func NewCategoryService(repository repository.ICategoryRepository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (service *CategoryService) CreateCategory(createCategoryDTO model.CreateCategoryDTO) (bool, error) {
	category := model.Category{
		Name: createCategoryDTO.Name,
	}
	categoryCustomFields := model.MapCCFToModel(0, createCategoryDTO.CustomFields)
	_, err := service.repository.CreateCategoryWithCustomFields(&category, categoryCustomFields)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (service *CategoryService) GetCategory(id uint) (model.CategoryOutputDTO, error) {
	categoryWithCustomFields, err := service.repository.GetCategoryByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.CategoryOutputDTO{}, fmt.Errorf("category with ID %d not found", id)
		}
		return model.CategoryOutputDTO{}, err
	}

	return *categoryWithCustomFields, nil
}

func (service *CategoryService) GetCategories() ([]model.CategoryDTO, error) {
	categories, err := service.repository.GetAllCategories()
	if err != nil {
		return nil, err
	}

	var categoriesDTO []model.CategoryDTO
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, model.CategoryDTO{
			ID:   category.ID,
			Name: category.Name,
		})
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

func (service *CategoryService) UpdateCategory(id uint, updatedCategory model.UpdateCategoryDTO) (bool, error) {
	category := model.Category{
		Name: updatedCategory.Name,
	}
	var categoryCustomFields []*model.CategoryCustomField
	for _, customField := range updatedCategory.CustomFields {
		categoryCustomFields = append(categoryCustomFields, &model.CategoryCustomField{
			ID:            customField.ID,
			CustomFieldID: customField.CustomFieldId,
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
