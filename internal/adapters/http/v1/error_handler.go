package v1

import (
	"encoding/json"
	"net/http"

	apierrs "github.com/horiondreher/go-boilerplate/internal/adapters/http/errors"
	"github.com/horiondreher/go-boilerplate/pkg/utils"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var encodeErr error

	switch e := err.(type) {
	case validator.ValidationErrors:
		errBody := apierrs.TransformValidatorError(e)
		encodeErr = encode(w, r, http.StatusBadRequest, errBody)
	case *json.UnmarshalTypeError:
		errBody := apierrs.TransformUnmarshalError(e)
		encodeErr = encode(w, r, http.StatusBadRequest, errBody)
	case *pgconn.PgError:
		errBody := apierrs.TransformPostgresError(e)
		encodeErr = encode(w, r, http.StatusBadRequest, errBody)
	case *utils.HashError:
		encodeErr = encode(w, r, http.StatusInternalServerError, apierrs.APIError{
			Code:   apierrs.InternalError,
			Errors: e.Error(),
		})
	case *SessionError:
		encodeErr = encode(w, r, http.StatusUnauthorized, apierrs.APIError{
			Code:   apierrs.UnauthorizedError,
			Errors: e.Error(),
		})
	default:
		httpCode, httpError := apierrs.MatchGenericError(e)
		encodeErr = encode(w, r, httpCode, httpError)
	}

	if encodeErr != nil {
		log.Err(encodeErr).Msg("Error encoding response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	httpError := apierrs.APIError{
		Code:   apierrs.NotFoundError,
		Errors: "The requested resource was not found",
	}

	encode(w, r, http.StatusNotFound, httpError)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	httpError := apierrs.APIError{
		Code:   apierrs.MehodNotAllowedError,
		Errors: "The request method is not allowed",
	}

	encode(w, r, http.StatusMethodNotAllowed, httpError)
}
