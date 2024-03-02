package repository

import (
	"github.com/kwiats/rate-all-things/internal/category/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	CreateCategory(*model.Category) (*model.Category, error)
	CreateCategoryWithCustomFields(*model.Category, *[]model.CategoryCustomField) (*model.Category, error)
	GetCategoryByID(uint) (*model.CategoryOutputDTO, error)
	GetAllCategories() ([]*model.Category, error)
	UpdateCategory(uint, *model.Category) (*model.Category, error)
	DeleteCategoryByID(uint, bool) error
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) CreateCategory(category *model.Category) (*model.Category, error) {
	if err := repo.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (repo *CategoryRepository) CreateCategoryWithCustomFields(category *model.Category, customFields *[]model.CategoryCustomField) (*model.Category, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(category).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range *customFields {
		(*customFields)[i].CategoryID = category.ID
	}

	if err := tx.Create(&customFields).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *CategoryRepository) GetCategoryByID(id uint) (*model.CategoryOutputDTO, error) {
	var categoryModel model.Category

	if err := repo.db.Preload("CustomFields").First(&categoryModel, id).Error; err != nil {
		return nil, err
	}

	categoryDTO := model.CategoryOutputDTO{
		ID:   categoryModel.ID,
		Name: categoryModel.Name,
	}

	var customFieldOutputs []model.CustomFieldOutputDTO
	if err := repo.db.Table("category_custom_fields").
		Select("custom_fields.id as custom_field_id, custom_fields.type, category_custom_fields.title").
		Joins("join custom_fields on custom_fields.id = category_custom_fields.custom_field_id").
		Where("category_custom_fields.category_id = ?", id).
		Scan(&customFieldOutputs).Error; err != nil {
		return nil, err
	}

	categoryDTO.CustomFields = customFieldOutputs

	return &categoryDTO, nil
}

func (repo *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	var categories []*model.Category
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) UpdateCategory(id uint, category *model.Category) (*model.Category, error) {
	var existingCategory model.Category
	if err := repo.db.First(&existingCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&existingCategory).Updates(category).Error; err != nil {
		return nil, err
	}

	return &existingCategory, nil
}

func (repo *CategoryRepository) DeleteCategoryByID(id uint, forceDelete bool) error {
	if forceDelete {
		if err := repo.db.Unscoped().Delete(&model.Category{}, id).Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.Delete(&model.Category{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}
