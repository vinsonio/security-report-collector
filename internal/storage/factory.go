package storage

import (
	"os"
	"strings"
)

// getEnv gets an environment variable and trims quotes.
func getEnv(key string) string {
	val := os.Getenv(key)
	return strings.Trim(val, "\"")
}

// NewStoreFromEnv creates a new storage backend based on environment variables.
func NewStoreFromEnv() (Store, error) {
	dbConnection := getEnv("DB_CONNECTION")
	if dbConnection == "" {
		dbConnection = "sqlite"
	}

	builder, err := GetStoreBuilder(dbConnection)
	if err != nil {
		return nil, err
	}

	return builder()
}
