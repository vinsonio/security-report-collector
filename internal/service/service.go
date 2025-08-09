package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/database"
	"github.com/vinsonio/security-report-collector/internal/types"
	"github.com/vinsonio/security-report-collector/internal/util"
)

// ReportService is the service for handling reports.
type ReportService struct {
	db           database.DB
	cache        cache.Cache
	cacheEnabled bool
}

// NewReportService creates a new ReportService.
func NewReportService(db database.DB, cache cache.Cache, cacheEnabled bool) *ReportService {
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

	if s.cacheEnabled {
		b, err := json.Marshal(report)
		if err != nil {
			return err
		}
		if err := s.cache.Set(hashStr, b, time.Hour); err != nil {
			return err
		}
	}

	return s.db.Save(reportType, report, userAgent, hashStr)
}
