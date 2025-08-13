package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vinsonio/security-report-collector/internal/database"
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

// CreateReport returns a new http.Handler for routing reports.
func CreateReport(reportService *service.ReportService, handlers map[string]ReportHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportType := chi.URLParam(r, "type")

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
		if err := reportService.SaveReport(reportType, report, userAgent); err != nil {
			if err == database.ErrDuplicateReport {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
