# Go Gin Microservices Boilerplate

A production-ready microservices architecture built with Go and Gin framework featuring authentication, user management, and an API gateway pattern.

## Architecture Overview

This boilerplate implements a clean, scalable microservices architecture with the following components:

```
┌─────────────────┐
│   API Gateway   │ (Port 8000) - Routes requests to downstream services
└────────┬────────┘
         │
    ┌────┴────┐
    │          │
┌───▼──┐  ┌──▼────┐
│ Auth │  │ Users  │ (Ports 8001-8002) - Domain services
│Service│  │Service │
└───┬──┘  └──┬────┘
    │        │
    └────┬───┘
         │
    ┌────▼──────┐
    │ PostgreSQL │ (via Supabase) - Shared data layer
    └───────────┘
```

## Services

### API Gateway
- **Port**: 8000
- **Purpose**: Single entry point for all requests, routes to downstream services
- **Routes**:
  - `POST /api/v1/auth/*` → Auth Service
  - `GET/POST/PUT/DELETE /api/v1/users/*` → User Service

### Auth Service
- **Port**: 8001
- **Purpose**: User authentication and JWT token management
- **Features**:
  - User registration (signup)
  - User login (signin)
  - Token refresh
  - Password hashing with bcrypt
  - JWT token generation

### User Service
- **Port**: 8002
- **Purpose**: User profile management
- **Features**:
  - Get user profile
  - List all users with pagination
  - Update user profile
  - Delete user (soft delete)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+ (or use Supabase)
- Docker and Docker Compose (optional)

### Environment Setup

1. Clone the repository and navigate to the project root:

```bash
cd /path/to/project
```

2. Copy the environment template:

```bash
cp .env.example .env
```

3. Update `.env` with your configuration:

```env
# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/microservices

# Ports
GATEWAY_PORT=8000
AUTH_SERVICE_PORT=8001
USER_SERVICE_PORT=8002

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRATION=3600
REFRESH_TOKEN_EXPIRATION=604800

# Application
ENVIRONMENT=development
LOG_LEVEL=debug
CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

### Database Setup

1. **Using PostgreSQL locally**:

```bash
createdb microservices
psql microservices < migrations/001_create_users_table.sql
```

2. **Using Docker Compose** (includes PostgreSQL):

```bash
docker-compose up -d postgres
docker exec microservices-db psql -U postgres -d microservices -f /docker-entrypoint-initdb.d/001_create_users_table.sql
```

3. **Using Supabase**:
   - Create a Supabase project
   - Get your database connection URL
   - Update `DATABASE_URL` in `.env`
   - Run migrations via Supabase dashboard or psql

### Running Services

#### Option 1: Local Development (3 terminals)

Terminal 1 - Auth Service:
```bash
make run-auth
```

Terminal 2 - User Service:
```bash
make run-user
```

Terminal 3 - API Gateway:
```bash
make run-gateway
```

#### Option 2: Docker Compose

```bash
docker-compose up
```

This will start all services with PostgreSQL.

#### Option 3: Build and Run Binaries

```bash
make build
./services/api-gateway/bin/gateway
./services/auth-service/bin/auth
./services/user-service/bin/user
```

### Verification

Test the API Gateway health endpoint:

```bash
curl http://localhost:8000/health
```

Expected response:
```json
{
  "success": true,
  "message": "API Gateway is healthy",
  "data": {
    "service": "api-gateway",
    "status": "healthy"
  }
}
```

### API Documentation with Swagger

Each service includes interactive Swagger/OpenAPI documentation:

**Generate documentation** (requires `swag` CLI):
```bash
# Install swag tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs for all services
make swagger
```

**Access the documentation:**
- Auth Service: http://localhost:8001/swagger/index.html
- User Service: http://localhost:8002/swagger/index.html

For detailed Swagger setup and usage, see [docs/SWAGGER.md](docs/SWAGGER.md).

## API Endpoints

### Authentication

#### Sign Up
```bash
POST /api/v1/auth/signup

{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

#### Sign In
```bash
POST /api/v1/auth/signin

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Refresh Token
```bash
POST /api/v1/auth/refresh

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### User Management

#### Get User Profile
```bash
GET /api/v1/users/:id

Authorization: Bearer <access_token>
```

#### List All Users
```bash
GET /api/v1/users?page=1&page_size=10

Authorization: Bearer <access_token>
```

Response:
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "bio": "Software developer",
      "avatar_url": "https://example.com/avatar.jpg",
      "is_active": true,
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_page": 5
  }
}
```

#### Update User Profile
```bash
PUT /api/v1/users/:id

Authorization: Bearer <access_token>

{
  "first_name": "Jane",
  "last_name": "Doe",
  "bio": "Senior developer",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

#### Delete User
```bash
DELETE /api/v1/users/:id

Authorization: Bearer <access_token>
```

## Project Structure

```
.
├── Makefile                    # Build and development commands
├── go.mod                      # Root module definition
├── docker-compose.yml          # Docker Compose configuration
├── .env.example                # Environment template
├── migrations/                 # Database migrations
│   └── 001_create_users_table.sql
├── shared/                     # Shared libraries
│   ├── go.mod
│   └── pkg/
│       ├── logger/            # Structured logging
│       ├── config/            # Configuration loader
│       ├── database/          # Database connection
│       ├── response/          # HTTP response helpers
│       ├── errors/            # Custom error types
│       ├── middleware/        # Reusable Gin middleware
│       └── validator/         # Input validation
├── services/
│   ├── api-gateway/
│   │   ├── cmd/main.go
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── internal/
│   │       ├── handler/
│   │       ├── router/
│   │       └── proxy/
│   ├── auth-service/
│   │   ├── cmd/main.go
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── internal/
│   │       ├── handler/
│   │       ├── service/
│   │       ├── repository/
│   │       ├── router/
│   │       └── model/
│   └── user-service/
│       ├── cmd/main.go
│       ├── Dockerfile
│       ├── go.mod
│       └── internal/
│           ├── handler/
│           ├── service/
│           ├── repository/
│           ├── router/
│           └── model/
└── docs/                      # Documentation
```

## Architecture Patterns

### Layered Architecture
Each service follows a clean layered architecture:
- **Handler**: HTTP request/response handling
- **Service**: Business logic
- **Repository**: Data access layer
- **Model**: Domain objects and DTOs

### Middleware Stack
All services include consistent middleware:
- Request ID generation
- Structured logging
- CORS support
- Panic recovery
- JWT authentication (protected routes)

### Error Handling
Custom error types provide consistent error responses across services with proper HTTP status codes and error messages.

## Development Commands

```bash
# Build all services
make build

# Run individual services
make run-gateway
make run-auth
make run-user

# Run tests
make test

# Lint code
make lint

# Generate Swagger API documentation
make swagger
make swagger-auth  # Auth service only
make swagger-user  # User service only

# Clean build artifacts
make clean

# Docker Compose commands
make docker-up
make docker-down
make docker-logs
```

## Security Considerations

1. **JWT Tokens**: Change `JWT_SECRET` in production
2. **Password Hashing**: Uses bcrypt with default cost (10 rounds)
3. **Soft Deletes**: Users are soft-deleted to preserve referential integrity
4. **CORS**: Configure `CORS_ORIGINS` for your frontend domains
5. **Database**: Use Supabase for managed PostgreSQL with built-in security features

## Scaling and Extending

### Adding a New Service

1. Create a new directory: `services/my-service/`
2. Copy the structure from `auth-service` or `user-service`
3. Update `services/my-service/go.mod` with correct module path
4. Implement your handlers, services, and repositories
5. Create a Dockerfile
6. Add routes in API Gateway
7. Update docker-compose.yml

### Database Schema Changes

1. Create a new migration: `migrations/002_your_change.sql`
2. Apply with: `psql $DATABASE_URL -f migrations/002_your_change.sql`
3. Or use Supabase dashboard migrations

## Troubleshooting

### Database Connection Errors
- Verify `DATABASE_URL` is correct and database is running
- Check PostgreSQL credentials
- Ensure migrations have been applied

### Service Won't Start
- Check port availability: `lsof -i :8000` (replace with service port)
- Review service logs for detailed error messages
- Verify environment variables are set correctly

### Cross-Service Communication
- Ensure all services are running
- Check firewall rules allowing service-to-service communication
- Verify service URLs in API Gateway configuration

## Next Steps

1. Add unit tests for service layers
2. Implement integration tests
3. Set up CI/CD pipeline (GitHub Actions)
4. Add API documentation (Swagger/OpenAPI)
5. Implement caching layer (Redis)
6. Add monitoring and observability (Prometheus, Grafana)
7. Implement gRPC for internal service communication
8. Add request validation and sanitization
9. Implement rate limiting per user
10. Add email verification for sign-up

## License

MIT