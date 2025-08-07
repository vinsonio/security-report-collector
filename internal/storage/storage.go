package storage

import "github.com/vinsonio/security-report-collector/internal/types"

// Store is the interface for a report storage.
type Store interface {
	Save(reportType string, report types.Report, userAgent, hash string) error
}