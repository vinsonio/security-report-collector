package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCache_Defaults(t *testing.T) {
	os.Unsetenv("CACHE_DRIVER")
	os.Unsetenv("FILE_CACHE_DIR")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("REDIS_DB")
	os.Unsetenv("MEMCACHED_SERVERS")

	cfg := NewCache()
	assert.Equal(t, "file", cfg.Driver)
	assert.Equal(t, "cache", cfg.File.Dir)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, 0, cfg.Redis.DB)
	assert.Equal(t, []string{"localhost:11211"}, cfg.Memcached.Servers)
}

func TestNewCache_FromEnv(t *testing.T) {
	os.Setenv("CACHE_DRIVER", "redis")
	os.Setenv("FILE_CACHE_DIR", "tmpcache")
	os.Setenv("REDIS_ADDR", "127.0.0.1:6380")
	os.Setenv("REDIS_PASSWORD", "secret")
	os.Setenv("REDIS_DB", "2")
	os.Setenv("MEMCACHED_SERVERS", "1.2.3.4:11211,5.6.7.8:11211")
	t.Cleanup(func() {
		os.Unsetenv("CACHE_DRIVER")
		os.Unsetenv("FILE_CACHE_DIR")
		os.Unsetenv("REDIS_ADDR")
		os.Unsetenv("REDIS_PASSWORD")
		os.Unsetenv("REDIS_DB")
		os.Unsetenv("MEMCACHED_SERVERS")
	})

	cfg := NewCache()
	assert.Equal(t, "redis", cfg.Driver)
	assert.Equal(t, "tmpcache", cfg.File.Dir)
	assert.Equal(t, "127.0.0.1:6380", cfg.Redis.Addr)
	assert.Equal(t, "secret", cfg.Redis.Password)
	assert.Equal(t, 2, cfg.Redis.DB)
	assert.Equal(t, []string{"1.2.3.4:11211", "5.6.7.8:11211"}, cfg.Memcached.Servers)
}