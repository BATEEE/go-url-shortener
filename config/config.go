package config

import "os"

type Config struct {
	BaseURL      string
	ServerPort   string
	DatabasePath string
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	if port != "" && port[0] != ':' {
		port = ":" + port
	}

	return &Config{
		BaseURL:      getEnv("BASE_URL", "http://localhost:8080"),
		ServerPort:   port,
		DatabasePath: getEnv("DB_PATH", "./shortener.db"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
