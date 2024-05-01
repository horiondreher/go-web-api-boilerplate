package apierrs

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func MapValidationTags(tag string) string {
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

func TransformValidatorError(err validator.ValidationErrors) APIError {
	errors := make(map[string]string)

	for _, e := range err {
		errors[e.Field()] = MapValidationTags(e.Tag())
	}

	return APIError{
		Code:   ValidationError,
		Errors: errors,
	}
}

func TransformUnmarshalError(err *json.UnmarshalTypeError) APIError {
	errors := make(map[string]string)

	errors[err.Field] = fmt.Sprintf("The field is invalid. Expected type %v", err.Type)

	return APIError{
		Code:   JsonDecodeError,
		Errors: errors,
	}
}
