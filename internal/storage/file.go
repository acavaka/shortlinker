package storage

import (
	"fmt"
	"sync"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/logger"
)

type FileStorage struct {
	InMemoryStorage
	filePath string
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

func NewFileStorage(cfg *config.Config) (*FileStorage, error) {
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
