package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	category_model "github.com/kwiats/rate-all-things/internal/category"
)

func InitializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("rat.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&category_model.Category{},
		&category_model.CustomField{},
		&category_model.CategoryCustomField{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
