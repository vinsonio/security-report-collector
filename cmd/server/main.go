package main

import (
	"log"
	"net/http"

	"github.com/vinsonio/security-report-collector/internal/bootstrap"
	"github.com/vinsonio/security-report-collector/internal/config"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/router"
	"github.com/vinsonio/security-report-collector/internal/service"
)

func buildRouter() (http.Handler, error) {
	db, cache, err := bootstrap.Init()
	if err != nil {
		return nil, err
	}

	appConfig := config.NewApp()
	reportService := service.NewReportService(db, cache, appConfig.CacheEnabled)

	reportHandlers := map[string]handler.ReportHandler{
		"csp": &handler.CSPReportHandler{},
	}

	r := router.New(reportService, reportHandlers)
	return r, nil
}

func main() {
	r, err := buildRouter()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
