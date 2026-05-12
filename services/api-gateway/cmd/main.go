package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/api-gateway/internal/router"
	"github.com/example/microservices/shared/pkg/config"
	"github.com/example/microservices/shared/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init(cfg.LogLevel)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	router.RegisterRoutes(engine, cfg)

	port := cfg.GatewayPort
	logger.Info("API Gateway starting", "port", port)

	if err := engine.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
