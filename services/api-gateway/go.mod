module github.com/example/microservices/api-gateway

go 1.21

require (
	github.com/example/microservices/shared v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.31.0
)

replace github.com/example/microservices/shared => ../../shared
