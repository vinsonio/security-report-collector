package storage

import "fmt"

// StoreBuilder is a function that creates a new storage backend.
type StoreBuilder func() (Store, error)

var storeBuilders = make(map[string]StoreBuilder)

// RegisterStore registers a new storage backend.
func RegisterStore(name string, builder StoreBuilder) {
	storeBuilders[name] = builder
}

// GetStoreBuilder returns a storage builder by name.
func GetStoreBuilder(name string) (StoreBuilder, error) {
	builder, ok := storeBuilders[name]
	if !ok {
		return nil, fmt.Errorf("unsupported storage type: %s", name)
	}
	return builder, nil
}