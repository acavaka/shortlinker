// Package service provides business logic for URL shortening service
package service

import (
	"fmt"
	"math/rand"

	"github.com/acavaka/shortlinker/internal/storage"
)

var ErrURLNotFound = fmt.Errorf("URL not found")

type Service struct {
	DB              storage.URLStorage
	FileStoragePath string
	BaseURL         string
}

func (svc *Service) SaveURL(longURL string) string {
	shortURL := svc.generateUniqueShortLink()
	svc.DB.Save(shortURL, longURL)
	return shortURL
}

func (svc *Service) GetURL(shortURL string) (string, error) {
	longURL, ok := svc.DB.Get(shortURL)
	if !ok {
		return "", fmt.Errorf("%w: short URL '%s' does not exist in storage", ErrURLNotFound, shortURL)
	}
	return longURL, nil
}

func (svc *Service) generateUniqueShortLink() string {
	const length = 8
	var uniqueShortURL string

	for {
		shortURLCandidate := generateRandomString(length)
		_, exists := svc.DB.Get(shortURLCandidate)
		if !exists {
			uniqueShortURL = shortURLCandidate
			break
		}
	}
	return uniqueShortURL
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}
