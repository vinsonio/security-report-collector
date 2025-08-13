package cache

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileCache is a file-based cache implementation.
type FileCache struct {
	Dir string
}

// NewFileCache creates a new file-based cache.
func NewFileCache(dir string) (*FileCache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileCache{Dir: dir}, nil
}

// Set stores a value in the cache.
func (c *FileCache) Set(key string, value []byte, ttl time.Duration) error {
	hash := sha1.New()
	hash.Write([]byte(key))
	filename := hex.EncodeToString(hash.Sum(nil))
	path := filepath.Join(c.Dir, filename)

	// Add expiration to the file content
	expiration := time.Now().Add(ttl).UnixNano()
	data := []byte(strconv.FormatInt(expiration, 10) + "\n" + string(value))

	return os.WriteFile(path, data, 0644)
}

// Get retrieves a value from the cache.
func (c *FileCache) Get(key string) ([]byte, error) {
	hash := sha1.New()
	hash.Write([]byte(key))
	filename := hex.EncodeToString(hash.Sum(nil))
	path := filepath.Join(c.Dir, filename)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	// Check expiration
	parts := strings.SplitN(string(data), "\n", 2)
	if len(parts) != 2 {
		return nil, nil // Invalid format
	}

	expiration, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return nil, nil // Invalid format
	}

	if time.Now().UnixNano() > expiration {
		_ = os.Remove(path) // Expired
		return nil, nil
	}

	return []byte(parts[1]), nil
}

// Delete removes a value from the cache.
func (c *FileCache) Delete(key string) error {
	hash := sha1.New()
	hash.Write([]byte(key))
	filename := hex.EncodeToString(hash.Sum(nil))
	path := filepath.Join(c.Dir, filename)
	return os.Remove(path)
}

// Close closes the cache store.
func (c *FileCache) Close() error {
	return nil
}
