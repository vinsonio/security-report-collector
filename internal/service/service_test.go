package service_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/vinsonio/security-report-collector/internal/util"

	"github.com/stretchr/testify/assert"
	"github.com/vinsonio/security-report-collector/internal/service"
	storagetesting "github.com/vinsonio/security-report-collector/internal/testing"
	"github.com/vinsonio/security-report-collector/internal/types"
)

func TestReportService_SaveReport(t *testing.T) {
	store := new(storagetesting.MockDB)
	reportService := service.NewReportService(store)

	report := types.CSPReport{
		ReportType: "csp-violation",
		URL:  "https://example.com",
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

	store.On("Save", "csp", report, "user-agent", hash).Return(nil)

	err = reportService.SaveReport(report, "user-agent")
	assert.NoError(t, err)

	store.AssertExpectations(t)
}