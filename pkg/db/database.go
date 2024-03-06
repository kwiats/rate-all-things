package database

import (
	categoryModel "github.com/kwiats/rate-all-things/internal/category/model"
	reviewModel "github.com/kwiats/rate-all-things/internal/review/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("rat.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&reviewModel.Review{},
		&categoryModel.Category{},
		&categoryModel.CustomField{},
		&categoryModel.CategoryCustomField{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
