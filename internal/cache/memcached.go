package cache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache is a Memcached cache implementation.
type MemcachedCache struct {
	client *memcache.Client
}

// NewMemcachedCache creates a new Memcached cache.
func NewMemcachedCache(servers ...string) (*MemcachedCache, error) {
	client := memcache.New(servers...)
	return &MemcachedCache{client: client}, nil
}

// Set stores a value in the cache.
func (c *MemcachedCache) Set(key string, value []byte, ttl time.Duration) error {
	return c.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(ttl.Seconds()),
	})
}

// Get retrieves a value from the cache.
func (c *MemcachedCache) Get(key string) ([]byte, error) {
	it, err := c.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}
	return it.Value, nil
}

// Delete removes a value from the cache.
func (c *MemcachedCache) Delete(key string) error {
	return c.client.Delete(key)
}

// Close closes the cache store.
func (c *MemcachedCache) Close() error {
	return nil // gomemcache does not have a Close method
}
