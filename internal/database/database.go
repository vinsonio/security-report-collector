package database

import (
	"errors"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// ErrDuplicateReport is returned when a report with the same hash already exists.
var ErrDuplicateReport = errors.New("duplicate report")

// DB is the interface for a report database.
type DB interface {
	Save(reportType string, report types.Report, userAgent, hash string) error
	Migrate() error
}
