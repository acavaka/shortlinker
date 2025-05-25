package storage

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStorage(t *testing.T) {
	storage := &InMemoryStorage{
		urlMappings: make(map[string]string),
		mutex:       &sync.RWMutex{},
	}

	tests := []struct {
		name     string
		shortURL string
		longURL  string
	}{
		{
			name:     "basic url",
			shortURL: "abc123",
			longURL:  "https://example.com",
		},
		{
			name:     "empty url",
			shortURL: "def456",
			longURL:  "",
		},
		{
			name:     "complex url",
			shortURL: "ghi789",
			longURL:  "https://example.com/path?param=value#fragment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Save
			storage.Save(tt.shortURL, tt.longURL)

			// Test Get
			gotURL, exists := storage.Get(tt.shortURL)
			assert.True(t, exists)
			assert.Equal(t, tt.longURL, gotURL)
		})
	}

	// Test non-existent URL
	_, exists := storage.Get("nonexistent")
	assert.False(t, exists)
}

func TestFileStorage(t *testing.T) {
	// Create temporary directory for test file
	tmpDir, err := os.MkdirTemp("", "shortlinker_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "urls.txt")

	tests := []struct {
		name     string
		shortURL string
		longURL  string
	}{
		{
			name:     "basic url",
			shortURL: "abc123",
			longURL:  "https://example.com",
		},
		{
			name:     "complex url",
			shortURL: "def456",
			longURL:  "https://example.com/path?param=value#fragment",
		},
	}

	t.Run("file storage operations", func(t *testing.T) {
		storage := &FileStorage{
			InMemoryStorage: InMemoryStorage{
				urlMappings: make(map[string]string),
				mutex:       &sync.RWMutex{},
			},
			filePath: filePath,
		}

		for _, tt := range tests {
			// Test Save
			storage.Save(tt.shortURL, tt.longURL)

			// Test Get
			gotURL, exists := storage.Get(tt.shortURL)
			assert.True(t, exists)
			assert.Equal(t, tt.longURL, gotURL)
		}

		// Test non-existent URL
		_, exists := storage.Get("nonexistent")
		assert.False(t, exists)
	})

	t.Run("restore from file", func(t *testing.T) {
		// Create new storage instance to test restore
		newStorage := &FileStorage{
			InMemoryStorage: InMemoryStorage{
				urlMappings: make(map[string]string),
				mutex:       &sync.RWMutex{},
			},
			filePath: filePath,
		}

		err := newStorage.restore()
		require.NoError(t, err)

		// Verify all URLs were restored
		for _, tt := range tests {
			gotURL, exists := newStorage.Get(tt.shortURL)
			assert.True(t, exists)
			assert.Equal(t, tt.longURL, gotURL)
		}
	})
}

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name           string
		fileStorageDir string
		wantType       URLStorage
		wantErr        bool
	}{
		{
			name:           "in-memory storage",
			fileStorageDir: "",
			wantType:       &InMemoryStorage{},
			wantErr:        false,
		},
		{
			name:           "file storage",
			fileStorageDir: "testdata",
			wantType:       &FileStorage{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				Service: config.ServiceConfig{
					FileStoragePath: tt.fileStorageDir,
				},
			}

			storage, err := NewStorage(cfg)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.IsType(t, tt.wantType, storage)
		})
	}
}
