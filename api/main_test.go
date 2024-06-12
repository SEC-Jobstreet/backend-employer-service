package api

import (
	"os"
	"testing"

	db "github.com/SEC-Jobstreet/backend-employer-service/db/sqlc"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Querier) *Server {
	config := utils.Config{}

	// require.NoError(t, err)
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
