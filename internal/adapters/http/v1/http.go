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

type HTTPHandler struct {
	service ports.Service

	config *utils.Config
	router *chi.Mux
	server *http.Server

	tokenMaker *token.PasetoMaker
}

func NewHTTPHandler(service ports.Service) (*HTTPHandler, error) {

	httpHandler := &HTTPHandler{
		service: service,
		config:  utils.GetConfig(),
	}

	setupValidator()
	httpHandler.setupRouter()
	err := httpHandler.setupServer()

	if err != nil {
		log.Err(err).Msg("Error setting up server")
		return nil, err
	}

	return httpHandler, nil
}

func (handler *HTTPHandler) Start() error {
	log.Info().Str("address", handler.server.Addr).Msg("Starting server")

	return handler.server.ListenAndServe()
}

func (handler *HTTPHandler) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := handler.server.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Error shutting down server")
	}
}

func (handler *HTTPHandler) setupRouter() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Post("/user", handler.createUser)
	router.Post("/login", handler.loginUser)

	handler.router = router
}

func (handler *HTTPHandler) setupServer() error {

	tokenMaker, err := token.NewPasetoMaker(handler.config.TokenSymmetricKey)

	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    handler.config.HTTPServerAddress,
		Handler: handler.router,
	}

	handler.tokenMaker = tokenMaker
	handler.server = server

	return nil
}

func encode[T any](w http.ResponseWriter, _ *http.Request, status int, v T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Err(err).Msg("Error encoding json")
		return err
	}

	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T

	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		log.Err(err).Msg("Error decoding JSON")
		return v, err
	}

	return v, nil
}
