package api

import (
	"os"
	"time"

	"github.com/SEC-Jobstreet/backend-employer-service/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "content-type", "accept", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRoutes := router.Group("/api/v1/authentications")

	authRoutes.POST("/sign_up", s.signUp)
	authRoutes.POST("/confirm_email", s.confirmSignUp)
	authRoutes.POST("/login", s.login)

	protectedRoutes := router.Group("/protected")
	protectedRoutes.Use(middleware.CognitoAuthMiddleware(os.Getenv("COGNITO_JWKS_URL")))

	protectedRoutes.POST("/test")

	s.router = router
}
