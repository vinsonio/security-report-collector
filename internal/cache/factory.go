package cache

import (
	"fmt"

	"github.com/vinsonio/security-report-collector/internal/config"
)

// New creates a new cache based on the provided configuration.
func New(cfg *config.Cache) (Cache, error) {
	switch cfg.Driver {
	case "file":
		return NewFileCache(cfg.File.Dir)
	case "redis":
		return NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	case "memcached":
		return NewMemcachedCache(cfg.Memcached.Servers...)
	default:
		return nil, fmt.Errorf("unsupported cache driver: %s", cfg.Driver)
	}
}
