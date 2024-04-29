package middleware

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customWriter := NewResponseWriter(w)

		next.ServeHTTP(customWriter, r)

		statusCode := strconv.Itoa(customWriter.statusCode)

		log.Info().Str("method", r.Method).Str("path", r.URL.Path).Str("response", statusCode).Msg("Request")
	})
}
