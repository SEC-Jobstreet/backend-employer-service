package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	config utils.Config
	store  *gorm.DB
	router *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config utils.Config, store *gorm.DB) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	server.setupRouter()

	return server, nil
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(ctx context.Context, waitGroup *errgroup.Group, address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}
	var err error
	waitGroup.Go(func() error {
		log.Info().Msgf("RESTFUL API server serve at %s", address)
		err = srv.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Fatal().Msg("RESTFUL API server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown RESTFUL API server")

		err = srv.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown RESTFUL API server")
			return err
		}
		log.Info().Msg("RESTFUL API server is stopped")
		return nil
	})

	return err
}
