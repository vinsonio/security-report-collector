package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vinsonio/security-report-collector/internal/handler"
	"github.com/vinsonio/security-report-collector/internal/service"
	databasetesting "github.com/vinsonio/security-report-collector/internal/testing"
	cachetesting "github.com/vinsonio/security-report-collector/internal/testing/cache"
	"github.com/vinsonio/security-report-collector/internal/types"
)

// simple handler implementation that returns a fixed payload
type okHandler struct{}

func (okHandler) Handle(r *http.Request) (types.Report, error) {
	return &mockReport{}, nil
}

type mockReport struct{}

func (m *mockReport) Type() string                   { return "csp" }
func (m *mockReport) JSON() ([]byte, error)          { return []byte("{}"), nil }
func (m *mockReport) HashData() (interface{}, error) { return map[string]string{"k": "v"}, nil }

func newTestServer(t *testing.T) http.Handler {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	svc := service.NewReportService(store, cache, false)

	handlers := map[string]handler.ReportHandler{
		"csp": okHandler{},
	}
	return New(svc, handlers)
}

func TestRouter_Healthz(t *testing.T) {
	r := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_CreateReport_CORSAllowed(t *testing.T) {
	t.Setenv("ALLOWED_DOMAINS", "example.com,sub.example.org,*.wild.test")

	// Build router with expectations on DB Save
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	store.On("Save", "csp", mock.Anything, mock.Anything, mock.AnythingOfType("string")).Return(nil)
	svc := service.NewReportService(store, cache, false)
	mux := New(svc, map[string]handler.ReportHandler{"csp": okHandler{}})

	req := httptest.NewRequest(http.MethodPost, "/reports/csp", strings.NewReader("{}"))
	req.Header.Set("Origin", "https://sub.example.org")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	store.AssertExpectations(t)
}

func TestRouter_CreateReport_CORSForbidden(t *testing.T) {
	t.Setenv("ALLOWED_DOMAINS", "example.com,*.allowed.test")

	r := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/reports/csp", strings.NewReader("{}"))
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRouter_CreateReport_MissingHandler(t *testing.T) {
	// No handler registered for type "x"
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	svc := service.NewReportService(store, cache, false)
	mux := New(svc, map[string]handler.ReportHandler{})

	req := httptest.NewRequest(http.MethodPost, "/reports/x", strings.NewReader("{}"))
	req.Header.Set("Origin", "https://example.com")
	t.Setenv("ALLOWED_DOMAINS", "example.com")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}
