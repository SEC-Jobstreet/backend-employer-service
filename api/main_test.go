package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestServer(t *testing.T) *Server {
	// require.NoError(t, err)

	return nil
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
