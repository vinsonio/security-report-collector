package database

import (
	"sync"

	"github.com/vinsonio/security-report-collector/internal/config"
)

var (
	once sync.Once
	db   DB
)

// Get returns a new or existing database backend.
func Get() (DB, error) {
	var err error
	once.Do(func() {
		db, err = New(config.NewDB())
	})
	return db, err
}