package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApp_Defaults(t *testing.T) {
	// Ensure env vars are not set
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("CACHE_ENABLED")

	cfg := NewApp()
	assert.Equal(t, "report-collector", cfg.Name)
	assert.Equal(t, "development", cfg.Env)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, false, cfg.CacheEnabled)
}

func TestNewApp_FromEnv(t *testing.T) {
	os.Setenv("APP_NAME", "urc")
	os.Setenv("APP_ENV", "production")
	os.Setenv("APP_PORT", "9000")
	os.Setenv("CACHE_ENABLED", "true")
	t.Cleanup(func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_ENV")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("CACHE_ENABLED")
	})

	cfg := NewApp()
	assert.Equal(t, "urc", cfg.Name)
	assert.Equal(t, "production", cfg.Env)
	assert.Equal(t, "9000", cfg.Port)
	assert.Equal(t, true, cfg.CacheEnabled)
}