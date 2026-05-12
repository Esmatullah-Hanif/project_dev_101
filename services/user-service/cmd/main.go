package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/user-service/internal/handler"
	"github.com/example/microservices/user-service/internal/repository"
	"github.com/example/microservices/user-service/internal/router"
	"github.com/example/microservices/user-service/internal/service"
	"github.com/example/microservices/shared/pkg/config"
	"github.com/example/microservices/shared/pkg/database"
	"github.com/example/microservices/shared/pkg/logger"
)

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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	router.RegisterRoutes(engine, userHandler, cfg.CORSOrigins, cfg.JWTSecret)

	port := cfg.UserServicePort
	logger.Info("User service starting", "port", port)

	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
