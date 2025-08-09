package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware_NoAllowedDomains(t *testing.T) {
	os.Unsetenv("ALLOWED_DOMAINS")

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORSMiddleware_MissingHeaders(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_DOMAINS") })

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCORSMiddleware_InvalidOrigin(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_DOMAINS") })

	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("Origin", "http://%zz")
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCORSMiddleware_ForbiddenOrigin(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_DOMAINS") })

	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("Origin", "http://notallowed.com")
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCORSMiddleware_AllowedExactMatch(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_DOMAINS") })

	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("Origin", "http://example.com")
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCORSMiddleware_AllowedWildcardSubdomain(t *testing.T) {
	os.Setenv("ALLOWED_DOMAINS", "*.example.com")
	t.Cleanup(func() { os.Unsetenv("ALLOWED_DOMAINS") })

	req := httptest.NewRequest("GET", "/healthz", nil)
	req.Header.Set("Origin", "http://api.example.com")
	rec := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Use(CORSMiddleware)
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}