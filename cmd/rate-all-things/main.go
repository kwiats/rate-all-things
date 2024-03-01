package main

import (
	"github.com/kwiats/rate-all-things/server"
)

func main() {
	server := server.RunAPIServer(":8000")
	server.Run()
}
