package apierrs

import (
	"errors"
	"io"
	"net/http"

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

	return http.StatusInternalServerError, APIError{
		Code:   UnexpectedError,
		Errors: err.Error(),
	}
}
