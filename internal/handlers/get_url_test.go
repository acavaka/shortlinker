package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/acavaka/shortlinker/internal/config"
	"github.com/acavaka/shortlinker/internal/logger"
	"github.com/acavaka/shortlinker/internal/service"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Get(shortLink string) (string, bool) {
	args := m.Called(shortLink)
	return args.String(0), args.Bool(1)
}

func (m *MockDB) Save(shortLink, longLink string) {
	m.Called(shortLink, longLink)
}

func TestGetHandler(t *testing.T) {
	cfg := config.LoadConfig()

	mockedDB := &MockDB{}
	svc := &service.Service{DB: mockedDB, BaseURL: cfg.Service.BaseURL}

	type want struct {
		contentType string
		statusCode  int
		success     bool
	}
	cases := []struct {
		name    string
		route   string
		longURL string
		method  string
		want    want
	}{
		{
			name:    "existing_link_redirects_to_ya_ru",
			route:   "ABC1234X",
			longURL: "https://ya.ru",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
				success:     true,
			},
		},
		{
			name:    "existing_link_redirects_to_ozon_ru",
			route:   "Ur0lH9i9",
			longURL: "https://ozon.ru",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusTemporaryRedirect,
				success:     true,
			},
		},
		{
			name:    "missing_link_returns_400",
			route:   "MissingRoute",
			longURL: "",
			method:  http.MethodGet,
			want: want{
				contentType: `"text/plain; charset=utf-8"`,
				statusCode:  http.StatusBadRequest,
				success:     false,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			mockedDB.On("Save", tt.route, tt.longURL).Return()
			mockedDB.On("Get", tt.route).Return(tt.longURL, tt.want.success)

			router := chi.NewRouter()
			router.Get("/{id}", GetHandler(svc, cfg.Service.BaseURL, logger.Initialize()))
			r := httptest.NewRequest(tt.method, "/"+tt.route, http.NoBody)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			res := w.Result()
			defer res.Body.Close() // лучше defer для 100% закрытия?

			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Ошибка при чтении тела ответа")

			if tt.want.success {
				assert.NotEmpty(t, resBody,
					"Для существующей ссылки тело ответа не должно быть пустым")
			} else {
				assert.Empty(t, resBody,
					"Для отсутствующей ссылки тело ответа должно быть пустым")
			}

			assert.Equal(t, tt.want.statusCode, res.StatusCode,
				"Ожидался статус %d (%s), получен %d (%s)",
				tt.want.statusCode, http.StatusText(tt.want.statusCode),
				res.StatusCode, http.StatusText(res.StatusCode))
		})
	}
}
