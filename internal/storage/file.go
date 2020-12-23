package storage

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
)

func init() {
	const storeType string = "file"
	Register(storeType, newFileStore)
}

const bucketName string = "root"

// ErrMissingFilepath happens when the filepath parameter is not passed to Open()
var ErrMissingFilepath error = errors.New("cannot find parameter \"path\"")

// IntegratedKV is an embedded KV store. It stores data in a file.
type file struct {
	config fileConfig
	db     *bolt.DB
}

type fileConfig struct {
	StoreType string `mapstructure:"type"`
	Path      string
}

func newFileStore(v *viper.Viper) (Store, error) {
	store := &file{config: fileConfig{}}
	err := v.Unmarshal(&store.config)
	if err != nil {
		panic(fmt.Errorf("cannot read storage config"))
	}
	store.Open()
	return store, err
}

// Open opens the integrated K/V storage
func (store *file) Open() error {
	if store.db != nil {
		return fmt.Errorf("the store has already been opened")
	}

	path := store.config.Path
	if path == "" {
		panic(fmt.Errorf("cannot open the store with an empty filepath db"))
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(fmt.Errorf("cannot open the store: %v", err))
	}

	store.db = db

	store.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}

// Close closes the integrated K/V storage
func (store *file) Close() error {
	return store.db.Close()
}

// ReadRedirection reads a redirection from the integrated K/V storage
func (store *file) ReadRedirection(key string) (string, error) {
	var result []byte
	store.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		result = b.Get([]byte(key))
		return nil
	})
	var err error
	if result == nil {
		err = ErrRedirectionNotfound
	}
	return string(result), err
}

// WriteRedirection writes a new redirection into the integrated K/V storage
func (store *file) WriteRedirection(key string, target string) (err error) {
	store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(target))
		return err
	})
	return err
}

// DeleteRedirection deletes a redirection from integrated K/V storage
func (store *file) DeleteRedirection(key string) (err error) {
	store.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(key))
		return err
	})
	return err
}
