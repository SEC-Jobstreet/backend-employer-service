package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/SEC-Jobstreet/backend-employer-service/api"
	_ "github.com/SEC-Jobstreet/backend-employer-service/docs"
	"github.com/SEC-Jobstreet/backend-employer-service/repository"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
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

	// sqlDB, err := sql.Open("pgx", config.DBSource)
	// if err != nil {
	// 	log.Fatal().Msg("cannot connect to db")
	// }

	// store, err := gorm.Open(postgres.New(postgres.Config{
	// 	Conn: sqlDB,
	// }), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal().Msg("cannot connect to db")
	// }

	// err = models.MigrateEnterprises(store)
	// if err != nil {
	// 	log.Fatal().Msg("could not migrate db")
	// }
	DBConnections := []string{
		"postgresql://admin:admin@34.126.133.75:5432/employer_service_jobstreet?sslmode=disable",
		"postgresql://admin:admin@3.106.126.218:5432/employer_service_jobstreet?sslmode=disable",
		"postgresql://admin:admin@34.124.137.133:5432/employer_service_jobstreet?sslmode=disable",
	}
	stores := repository.NewEmployerRepo(DBConnections)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runGinServer(ctx, waitGroup, config, stores)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config utils.Config, stores *repository.EmployerRepo) {
	ginServer, err := api.NewServer(config, stores)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = ginServer.Start(ctx, waitGroup, config.RESTfulServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
