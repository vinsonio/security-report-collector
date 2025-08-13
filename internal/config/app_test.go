package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp_Defaults(t *testing.T) {
	// Explicitly set defaults to simulate unset env behavior
	t.Setenv("APP_NAME", "report-collector")
	t.Setenv("APP_ENV", "development")
	t.Setenv("APP_PORT", "8080")
	t.Setenv("CACHE_ENABLED", "false")

	cfg := NewApp()
	assert.Equal(t, "report-collector", cfg.Name)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, false, cfg.CacheEnabled)
}

func TestNewApp_FromEnv(t *testing.T) {
	t.Setenv("APP_NAME", "urc")
	t.Setenv("APP_ENV", "production")
	t.Setenv("APP_PORT", "9000")
	t.Setenv("CACHE_ENABLED", "true")

	cfg := NewApp()
	assert.Equal(t, "urc", cfg.Name)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, true, cfg.CacheEnabled)
}
