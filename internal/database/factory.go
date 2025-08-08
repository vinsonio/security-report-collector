package database

import (
	"fmt"

	"github.com/vinsonio/security-report-collector/internal/config"
)

// New creates a new database backend based on the provided configuration.
func New(cfg *config.DB) (DB, error) {
	switch cfg.Connection {
	case "sqlite":
		return NewSQLiteDB(cfg.SQLite)
	case "mysql":
		return NewMySQLDB(cfg.MySQL)
	default:
		return nil, fmt.Errorf("unsupported database connection: %s", cfg.Connection)
	}
}
