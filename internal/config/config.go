package config

import (
	"os"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Server config
	ServerPort string

	// Database config
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT config
	JWTSecret    string
	JWTExpiresIn time.Duration
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		// Server config
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Database config
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "myapp"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT config
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiresIn: time.Duration(getEnvAsInt("JWT_EXPIRES_IN", 24)) * time.Hour,
	}
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
