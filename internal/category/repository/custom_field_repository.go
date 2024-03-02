package repository

import (
	category_model "github.com/kwiats/rate-all-things/internal/category/model"
	"gorm.io/gorm"
)

type CustomFieldRepository struct {
	db *gorm.DB
}

type ICustomFieldRepository interface {
	CreateCustomField(*category_model.CustomField) (*category_model.CustomField, error)
	GetCustomFieldByID(uint) (*category_model.CustomField, error)
	GetAllCustomFields() ([]*category_model.CustomField, error)
	UpdateCustomField(uint, *category_model.CustomField) (*category_model.CustomField, error)
	DeleteCustomFieldByID(uint, bool) error
}

func NewCustomFieldRepository(db *gorm.DB) *CustomFieldRepository {
	return &CustomFieldRepository{db: db}
}

func (repo *CustomFieldRepository) CreateCustomField(customField *category_model.CustomField) (*category_model.CustomField, error) {
	if err := repo.db.Create(customField).Error; err != nil {
		return nil, err
	}
	return customField, nil
}

func (repo *CustomFieldRepository) GetCustomFieldByID(id uint) (*category_model.CustomField, error) {
	var customField category_model.CustomField
	if err := repo.db.First(&customField, id).Error; err != nil {
		return nil, err
	}
	return &customField, nil
}

func (repo *CustomFieldRepository) GetAllCustomFields() ([]*category_model.CustomField, error) {
	var categories []*category_model.CustomField
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CustomFieldRepository) UpdateCustomField(id uint, customField *category_model.CustomField) (*category_model.CustomField, error) {
	var existingCustomField category_model.CustomField
	if err := repo.db.First(&existingCustomField, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&existingCustomField).Updates(customField).Error; err != nil {
		return nil, err
	}

	return &existingCustomField, nil
}

func (repo *CustomFieldRepository) DeleteCustomFieldByID(id uint, forceDelete bool) error {
	if forceDelete {
		if err := repo.db.Unscoped().Delete(&category_model.CustomField{}, id).Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.Delete(&category_model.CustomField{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}
