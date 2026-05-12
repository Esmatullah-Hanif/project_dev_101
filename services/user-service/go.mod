module github.com/example/microservices/user-service

go 1.21

require (
	github.com/example/microservices/shared v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/jackc/pgx/v5 v5.5.0
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.17.0
	github.com/rs/zerolog v1.31.0
)

replace github.com/example/microservices/shared => ../../shared
