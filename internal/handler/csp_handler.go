package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vinsonio/security-report-collector/internal/types"
)

// CSPReportHandler handles CSP violation reports.
type CSPReportHandler struct{}

// Handle decodes a CSP report from the request body.
func (h *CSPReportHandler) Handle(r *http.Request) (types.Report, error) {
	var report types.CSPReport
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		return nil, err
	}
	return &report, nil
}
