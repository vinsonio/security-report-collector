package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/vinsonio/security-report-collector/internal/types"
	"github.com/vinsonio/security-report-collector/internal/util"
)

// Database is the interface for database operations.
type Database interface {
	Save(reportType string, report types.Report, userAgent string, hash string) error
}

// Cacher is the interface for cache operations.
type Cacher interface {
	Set(key string, value []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}

// ReportService is the service for handling reports.
type ReportService struct {
	db           Database
	cache        Cacher
	cacheEnabled bool
}

// NewReportService creates a new ReportService.
func NewReportService(db Database, cache Cacher, cacheEnabled bool) *ReportService {
	return &ReportService{db: db, cache: cache, cacheEnabled: cacheEnabled}
}

// SaveReport saves a report.
func (s *ReportService) SaveReport(reportType string, report types.Report, userAgent string) error {
	hashData, err := report.HashData()
	if err != nil {
		return err
	}

	data, err := util.StableMarshal(hashData)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// If cache is enabled, short-circuit on cache hit
	if s.cacheEnabled {
		b, err := s.cache.Get(hashStr)
		if err != nil {
			return err
		}
		if b != nil {
			// Report already cached; treat as success without hitting DB
			return nil
		}
	}

	// Persist to database first
	if err := s.db.Save(reportType, report, userAgent, hashStr); err != nil {
		return err
	}

	// After successful DB save, populate cache if enabled
	if s.cacheEnabled {
		fmt.Printf("caching report, hash: %s, type: %s\n", hashStr, reportType)

		b, err := json.Marshal(report)
		if err != nil {
			return err
		}
		if err := s.cache.Set(hashStr, b, time.Hour); err != nil {
			return err
		}
	}

	return nil
}
