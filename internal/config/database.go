package config

// DB holds the database configuration.
type DB struct {
	Connection string
	SQLite     SQLite
	MySQL      MySQL
}

// SQLite holds the SQLite database configuration.
type SQLite struct {
	Database string
}

// MySQL holds the MySQL database configuration.
type MySQL struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// NewDB creates a new DB configuration.
func NewDB() *DB {
	return &DB{
		Connection: getEnv("DB_CONNECTION", "sqlite"),
		SQLite: SQLite{
			Database: getEnv("DB_DATABASE", "reports.db"),
		},
		MySQL: MySQL{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_DATABASE", "reports"),
		},
	}
}
