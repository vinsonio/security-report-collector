package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestRedisCache(t *testing.T) {
	t.Skip("skipping test; redis not available in test environment")
	c, err := NewRedisCache("localhost:6379", "", 0)
	assert.NoError(t, err)
	defer c.Close()
	testCache(t, c, 1*time.Millisecond, 2*time.Millisecond)
}

func TestMemcachedCache(t *testing.T) {
	t.Skip("skipping test; memcached not available in test environment")
	c, err := NewMemcachedCache("localhost:11211")
	assert.NoError(t, err)
	defer c.Close()
	testCache(t, c, 1*time.Second, 2*time.Second)
}