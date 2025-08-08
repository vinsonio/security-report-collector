package config

import "os"

// App holds the application configuration.

type App struct {
	Name string
	Env  string
	Port string
}

// NewApp creates a new App configuration.

func NewApp() *App {
	return &App{
		Name: getEnv("APP_NAME", "report-collector"),
		Env:  getEnv("APP_ENV", "development"),
		Port: getEnv("APP_PORT", "8080"),
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}