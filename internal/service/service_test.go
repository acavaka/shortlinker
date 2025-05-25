package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockStorage is a mock implementation of URLStorage for testing
type MockStorage struct {
	urls map[string]string
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		urls: make(map[string]string),
	}
}

func (m *MockStorage) Get(shortURL string) (string, bool) {
	longURL, exists := m.urls[shortURL]
	return longURL, exists
}

func (m *MockStorage) Save(shortURL, longURL string) {
	m.urls[shortURL] = longURL
}

func TestService_SaveURL(t *testing.T) {
	tests := []struct {
		name    string
		longURL string
	}{
		{
			name:    "basic url",
			longURL: "https://example.com",
		},
		{
			name:    "complex url",
			longURL: "https://example.com/path?param=value#fragment",
		},
		{
			name:    "empty url",
			longURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMockStorage()
			svc := &Service{DB: storage}

			shortURL := svc.SaveURL(tt.longURL)

			// Verify the short URL is not empty
			assert.NotEmpty(t, shortURL)

			// Verify the URL was saved
			savedURL, exists := storage.Get(shortURL)
			require.True(t, exists)
			assert.Equal(t, tt.longURL, savedURL)
		})
	}
}

func TestService_GetURL(t *testing.T) {
	tests := []struct {
		name        string
		setupURLs   map[string]string
		shortURL    string
		wantURL     string
		wantErr     error
		shouldExist bool
	}{
		{
			name: "existing url",
			setupURLs: map[string]string{
				"abc123": "https://example.com",
			},
			shortURL:    "abc123",
			wantURL:     "https://example.com",
			shouldExist: true,
		},
		{
			name:        "non-existing url",
			setupURLs:   map[string]string{},
			shortURL:    "nonexistent",
			wantErr:     ErrURLNotFound,
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMockStorage()
			for k, v := range tt.setupURLs {
				storage.Save(k, v)
			}

			svc := &Service{DB: storage}

			gotURL, err := svc.GetURL(tt.shortURL)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantURL, gotURL)
		})
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "length 8",
			length: 8,
		},
		{
			name:   "length 16",
			length: 16,
		},
		{
			name:   "length 0",
			length: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateRandomString(tt.length)
			assert.Len(t, result, tt.length)

			// Verify the string contains only valid characters
			for _, char := range result {
				assert.Contains(t,
					"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
					string(char))
			}
		})
	}
}

func TestService_GenerateUniqueShortLink(t *testing.T) {
	storage := NewMockStorage()
	svc := &Service{DB: storage}

	// Generate multiple short links and verify they're unique
	generated := make(map[string]bool)
	for i := 0; i < 100; i++ {
		shortURL := svc.generateUniqueShortLink()

		// Verify the URL is not empty and has correct length
		assert.NotEmpty(t, shortURL)
		assert.Len(t, shortURL, 8)

		// Verify it's unique
		assert.False(t, generated[shortURL], "Generated duplicate short URL: %s", shortURL)
		generated[shortURL] = true
	}
}
