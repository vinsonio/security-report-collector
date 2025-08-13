package database

import "fmt"

// DBBuilder is a function that creates a new database backend.
type DBBuilder func() (DB, error)

var dbBuilders = make(map[string]DBBuilder)

// RegisterDB registers a new database backend.
func RegisterDB(name string, builder DBBuilder) {
	dbBuilders[name] = builder
}

// GetDBBuilder returns a database builder by name.
func GetDBBuilder(name string) (DBBuilder, error) {
	builder, ok := dbBuilders[name]
	if !ok {
		return nil, fmt.Errorf("unsupported database type: %s", name)
	}
	return builder, nil
}
