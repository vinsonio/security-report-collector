package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSPReport_Type(t *testing.T) {
	report := CSPReport{}
	assert.Equal(t, "csp", report.Type())
}

func TestCSPReport_JSON(t *testing.T) {
	report := CSPReport{
		URL:        "https://example.com",
		ReportType: "csp-violation",
		Body: CSPReportBody{
			DocumentURL:        "https://example.com/page.html",
			EffectiveDirective: "script-src",
			BlockedURL:         "https://malicious.com/script.js",
		},
	}

	b, err := report.JSON()
	assert.NoError(t, err)
	assert.True(t, json.Valid(b))

	var unmarshaled CSPReport
	err = json.Unmarshal(b, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, report.URL, unmarshaled.URL)
	assert.Equal(t, report.ReportType, unmarshaled.ReportType)
	assert.Equal(t, report.Body.DocumentURL, unmarshaled.Body.DocumentURL)
}

func TestCSPReport_HashData(t *testing.T) {
	report := CSPReport{
		Body: CSPReportBody{
			DocumentURL:        "https://example.com/page.html",
			EffectiveDirective: "script-src",
			BlockedURL:         "https://malicious.com/script.js",
			SourceFile:         "https://example.com/app.js",
			LineNumber:         42,
			ColumnNumber:       10,
			// Extra fields not included in hash
			Disposition:    "enforce",
			Referrer:       "https://example.com",
			OriginalPolicy: "default-src 'self'",
			StatusCode:     200,
			Sample:         "sample code",
		},
	}

	hashData, err := report.HashData()
	assert.NoError(t, err)

	expected := CSPReportHashData{
		DocumentURL:        "https://example.com/page.html",
		EffectiveDirective: "script-src",
		BlockedURL:         "https://malicious.com/script.js",
		SourceFile:         "https://example.com/app.js",
		LineNumber:         42,
		ColumnNumber:       10,
	}

	assert.Equal(t, expected, hashData)
}

func TestCSPReport_EmptyFields(t *testing.T) {
	report := CSPReport{}

	// Test Type
	assert.Equal(t, "csp", report.Type())

	// Test JSON with empty fields
	b, err := report.JSON()
	assert.NoError(t, err)
	assert.True(t, json.Valid(b))

	// Test HashData with empty fields
	hashData, err := report.HashData()
	assert.NoError(t, err)
	expected := CSPReportHashData{}
	assert.Equal(t, expected, hashData)
}
