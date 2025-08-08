package service

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/vinsonio/security-report-collector/internal/database"
	"github.com/vinsonio/security-report-collector/internal/types"
	"github.com/vinsonio/security-report-collector/internal/util"
)

// ReportService is the service for handling reports.
type ReportService struct {
	db database.DB
}

// NewReportService creates a new ReportService.
func NewReportService(db database.DB) *ReportService {
	return &ReportService{db: db}
}

// SaveReport saves a report.
func (s *ReportService) SaveReport(report types.Report, userAgent string) error {
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

	return s.db.Save(report.Type(), report, userAgent, hashStr)
}
