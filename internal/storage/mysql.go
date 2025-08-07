package storage

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
	migrate_mysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oklog/ulid/v2"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// MySQLStore is a MySQL-backed storage implementation.
type MySQLStore struct {
	DB *sql.DB
}

func init() {
	RegisterStore("mysql", newMySQLStore)
}

// newMySQLStore creates a new MySQLStore.
func newMySQLStore() (Store, error) {
	dbUser := getEnv("DB_USER")
	dbPassword := getEnv("DB_PASSWORD")
	dbHost := getEnv("DB_HOST")
	dbPort := getEnv("DB_PORT")
	dbName := getEnv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		return nil, fmt.Errorf("DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, and DB_NAME must be set for mysql connection")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	driver, err := migrate_mysql.WithInstance(db, &migrate_mysql.Config{})
	if err != nil {
		return nil, err
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationsPath := "file://" + filepath.Join(basepath, "..", "..", "database", "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &MySQLStore{DB: db}, nil
}

// Save saves a report to the database.
func (s *MySQLStore) Save(reportType string, report types.Report, userAgent, hash string) error {
	data, err := report.JSON()
	if err != nil {
		return err
	}

	ms := ulid.Timestamp(time.Now())
	id, err := ulid.New(ms, rand.Reader)
	if err != nil {
		return err
	}

	res, err := s.DB.Exec("INSERT INTO reports (id, report_type, data, user_agent, hash) VALUES (?, ?, ?, ?, ?)", id.String(), reportType, data, userAgent, hash)
	if err != nil {
		mysqlErr, ok := err.(*mysql.MySQLError)
		if ok && mysqlErr.Number == 1062 {
			return nil
		}
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were affected")
	}

	return nil
}
