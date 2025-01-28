package domainerr

import (
	"fmt"
	"net/http"
)

var (
	UnauthorizedError    = "auth/unauthorized"
	InvalidToken         = "auth/invalid-token"
	ExpiredToken         = "auth/expired-token"
	InvalidPasswordError = "auth/invalid-password"
	ValidationError      = "request/invalid-fields"
	JsonDecodeError      = "request/invalid-json"
	DuplicateError       = "data/duplicate"
	QueryError           = "data/invalid-query"
	UnexpectedError      = "server/internal-error"
	InternalError        = "server/internal-error"
	NotFoundError        = "server/not-found"
	MehodNotAllowedError = "server/method-not-allowed"
)

type HTTPErrorBody struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}

type DomainError struct {
	HTTPCode      int
	HTTPErrorBody HTTPErrorBody
	OriginalError string
}

func NewDomainError(httpCode int, errorCode string, errorMsg any, err error) *DomainError {
	return &DomainError{
		HTTPCode:      httpCode,
		OriginalError: err.Error(),
		HTTPErrorBody: HTTPErrorBody{
			Code:   errorCode,
			Errors: errorMsg,
		},
	}
}

func NewInternalError(err error) *DomainError {
	return &DomainError{
		HTTPCode:      http.StatusInternalServerError,
		OriginalError: err.Error(),
		HTTPErrorBody: HTTPErrorBody{
			Code:   UnexpectedError,
			Errors: "Internal server error",
		},
	}
}

func (e DomainError) Error() string {
	return fmt.Sprintf("api error: %d", e.HTTPCode)
}
