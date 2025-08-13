package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisCache is a Redis cache implementation.
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache.
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Set stores a value in the cache.
func (c *RedisCache) Set(key string, value []byte, ttl time.Duration) error {
	return c.client.Set(c.ctx, key, value, ttl).Err()
}

// Get retrieves a value from the cache.
func (c *RedisCache) Get(key string) ([]byte, error) {
	val, err := c.client.Get(c.ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Not found
	}
	return val, err
}

// Delete removes a value from the cache.
func (c *RedisCache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Close closes the cache store.
func (c *RedisCache) Close() error {
	return c.client.Close()
}
