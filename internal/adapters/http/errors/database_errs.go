package apierrs

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapDuplicateError(constraintName string) string {
	var errorMessage string

	switch constraintName {
	case "user_email_idx":
		errorMessage = "The email is already in use"
	}

	return errorMessage
}

func TransformPostgresError(err *pgconn.PgError) error {
	httpError := APIError{
		HTTPCode: http.StatusBadRequest,
		Body: APIErrorBody{
			Code:   QueryError,
			Errors: err.ConstraintName,
		},
	}

	switch err.Code {
	case "23505":
		httpError.HTTPCode = http.StatusConflict
		httpError.Body.Code = DuplicateError
		httpError.Body.Errors = MapDuplicateError(err.ConstraintName)
	}

	return httpError
}
