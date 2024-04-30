package v1

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/horiondreher/go-boilerplate/internal/application/service"
	"github.com/horiondreher/go-boilerplate/internal/infrastructure/persistence/pgsqlc"
	"github.com/horiondreher/go-boilerplate/pkg/utils"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testService *service.Service

func TestMain(m *testing.M) {
	ctx := context.Background()

	utils.SetConfigPath("../../../../")
	config := utils.GetConfig()

	migrationsPath := filepath.Join("..", "..", "..", "..", "db", "postgres", "migration", "*.up.sql")
	upMigrations, err := filepath.Glob(migrationsPath)

	if err != nil {
		log.Fatalf("cannot find up migrations: %v", err)
	}

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16.2"),
		postgres.WithInitScripts(upMigrations...),
		postgres.WithDatabase(config.DBName),
		postgres.WithUsername(config.DBUser),
		postgres.WithPassword(config.DBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		log.Fatalf("cannot start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		log.Fatalf("cannot get connection string: %v", err)
	}

	conn, err := pgxpool.New(ctx, connStr)

	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	testStore := pgsqlc.New(conn)
	testService = service.NewService(testStore)

	os.Exit(m.Run())
}
