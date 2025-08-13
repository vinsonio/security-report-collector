package bootstrap

import (
	"testing"

	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/database"
)

// TestInit_SucceedsWithSQLiteAndFileCache verifies that Init can initialize when defaults are used.
func TestInit_SucceedsWithSQLiteAndFileCache(t *testing.T) {
	// Reset singletons for deterministic tests
	database.ResetSingletonForTest()
	cache.ResetSingletonForTest()

	// Explicitly set defaults to simulate "no env" behavior without relying on empty-string semantics
	t.Setenv("DB_CONNECTION", "sqlite")
	t.Setenv("CACHE_DRIVER", "file")
	// SQLite DB path must be writable; default is reports.db in repo root; isolate per test
	t.Setenv("DB_DATABASE", t.TempDir()+"/test.db")

	db, c, err := Init()
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	if db == nil || c == nil {
		t.Fatalf("expected non-nil db and cache")
	}
}

// TestInit_InvalidDBDriver ensures Init fails fast when DB driver unsupported.
func TestInit_InvalidDBDriver(t *testing.T) {
	// Reset singletons first so Get() runs again
	database.ResetSingletonForTest()
	cache.ResetSingletonForTest()

	t.Setenv("DB_CONNECTION", "invalid")

	if _, _, err := Init(); err == nil {
		t.Fatalf("expected error for invalid DB driver")
	}
}
