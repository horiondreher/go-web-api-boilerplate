package v1

import (
	"net/http"

	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/domainerr"
)

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	httpError := domainerr.HTTPErrorBody{
		Code:   domainerr.NotFoundError,
		Errors: "The requested resource was not found",
	}

	_ = httputils.Encode(w, r, http.StatusNotFound, httpError)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	httpError := domainerr.HTTPErrorBody{
		Code:   domainerr.MehodNotAllowedError,
		Errors: "The request method is not allowed",
	}

	_ = httputils.Encode(w, r, http.StatusMethodNotAllowed, httpError)
}
