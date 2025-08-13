package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	databasetesting "github.com/vinsonio/security-report-collector/internal/testing"
	cachetesting "github.com/vinsonio/security-report-collector/internal/testing/cache"
	"github.com/vinsonio/security-report-collector/internal/types"
)

func TestSaveReport_CacheHitSkipsDB(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	service := NewReportService(store, cache, true)

	report := types.CSPReport{Body: types.CSPReportBody{DocumentURL: "https://example.com"}}

	// Arrange cache hit
	cache.On("Get", mock.AnythingOfType("string")).Return([]byte("1"), nil)

	err := service.SaveReport("csp", report, "UA")
	assert.NoError(t, err)

	// DB should not be called
	store.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestSaveReport_CacheMissSavesAndSets(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	service := NewReportService(store, cache, true)

	report := types.CSPReport{Body: types.CSPReportBody{DocumentURL: "https://example.com"}}

	cache.On("Get", mock.AnythingOfType("string")).Return([]byte(nil), nil)
	store.On("Save", "csp", mock.AnythingOfType("types.CSPReport"), "UA", mock.AnythingOfType("string")).Return(nil)
	cache.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)

	err := service.SaveReport("csp", report, "UA")
	assert.NoError(t, err)
	store.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestSaveReport_CacheDisabled_NoGet(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	service := NewReportService(store, cache, false)

	report := types.CSPReport{Body: types.CSPReportBody{DocumentURL: "https://example.com"}}

	// Expect only Save
	store.On("Save", "csp", mock.AnythingOfType("types.CSPReport"), "UA", mock.AnythingOfType("string")).Return(nil)

	err := service.SaveReport("csp", report, "UA")
	assert.NoError(t, err)
	store.AssertExpectations(t)
	cache.AssertNotCalled(t, "Get", mock.Anything)
}

func TestSaveReport_CacheGetError(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	service := NewReportService(store, cache, true)

	report := types.CSPReport{Body: types.CSPReportBody{DocumentURL: "https://example.com"}}

	// Cache Get returns error
	cache.On("Get", mock.AnythingOfType("string")).Return([]byte(nil), errors.New("cache error"))

	err := service.SaveReport("csp", report, "UA")
	assert.Error(t, err)
	assert.Equal(t, "cache error", err.Error())

	// DB and Set should not be called
	store.AssertNotCalled(t, "Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestSaveReport_CacheSetError(t *testing.T) {
	store := new(databasetesting.MockDB)
	cache := new(cachetesting.MockCache)
	service := NewReportService(store, cache, true)

	report := types.CSPReport{Body: types.CSPReportBody{DocumentURL: "https://example.com"}}

	cache.On("Get", mock.AnythingOfType("string")).Return([]byte(nil), nil)
	store.On("Save", "csp", mock.AnythingOfType("types.CSPReport"), "UA", mock.AnythingOfType("string")).Return(nil)
	cache.On("Set", mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(errors.New("cache set error"))

	err := service.SaveReport("csp", report, "UA")
	assert.Error(t, err)
	assert.Equal(t, "cache set error", err.Error())
	store.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestReport_JSONAndHashData(t *testing.T) {
	report := types.CSPReport{
		URL:        "https://example.com",
		ReportType: "csp-violation",
		Body: types.CSPReportBody{
			DocumentURL:        "https://example.com",
			EffectiveDirective: "script-src",
			BlockedURL:         "https://cdn.example.com/script.js",
			SourceFile:         "https://example.com/app.js",
			LineNumber:         10,
			ColumnNumber:       20,
		},
	}

	b, err := report.JSON()
	assert.NoError(t, err)
	assert.True(t, json.Valid(b))

	h, err := report.HashData()
	assert.NoError(t, err)
	m, ok := h.(types.CSPReportHashData)
	assert.True(t, ok)
	assert.Equal(t, "https://example.com", m.DocumentURL)
	assert.Equal(t, "script-src", m.EffectiveDirective)
	assert.Equal(t, "https://cdn.example.com/script.js", m.BlockedURL)
	assert.Equal(t, "https://example.com/app.js", m.SourceFile)
	assert.Equal(t, 10, m.LineNumber)
	assert.Equal(t, 20, m.ColumnNumber)
}
