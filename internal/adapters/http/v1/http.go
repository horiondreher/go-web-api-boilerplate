package v1

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	apierrs "github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/errors"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/httputils"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/ports"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validate *validator.Validate

func setupValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
}

type HTTPAdapter struct {
	userService ports.UserService

	config *utils.Config
	router *chi.Mux
	server *http.Server

	tokenMaker *token.PasetoMaker
}

func NewHTTPAdapter(userService ports.UserService) (*HTTPAdapter, error) {

	httpAdapter := &HTTPAdapter{
		userService: userService,
		config:      utils.GetConfig(),
	}

	setupValidator()

	err := httpAdapter.setupTokenMaker()

	if err != nil {
		log.Err(err).Msg("error setting up server")
		return nil, err
	}

	httpAdapter.setupRouter()
	httpAdapter.setupServer()

	return httpAdapter, nil
}

func (adapter *HTTPAdapter) Start() error {
	log.Info().Str("address", adapter.server.Addr).Msg("starting server")

	chi.Walk(adapter.router, adapter.printRoutes)

	return adapter.server.ListenAndServe()
}

func (adapter *HTTPAdapter) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := adapter.server.Shutdown(ctx); err != nil {
		log.Err(err).Msg("error shutting down server")
	}
}

func (adapter *HTTPAdapter) setupRouter() {
	router := chi.NewRouter()

	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.RedirectSlashes)

	router.NotFound(notFoundResponse)
	router.MethodNotAllowed(methodNotAllowedResponse)

	v1Router := chi.NewRouter()
	v1Router.Use(middleware.RequestID)
	v1Router.Use(middleware.Logger)

	v1Router.Post("/users", adapter.handlerWrapper(adapter.createUser))
	v1Router.Post("/login", adapter.handlerWrapper(adapter.loginUser))
	v1Router.Post("/renew-token", adapter.handlerWrapper(adapter.renewAccessToken))

	// private routes
	v1Router.Group(func(r chi.Router) {
		r.Use(middleware.Authentication(adapter.tokenMaker))
		r.Get("/user/{uid}", adapter.handlerWrapper(adapter.getUserByUID))
	})

	router.Mount("/api/v1", v1Router)

	adapter.router = router
}

type HandlerWrapper func(w http.ResponseWriter, r *http.Request) error

func (adapter *HTTPAdapter) handlerWrapper(handlerFn HandlerWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if apiErr := handlerFn(w, r); apiErr != nil {
			var apiErrIntf apierrs.APIError
			var err error

			requestID := middleware.GetRequestID(r.Context())

			if errors.As(apiErr, &apiErrIntf) {
				log.Info().Str("id", requestID).Str("error message", apiErrIntf.OriginalError).Msg("request error")
				err = httputils.Encode(w, r, apiErrIntf.HTTPCode, apiErrIntf.Body)

			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			if err != nil {
				log.Err(err).Msg("error encoding response")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}
	}
}

func (adapter *HTTPAdapter) printRoutes(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	log.Info().Str("method", method).Str("route", route).Msg("route registered")
	return nil
}

func (adapter *HTTPAdapter) setupTokenMaker() error {
	tokenMaker, err := token.NewPasetoMaker(adapter.config.TokenSymmetricKey)

	if err != nil {
		return err
	}

	adapter.tokenMaker = tokenMaker

	return nil
}

func (adapter *HTTPAdapter) setupServer() {
	server := &http.Server{
		Addr:    adapter.config.HTTPServerAddress,
		Handler: adapter.router,
	}

	adapter.server = server
}
