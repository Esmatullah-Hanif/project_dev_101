.PHONY: help build run test clean lint migrate-up

GATEWAY_BIN := services/api-gateway/bin/gateway
AUTH_BIN := services/auth-service/bin/auth
USER_BIN := services/user-service/bin/user

help:
	@echo "Available commands:"
	@echo "  make build              - Build all services"
	@echo "  make run                - Run all services (use in separate terminals)"
	@echo "  make run-gateway        - Run API gateway"
	@echo "  make run-auth           - Run auth service"
	@echo "  make run-user           - Run user service"
	@echo "  make test               - Run tests for all services"
	@echo "  make lint               - Run linter"
	@echo "  make clean              - Clean build artifacts"
	@echo "  make migrate-up         - Run database migrations"
	@echo "  make swagger            - Generate Swagger documentation for all services"
	@echo "  make swagger-auth       - Generate Swagger docs for auth service"
	@echo "  make swagger-user       - Generate Swagger docs for user service"
	@echo "  make docker-up          - Start services with docker-compose"
	@echo "  make docker-down        - Stop docker-compose services"

build:
	@echo "Building API Gateway..."
	@mkdir -p $(GATEWAY_BIN)
	@cd services/api-gateway && go build -o ../../$(GATEWAY_BIN) ./cmd
	@echo "Building Auth Service..."
	@mkdir -p $(AUTH_BIN)
	@cd services/auth-service && go build -o ../../$(AUTH_BIN) ./cmd
	@echo "Building User Service..."
	@mkdir -p $(USER_BIN)
	@cd services/user-service && go build -o ../../$(USER_BIN) ./cmd
	@echo "Build complete!"

run-gateway:
	@cd services/api-gateway && go run ./cmd

run-auth:
	@cd services/auth-service && go run ./cmd

run-user:
	@cd services/user-service && go run ./cmd

test:
	@echo "Running tests..."
	@cd services/auth-service && go test -v -cover ./...
	@cd services/user-service && go test -v -cover ./...

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf services/*/bin
	@go clean -cache
	@echo "Clean complete!"

migrate-up:
	@echo "Running database migrations..."
	@psql $(DATABASE_URL) -f migrations/001_create_users_table.sql
	@echo "Migrations complete!"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

swagger:
	@echo "Generating Swagger documentation..."
	@command -v swag >/dev/null 2>&1 || { echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	@echo "Generating Auth Service docs..."
	@cd services/auth-service && swag init -g cmd/main.go -o ../../services/auth-service/docs
	@echo "Generating User Service docs..."
	@cd services/user-service && swag init -g cmd/main.go -o ../../services/user-service/docs
	@echo "Swagger docs generated successfully!"
	@echo ""
	@echo "Access Swagger UI at:"
	@echo "  Auth Service:  http://localhost:8001/swagger/index.html"
	@echo "  User Service:  http://localhost:8002/swagger/index.html"

swagger-auth:
	@echo "Generating Auth Service Swagger docs..."
	@command -v swag >/dev/null 2>&1 || { echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	@cd services/auth-service && swag init -g cmd/main.go -o ../../services/auth-service/docs
	@echo "Auth Service Swagger docs generated!"
	@echo "Access at: http://localhost:8001/swagger/index.html"

swagger-user:
	@echo "Generating User Service Swagger docs..."
	@command -v swag >/dev/null 2>&1 || { echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; exit 1; }
	@cd services/user-service && swag init -g cmd/main.go -o ../../services/user-service/docs
	@echo "User Service Swagger docs generated!"
	@echo "Access at: http://localhost:8002/swagger/index.html"
