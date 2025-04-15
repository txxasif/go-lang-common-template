package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Server holds server configuration
type Server struct {
	Address string
}

// Database holds database configuration
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWT holds JWT configuration
type JWT struct {
	Secret           string
	AccessExpiresIn  time.Duration
	RefreshExpiresIn time.Duration
}

// Config holds all application configuration
type Config struct {
	Server   Server
	Database Database
	JWT      JWT
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Server: Server{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
		Database: Database{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "myapp"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWT{
			Secret:           getEnv("JWT_SECRET", "your-secret-key"),
			AccessExpiresIn:  time.Duration(getEnvAsInt("JWT_ACCESS_EXPIRES_IN", 24)) * time.Hour,
			RefreshExpiresIn: time.Duration(getEnvAsInt("JWT_REFRESH_EXPIRES_IN", 7*24)) * time.Hour,
		},
	}, nil
}

// getEnv retrieves environment variables with fallback values
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves environment variables as integers with fallback values
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return int(value.Hours())
	}
	return fallback
}
