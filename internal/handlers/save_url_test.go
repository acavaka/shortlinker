package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/service"
	"github.com/acavaka/shortlinker/internal/storage"
)

func TestSaveHandler(t *testing.T) {
	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}

	type want struct {
		statusCode int
	}

	cases := []struct {
		name   string
		route  string
		method string
		body   string
		want   want
	}{
		{
			name:   "valid_url_ya_ru",
			route:  "/",
			method: http.MethodPost,
			body:   "https://ya.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "valid_url_ozon_ru",
			route:  "/",
			method: http.MethodPost,
			body:   "https://ozon.ru",
			want: want{
				statusCode: http.StatusCreated,
			},
		},
		{
			name:   "invalid_url_returns_400",
			route:  "/",
			method: http.MethodPost,
			body:   "not-a-valid-url",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name:   "empty_body_returns_400",
			route:  "/",
			method: http.MethodPost,
			body:   "",
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post("/", SaveHandler(svc))

			reqBody := strings.NewReader(tt.body)
			r := httptest.NewRequest(tt.method, tt.route, reqBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			res := w.Result()
			defer func() {
				if err := res.Body.Close(); err != nil {
					log.Printf("failed to close response body: %v", err)
				}
			}()

			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				log.Printf("failed to read response body: %v", err)
				require.FailNow(t, "Failed to read response body", err)
			}

			if tt.want.statusCode == http.StatusCreated {
				assert.NotEmpty(t, resBody, "Response body should not be empty for successful creation")
			} else {
				assert.Equal(t, tt.want.statusCode, res.StatusCode,
					"Expected status code %d for case '%s', got %d",
					tt.want.statusCode, tt.name, res.StatusCode)
			}
		})
	}
}
