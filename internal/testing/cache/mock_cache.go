package cache

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockCache is a mock of cache.Cache
type MockCache struct {
	mock.Mock
}

// Get is a mock of cache.Cache.Get
func (m *MockCache) Get(key string) ([]byte, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// Set is a mock of cache.Cache.Set
func (m *MockCache) Set(key string, value []byte, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

// Delete is a mock of cache.Cache.Delete
func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

// Close is a mock of cache.Cache.Close
func (m *MockCache) Close() error {
	args := m.Called()
	return args.Error(0)
}
