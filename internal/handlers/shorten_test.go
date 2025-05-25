package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/models"
	"github.com/acavaka/shortlinker/internal/service"
	"github.com/acavaka/shortlinker/internal/storage"
)

func TestShortenHandler(t *testing.T) {
	const (
		ct    = "application/json"
		route = "/api/shorten"
	)

	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}
	type want struct {
		statusCode  int
		contentType string
	}
	cases := []struct {
		name   string
		method string
		body   string
		want   want
	}{
		{
			name:   "Valid URL",
			method: http.MethodPost,
			body:   `{"url": "https://practicum.yandex.ru"}`,
			want: want{
				statusCode:  http.StatusCreated,
				contentType: ct,
			},
		},
		{
			name:   "Invalid JSON",
			method: http.MethodPost,
			body:   `{"url": invalid}`,
			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: ct,
			},
		},
		{
			name:   "Empty URL",
			method: http.MethodPost,
			body:   `{"url": ""}`,
			want: want{
				statusCode:  http.StatusInternalServerError,
				contentType: ct,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			router := chi.NewRouter()
			router.Post(route, ShortenHandler(svc))

			reqBody := strings.NewReader(tc.body)
			r := httptest.NewRequest(tc.method, route, reqBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)
			res := w.Result()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			defer res.Body.Close()

			assert.Equal(t, tc.want.statusCode, res.StatusCode)
			assert.Equal(t, tc.want.contentType, res.Header.Get("Content-Type"))

			if res.StatusCode == http.StatusCreated {
				var response models.Response
				err = json.Unmarshal(resBody, &response)
				require.NoError(t, err)
				assert.NotEmpty(t, response.Result)
				assert.Contains(t, response.Result, cfg.Server.BaseURL)
			}
		})
	}
}
