package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SEC-Jobstreet/backend-employer-service/api"
	db "github.com/SEC-Jobstreet/backend-employer-service/db/sqlc"
	_ "github.com/SEC-Jobstreet/backend-employer-service/docs"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

//	@title			employer Service API
//	@version		1.0
//	@description	This is a employer Service Server.

// @host		localhost:4000
// @BasePath	/api/v1
func main() {

	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	runDBMigration(config)

	store := db.NewStore(connPool)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGinServer(ctx, waitGroup, config, store)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runDBMigration(config utils.Config) {
	migration, err := migrate.New(config.MigrationURL, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, store db.Store) {
	ginServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = ginServer.Start(ctx, waitGroup, config.RESTfulServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
