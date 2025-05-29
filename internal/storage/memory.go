// Package storage provides interfaces and implementations for URL storage
package storage

import (
	"fmt"
	"sync"

	"github.com/acavaka/shortlinker/internal/config"
)

var ErrStorageRestore = fmt.Errorf("failed to restore storage from file")

type InMemoryStorage struct {
	urlMappings map[string]string
	mutex       *sync.RWMutex
	counter     uint64
}

type URLStorage interface {
	Get(shortURL string) (string, bool)
	Save(shortURL, longURL string)
}

func (storage *InMemoryStorage) Get(shortURL string) (string, bool) {
	storage.mutex.RLock()
	longURL, exists := storage.urlMappings[shortURL]
	storage.mutex.RUnlock()
	return longURL, exists
}

func (storage *InMemoryStorage) Save(shortURL, longURL string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.urlMappings[shortURL] = longURL
	storage.counter++
}

func NewMemoryStorage(cfg *config.Config) *InMemoryStorage {
	return &InMemoryStorage{
		urlMappings: make(map[string]string),
		mutex:       &sync.RWMutex{},
	}
}
