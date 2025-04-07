package validation

import (
	"os"
	"strconv"
	"strings"
)

// Config holds validation configuration
type Config struct {
	Password struct {
		MinLength        int      `json:"min_length"`
		RequireUppercase bool     `json:"require_uppercase"`
		RequireLowercase bool     `json:"require_lowercase"`
		RequireNumbers   bool     `json:"require_numbers"`
		RequireSpecial   bool     `json:"require_special"`
		SpecialChars     string   `json:"special_chars"`
		Disallowed       []string `json:"disallowed"`
	} `json:"password"`

	Username struct {
		MinLength    int      `json:"min_length"`
		MaxLength    int      `json:"max_length"`
		Reserved     []string `json:"reserved"`
		ProfaneWords []string `json:"profane_words"`
		AllowedChars string   `json:"allowed_chars"`
	} `json:"username"`

	Name struct {
		MinLength      int    `json:"min_length"`
		MaxLength      int    `json:"max_length"`
		AllowedChars   string `json:"allowed_chars"`
		SpecialChars   string `json:"special_chars"`
		MaxConsecutive int    `json:"max_consecutive"`
	} `json:"name"`
}

// DefaultConfig returns the default validation configuration
func DefaultConfig() *Config {
	return &Config{
		Password: struct {
			MinLength        int      `json:"min_length"`
			RequireUppercase bool     `json:"require_uppercase"`
			RequireLowercase bool     `json:"require_lowercase"`
			RequireNumbers   bool     `json:"require_numbers"`
			RequireSpecial   bool     `json:"require_special"`
			SpecialChars     string   `json:"special_chars"`
			Disallowed       []string `json:"disallowed"`
		}{
			MinLength:        8,
			RequireUppercase: true,
			RequireLowercase: true,
			RequireNumbers:   true,
			RequireSpecial:   true,
			SpecialChars:     "!@#$%^&*()_+-=[]{}|;:,.<>?",
			Disallowed:       []string{"password", "123456", "qwerty", "admin", "welcome"},
		},
		Username: struct {
			MinLength    int      `json:"min_length"`
			MaxLength    int      `json:"max_length"`
			Reserved     []string `json:"reserved"`
			ProfaneWords []string `json:"profane_words"`
			AllowedChars string   `json:"allowed_chars"`
		}{
			MinLength:    3,
			MaxLength:    20,
			Reserved:     []string{"admin", "root", "system", "support", "help", "info", "contact"},
			ProfaneWords: []string{"fuck", "shit", "ass", "bitch", "cunt"},
			AllowedChars: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_",
		},
		Name: struct {
			MinLength      int    `json:"min_length"`
			MaxLength      int    `json:"max_length"`
			AllowedChars   string `json:"allowed_chars"`
			SpecialChars   string `json:"special_chars"`
			MaxConsecutive int    `json:"max_consecutive"`
		}{
			MinLength:      2,
			MaxLength:      50,
			AllowedChars:   "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
			SpecialChars:   " -'",
			MaxConsecutive: 2,
		},
	}
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	if v := os.Getenv("VALIDATION_PASSWORD_MIN_LENGTH"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			config.Password.MinLength = n
		}
	}

	if v := os.Getenv("VALIDATION_PASSWORD_REQUIRE_UPPERCASE"); v != "" {
		config.Password.RequireUppercase = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("VALIDATION_PASSWORD_REQUIRE_LOWERCASE"); v != "" {
		config.Password.RequireLowercase = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("VALIDATION_PASSWORD_REQUIRE_NUMBERS"); v != "" {
		config.Password.RequireNumbers = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("VALIDATION_PASSWORD_REQUIRE_SPECIAL"); v != "" {
		config.Password.RequireSpecial = strings.ToLower(v) == "true"
	}

	if v := os.Getenv("VALIDATION_PASSWORD_SPECIAL_CHARS"); v != "" {
		config.Password.SpecialChars = v
	}

	if v := os.Getenv("VALIDATION_PASSWORD_DISALLOWED"); v != "" {
		config.Password.Disallowed = strings.Split(v, ",")
	}

	return config
}
