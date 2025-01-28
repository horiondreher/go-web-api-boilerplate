package domainerr

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func MatchHashError(err error) *DomainError {
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return &DomainError{
			HTTPCode:      http.StatusUnauthorized,
			OriginalError: err.Error(),
			HTTPErrorBody: HTTPErrorBody{
				Code:   InvalidPasswordError,
				Errors: "The password is invalid",
			},
		}
	}

	return NewInternalError(err)
}
