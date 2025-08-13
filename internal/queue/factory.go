package queue

import (
	"fmt"

	"github.com/vinsonio/security-report-collector/internal/config"
)

// New creates a new queue based on the provided configuration.
func New(cfg *config.Cache, queueName string) (Queue, error) {
	switch cfg.Driver {
	case "redis":
		return NewRedisQueue(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, queueName)
	case "file", "memcached":
		// For file and memcached, we use in-memory queue (not persistent)
		return NewInMemoryQueue(), nil
	default:
		return nil, fmt.Errorf("unsupported queue driver: %s", cfg.Driver)
	}
}
