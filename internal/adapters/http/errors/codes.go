package apierrs

import "fmt"

var (
	ValidationError      = "request/invalid-fields"
	JsonDecodeError      = "request/invalid-json"
	UnexpectedError      = "server/internal-error"
	InvalidPasswordError = "auth/invalid-password"
	InternalError        = "server/internal-error"
	NotFoundError        = "server/not-found"
	MehodNotAllowedError = "server/method-not-allowed"
	DuplicateError       = "data/duplicate"
	QueryError           = "data/invalid-query"
	UnauthorizedError    = "auth/unauthorized"
)

type APIError struct {
	HTTPCode int
	Body     APIErrorBody
}

type APIErrorBody struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.HTTPCode)
}
