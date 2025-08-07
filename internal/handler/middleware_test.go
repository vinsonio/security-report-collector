package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		allowedDomains string
		origin         string
		referer        string
		expectedStatus int
	}{
		{"Allowed Origin", "example.com", "http://example.com", "", http.StatusOK},
		{"Allowed Referer", "example.com", "", "http://example.com", http.StatusOK},
		{"Disallowed Origin", "example.com", "http://disallowed.com", "", http.StatusForbidden},
		{"No Origin or Referer", "example.com", "", "", http.StatusBadRequest},
		{"Empty Allowed Domains", "", "http://example.com", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("ALLOWED_DOMAINS", tt.allowedDomains)
			defer os.Unsetenv("ALLOWED_DOMAINS")

			req, err := http.NewRequest("POST", "/csp-report", nil)
			if err != nil {
				t.Fatal(err)
			}
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.referer != "" {
				req.Header.Set("Referer", tt.referer)
			}

			rr := httptest.NewRecorder()
			handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}
}