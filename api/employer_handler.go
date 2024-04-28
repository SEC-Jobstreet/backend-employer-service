package api

import (
	"net/http"

	"github.com/SEC-Jobstreet/backend-employer-service/models"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) example(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "OK")
}

type enterpriseCreationRequest struct {
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Field     string `json:"field" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Url       string `json:"url"`
	License   string `json:"license"`

	EmployerRole string `json:"employer_role" binding:"required"`
}

func (s *Server) CreateEnterprise(ctx *gin.Context) {

	var request enterpriseCreationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	enterprise := &models.Enterprises{
		ID:        id,
		Name:      request.Name,
		Country:   request.Country,
		Address:   request.Address,
		Latitude:  request.Latitude,
		Longitude: request.Longitude,
		Field:     request.Field,
		Size:      request.Size,
		Url:       request.Url,

		EmployerID:   currentUser.Username,
		EmployerRole: request.EmployerRole,
	}
	s.store.Create(enterprise)

	ctx.JSON(http.StatusOK, enterprise)
}

func (s *Server) GetEnterpriseByEmployer(ctx *gin.Context) {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	enterprises := &[]models.Enterprises{}

	s.store.Where("employer_id = ?", currentUser.Username).Find(enterprises)

	ctx.JSON(http.StatusOK, enterprises)
}

type enterpriseIdRequest struct {
	EnterpriseID string `uri:"id" binding:"required"`
}

func (s *Server) GetEnterpriseByID(ctx *gin.Context) {
	var request enterpriseIdRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	id, err := uuid.Parse(request.EnterpriseID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	enterprise := &models.Enterprises{}

	s.store.Where("id = ?", id).Find(enterprise)

	ctx.JSON(http.StatusOK, enterprise)
}
