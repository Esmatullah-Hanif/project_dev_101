package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	_ "github.com/example/microservices/auth-service/docs"
	"github.com/example/microservices/auth-service/internal/handler"
	"github.com/example/microservices/shared/pkg/middleware"
)

func RegisterRoutes(engine *gin.Engine, authHandler *handler.AuthHandler, corsOrigins string) {
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS(corsOrigins))
	engine.Use(middleware.Recovery())

	engine.GET("/health", authHandler.HealthCheck)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authGroup := engine.Group("/api/v1/auth")
	{
		authGroup.POST("/signup", authHandler.SignUp)
		authGroup.POST("/signin", authHandler.SignIn)
		authGroup.POST("/refresh", authHandler.RefreshToken)
	}
}
