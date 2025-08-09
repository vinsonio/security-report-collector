package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vinsonio/security-report-collector/internal/cache"
	"github.com/vinsonio/security-report-collector/internal/database"
)

func TestBuildRouter_Succeeds(t *testing.T) {
	// Reset singletons to ensure a fresh init
	database.ResetSingletonForTest()
	cache.ResetSingletonForTest()

	os.Setenv("DB_DATABASE", t.TempDir()+"/srv.db")
	t.Cleanup(func(){ os.Unsetenv("DB_DATABASE") })

	r, err := buildRouter()
	require.NoError(t, err)
	require.NotNil(t, r)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestBuildRouter_InitFailure(t *testing.T) {
	// Reset singletons so that invalid driver is re-evaluated
	database.ResetSingletonForTest()
	cache.ResetSingletonForTest()

	os.Setenv("DB_CONNECTION", "invalid")
	t.Cleanup(func(){ os.Unsetenv("DB_CONNECTION") })

	_, err := buildRouter()
	require.Error(t, err)
}