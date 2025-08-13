package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/service"
)

func New(reportService *service.ReportService, reportHandlers map[string]handler.ReportHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", handler.HealthCheck)

	r.Group(func(r chi.Router) {
		r.Use(CORSMiddleware)
		r.Post("/reports/{type}", handler.CreateReport(reportService, reportHandlers))
	})

	return r
}
