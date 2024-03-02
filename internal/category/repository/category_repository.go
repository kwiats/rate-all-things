package repository

import (
	category_model "github.com/kwiats/rate-all-things/internal/category/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

type ICategoryRepository interface {
	CreateCategory(*category_model.Category) (*category_model.Category, error)
	GetCategoryByID(uint) (*category_model.Category, error)
	GetAllCategories() ([]*category_model.Category, error)
	UpdateCategory(uint, *category_model.Category) (*category_model.Category, error)
	DeleteCategoryByID(uint, bool) error
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (repo *CategoryRepository) CreateCategory(category *category_model.Category) (*category_model.Category, error) {
	if err := repo.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (repo *CategoryRepository) GetCategoryByID(id uint) (*category_model.Category, error) {
	var category category_model.Category
	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (repo *CategoryRepository) GetAllCategories() ([]*category_model.Category, error) {
	var categories []*category_model.Category
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) UpdateCategory(id uint, category *category_model.Category) (*category_model.Category, error) {
	var existingCategory category_model.Category
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
		if err := repo.db.Unscoped().Delete(&category_model.Category{}, id).Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.Delete(&category_model.Category{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}
