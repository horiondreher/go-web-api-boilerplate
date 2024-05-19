package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type key int

const (
	KeyRequestID key = iota
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
