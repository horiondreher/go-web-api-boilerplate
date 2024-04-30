package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpV1 "github.com/horiondreher/go-boilerplate/internal/adapters/http/v1"
	"github.com/horiondreher/go-boilerplate/internal/application/service"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	utils.StartLogger()

	// creates a new context with a cancel function that is called when the interrupt signal is received
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	config := utils.GetConfig()

	conn, err := pgxpool.New(ctx, config.DBSource)

	if err != nil {
		log.Err(err).Msg("Error connecting to database")
	}

	store := pgsqlc.New(conn)
	service := service.NewService(store)
	server, err := httpV1.NewHTTPAdapter(service)

	if err != nil {
		log.Err(err).Msg("Error creating server")
		stop()
	}

	// starts the server in a goroutine to let the main goroutine listen for the interrupt signal
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("Error starting server")
		}
	}()

	<-ctx.Done()

	// gracefully shutdown the server
	server.Shutdown()

	log.Info().Msg("Server stopped")
}
