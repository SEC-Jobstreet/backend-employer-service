package api

import (
	"net/http"

	"github.com/SEC-Jobstreet/backend-employer-service/api/dto"
	"github.com/SEC-Jobstreet/backend-employer-service/api/middleware"
	db "github.com/SEC-Jobstreet/backend-employer-service/db/sqlc"
	"github.com/gin-gonic/gin"
)

type SignUpInput struct {
	Email     string `json:"email" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

type ConfirmSignUpInput struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var cognitoClient middleware.CognitoClient

func init() {
	clientId := "27qr5vcjamku7t7jle3lg7poc8"
	cognitoClient = middleware.NewCognitoClient("ap-southeast-1", clientId)
}

func (s *Server) signUp(ctx *gin.Context) {
	var input SignUpInput

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, result := cognitoClient.SignUp(dto.Employer{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Address:   input.Address,
		Password:  input.Password,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employerProfileParams := db.CreateEmployerProfileParams{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.Phone,
		Address:   input.Address,
	}

	_, err = s.store.CreateEmployerProfile(ctx, employerProfileParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) confirmSignUp(ctx *gin.Context) {
	var input ConfirmSignUpInput

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, result := cognitoClient.ConfirmSignUp(input.Email, input.Code)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	confirm_err := s.store.UpdateEmailConfirmed(ctx, input.Email)

	if confirm_err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": confirm_err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (s *Server) login(ctx *gin.Context) {
	var input LoginInput

	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, accessToken, refreshToken := cognitoClient.Login(input.Email, input.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	confirm_err := s.store.UpdateEmailConfirmed(ctx, input.Email)

	if confirm_err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": confirm_err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (s *Server) jwt_test(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"jwt_test": "success",
	})
}
