package main

import (
	"log"

	"github.com/raghavendrah25/golang-backend/internal/server"
)

func main() {
	server, err := server.NewServer(server.Config{Port: "8080"})
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
