package main

import (
	api "go-onvif/json_apis"
	"log"
)

func main() {
	// Load configuration
	config := api.LoadConfig()

	// Create and run server
	server := api.NewAPIServer(config)
	server.SetupRoutes()

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
