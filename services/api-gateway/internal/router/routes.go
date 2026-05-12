package router

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/api-gateway/internal/handler"
	"github.com/example/microservices/api-gateway/internal/proxy"
	"github.com/example/microservices/shared/pkg/config"
	"github.com/example/microservices/shared/pkg/middleware"
)

func RegisterRoutes(engine *gin.Engine, cfg *config.Config) {
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS(cfg.CORSOrigins))
	engine.Use(middleware.Recovery())

	healthHandler := handler.NewHealthHandler()
	engine.GET("/health", healthHandler.HealthCheck)

	authServiceURL := fmt.Sprintf("http://localhost:%s", cfg.AuthServicePort)
	userServiceURL := fmt.Sprintf("http://localhost:%s", cfg.UserServicePort)

	authGroup := engine.Group("/api/v1/auth")
	{
		authGroup.POST("/signup", proxy.ForwardToAuthService(authServiceURL))
		authGroup.POST("/signin", proxy.ForwardToAuthService(authServiceURL))
		authGroup.POST("/refresh", proxy.ForwardToAuthService(authServiceURL))
	}

	userGroup := engine.Group("/api/v1/users")
	{
		userGroup.GET("", proxy.ForwardToUserService(userServiceURL))
		userGroup.GET("/:id", proxy.ForwardToUserService(userServiceURL))
		userGroup.PUT("/:id", proxy.ForwardToUserService(userServiceURL))
		userGroup.DELETE("/:id", proxy.ForwardToUserService(userServiceURL))
	}
}
