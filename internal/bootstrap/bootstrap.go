package bootstrap

import (
	"log"

	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/database"
)

// Init initializes the application's dependencies.
func Init() (database.DB, cache.Cache, error) {
	cfg := config.NewDB()
	db, err := database.New(cfg)
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
