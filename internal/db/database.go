package database

import (
	"log"
	"tit/internal/db/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(connectionUri string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: connectionUri}), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		log.Printf("Failed to auto-migrate database models: %v\n", err)
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
