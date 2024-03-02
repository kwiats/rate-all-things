package service

import (
	"github.com/kwiats/rate-all-things/internal/category/model"
	"github.com/kwiats/rate-all-things/internal/category/repository"
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

func (service *CategoryService) CreateCategory(categoryDTO model.CategoryDTO) (bool, error) {
	customFields := model.MapCustomFieldsDTOToModel(categoryDTO.CustomFields)

	category := model.Category{
		Name:         categoryDTO.Name,
		CustomFields: customFields,
	}

	if _, err := service.repository.CreateCategory(&category); err != nil {
		return false, err
	}
	return true, nil
}

func (service *CategoryService) GetCategory(id uint) (model.CategoryDTO, error) {
	category, err := service.repository.GetCategoryByID(id)
	if err != nil {
		return model.CategoryDTO{}, err
	}

	categoryDTO := model.CategoryDTO{
		ID:           category.ID,
		Name:         category.Name,
		CustomFields: model.MapCustomFieldsModelToDTO(category.CustomFields),
	}

	return categoryDTO, nil
}

func (service *CategoryService) GetCategories() ([]model.CategoryDTO, error) {
	categories, err := service.repository.GetAllCategories()
	if err != nil {
		return nil, err
	}

	var categoriesDTO []model.CategoryDTO
	for _, category := range categories {
		categoriesDTO = append(categoriesDTO, model.CategoryDTO{
			ID:           category.ID,
			Name:         category.Name,
			CustomFields: model.MapCustomFieldsModelToDTO(category.CustomFields),
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

func (service *CategoryService) UpdateCategory(id uint, categoryDTO model.CategoryDTO) (bool, error) {
	customFields := model.MapCustomFieldsDTOToModel(categoryDTO.CustomFields)

	category := model.Category{
		Name:         categoryDTO.Name,
		CustomFields: customFields,
	}

	if _, err := service.repository.UpdateCategory(id, &category); err != nil {
		return false, err
	}

	return true, nil
}
