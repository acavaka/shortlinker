// Package storage provides interfaces and implementations for URL storage
package storage

import (
	"fmt"
	"sync"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/logger"
)

var ErrStorageRestore = fmt.Errorf("failed to restore storage from file")

type InMemoryStorage struct {
	urlMappings map[string]string
	mutex       *sync.RWMutex
	counter     uint64
}

type FileStorage struct {
	InMemoryStorage
	filePath string
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

func (storage *FileStorage) Save(shortURL, longURL string) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()
	storage.urlMappings[shortURL] = longURL
	if err := AppendToFile(storage.filePath, shortURL, longURL, storage.counter); err != nil {
		logger.Error("failed to persist URL mapping to file", err)
	}
	storage.counter++
}

func (storage *FileStorage) restore() error {
	if storage.filePath == "" {
		return nil
	}

	mapping, err := ReadFileStorage(storage.filePath)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrStorageRestore, err)
	}

	storage.mutex.Lock()
	storage.urlMappings = mapping
	storage.counter = uint64(len(mapping))
	storage.mutex.Unlock()

	return nil
}

func NewStorage(cfg *config.Config) (URLStorage, error) {
	if cfg.Service.FileStoragePath == "" {
		return &InMemoryStorage{
			urlMappings: make(map[string]string),
			mutex:       &sync.RWMutex{},
		}, nil
	}

	storage := &FileStorage{
		InMemoryStorage: InMemoryStorage{
			urlMappings: make(map[string]string),
			mutex:       &sync.RWMutex{},
		},
		filePath: cfg.Service.FileStoragePath,
	}

	if err := storage.restore(); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	return storage, nil
}
