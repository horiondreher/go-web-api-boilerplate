package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	apierrs "github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/errors"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/rs/zerolog/log"
)

const (
	bearerAuth = "bearer"
)

func Authentication(tokenMaker *token.PasetoMaker) func(next http.Handler) http.Handler {
	if tokenMaker == nil {
		fmt.Println("PasetoMaker is not initialized")
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")

			requestID := GetRequestID(r.Context())

			if len(auth) == 0 {
				log.Info().Str("id", requestID).Str("error message", "empty authorization header").Msg("request error")
				httputils.Encode(w, r, http.StatusUnauthorized, apierrs.APIErrorBody{
					Code:   apierrs.UnauthorizedError,
					Errors: "Empty Authorization Header",
				})
				return
			}

			fields := strings.Fields(auth)

			if len(fields) < 2 {
				log.Info().Str("id", requestID).Str("error message", "invalid authorization header").Msg("request error")
				httputils.Encode(w, r, http.StatusUnauthorized, apierrs.APIErrorBody{
					Code:   apierrs.UnauthorizedError,
					Errors: "Invalid Authorization Header",
				})
				return
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != bearerAuth {
				log.Info().Str("id", requestID).Str("error message", "invalid authorization type").Msg("request error")
				httputils.Encode(w, r, http.StatusUnauthorized, apierrs.APIErrorBody{
					Code:   apierrs.UnauthorizedError,
					Errors: "Invalid Authorization Type",
				})
				return
			}

			accessToken := fields[1]
			fmt.Println(accessToken)
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				log.Info().Str("id", requestID).Str("error message", err.Error()).Msg("request error")
				apiErr := apierrs.MatchGenericError(err).(apierrs.APIError)

				httputils.Encode(w, r, http.StatusUnauthorized, apiErr.Body)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, KeyAuthUser, payload)

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
