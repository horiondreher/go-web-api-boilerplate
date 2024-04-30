package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	ValidationError = "validation-error"
	JsonDecodeError = "json-decode-error"
	UnexpectedError = "unexpected-error"
	InvalidPassword = "invalid-password"
	InternalError   = "internal-error"
	NotFound        = "not-found"
	MehodNotAllowed = "method-not-allowed"
)

type HttpError struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}

func mapValidationTags(tag string) string {
	var tagMessage string

	switch tag {
	case "required":
		tagMessage = "The field is required"
	case "email":
		tagMessage = "The field must be a valid email address"
	default:
		tagMessage = "The field is invalid"
	}

	return tagMessage
}

func transformValidatorError(err validator.ValidationErrors) HttpError {
	errors := make(map[string]string)

	for _, e := range err {
		errors[e.Field()] = mapValidationTags(e.Tag())
	}

	return HttpError{
		Code:   ValidationError,
		Errors: errors,
	}
}

func transformUnmarshalError(err *json.UnmarshalTypeError) HttpError {
	errors := make(map[string]string)

	errors[err.Field] = fmt.Sprintf("The field is invalid. Expected type %v", err.Type)

	return HttpError{
		Code:   JsonDecodeError,
		Errors: errors,
	}
}

func matchGenericError(err error) (int, HttpError) {
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return http.StatusBadRequest, HttpError{
			Code:   JsonDecodeError,
			Errors: "The request body is invalid",
		}
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return http.StatusBadRequest, HttpError{
			Code:   InvalidPassword,
			Errors: "The password is invalid",
		}
	}

	return http.StatusInternalServerError, HttpError{
		Code:   UnexpectedError,
		Errors: err.Error(),
	}
}

func errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	var encodeErr error

	switch e := err.(type) {
	case validator.ValidationErrors:
		errBody := transformValidatorError(e)
		encodeErr = encode(w, r, http.StatusBadRequest, errBody)
	case *json.UnmarshalTypeError:
		errBody := transformUnmarshalError(e)
		encodeErr = encode(w, r, http.StatusBadRequest, errBody)
	case utils.PasswordError:
		encodeErr = encode(w, r, http.StatusInternalServerError, HttpError{
			Code:   InternalError,
			Errors: e.Error(),
		})
	default:
		httpCode, httpError := matchGenericError(e)
		encodeErr = encode(w, r, httpCode, httpError)
	}

	if encodeErr != nil {
		log.Err(encodeErr).Msg("Error encoding response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	httpError := HttpError{
		Code:   NotFound,
		Errors: "The requested resource was not found",
	}

	encode(w, r, http.StatusNotFound, httpError)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	httpError := HttpError{
		Code:   MehodNotAllowed,
		Errors: "The request method is not allowed",
	}

	encode(w, r, http.StatusMethodNotAllowed, httpError)
}
