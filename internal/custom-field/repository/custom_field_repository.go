package repository

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/gorm"
)

type CustomFieldRepository struct {
	db *gorm.DB
}

func NewCustomFieldRepository(db *gorm.DB) *CustomFieldRepository {
	return &CustomFieldRepository{db: db}
}

func (repo *CustomFieldRepository) CreateCustomField(customField *model.CustomField) (*model.CustomField, error) {
	if err := repo.db.Create(customField).Error; err != nil {
		return nil, err
	}
	return customField, nil
}

func (repo *CustomFieldRepository) GetCustomFieldByID(id uint) (*model.CustomField, error) {
	var customField model.CustomField
	if err := repo.db.First(&customField, id).Error; err != nil {
		return nil, err
	}
	return &customField, nil
}

func (repo *CustomFieldRepository) GetAllCustomFields() ([]*model.CustomField, error) {
	var categories []*model.CustomField
	if err := repo.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CustomFieldRepository) UpdateCustomField(id uint, customField *model.CustomField) (*model.CustomField, error) {
	var existingCustomField model.CustomField
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
		if err := repo.db.Unscoped().Delete(&model.CustomField{}, id).Error; err != nil {
			return err
		}
	} else {
		if err := repo.db.Delete(&model.CustomField{}, id).Error; err != nil {
			return err
		}
	}
	return nil
}
