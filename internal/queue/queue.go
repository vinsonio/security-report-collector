package queue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vinsonio/security-report-collector/internal/types"
)

// ReportEnvelope contains all data needed to persist a report to the database.
type ReportEnvelope struct {
	Type      string       `json:"type"`
	UserAgent string       `json:"user_agent"`
	Hash      string       `json:"hash"`
	Report    types.Report `json:"report"`
	Timestamp time.Time    `json:"timestamp"`
}

// Queue is the interface for a report queue.
type Queue interface {
	// Enqueue adds a report envelope to the queue.
	Enqueue(envelope *ReportEnvelope) error
	// DequeueN retrieves and removes up to n envelopes from the queue.
	DequeueN(n int) ([]*ReportEnvelope, error)
	// Size returns the approximate number of items in the queue.
	Size() (int, error)
	// Contains checks if a hash exists in the queue (for deduplication).
	Contains(hash string) (bool, error)
	// Close closes the queue.
	Close() error
}

// MarshalEnvelope serializes a report envelope to JSON bytes.
func MarshalEnvelope(envelope *ReportEnvelope) ([]byte, error) {
	return json.Marshal(envelope)
}

// UnmarshalEnvelope deserializes JSON bytes to a report envelope.
func UnmarshalEnvelope(data []byte) (*ReportEnvelope, error) {
	// First unmarshal into an alias that keeps report as raw JSON
	var alias struct {
		Type      string          `json:"type"`
		UserAgent string          `json:"user_agent"`
		Hash      string          `json:"hash"`
		Report    json.RawMessage `json:"report"`
		Timestamp time.Time       `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &alias); err != nil {
		return nil, err
	}

	// Decode the concrete report based on the type
	var rep types.Report
	switch alias.Type {
	case "csp":
		var r types.CSPReport
		if err := json.Unmarshal(alias.Report, &r); err != nil {
			return nil, err
		}
		rep = r
	default:
		return nil, fmt.Errorf("unsupported report type for envelope unmarshal: %s", alias.Type)
	}

	return &ReportEnvelope{
		Type:      alias.Type,
		UserAgent: alias.UserAgent,
		Hash:      alias.Hash,
		Report:    rep,
		Timestamp: alias.Timestamp,
	}, nil
}
