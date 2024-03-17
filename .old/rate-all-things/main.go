package main

import (
	"gorm.io/gorm"
	"log"

	"github.com/kwiats/rate-all-things/pkg/config"
	"github.com/kwiats/rate-all-things/pkg/db"
	"github.com/kwiats/rate-all-things/server"
)

func main() {
	settings, db := Initialize()

	srv := server.NewAPIServer(":8000", db, settings)
	srv.Run()
}

func Initialize() (*config.Config, *gorm.DB) {
	configDone := make(chan *config.Config)
	dbDone := make(chan *gorm.DB)
	go initializeConfigurations(configDone)
	settings := <-configDone
	go initializeDatabase(settings, dbDone)
	db := <-dbDone
	return settings, db
}

func initializeConfigurations(cfg chan<- *config.Config) {
	settings, err := config.NewConfiguration()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	cfg <- settings
}

func initializeDatabase(settings *config.Config, dbDone chan<- *gorm.DB) {
	db, err := database.InitializeDB(settings.GetDBConnectionUri())
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	dbDone <- db
}
