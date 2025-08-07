package testing

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vinsonio/security-report-collector/internal/storage"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// MockStore is a mock of storage.Store.
type MockStore struct {
	mock.Mock
}

// Save is a mock of the Save method.
func (m *MockStore) Save(reportType string, report types.Report, userAgent, hash string) error {
	args := m.Called(reportType, report, userAgent, hash)
	return args.Error(0)
}

// Store is an interface that extends the storage.Store interface with testing-specific methods.
type Store interface {
	storage.Store
	Count(t *testing.T) int
}

// getEnv gets an environment variable and trims quotes.
func getEnv(key string) string {
	val := os.Getenv(key)
	return strings.Trim(val, "\"")
}

// GetStoreForTest returns a Store implementation based on the current environment settings.
func GetStoreForTest(t *testing.T) Store {
	t.Helper()
	dbConnection := getEnv("DB_CONNECTION")
	// Default to in-memory sqlite for tests if no specific connection is set
	if dbConnection == "" || dbConnection == "sqlite" {
		// Use a unique name for each test to avoid conflicts if tests run in parallel.
		// The cache=shared is important to allow multiple connections to the same in-memory db.
		dbURL := "file:" + t.Name() + "?mode=memory&cache=shared"
		t.Setenv("DB_URL", dbURL)
		t.Setenv("DB_CONNECTION", "sqlite")
	}
	store, err := storage.NewStoreFromEnv()
	if err != nil {
		t.Fatalf("Failed to create store for testing: %v", err)
	}

	switch s := store.(type) {
	case *storage.SQLiteStore:
		return &sqliteStore{s}
	case *storage.MySQLStore:
		return &mysqlStore{s}
	default:
		t.Fatalf("Unsupported store type for testing: %T", s)
		return nil
	}
}

type sqliteStore struct {
	*storage.SQLiteStore
}

func (s *sqliteStore) Count(t *testing.T) int {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM reports").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count reports: %v", err)
	}
	return count
}

type mysqlStore struct {
	*storage.MySQLStore
}

func (s *mysqlStore) Count(t *testing.T) int {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM reports").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count reports: %v", err)
	}
	return count
}
