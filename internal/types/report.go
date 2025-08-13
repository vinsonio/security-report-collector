package types

// Report is an interface that all report types must implement.
type Report interface {
	// Type returns the type of the report (e.g., "csp", "hsts").
	Type() string
	// JSON returns the JSON representation of the report.
	JSON() ([]byte, error)
	// HashData returns the data used to generate the report's hash.
	HashData() (interface{}, error)
}
