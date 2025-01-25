package httperr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

func MatchEncodingError(err error) *domainerr.DomainError {
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return &domainerr.DomainError{
			HTTPCode:      http.StatusBadRequest,
			OriginalError: err.Error(),
			HTTPErrorBody: domainerr.HTTPErrorBody{
				Code:   domainerr.JsonDecodeError,
				Errors: "The request body is invalid",
			},
		}
	}

	jsonErr, ok := err.(*json.UnmarshalTypeError)
	if ok {
		return transformUnmarshalError(jsonErr)
	}

	return domainerr.NewInternalError(err)
}

func transformUnmarshalError(err *json.UnmarshalTypeError) *domainerr.DomainError {
	errors := make(map[string]string)
	errors[err.Field] = fmt.Sprintf("The field is invalid. Expected type %v", err.Type)

	return &domainerr.DomainError{
		HTTPCode:      http.StatusUnprocessableEntity,
		OriginalError: err.Error(),
		HTTPErrorBody: domainerr.HTTPErrorBody{
			Code:   domainerr.JsonDecodeError,
			Errors: errors,
		},
	}
}
