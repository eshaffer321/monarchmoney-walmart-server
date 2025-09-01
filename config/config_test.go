package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Arrange - Clear environment variables
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("GIN_MODE")
	_ = os.Unsetenv("SENTRY_DSN")
	_ = os.Unsetenv("EXTENSION_SECRET_KEY")

	// Act
	cfg := LoadConfig()

	// Assert
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "debug", cfg.GinMode)
	assert.Equal(t, "", cfg.SentryDSN)
	assert.Equal(t, "test-secret", cfg.ExtensionKey)
}

func TestLoadConfig_FromEnvironment(t *testing.T) {
	// Arrange
	_ = os.Setenv("PORT", "3000")
	_ = os.Setenv("GIN_MODE", "release")
	_ = os.Setenv("SENTRY_DSN", "https://test@sentry.io/123")
	_ = os.Setenv("EXTENSION_SECRET_KEY", "my-secret")
	defer func() {
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("GIN_MODE")
		_ = os.Unsetenv("SENTRY_DSN")
		_ = os.Unsetenv("EXTENSION_SECRET_KEY")
	}()

	// Act
	cfg := LoadConfig()

	// Assert
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "release", cfg.GinMode)
	assert.Equal(t, "https://test@sentry.io/123", cfg.SentryDSN)
	assert.Equal(t, "my-secret", cfg.ExtensionKey)
}

func TestConfig_IsSentryEnabled(t *testing.T) {
	tests := []struct {
		name      string
		sentryDSN string
		expected  bool
	}{
		{
			name:      "Enabled with DSN",
			sentryDSN: "https://test@sentry.io/123",
			expected:  true,
		},
		{
			name:      "Disabled without DSN",
			sentryDSN: "",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{SentryDSN: tt.sentryDSN}
			assert.Equal(t, tt.expected, cfg.IsSentryEnabled())
		})
	}
}
