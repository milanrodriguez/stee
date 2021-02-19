package stee

import (
	"errors"
	"math/rand"

	"github.com/milanrodriguez/stee/internal/storage"
)

type redirections map[string]string

var (
	// ErrRedirectionNotfound is used when the redirection was not found
	ErrRedirectionNotfound = errors.New("no redirection found for this key")
	// ErrRedirectionAlreadyExists is used when the redirection already exists
	ErrRedirectionAlreadyExists = errors.New("a redirection is already associated with this key")
	// ErrTargetIsNotAValidURL is used when the provided target is not a valid URL
	ErrTargetIsNotAValidURL = errors.New("the redirection target is not a valid URL")
)

// GetRedirection gets a redirection based on its key
func (c *Core) GetRedirection(key string) (target string, err error) {
	target, err = c.store.ReadRedirection(key)
	return
}

// AddRedirectionWithKey adds a redirection. It takes both the key and the target of the redirection.
func (c *Core) AddRedirectionWithKey(key string, target string) (err error) {
	exists, err := c.keyExists(key)
	if err != nil {
		return err
	}
	if exists {
		return ErrRedirectionAlreadyExists
	}
	err = c.store.WriteRedirection(key, target)
	return
}

// AddRedirectionWithoutKey adds a redirection. It takes only the target and generate a key.
func (c *Core) AddRedirectionWithoutKey(target string) (key string, err error) {
	var exists bool
	for {
		key = generateRedirectionKey()
		exists, err = c.keyExists(key)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
	}
	err = c.store.WriteRedirection(key, target)
	return key, err
}

// DeleteRedirection deletes a redirection based on its key.
func (c *Core) DeleteRedirection(key string) (err error) {
	err = c.store.DeleteRedirection(key)
	return err
}

func generateRedirectionKey() string {
	const base64CharSet string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	const keylength = 6
	var charSet string = base64CharSet

	var key string
	for i := 0; i < keylength; i++ {
		key = key + string(charSet[rand.Intn(len(charSet))])
	}

	return key
}

func (c *Core) keyExists(key string) (exists bool, err error) {
	target, err := c.store.ReadRedirection(key)
	if err != nil && !errors.Is(err, storage.ErrRedirectionNotfound) {
		return false, err
	}
	exists = target != ""
	return exists, nil
}
