package main

import (
	"log"

	database "github.com/kwiats/rate-all-things/pkg/db"
	"github.com/kwiats/rate-all-things/server"
)

func main() {
	_, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	server := server.RunAPIServer(":8000")
	server.Run()
}
