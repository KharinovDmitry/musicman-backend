package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/musicman-backend/internal/http/handler/music"
	"github.com/musicman-backend/internal/http/handler/payment"
	"github.com/musicman-backend/internal/http/handler/purchase"
	swaggerFiles "github.com/swaggo/files"
	"log/slog"
	"time"

	"github.com/musicman-backend/internal/di"
	"github.com/musicman-backend/internal/http/handler/auth"
	"github.com/musicman-backend/internal/http/handler/health"
	"github.com/musicman-backend/internal/http/handler/profile"
	"github.com/musicman-backend/internal/http/middleware"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(container *di.Container) *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	slog.Info("cors applied")

	healthHandler := health.NewHandler()
	router.GET("/health", healthHandler.Health)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	apiV1.Use(
		gin.Recovery(),
		gin.Logger(),
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

	musicHandler := music.New(container.Service.Music, container.Service.Purchase)

	apiV1.Group("/samples").
		//Use(authMiddleware).
		GET("", musicHandler.GetSamples).
		GET("/:id", musicHandler.GetSample).
		PUT("/:id", musicHandler.UpdateSample).
		POST("/:id", musicHandler.UploadAudio).
		DELETE("/:id", musicHandler.DeleteSample).
		POST("", musicHandler.CreateSample)

	apiV1.Group("/packs").
		//Use(authMiddleware).
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

	purchaseHandler := purchase.New(container.Service.Purchase, container.Service.Music)

	// Покупка семпла
	apiV1.Group("/samples").
		Use(authMiddleware).
		POST("/:id/purchase", purchaseHandler.PurchaseSample)

	// Получение всех покупок
	purchasesGroup := apiV1.Group("/purchases")
	purchasesGroup.Use(authMiddleware)
	{
		purchasesGroup.GET("", purchaseHandler.GetUserPurchases)
	}

	return router
}
