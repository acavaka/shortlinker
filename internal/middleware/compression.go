package middleware

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func allowedContentType(header string) bool {
	allowed := map[string]struct{}{
		"text/html":        {},
		"application/json": {},
	}
	_, ok := allowed[header]
	return ok
}

func GzipMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ow := w
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				if allowedContentType(r.Header.Get("Content-Type")) {
					cw := newCompressWriter(w)
					ow = cw
					defer func() {
						err := cw.Close()
						if err != nil {
							logger.Error("failed to close compress writer", zap.Error(err))
							http.Error(w, "", http.StatusInternalServerError)
							return
						}
					}()
				}
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					logger.Error("failed to read compressed body", zap.Error(err))
					http.Error(w, "check if gzip data is valid", http.StatusBadRequest)
					return
				} else {
					r.Body = cr
					defer func() {
						err = cr.Close()
						if err != nil {
							logger.Error("failed to close compress reader", zap.Error(err))
							http.Error(w, "", http.StatusInternalServerError)
							return
						}
					}()
				}
			}
			next.ServeHTTP(ow, r)
		})
	}
}
