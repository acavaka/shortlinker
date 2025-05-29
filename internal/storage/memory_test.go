package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
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
			name:     "basic_url",
			shortURL: "abc123",
			longURL:  "https://example.com",
		},
		{
			name:     "empty_url",
			shortURL: "def456",
			longURL:  "",
		},
		{
			name:     "complex_url",
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
