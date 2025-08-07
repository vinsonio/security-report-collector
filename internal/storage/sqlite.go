package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// SQLiteStore is a store that uses SQLite.
type SQLiteStore struct {
	DB *sql.DB
}

var (
	sqliteOnce sync.Once
	sqliteInst *SQLiteStore
	sqliteErr  error
)

func init() {
	RegisterStore("sqlite", newSQLiteStore)
}

// newSQLiteStore creates a new SQLiteStore.
func newSQLiteStore() (Store, error) {
	sqliteOnce.Do(func() {
		dbURL := getEnv("DB_URL")
		if dbURL == "" {
			dbURL = "./reports.db"
		}

		var db *sql.DB
		db, sqliteErr = sql.Open("sqlite3", dbURL)
		if sqliteErr != nil {
			return
		}

		driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			sqliteErr = err
			return
		}

		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		migrationsPath := "file://" + filepath.Join(basepath, "..", "..", "database", "migrations")

		m, err := migrate.NewWithDatabaseInstance(
			migrationsPath,
			"sqlite3",
			driver,
		)
		if err != nil {
			sqliteErr = err
			return
		}

		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			if dirtyErr, ok := err.(migrate.ErrDirty); ok {
				log.Printf("Database is dirty at version %d. Forcing to version %d and retrying.\n", dirtyErr.Version, dirtyErr.Version-1)
				if err = m.Force(int(dirtyErr.Version - 1)); err != nil {
					sqliteErr = fmt.Errorf("failed to force migration version: %w", err)
					return
				}
				if err = m.Up(); err != nil && err != migrate.ErrNoChange {
					sqliteErr = fmt.Errorf("failed to migrate up after forcing version: %w", err)
					return
				}
			} else {
				sqliteErr = err
				return
			}
		}

		sqliteInst = &SQLiteStore{DB: db}
	})

	return sqliteInst, sqliteErr
}

// Save saves a report to the database.
func (s *SQLiteStore) Save(reportType string, report types.Report, userAgent, hash string) error {
	data, err := report.JSON()
	if err != nil {
		return err
	}

	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, rand.Reader)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec("INSERT INTO reports (id, report_type, data, user_agent, hash) VALUES (?, ?, ?, ?, ?)", id.String(), reportType, string(data), userAgent, hash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: reports.hash") {
			return nil
		}
		return err
	}

	return nil
}
