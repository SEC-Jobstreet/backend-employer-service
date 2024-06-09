package api

import (
	"net/http"
	"os"

	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
)

type publicRequest struct {
	UserId int32  `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type publicResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	InstanceId string      `json:"instanceId"`
}

func (s *Server) PublicData(ctx *gin.Context) {

	var request publicRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	instanceId := os.Getenv("CONTAINER_ID")
	if instanceId == "" {
		instanceId = "Instance-Default"
	}

	var responseData publicResponse
	responseData.Status = "success"
	responseData.Message = "Data received"
	responseData.Data = request
	responseData.InstanceId = instanceId

	ctx.JSON(http.StatusOK, responseData)
}
