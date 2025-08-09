package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Defaults(t *testing.T) {
	// Ensure unrelated env vars do not affect defaults
	os.Unsetenv("APP_NAME")
	os.Unsetenv("DB_CONNECTION")

	cfg := New()
	assert.NotNil(t, cfg.App)
	assert.NotNil(t, cfg.Cache)
	assert.NotNil(t, cfg.DB)
	assert.Equal(t, "sqlite", cfg.DB.Connection)
}

func TestNewDB_FromEnv(t *testing.T) {
	os.Setenv("DB_CONNECTION", "mysql")
	os.Setenv("DB_DATABASE", "urc")
	os.Setenv("DB_HOST", "db.example")
	os.Setenv("DB_PORT", "3307")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	t.Cleanup(func() {
		os.Unsetenv("DB_CONNECTION")
		os.Unsetenv("DB_DATABASE")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
	})

	db := NewDB()
	assert.Equal(t, "mysql", db.Connection)
	assert.Equal(t, "urc", db.MySQL.Database)
	assert.Equal(t, "db.example", db.MySQL.Host)
	assert.Equal(t, 3307, db.MySQL.Port)
	assert.Equal(t, "user", db.MySQL.User)
	assert.Equal(t, "pass", db.MySQL.Password)
}