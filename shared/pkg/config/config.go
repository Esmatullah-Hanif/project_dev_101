package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL          string
	GatewayPort          string
	AuthServicePort      string
	UserServicePort      string
	JWTSecret            string
	JWTExpiration        string
	RefreshTokenExpiration string
	Environment          string
	LogLevel             string
	CORSOrigins          string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		GatewayPort:         getEnv("GATEWAY_PORT", "8000"),
		AuthServicePort:     getEnv("AUTH_SERVICE_PORT", "8001"),
		UserServicePort:     getEnv("USER_SERVICE_PORT", "8002"),
		JWTSecret:           getEnv("JWT_SECRET", "change-this-secret-in-production"),
		JWTExpiration:       getEnv("JWT_EXPIRATION", "3600"),
		RefreshTokenExpiration: getEnv("REFRESH_TOKEN_EXPIRATION", "604800"),
		Environment:         getEnv("ENVIRONMENT", "development"),
		LogLevel:            getEnv("LOG_LEVEL", "info"),
		CORSOrigins:         getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
