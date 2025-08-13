package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCache_Defaults(t *testing.T) {
	// Explicitly set defaults to simulate unset env behavior
	t.Setenv("CACHE_DRIVER", "file")
	t.Setenv("FILE_CACHE_DIR", "cache")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "0")
	t.Setenv("MEMCACHED_SERVERS", "localhost:11211")

	cfg := NewCache()
	assert.Equal(t, "file", cfg.Driver)
	assert.Equal(t, "cache", cfg.File.Dir)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, []string{"localhost:11211"}, cfg.Memcached.Servers)
}

func TestNewCache_FromEnv(t *testing.T) {
	t.Setenv("CACHE_DRIVER", "redis")
	t.Setenv("FILE_CACHE_DIR", "tmpcache")
	t.Setenv("REDIS_ADDR", "127.0.0.1:6380")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "2")
	t.Setenv("MEMCACHED_SERVERS", "1.2.3.4:11211,5.6.7.8:11211")

	cfg := NewCache()
	assert.Equal(t, "redis", cfg.Driver)
	assert.Equal(t, "tmpcache", cfg.File.Dir)
	assert.Equal(t, "127.0.0.1:6380", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, []string{"1.2.3.4:11211", "5.6.7.8:11211"}, cfg.Memcached.Servers)
}
