package cache

import (
	"sync"

	"github.com/vinsonio/security-report-collector/internal/config"
)

var (
	once  sync.Once
	cache Cache
	err   error
)

// Get returns the singleton instance of the cache.
func Get() (Cache, error) {
	once.Do(func() {
		cfg := config.NewCache()
		cache, err = New(cfg)
	})
	return cache, err
}