package database

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/oklog/ulid/v2"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// SQLiteDB is a db that uses SQLite.
type SQLiteDB struct {
	DB *sql.DB
}

// NewSQLiteDB creates a new SQLiteDB.
func NewSQLiteDB(cfg config.SQLite) (DB, error) {
	db, err := sql.Open("sqlite3", cfg.Database)
	if err != nil {
		return nil, err
	}

	return &SQLiteDB{DB: db}, nil
}

// Migrate runs the database migrations.
func (s *SQLiteDB) Migrate() error {
	var driver database.Driver
	driver, err := sqlite3.WithInstance(s.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	// get the path to the migrations directory
	_, b, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(b), "..", "..", "database", "migrations")

	var m *migrate.Migrate
	m, err = migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// Save saves a report to the database.
func (s *SQLiteDB) Save(reportType string, report types.Report, userAgent, hash string) error {
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
			return ErrDuplicateReport
		}
		return err
	}

	return nil
}
