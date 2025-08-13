package testing

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/database"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// MockDB is a mock of database.DB.
type MockDB struct {
	mock.Mock
}

// Save is a mock of the Save method.
func (m *MockDB) Save(reportType string, report types.Report, userAgent, hash string) error {
	args := m.Called(reportType, report, userAgent, hash)
	return args.Error(0)
}

// DB is an interface that extends the database.DB interface with testing-specific methods.
type DB interface {
	database.DB
	Count(t *testing.T) int
}

// GetDBForTest returns a DB implementation based on the current environment settings.
func GetDBForTest(t *testing.T) DB {
	t.Helper()
	cfg := config.NewDB()
	db, err := database.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create db for testing: %v", err)
	}

	// Ensure database schema is migrated before attempting to modify tables
	if err := db.Migrate(); err != nil {
		t.Fatalf("Failed to migrate db for testing: %v", err)
	}

	switch d := db.(type) {
	case *database.SQLiteDB:
		_, err := d.DB.Exec("DELETE FROM reports")
		if err != nil {
			t.Fatalf("failed to truncate reports table: %v", err)
		}
		return &sqliteDB{d}
	case *database.MySQLDB:
		_, err := d.DB.Exec("TRUNCATE TABLE reports")
		if err != nil {
			t.Fatalf("failed to truncate reports table: %v", err)
		}
		return &mysqlDB{d}
	default:
		t.Fatalf("Unsupported db type for testing: %T", d)
		return nil
	}
}

type sqliteDB struct {
	*database.SQLiteDB
}

func (s *sqliteDB) Count(t *testing.T) int {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM reports").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count reports: %v", err)
	}
	return count
}

type mysqlDB struct {
	*database.MySQLDB
}

func (s *mysqlDB) Count(t *testing.T) int {
	var count int
	err := s.DB.QueryRow("SELECT COUNT(*) FROM reports").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count reports: %v", err)
	}
	return count
}
