package cache

import "time"

// Cache is the interface for a cache store.
type Cache interface {
	// Set stores a value in the cache.
	Set(key string, value []byte, ttl time.Duration) error
	// Get retrieves a value from the cache.
	Get(key string) ([]byte, error)
	// Delete removes a value from the cache.
	Delete(key string) error
	// Close closes the cache store.
	Close() error
}
