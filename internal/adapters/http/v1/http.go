package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/horiondreher/go-boilerplate/internal/adapters/http/middleware"
	"github.com/horiondreher/go-boilerplate/internal/adapters/http/token"
	"github.com/horiondreher/go-boilerplate/internal/application/ports"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/go-chi/chi/v5"
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
	userService ports.Service

	config *utils.Config
	router *chi.Mux
	server *http.Server

	tokenMaker *token.PasetoMaker
}

func NewHTTPAdapter(userService ports.Service) (*HTTPAdapter, error) {

	httpAdapter := &HTTPAdapter{
		userService: userService,
		config:      utils.GetConfig(),
	}

	setupValidator()
	httpAdapter.setupRouter()
	err := httpAdapter.setupServer()

	if err != nil {
		log.Err(err).Msg("error setting up server")
		return nil, err
	}

	return httpAdapter, nil
}

func (adapter *HTTPAdapter) Start() error {
	log.Info().Str("address", adapter.server.Addr).Msg("Starting server")

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

	router.NotFound(notFoundResponse)
	router.MethodNotAllowed(methodNotAllowedResponse)

	v1Router := chi.NewRouter()
	v1Router.Use(middleware.Logger)

	v1Router.Post("/users", adapter.createUser)
	v1Router.Post("/login", adapter.loginUser)
	v1Router.Post("/renew-token", adapter.renewAccessToken)

	router.Mount("/api/v1", v1Router)

	adapter.router = router
}

func (adapter *HTTPAdapter) printRoutes(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	log.Info().Str("method", method).Str("route", route).Msg("Route registered")
	return nil
}

func (adapter *HTTPAdapter) setupServer() error {

	tokenMaker, err := token.NewPasetoMaker(adapter.config.TokenSymmetricKey)

	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    adapter.config.HTTPServerAddress,
		Handler: adapter.router,
	}

	adapter.tokenMaker = tokenMaker
	adapter.server = server

	return nil
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Err(err).Msg("error encoding json")
		return err
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Err(err).Msg("error decoding JSON")
		return v, err
	}

	return v, nil
}
