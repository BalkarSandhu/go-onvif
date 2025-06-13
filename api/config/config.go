package config

import (
	"log"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port           string
	LogLevel       string
	RateLimitReqs  int
	RateLimitBurst int
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables directly...")
	}

	getEnv := func(key, fallback string) string {
		if val := os.Getenv(key); val != "" {
			return val
		}
		return fallback
	}

	toInt := func(val string, fallback int) int {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
		return fallback
	}

	return Config{
		Port:           getEnv("PORT", "8084"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		RateLimitReqs:  toInt(getEnv("RATE_LIMIT_REQS", "10"), 10),
		RateLimitBurst: toInt(getEnv("RATE_LIMIT_BURST", "20"), 20),
	}
}
