package database

import (
	"github.com/kwiats/rate-all-things/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitializeDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("rat.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.Review{},
		&model.Category{},
		&model.CustomField{},
		&model.CategoryCustomField{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
