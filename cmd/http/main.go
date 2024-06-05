package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpV1 "github.com/horiondreher/go-web-api-boilerplate/internal/adapters/http/v1"
	"github.com/horiondreher/go-web-api-boilerplate/internal/domain/services"
	"github.com/horiondreher/go-web-api-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-web-api-boilerplate/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	os.Setenv("TZ", "UTC")

	utils.StartLogger()

	// creates a new context with a cancel function that is called when the interrupt signal is received
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	config := utils.GetConfig()

	conn, err := pgxpool.New(ctx, config.DBSource)

	if err != nil {
		log.Err(err).Msg("error connecting to database")
	}

	store := pgsqlc.New(conn)
	userService := services.NewUserManager(store)
	server, err := httpV1.NewHTTPAdapter(userService)

	if err != nil {
		log.Err(err).Msg("error creating server")
		stop()
	}

	// starts the server in a goroutine to let the main goroutine listen for the interrupt signal
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("error starting server")
		}
	}()

	<-ctx.Done()

	// gracefully shutdown the server
	server.Shutdown()

	log.Info().Msg("server stopped")
}
