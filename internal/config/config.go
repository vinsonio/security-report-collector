package config

// Config holds all configurations for the application.
type Config struct {
	App   *App
	Cache *Cache
	DB    *DB
}

// New creates a new Config instance.
func New() *Config {
	return &Config{
		App:   NewApp(),
		Cache: NewCache(),
		DB:    NewDB(),
	}
}
