package config

import "os"

type Config struct {
	BaseURL      string
	ServerPort   string
	DatabasePath string
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		ServerPort:   getEnv("PORT", ":8080"),
		DatabasePath: getEnv("DB_PATH", "./shortener.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
