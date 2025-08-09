package config

import (
	"os"
	"strconv"
)

// App holds the application configuration.
type App struct {
	Name         string
	Env          string
	Port         string
	CacheEnabled bool
}

// NewApp creates a new App configuration.
func NewApp() *App {
	return &App{
		Name:         getEnv("APP_NAME", "report-collector"),
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("APP_PORT", "8080"),
		CacheEnabled: getEnvAsBool("CACHE_ENABLED", false),
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsBool returns the boolean value of an environment variable or a default value.
func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}