package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/auth-service/internal/handler"
	"github.com/example/microservices/auth-service/internal/repository"
	"github.com/example/microservices/auth-service/internal/router"
	"github.com/example/microservices/auth-service/internal/service"
	"github.com/example/microservices/shared/pkg/config"
	"github.com/example/microservices/shared/pkg/database"
	"github.com/example/microservices/shared/pkg/logger"
)

// @title Auth Service API
// @version 1.0
// @description Authentication and token management service for the microservices architecture
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8001
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Type 'Bearer TOKEN' to correctly set the API Key"

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.LogLevel)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	logger.Info("Database connected", "url", cfg.DatabaseURL)

	tokenExp := 3600
	refreshExp := 604800

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, tokenExp, refreshExp)
	authHandler := handler.NewAuthHandler(authService)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	router.RegisterRoutes(engine, authHandler, cfg.CORSOrigins)

	port := cfg.AuthServicePort
	logger.Info("Auth service starting", "port", port)

	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
