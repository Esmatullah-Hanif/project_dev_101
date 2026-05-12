package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	_ "github.com/example/microservices/user-service/docs"
	"github.com/example/microservices/user-service/internal/handler"
	"github.com/example/microservices/shared/pkg/middleware"
)

func RegisterRoutes(engine *gin.Engine, userHandler *handler.UserHandler, corsOrigins string, jwtSecret string) {
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS(corsOrigins))
	engine.Use(middleware.Recovery())

	engine.GET("/health", userHandler.HealthCheck)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userGroup := engine.Group("/api/v1/users")
	{
		userGroup.GET("", userHandler.ListUsers)
		userGroup.GET("/:id", userHandler.GetUser)
		userGroup.PUT("/:id", middleware.AuthMiddleware(jwtSecret), userHandler.UpdateUser)
		userGroup.DELETE("/:id", middleware.AuthMiddleware(jwtSecret), userHandler.DeleteUser)
	}
}
