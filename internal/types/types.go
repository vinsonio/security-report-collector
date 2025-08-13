package types

import "encoding/json"

// CSPReport is a wrapper for the CSP report body
type CSPReport struct {
	URL        string        `json:"url,omitempty"`
	ReportType string        `json:"type,omitempty"`
	Body       CSPReportBody `json:"body"`
}

// Type returns the type of the report.
func (r CSPReport) Type() string {
	return "csp"
}

// JSON returns the JSON representation of the report.
func (r CSPReport) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// HashData returns the data used to generate the report's hash.
func (r CSPReport) HashData() (interface{}, error) {
	return CSPReportHashData{
		DocumentURL:        r.Body.DocumentURL,
		EffectiveDirective: r.Body.EffectiveDirective,
		BlockedURL:         r.Body.BlockedURL,
		SourceFile:         r.Body.SourceFile,
		LineNumber:         r.Body.LineNumber,
		ColumnNumber:       r.Body.ColumnNumber,
	}, nil
}

// CSPReportHashData defines the structure of the data used to generate the report hash.
// This is a subset of the full CSP report, containing only the fields that uniquely
// identify a specific violation.
// This helps in deduplicating reports and tracking the frequency of specific issues.
type CSPReportHashData struct {
	DocumentURL        string `json:"documentURL,omitempty"`
	EffectiveDirective string `json:"effectiveDirective,omitempty"`
	BlockedURL         string `json:"blockedURL,omitempty"`
	SourceFile         string `json:"sourceFile,omitempty"`
	LineNumber         int    `json:"lineNumber,omitempty"`
	ColumnNumber       int    `json:"columnNumber,omitempty"`
}

// CSPReportBody defines the structure of a CSP report body.
type CSPReportBody struct {
	DocumentURL        string `json:"documentURL,omitempty"`
	Disposition        string `json:"disposition,omitempty"`
	Referrer           string `json:"referrer,omitempty"`
	EffectiveDirective string `json:"effectiveDirective,omitempty"`
	BlockedURL         string `json:"blockedURL,omitempty"`
	OriginalPolicy     string `json:"originalPolicy,omitempty"`
	StatusCode         int    `json:"statusCode,omitempty"`
	Sample             string `json:"sample,omitempty"`
	SourceFile         string `json:"sourceFile,omitempty"`
	LineNumber         int    `json:"lineNumber,omitempty"`
	ColumnNumber       int    `json:"columnNumber,omitempty"`
}
