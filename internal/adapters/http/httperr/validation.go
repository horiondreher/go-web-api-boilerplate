package httperr

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

func MatchValidationError(err error) *domainerr.DomainError {
	validationErr, ok := err.(validator.ValidationErrors)
	if ok {
		return TransformValidatorError(validationErr)
	}

	return domainerr.NewInternalError(err)
}

func TransformValidatorError(err validator.ValidationErrors) *domainerr.DomainError {
	errors := make(map[string]string)

	for _, e := range err {
		errors[e.Field()] = mapValidationTags(e.Tag())
	}

	return &domainerr.DomainError{
		HTTPCode:      http.StatusUnprocessableEntity,
		OriginalError: err.Error(),
		HTTPErrorBody: domainerr.HTTPErrorBody{
			Code:   domainerr.ValidationError,
			Errors: errors,
		},
	}
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
