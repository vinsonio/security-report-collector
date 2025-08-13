package bootstrap

import (
	"log"

	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/database"
	"github.com/vinsonio/security-report-collector/internal/service"
)

// Init initializes the application's dependencies.
func Init() (service.Database, service.Cacher, error) {
	db, err := database.Get()
	if err != nil {
		return nil, nil, err
	}

	log.Println("Database connected successfully")

	if err := db.Migrate(); err != nil {
		return nil, nil, err
	}

	log.Println("Database migration completed")

	cacheInstance, err := cache.Get()
	if err != nil {
		return nil, nil, err
	}
	log.Println("Cache connected successfully")

	return db, cacheInstance, nil
}
