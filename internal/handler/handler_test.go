package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vinsonio/security-report-collector/internal/database"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/service"
	databasetesting "github.com/vinsonio/security-report-collector/internal/testing"
	cachetesting "github.com/vinsonio/security-report-collector/internal/testing/cache"
)

func TestCreateReport_DuplicateHandled(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	reportService := service.NewReportService(store, cache, false)

	jsonStr := []byte(`{"csp-report":{}}`)
	req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	req.Header.Set("User-Agent", "test-agent")

	rr := httptest.NewRecorder()

	reportHandlers := map[string]handler.ReportHandler{
		"csp": &handler.CSPReportHandler{},
	}

	store.On("Save", "csp", mock.AnythingOfType("*types.CSPReport"), "test-agent", mock.AnythingOfType("string")).Return(database.ErrDuplicateReport)

	router := chi.NewRouter()
	router.Post("/reports/{type}", handler.CreateReport(reportService, reportHandlers))
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	store.AssertExpectations(t)
}

func TestCreateReport(t *testing.T) {
	t.Run("handles valid report", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		cache := new(cachetesting.MockCache)
		reportService := service.NewReportService(store, cache, false)

		jsonStr := []byte(`{"csp-report":{}}`)
		req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)
		req.Header.Set("User-Agent", "test-agent")

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		store.On("Save", "csp", mock.AnythingOfType("*types.CSPReport"), "test-agent", mock.AnythingOfType("string")).Return(nil)

		router := chi.NewRouter()
		router.Post("/reports/{type}", handler.CreateReport(reportService, reportHandlers))
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		store.AssertExpectations(t)
	})

	t.Run("handles unknown report type", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		cache := new(cachetesting.MockCache)
		reportService := service.NewReportService(store, cache, false)

		jsonStr := []byte(`{"type":"unknown"}`)
		req, err := http.NewRequest("POST", "/reports/unknown", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		router := chi.NewRouter()
		router.Post("/reports/{type}", handler.CreateReport(reportService, reportHandlers))
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("handles invalid json", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		cache := new(cachetesting.MockCache)
		reportService := service.NewReportService(store, cache, false)

		jsonStr := []byte(`invalid-json`)
		req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		router := chi.NewRouter()
		router.Post("/reports/{type}", handler.CreateReport(reportService, reportHandlers))
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
