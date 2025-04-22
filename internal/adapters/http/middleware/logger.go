package middleware

import (
	"net/http"

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
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(KeyRequestID).(string)

		log.Info().Str("id", requestID).Str("method", r.Method).Str("path", r.URL.Path).Msg("request received")

		customWriter := NewResponseWriter(w)
		next.ServeHTTP(customWriter, r)

		log.Info().Str("id", requestID).Str("method", r.Method).Str("path", r.URL.Path).Int("response", customWriter.statusCode).Msg("request response")
	}

	return http.HandlerFunc(fn)
}
