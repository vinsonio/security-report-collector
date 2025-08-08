package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/service"
	databasetesting "github.com/vinsonio/security-report-collector/internal/testing"
)

func TestReportMux(t *testing.T) {
	t.Run("handles valid report", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		reportService := service.NewReportService(store)

		jsonStr := []byte(`{"body":{}}`)
		req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		store.On("Save", "csp", mock.Anything, "", mock.AnythingOfType("string")).Return(nil)

		mux := handler.ReportMux(reportService, reportHandlers)
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		store.AssertExpectations(t)
	})

	t.Run("handles unknown report type", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		reportService := service.NewReportService(store)

		jsonStr := []byte(`{"type":"unknown"}`)
		req, err := http.NewRequest("POST", "/reports/unknown", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		mux := handler.ReportMux(reportService, reportHandlers)
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		store.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("handles invalid json", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		reportService := service.NewReportService(store)

		jsonStr := []byte(`invalid-json`)
		req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		mux := handler.ReportMux(reportService, reportHandlers)
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		store.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}