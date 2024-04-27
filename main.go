package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/SEC-Jobstreet/backend-employer-service/api"
	_ "github.com/SEC-Jobstreet/backend-employer-service/docs"
	"github.com/SEC-Jobstreet/backend-employer-service/models"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	sqlDB, err := sql.Open("pgx", config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	store, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	err = models.MigrateEnterprises(store)
	if err != nil {
		log.Fatal().Msg("could not migrate db")
	}

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGinServer(ctx, waitGroup, config, store)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, store *gorm.DB) {
	ginServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = ginServer.Start(ctx, waitGroup, config.RESTfulServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
