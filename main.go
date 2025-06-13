package main

import (
	api "go-onvif/json_apis"
	"go-onvif/json_apis/config"
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
