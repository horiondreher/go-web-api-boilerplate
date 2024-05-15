package apierrs

import (
	"errors"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func MatchGenericError(err error) error {
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return APIError{
			HTTPCode: http.StatusBadRequest,
			Body: APIErrorBody{
				Code:   JsonDecodeError,
				Errors: "The request body is invalid",
			},
		}
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return APIError{
			HTTPCode: http.StatusUnauthorized,
			Body: APIErrorBody{
				Code:   InvalidPasswordError,
				Errors: "The password is invalid",
			},
		}
	}

	log.Err(err).Msg("unhandled http error")

	return APIError{
		HTTPCode: http.StatusInternalServerError,
		Body: APIErrorBody{
			Code:   UnexpectedError,
			Errors: "Internal server error",
		},
	}
}
