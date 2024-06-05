package apierrs

import (
	"errors"
	"io"
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"golang.org/x/crypto/bcrypt"
)

func MatchGenericError(err error) error {
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return APIError{
			HTTPCode:      http.StatusBadRequest,
			OriginalError: err.Error(),
			Body: APIErrorBody{
				Code:   JsonDecodeError,
				Errors: "The request body is invalid",
			},
		}
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return APIError{
			HTTPCode:      http.StatusUnauthorized,
			OriginalError: err.Error(),
			Body: APIErrorBody{
				Code:   InvalidPasswordError,
				Errors: "The password is invalid",
			},
		}
	}

	if errors.Is(err, token.ErrInvalidToken) {
		return APIError{
			HTTPCode:      http.StatusUnauthorized,
			OriginalError: err.Error(),
			Body: APIErrorBody{
				Code:   InvalidToken,
				Errors: "Invalid token",
			},
		}
	}

	if errors.Is(err, token.ErrExpiredToken) {
		return APIError{
			HTTPCode:      http.StatusUnauthorized,
			OriginalError: err.Error(),
			Body: APIErrorBody{
				Code:   ExpiredToken,
				Errors: "Expired token",
			},
		}
	}

	return APIError{
		HTTPCode:      http.StatusInternalServerError,
		OriginalError: err.Error(),
		Body: APIErrorBody{
			Code:   UnexpectedError,
			Errors: "Internal server error",
		},
	}
}
