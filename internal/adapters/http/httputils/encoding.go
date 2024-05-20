package httputils

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func Encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Err(err).Msg("error encoding json")
		return err
	}

	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Err(err).Msg("error decoding JSON")
		return v, err
	}

	return v, nil
}
