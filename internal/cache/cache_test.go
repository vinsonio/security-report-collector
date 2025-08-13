package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vinsonio/security-report-collector/internal/config"
)

func testCache(t *testing.T, c Cache, expirationTTL time.Duration, expirationSleep time.Duration) {
	key := "test-key"
	value := []byte("test-value")

	// Test Set and Get
	err := c.Set(key, value, 1*time.Minute)
	assert.NoError(t, err)

	retrievedValue, err := c.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// Test Delete
	err = c.Delete(key)
	assert.NoError(t, err)

	retrievedValue, err = c.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, retrievedValue)

	// Test expiration
	err = c.Set(key, value, expirationTTL)
	assert.NoError(t, err)

	time.Sleep(expirationSleep)

	retrievedValue, err = c.Get(key)
	assert.NoError(t, err)
	assert.Nil(t, retrievedValue)
}

func TestFileCache(t *testing.T) {
	tmpDir := t.TempDir()

	c, err := NewFileCache(tmpDir)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()

	testCache(t, c, 1*time.Millisecond, 2*time.Millisecond)
}

func TestFileCache_InvalidPath(t *testing.T) {
	_, err := NewFileCache("/invalid/path/that/should/not/exist")
	assert.Error(t, err)
}

func TestFileCache_GetNonExistentKey(t *testing.T) {
	tmpDir := t.TempDir()

	c, err := NewFileCache(tmpDir)
	assert.NoError(t, err)
	defer func() { _ = c.Close() }()

	value, err := c.Get("non-existent-key")
	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestFactory_FileCache(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Cache{
		Driver: "file",
		File:   config.FileCache{Dir: tmpDir},
	}

	cache, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, cache)
	assert.IsType(t, &FileCache{}, cache)
}

func TestFactory_UnsupportedDriver(t *testing.T) {
	cfg := &config.Cache{
		Driver: "unsupported",
	}

	cache, err := New(cfg)
	assert.Error(t, err)
	assert.Nil(t, cache)
	assert.Contains(t, err.Error(), "unsupported cache driver")
}

func TestRedisCache(t *testing.T) {
	t.Skip("skipping test; redis not available in test environment")
	c, err := NewRedisCache("localhost:6379", "", 0)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()
	testCache(t, c, 1*time.Millisecond, 2*time.Millisecond)
}

func TestMemcachedCache(t *testing.T) {
	t.Skip("skipping test; memcached not available in test environment")
	c, err := NewMemcachedCache("localhost:11211")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, c.Close())
	}()
	testCache(t, c, 1*time.Second, 2*time.Second)
}
