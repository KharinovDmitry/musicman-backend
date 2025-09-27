package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/internal/di"
	"github.com/musicman-backend/internal/http/handler/health"
)

func SetupRouter(container *di.Container) *gin.Engine {
	router := gin.New()

	healthHandler := health.NewHandler()
	router.GET("/health", healthHandler.Health)

	api := router.Group("/api")
	api.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.Default(),
	)

	return router
}
