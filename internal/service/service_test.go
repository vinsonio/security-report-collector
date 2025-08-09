package service_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vinsonio/security-report-collector/internal/service"
	databasetesting "github.com/vinsonio/security-report-collector/internal/testing"
	testingcache "github.com/vinsonio/security-report-collector/internal/testing/cache"
	"github.com/vinsonio/security-report-collector/internal/types"
	"github.com/vinsonio/security-report-collector/internal/util"
)

func TestReportService_SaveReport(t *testing.T) {
	report := &types.CSPReport{
		Body: types.CSPReportBody{
			DocumentURL:        "http://example.com/signup.html",
			Referrer:           "",
			BlockedURL:         "http://example.com/css/style.css",
			EffectiveDirective: "style-src cdn.example.com",
			OriginalPolicy:     "default-src 'none'; style-src cdn.example.com; report-uri /_/csp-reports",
		},
	}

	hashData, err := report.HashData()
	assert.NoError(t, err)
	hashBytes, err := util.StableMarshal(hashData)
	assert.NoError(t, err)
	hashSum := sha256.Sum256(hashBytes)
	hash := hex.EncodeToString(hashSum[:])

	t.Run("cache enabled", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		cache := new(testingcache.MockCache)
		reportService := service.NewReportService(store, cache, true)
		reportJSON, err := json.Marshal(report)
		assert.NoError(t, err)

		cache.On("Set", hash, reportJSON, time.Hour).Return(nil)
		store.On("Save", "csp", report, "user-agent", hash).Return(nil)

		err = reportService.SaveReport("csp", report, "user-agent")
		assert.NoError(t, err)

		store.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("cache disabled", func(t *testing.T) {
		store := new(databasetesting.MockDB)
		cache := new(testingcache.MockCache)
		reportService := service.NewReportService(store, cache, false)

		store.On("Save", "csp", report, "user-agent", hash).Return(nil)

		err = reportService.SaveReport("csp", report, "user-agent")
		assert.NoError(t, err)

		store.AssertExpectations(t)
		cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
	})
}