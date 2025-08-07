package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/vinsonio/security-report-collector/internal/service"
	storagetesting "github.com/vinsonio/security-report-collector/internal/testing"
	"github.com/vinsonio/security-report-collector/internal/types"
	"github.com/vinsonio/security-report-collector/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/vinsonio/security-report-collector/internal/handler"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handler.HealthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestReportMux(t *testing.T) {
	t.Run("handles valid report", func(t *testing.T) {
		store := new(storagetesting.MockStore)
		reportService := service.NewReportService(store)

		jsonStr := []byte(`{
			"type": "csp-violation",
			"url": "https://example.com",
			"body": {
				"documentURL": "https://example.com",
				"disposition": "report",
				"referrer": "",
				"effectiveDirective": "font-src",
				"blockedURL": "https://fonts.gstatic.com/s/inter/v19/UcC73FwrK3iLTeHuS_nVMrMxCp50SjIa2ZL7W0Q5n-wU.woff2",
				"originalPolicy": "default-src 'self'; script-src 'self'; img-src: 'self'; report-to csp",
				"statusCode": 200,
				"sample": "",
				"sourceFile": "https://example.com",
				"lineNumber": 0,
				"columnNumber": 1
			}
		}`)
		req, err := http.NewRequest("POST", "/reports/csp", bytes.NewBuffer(jsonStr))
		assert.NoError(t, err)

		rr := httptest.NewRecorder()

		reportHandlers := map[string]handler.ReportHandler{
			"csp": &handler.CSPReportHandler{},
		}

		report := types.CSPReport{
			ReportType: "csp-violation",
			URL:        "https://example.com",
			Body: types.CSPReportBody{
				DocumentURL:        "https://example.com",
				Disposition:        "report",
				Referrer:           "",
				EffectiveDirective: "font-src",
				BlockedURL:         "https://fonts.gstatic.com/s/inter/v19/UcC73FwrK3iLTeHuS_nVMrMxCp50SjIa2ZL7W0Q5n-wU.woff2",
				OriginalPolicy:     "default-src 'self'; script-src 'self'; img-src: 'self'; report-to csp",
				StatusCode:         200,
				Sample:             "",
				SourceFile:         "https://example.com",
				LineNumber:         0,
				ColumnNumber:       1,
			},
		}
		hashData := types.CSPReportHashData{
			DocumentURL:        report.Body.DocumentURL,
			EffectiveDirective: report.Body.EffectiveDirective,
			BlockedURL:         report.Body.BlockedURL,
			SourceFile:         report.Body.SourceFile,
			LineNumber:         report.Body.LineNumber,
			ColumnNumber:       report.Body.ColumnNumber,
		}

		data, err := util.StableMarshal(hashData)
		assert.NoError(t, err)
		hashBytes := sha256.Sum256(data)
		hash := hex.EncodeToString(hashBytes[:])

		store.On("Save", "csp", &report, "", hash).Return(nil)

		mux := handler.ReportMux(reportService, reportHandlers)
		mux.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		store.AssertExpectations(t)
	})
}