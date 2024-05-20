package apierrs

import "fmt"

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

type APIError struct {
	HTTPCode      int
	OriginalError string
	Body          APIErrorBody
}

type APIErrorBody struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.HTTPCode)
}
