package v1

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/horiondreher/go-boilerplate/internal/application/service"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/postgres"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testService *service.Service

func TestMain(m *testing.M) {
	utils.SetConfigPath("../../../../")
	config := utils.GetConfig()

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DBSource)

	testStore := postgres.New(conn)
	testService = service.NewService(testStore)

	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	os.Exit(m.Run())
}
