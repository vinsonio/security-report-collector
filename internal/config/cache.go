package config

import (
	"os"
	"strconv"
	"strings"
)

// Cache holds the cache configuration.
type Cache struct {
	Driver    string
	File      FileCache
	Redis     Redis
	Memcached Memcached
}

// FileCache holds the file cache configuration.
type FileCache struct {
	Dir string
}

// Redis holds the Redis configuration.
type Redis struct {
	Addr     string
	Password string
	DB       int
}

// Memcached holds the Memcached configuration.
type Memcached struct {
	Servers []string
}

// NewCache creates a new Cache configuration.
func NewCache() *Cache {
	return &Cache{
		Driver: getEnv("CACHE_DRIVER", "file"),
		File: FileCache{
			Dir: getEnv("FILE_CACHE_DIR", "cache"),
		},
		Redis: Redis{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Memcached: Memcached{
			Servers: getEnvAsSlice("MEMCACHED_SERVERS", []string{"localhost:11211"}, ","),
		},
	}
}

// getEnvAsInt returns the value of an environment variable as an integer or a default value.
func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

// getEnvAsSlice returns the value of an environment variable as a slice of strings or a default value.
func getEnvAsSlice(key string, fallback []string, sep string) []string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.Split(value, sep)
	}
	return fallback
}
