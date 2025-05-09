package main

import (
	"log"
	// "github.com/BalkarSandhu/go-onvif/xml_apis"
	api "github.com/BalkarSandhu/go-onvif/json_apis"
)

func main() {
	// Load configuration (could be from environment, flags, or config file)
	config := api.Config{
		Port:           "8081",
		LogLevel:       "info",
		RateLimitReqs:  10, // 10 requests per second
		RateLimitBurst: 20, // Allow bursts of up to 20 requests
	}

	// Create and run server
	server := api.NewAPIServer(config)
	server.SetupRoutes()

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
