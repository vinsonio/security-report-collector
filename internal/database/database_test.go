package database_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vinsonio/security-report-collector/internal/database"
	dbtesting "github.com/vinsonio/security-report-collector/internal/testing"
	"github.com/vinsonio/security-report-collector/internal/types"
	"log"
)

func TestMain(m *testing.M) {
	// Set the environment variable to use the in-memory SQLite database for tests
	if err := os.Setenv("DB_CONNECTION", "sqlite"); err != nil {
		log.Fatal(err)
	}
	if err := os.Setenv("DB_URL", "file:test.db?mode=memory&cache=shared"); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestSaveDuplicateReport(t *testing.T) {
	db := dbtesting.GetDBForTest(t)

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

	hash := "d1692b293b40495a372cf2473551125d5635393da55b6942647b013b0c2a2a59"

	// Save the report for the first time
	err := db.Save("csp", report, "test-agent", hash)
	assert.NoError(t, err)
	assert.Equal(t, 1, db.Count(t), "Report count should be 1 after first save")

	// Save the same report again
	err = db.Save("csp", report, "test-agent", hash)
	assert.Equal(t, database.ErrDuplicateReport, err)
	assert.Equal(t, 1, db.Count(t), "Report count should still be 1 after saving a duplicate")
}
