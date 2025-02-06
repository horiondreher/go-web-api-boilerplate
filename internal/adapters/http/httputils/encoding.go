package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httperr"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
	"github.com/rs/zerolog/log"
)

func Encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) *domainerr.DomainError {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Err(err).Msg("error encoding json")
		return httperr.MatchEncodingError(err)
	}

	return nil
}

func Decode[T any](r *http.Request) (T, *domainerr.DomainError) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Err(err).Msg("error decoding JSON")
		return v, httperr.MatchEncodingError(err)
	}

	return v, nil
}
