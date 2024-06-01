package v1

import (
	"encoding/json"
	"net/http"

	apierrs "github.com/horiondreher/go-boilerplate/internal/adapters/http/errors"
	"github.com/horiondreher/go-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-boilerplate/internal/utils"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/go-playground/validator/v10"
)

func errorResponse(err error) error {
	switch e := err.(type) {
	case validator.ValidationErrors:
		return apierrs.TransformValidatorError(e)
	case *json.UnmarshalTypeError:
		return apierrs.TransformUnmarshalError(e)
	case *pgconn.PgError:
		return apierrs.TransformPostgresError(e)
	case *utils.HashError:
		return apierrs.APIError{
			HTTPCode:      http.StatusInternalServerError,
			OriginalError: err.Error(),
			Body: apierrs.APIErrorBody{
				Code:   apierrs.InternalError,
				Errors: e.Error(),
			},
		}
	case *SessionError:
		return apierrs.APIError{
			HTTPCode:      http.StatusUnauthorized,
			OriginalError: err.Error(),
			Body: apierrs.APIErrorBody{
				Code:   apierrs.UnauthorizedError,
				Errors: e.Error(),
			},
		}
	default:
		return apierrs.MatchGenericError(e)
	}
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	httpError := apierrs.APIErrorBody{
		Code:   apierrs.NotFoundError,
		Errors: "The requested resource was not found",
	}

	httputils.Encode(w, r, http.StatusNotFound, httpError)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	httpError := apierrs.APIErrorBody{
		Code:   apierrs.MehodNotAllowedError,
		Errors: "The request method is not allowed",
	}

	httputils.Encode(w, r, http.StatusMethodNotAllowed, httpError)
}
