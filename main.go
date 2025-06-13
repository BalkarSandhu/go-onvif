package main

import (
	"go-onvif/api"
	"go-onvif/api/config"
	"log"
)

func main() {
	// Load configuration
	config := config.LoadConfig()

	// Create and run server
	server := api.NewAPIServer(config)
	server.SetupRoutes()

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
