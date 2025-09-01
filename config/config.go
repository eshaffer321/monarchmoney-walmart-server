// Package config provides configuration management for the Walmart-Monarch sync backend.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values for the application.
type Config struct {
	Port           string
	GinMode        string
	SentryDSN      string
	ExtensionKey   string
	MonarchAPIKey  string
	OllamaEndpoint string
	OpenAIAPIKey   string
	ClaudeAPIKey   string
}

// LoadConfig loads configuration from environment variables with fallback to defaults.
func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load() // Ignore error as it's OK if .env doesn't exist

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		GinMode:        getEnv("GIN_MODE", "debug"),
		SentryDSN:      getEnv("SENTRY_DSN", ""),
		ExtensionKey:   getEnv("EXTENSION_SECRET_KEY", "test-secret"),
		MonarchAPIKey:  getEnv("MONARCH_API_KEY", ""),
		OllamaEndpoint: getEnv("OLLAMA_ENDPOINT", "http://localhost:11434"),
		OpenAIAPIKey:   getEnv("OPENAI_API_KEY", ""),
		ClaudeAPIKey:   getEnv("CLAUDE_API_KEY", ""),
	}

	return cfg
}

// IsSentryEnabled returns true if Sentry error tracking is configured.
func (c *Config) IsSentryEnabled() bool {
	return c.SentryDSN != ""
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
