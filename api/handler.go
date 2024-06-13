package api

import (
	"net/http"

	db "github.com/SEC-Jobstreet/backend-employer-service/db/sqlc"
	"github.com/SEC-Jobstreet/backend-employer-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type enterpriseCreationRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Field     string `json:"field" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Url       string `json:"url"`
	License   string `json:"license"`

	EmployerId   string `json:"employer_id" binding:"required"`
	EmployerRole string `json:"employer_role" binding:"required"`
}

func (s *Server) CreateEnterprise(ctx *gin.Context) {

	var request enterpriseCreationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	var enterpriseId uuid.UUID
	if request.Id == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}
		enterpriseId = id
	} else {
		id, err := uuid.Parse(request.Id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}
		enterpriseId = id
	}

	enterprise := db.CreateEnterpriseParams{
		ID: enterpriseId,
		Name: pgtype.Text{
			String: request.Name,
			Valid:  request.Name != "",
		},
		Country: pgtype.Text{
			String: request.Country,
			Valid:  request.Country != "",
		},
		Address: pgtype.Text{
			String: request.Address,
			Valid:  request.Address != "",
		},
		Latitude: pgtype.Text{
			String: request.Latitude,
			Valid:  request.Latitude != "",
		},
		Longitude: pgtype.Text{
			String: request.Longitude,
			Valid:  request.Longitude != "",
		},
		Field: pgtype.Text{
			String: request.Field,
			Valid:  request.Field != "",
		},
		Size: pgtype.Text{
			String: request.Size,
			Valid:  request.Size != "",
		},
		Url: pgtype.Text{
			String: request.Url,
			Valid:  request.Url != "",
		},
		License: pgtype.Text{
			String: request.License,
			Valid:  request.License != "",
		},
		EmployerID: pgtype.Text{
			String: request.EmployerId,
			Valid:  request.EmployerId != "",
		},
		EmployerRole: pgtype.Text{
			String: request.EmployerRole,
			Valid:  request.EmployerRole != "",
		},
	}

	res, err := s.store.CreateEnterprise(ctx, enterprise)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

type enterpriseUpdatingRequest struct {
	Id        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Country   string `json:"country" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Field     string `json:"field" binding:"required"`
	Size      string `json:"size" binding:"required"`
	Url       string `json:"url"`
	License   string `json:"license"`

	EmployerRole string `json:"employer_role"`
}

func (s *Server) UpdateEnterprise(ctx *gin.Context) {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	var request enterpriseUpdatingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	enterpriseId, err := uuid.Parse(request.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	enterprise := db.UpdateEnterpriseParams{
		ID: enterpriseId,
		EmployerID: pgtype.Text{
			String: currentUser.Username,
			Valid:  true,
		},
		Name: pgtype.Text{
			String: request.Name,
			Valid:  request.Name != "",
		},
		Country: pgtype.Text{
			String: request.Country,
			Valid:  request.Country != "",
		},
		Address: pgtype.Text{
			String: request.Address,
			Valid:  request.Address != "",
		},
		Latitude: pgtype.Text{
			String: request.Latitude,
			Valid:  request.Latitude != "",
		},
		Longitude: pgtype.Text{
			String: request.Longitude,
			Valid:  request.Longitude != "",
		},
		Field: pgtype.Text{
			String: request.Field,
			Valid:  request.Field != "",
		},
		Size: pgtype.Text{
			String: request.Size,
			Valid:  request.Size != "",
		},
		Url: pgtype.Text{
			String: request.Url,
			Valid:  request.Url != "",
		},
		License: pgtype.Text{
			String: request.License,
			Valid:  request.License != "",
		},
		EmployerRole: pgtype.Text{
			String: request.EmployerRole,
			Valid:  request.EmployerRole != "",
		},
	}

	res, err := s.store.UpdateEnterprise(ctx, enterprise)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) GetEnterpriseByEmployer(ctx *gin.Context) {
	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	res, err := s.store.GetEnterpriseByEmployerId(ctx, pgtype.Text{
		String: currentUser.Username,
		Valid:  true,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)
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

	res, err := s.store.GetEnterpriseById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}
