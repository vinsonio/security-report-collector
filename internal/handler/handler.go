package handler

import (
	"strings"
	"net/http"
)

import (
	"github.com/vinsonio/security-report-collector/internal/service"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// ReportHandler defines the interface for handling a specific type of report.
type ReportHandler interface {
	Handle(r *http.Request) (types.Report, error)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ReportMux returns a new http.Handler for routing reports.
func ReportMux(reportService *service.ReportService, handlers map[string]ReportHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Extract report type from URL path, e.g., /reports/csp
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reportType := pathParts[1]

		handler, ok := handlers[reportType]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		report, err := handler.Handle(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userAgent := r.Header.Get("User-Agent")
		if err := reportService.SaveReport(report, userAgent); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}