package database

import (
	"sync"

	"github.com/vinsonio/security-report-collector/internal/config"
)

var (
	once sync.Once
	db   DB
	err  error
)

// Get returns a new or existing database backend.
func Get() (DB, error) {
	once.Do(func() {
		cfg := config.NewDB()
		db, err = New(cfg)
	})
	return db, err
}

// ResetSingletonForTest resets the singleton for testing purposes.
func ResetSingletonForTest() {
	once = sync.Once{}
	db = nil
	err = nil
}
