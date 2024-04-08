package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) example(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "OK")
}
