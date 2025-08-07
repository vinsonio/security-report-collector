package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/service"
	"github.com/vinsonio/security-report-collector/internal/storage"
)

func main() {
	godotenv.Load()

	store, err := storage.NewStoreFromEnv()
	if err != nil {
		log.Fatalf("failed to initialize storage: %v", err)
	}

	reportService := service.NewReportService(store)

	http.HandleFunc("/healthz", handler.HealthCheck)

	reportHandlers := map[string]handler.ReportHandler{
		"csp": &handler.CSPReportHandler{},
	}

	reportMux := handler.ReportMux(reportService, reportHandlers)
	http.Handle("/reports/", handler.CORSMiddleware(http.HandlerFunc(reportMux)))

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
