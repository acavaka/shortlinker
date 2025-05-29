package storage

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			name:     "basic_url",
			shortURL: "abc123",
			longURL:  "https://example.com",
		},
		{
			name:     "complex_url",
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
