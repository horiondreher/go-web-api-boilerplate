package apierrs

import (
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func MatchGenericError(err error) (int, APIError) {
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return http.StatusBadRequest, APIError{
			Code:   JsonDecodeError,
			Errors: "The request body is invalid",
		}
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return http.StatusBadRequest, APIError{
			Code:   InvalidPasswordError,
			Errors: "The password is invalid",
		}
	}

	log.Err(err).Msg("unhandled http error")

	return http.StatusInternalServerError, APIError{
		Code:   UnexpectedError,
		Errors: "internal server error",
	}
}
