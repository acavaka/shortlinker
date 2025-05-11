package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name            string
		setupEnv        func()
		expectedBaseURL string
		expectedAddr    string
	}{
		{
			name: "env_vars_set",
			setupEnv: func() {
				os.Setenv("BASE_URL", "http://test.ru")
				os.Setenv("SERVER_ADDRESS", ":9090")
			},
			expectedBaseURL: "http://test.ru",
			expectedAddr:    ":9090",
		},
		{
			name: "env_vars_not_set",
			setupEnv: func() {
				os.Unsetenv("BASE_URL")
				os.Unsetenv("SERVER_ADDRESS")
			},
			// Потом можно переделать parseFlags() на моки
			expectedBaseURL: "",
			expectedAddr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origBaseURL := os.Getenv("BASE_URL")
			origAddr := os.Getenv("SERVER_ADDRESS")

			tt.setupEnv()

			defer func() {
				os.Setenv("BASE_URL", origBaseURL)
				os.Setenv("SERVER_ADDRESS", origAddr)
			}()

			cfg := LoadConfig()

			assert.NotNil(t, cfg, "Config should not be nil")
			assert.Equal(t, 8, cfg.URL.Length, "URL length should be 8")
			
			if tt.name == "env_vars_set" {
				assert.Equal(t, tt.expectedBaseURL, cfg.Server.BaseURL)
				assert.Equal(t, tt.expectedAddr, cfg.Server.ServerAddress)
			}
		})
	}
}
