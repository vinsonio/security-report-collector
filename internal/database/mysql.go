package database

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oklog/ulid/v2"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// MySQLDB is a MySQL-backed database implementation.
type MySQLDB struct {
	DB *sql.DB
}

// NewMySQLDB creates a new MySQLDB.
func NewMySQLDB(cfg config.MySQL) (DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &MySQLDB{DB: db}, nil
}

// Migrate runs the database migrations.
func (s *MySQLDB) Migrate() error {
	var driver database.Driver
	driver, err := migrate_mysql.WithInstance(s.DB, &migrate_mysql.Config{})
	if err != nil {
		return err
	}

	// get the path to the migrations directory
	_, b, _, _ := runtime.Caller(0)
	migrationsPath := filepath.Join(filepath.Dir(b), "..", "..", "database", "migrations")

	var m *migrate.Migrate
	m, err = migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

// Save saves a report to the database.
func (s *MySQLDB) Save(reportType string, report types.Report, userAgent, hash string) error {
	data, err := report.JSON()
	if err != nil {
		return err
	}

	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, rand.Reader)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec("INSERT INTO reports (id, report_type, data, user_agent, hash) VALUES (?, ?, ?, ?, ?)", id.String(), reportType, data, userAgent, hash)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			return ErrDuplicateReport
		}
		return err
	}

	return nil
}
