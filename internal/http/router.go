package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/internal/http/handler/music"
	swaggerFiles "github.com/swaggo/files"
	"github.com/musicman-backend/internal/http/handler/payment"

	"github.com/musicman-backend/internal/di"
	"github.com/musicman-backend/internal/http/handler/auth"
	"github.com/musicman-backend/internal/http/handler/health"
	"github.com/musicman-backend/internal/http/handler/profile"
	"github.com/musicman-backend/internal/http/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(container *di.Container) *gin.Engine {
	router := gin.New()

	healthHandler := health.NewHandler()
	router.GET("/health", healthHandler.Health)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	apiV1.Use(
		gin.Recovery(),
		gin.Logger(),
		cors.Default(),
	)

	authGroup := apiV1.Group("/auth")
	{
		authHandler := auth.NewHandler(container.Service.Auth)
		authGroup.POST("/sign-up", authHandler.Register)
		authGroup.POST("/sign-in", authHandler.Login)
	}

	authMiddleware := middleware.AuthMiddleware(container.Service.Token)

	profileGroup := apiV1.Group("/profile")
	profileGroup.Use(authMiddleware)
	{
		profileHandler := profile.NewHandler(container.Repository.UserRepository)
		profileGroup.GET("/me", profileHandler.GetMyProfile)
	}

	musicHandler := music.New(container.Service.Music)

	apiV1.Group("/samples").
		Use(authMiddleware).
		GET("", musicHandler.GetSamples).
		GET("/:id", musicHandler.GetSample).
		PUT("", musicHandler.UpdateSample).
		DELETE("/:id", musicHandler.DeleteSample).
		POST("", musicHandler.CreateSample)

	apiV1.Group("/pack").
		Use(authMiddleware).
		GET("", musicHandler.GetPacks).
		GET("/:id", musicHandler.GetPack).
		PUT("", musicHandler.UpdatePack).
		DELETE("/:id", musicHandler.DeletePack).
		POST("", musicHandler.CreatePack)

	paymentsGroup := apiV1.Group("/payments")
	paymentsGroup.Use(authMiddleware)
	{
		paymentHandler := payment.NewHandler(container.Service.Payment, container.Repository.PaymentRepository)
		paymentsGroup.POST("/new", paymentHandler.NewPayment)
		paymentsGroup.GET("/history", paymentHandler.GetPayments)
	}

	return router
}
