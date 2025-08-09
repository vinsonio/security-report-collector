package bootstrap

import (
	"log"

	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/database"
)

// Init initializes the application's dependencies.
func Init() (database.DB, cache.Cache, error) {
	db, err := database.Get()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Database connected successfully")

	cacheInstance, err := cache.Get()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Cache connected successfully")

	return db, cacheInstance, nil
}
