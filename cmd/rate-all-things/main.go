package main

import (
	"log"

	database "github.com/kwiats/rate-all-things/pkg/db"
	"github.com/kwiats/rate-all-things/server"
)

func main() {
	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	srv := server.NewAPIServer(":8000", db)
	srv.Run()
}
