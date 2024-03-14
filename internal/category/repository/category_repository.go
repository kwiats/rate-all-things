package repository

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
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

func (repo *CategoryRepository) CreateCategoryWithCustomFields(category *model.Category, categoryCustomFields []*model.CategoryCustomField) (*model.Category, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(category).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, cf := range categoryCustomFields {
		cf.CategoryID = category.ID

		var customField model.CustomField
		if err := tx.Where("id = ?", cf.CustomFieldID).First(&customField).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if cf.Settings == nil {
			cf.Settings = customField.DefaultSettings
		}
	}

	if err := tx.Create(&categoryCustomFields).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *CategoryRepository) GetCategoryByID(id uint) (*model.Category, error) {
	var categoryModel model.Category

	if err := repo.db.First(&categoryModel, id).Error; err != nil {
		return nil, err
	}

	return &categoryModel, nil
}

func (repo *CategoryRepository) GetCategoryCustomField(id uint) ([]*model.CategoryCustomField, error) {
	var categoryCustomField []*model.CategoryCustomField
	if err := repo.db.Table("category_custom_fields").
		Select("category_custom_fields.id, custom_fields.id as custom_field_id, custom_fields.type, category_custom_fields.title, category_custom_fields.settings").
		Joins("join custom_fields on custom_fields.id = category_custom_fields.custom_field_id").
		Where("category_custom_fields.category_id = ?", id).
		Scan(&categoryCustomField).Error; err != nil {
		return nil, err
	}
	return categoryCustomField, nil
}

func (repo *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	var categories []*model.Category
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) UpdateCategory(id uint, category *model.Category, ccf []*model.CategoryCustomField) (*model.Category, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	var existingCategory model.Category
	if err := tx.First(&existingCategory, id).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&existingCategory).Updates(category).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, field := range ccf {
		if err := tx.Model(&model.CategoryCustomField{}).Where("id = ?", field.ID).Updates(field).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &existingCategory, nil
}

func (repo *CategoryRepository) DeleteCategoryByID(id uint, force bool) error {
	tx := repo.db.Begin()

	if force {
		if err := tx.Unscoped().Delete(&model.Category{}, id).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Where("category_id = ?", id).Delete(&model.CategoryCustomField{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Delete(&model.Category{}, id).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
