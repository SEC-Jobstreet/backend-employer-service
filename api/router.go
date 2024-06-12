package api

import (
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

	authRoutes := router.Group("/api/v1")

	authRoutes.POST("/create_enterprise", s.CreateEnterprise)
	authRoutes.POST("/update_enterprise",
		middleware.AuthMiddleware(s.config, []string{"employer"}),
		s.UpdateEnterprise)

	authRoutes.GET("/get_enterprise_by_employer", middleware.AuthMiddleware(s.config, []string{"employers", "admin"}), s.GetEnterpriseByEmployer)
	authRoutes.GET("/get_enterprise_by_id/:id", middleware.AuthMiddleware(s.config, []string{"employers", "admin"}), s.GetEnterpriseByID)

	// public api
	apiPublic := router.Group("/api/public")
	apiPublic.POST("/loadbalance", s.PublicData)

	s.router = router
}
