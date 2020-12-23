package storage

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// Store is an interface representing the possible interactions with a storage.
// When working with a Store, you should always open it before use and close it after use.
type Store interface {
	// Open initialize the store. (Could be opening a file or initiating a connection to a remote database)
	Open() error
	// Close closes the store. (Could be closing a file or closing a remote connection.)
	Close() error
	// ReadRedirection takes a redirection key and returns the redirection's target.
	ReadRedirection(key string) (target string, err error)
	// WriteRedirection writes a redirection to the store. Parameters are the key and the target.
	WriteRedirection(key string, target string) error
	// DeleteRedirection deletes a redirection from the store. it Takes the key as argument.
	DeleteRedirection(key string) error
}

// StoreFactory if a factory function to create a new store.
type StoreFactory func(v *viper.Viper) (Store, error)

var availableStores map[string]StoreFactory = make(map[string]StoreFactory)

var (
	ErrRedirectionNotfound error = errors.New("redirection not found")
)

// Register will register a new store type with its associated factory function.
func Register(name string, factory StoreFactory) error {
	_, ok := availableStores[name]
	if ok {
		return fmt.Errorf("could not register new store type \"%s\" because it already exists", name)
	}
	availableStores[name] = factory
	return nil
}

// GetFactory returns the factory function associated with a store type.
func GetFactory(name string) (StoreFactory, error) {
	factory, exists := availableStores[name]
	if !exists {
		return nil, fmt.Errorf("could not register new store type \"%s\" because it already exists", name)
	}
	return factory, nil
}
