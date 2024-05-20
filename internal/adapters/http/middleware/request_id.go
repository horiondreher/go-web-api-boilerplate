package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}

func GetRequestID(httpCtx context.Context) (requestID string) {
	if requestIdVal, ok := httpCtx.Value(KeyRequestID).(string); ok {
		requestID = requestIdVal
	}

	return
}
