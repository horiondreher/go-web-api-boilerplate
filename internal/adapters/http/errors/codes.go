package apierrs

var (
	ValidationError      = "validation-error"
	JsonDecodeError      = "json-decode-error"
	UnexpectedError      = "unexpected-error"
	InvalidPasswordError = "invalid-password"
	InternalError        = "internal-error"
	NotFoundError        = "not-found"
	MehodNotAllowedError = "method-not-allowed"
	DuplicateError       = "duplicate-error"
	QueryError           = "query-error"
)

type APIError struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}
