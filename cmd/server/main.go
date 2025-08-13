package main

import (
	"log"
	"net/http"
	"time"

	"github.com/vinsonio/security-report-collector/internal/bootstrap"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/queue"
	"github.com/vinsonio/security-report-collector/internal/router"
	"github.com/vinsonio/security-report-collector/internal/scheduler"
	"github.com/vinsonio/security-report-collector/internal/service"
)

func buildRouter() (http.Handler, error) {
	db, cache, err := bootstrap.Init()
	if err != nil {
		return nil, err
	}

	appConfig := config.NewApp()
	reportService := service.NewReportService(db, cache, appConfig.CacheEnabled)

	return buildRouterWithService(reportService)
}

// buildRouterWithService constructs the HTTP router using the provided service.
func buildRouterWithService(reportService *service.ReportService) (http.Handler, error) {
	reportHandlers := map[string]handler.ReportHandler{
		"csp": &handler.CSPReportHandler{},
	}

	r := router.New(reportService, reportHandlers)
	return r, nil
}

func main() {
	// Application bootstrap
	db, cache, err := bootstrap.Init()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	appConfig := config.NewApp()
	reportService := service.NewReportService(db, cache, appConfig.CacheEnabled)

	// Start background flusher (queue + scheduler) as part of app lifecycle, not router construction
	if appConfig.CacheEnabled {
		cacheCfg := config.NewCache()
		q, err := queue.New(cacheCfg, "reports")
		if err != nil {
			log.Fatalf("failed to initialize queue: %v", err)
		}
		reportService.AttachQueue(q)

		flusher := scheduler.NewBatchFlusher(q, db, appConfig.BatchSize)
		stop := make(chan struct{})
		interval := time.Duration(appConfig.FlushIntervalMinutes) * time.Minute
		go scheduler.Scheduler(interval, stop, flusher.Flush)
		log.Printf("Batch flusher scheduler started (interval: %v, batchSize: %d)", interval, appConfig.BatchSize)
	}

	r, err := buildRouterWithService(reportService)
	if err != nil {
		log.Fatalf("failed to build router: %v", err)
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
