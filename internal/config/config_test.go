package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Defaults(t *testing.T) {
	// Explicitly set defaults that New() would use when env is unset
	t.Setenv("APP_NAME", "report-collector")
	t.Setenv("DB_CONNECTION", "sqlite")

	cfg := New()
	assert.NotNil(t, cfg.App)
	assert.NotNil(t, cfg.Cache)
	assert.NotNil(t, cfg.DB)
	assert.Equal(t, "sqlite", cfg.DB.Connection)
}

func TestNewDB_FromEnv(t *testing.T) {
	t.Setenv("DB_CONNECTION", "mysql")
	t.Setenv("DB_DATABASE", "urc")
	t.Setenv("DB_HOST", "db.example")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_USER", "user")
	t.Setenv("DB_PASSWORD", "pass")

	db := NewDB()
	assert.Equal(t, "mysql", db.Connection)
	assert.Equal(t, "urc", db.MySQL.Database)
	assert.Equal(t, "db.example", db.MySQL.Host)
	assert.Equal(t, 3307, db.MySQL.Port)
	assert.Equal(t, "user", db.MySQL.User)
	assert.Equal(t, "pass", db.MySQL.Password)
}
