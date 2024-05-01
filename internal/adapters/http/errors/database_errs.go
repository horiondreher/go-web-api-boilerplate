package apierrs

import "github.com/jackc/pgx/v5/pgconn"

func MapDuplicateError(constraintName string) string {
	var errorMessage string

	switch constraintName {
	case "user_email_idx":
		errorMessage = "The email is already in use"
	}

	return errorMessage
}

func TransformPostgresError(err *pgconn.PgError) APIError {
	httpError := APIError{
		Code:   QueryError,
		Errors: err.ConstraintName,
	}

	switch err.Code {
	case "23505":
		httpError.Code = DuplicateError
		httpError.Errors = MapDuplicateError(err.ConstraintName)
	}

	return httpError
}
