# Implementation Guide - Go Gin Microservices Boilerplate

## Project Completion Summary

This document provides an overview of the complete microservices boilerplate implementation and how to get started.

## What Has Been Implemented

### 1. Core Infrastructure

- **Root-level Go module** (`go.mod`) with all necessary dependencies
- **Makefile** with targets for building, running, testing, and linting
- **Docker & Docker Compose** configuration for containerized deployment
- **Environment configuration** system with `.env` support
- **Git configuration** (.gitignore) for Go projects

### 2. Shared Libraries (`shared/`)

The `shared/` directory contains reusable packages used across all services:

#### Logger Package
- **File**: `shared/pkg/logger/logger.go`
- **Features**: Structured logging using zerolog
- **Usage**: Consistent log formatting and levels (debug, info, warn, error)

#### Config Package
- **File**: `shared/pkg/config/config.go`
- **Features**: Load configuration from environment variables
- **Usage**: Centralized config management across services

#### Database Package
- **File**: `shared/pkg/database/database.go`
- **Features**: PostgreSQL connection pooling with pgx
- **Usage**: Database initialization and connection management

#### Response Package
- **File**: `shared/pkg/response/response.go`
- **Features**: Standardized HTTP response helpers
- **Usage**: Consistent JSON responses with success/error states and pagination

#### Error Package
- **File**: `shared/pkg/errors/errors.go`
- **Features**: Custom error types and error codes
- **Usage**: Unified error handling across services

#### Validator Package
- **File**: `shared/pkg/validator/validator.go`
- **Features**: Input validation helpers (email, password, etc.)
- **Usage**: Request payload validation

#### Middleware Package
- **File**: `shared/pkg/middleware/middleware.go`
- **Features**: Reusable Gin middleware
- **Includes**:
  - Request ID generation
  - Structured logging
  - CORS handling
  - Panic recovery
  - JWT authentication

### 3. Auth Service (`services/auth-service/`)

Complete user authentication service with JWT tokens.

#### Structure
```
auth-service/
├── cmd/main.go                        # Service entrypoint
├── go.mod                             # Service module
├── Dockerfile                         # Container image
└── internal/
    ├── handler/auth_handler.go        # HTTP handlers
    ├── service/auth_service.go        # Business logic
    ├── repository/user_repository.go  # Database access
    ├── router/routes.go               # Route registration
    └── model/user.go                  # Domain models
```

#### Features
- User registration (POST `/api/v1/auth/signup`)
- User login (POST `/api/v1/auth/signin`)
- Token refresh (POST `/api/v1/auth/refresh`)
- Password hashing with bcrypt
- JWT token generation and validation

#### Database
- Creates users in PostgreSQL
- Stores hashed passwords
- Tracks creation and update timestamps

#### Error Handling
- Validation errors for missing/invalid input
- Auth errors for invalid credentials
- Conflict errors for duplicate emails
- Proper HTTP status codes

### 4. User Service (`services/user-service/`)

User profile management and listing.

#### Structure
```
user-service/
├── cmd/main.go                        # Service entrypoint
├── go.mod                             # Service module
├── Dockerfile                         # Container image
└── internal/
    ├── handler/user_handler.go        # HTTP handlers
    ├── service/user_service.go        # Business logic
    ├── repository/user_repository.go  # Database access
    ├── router/routes.go               # Route registration
    └── model/user.go                  # Domain models
```

#### Features
- Get user profile (GET `/api/v1/users/:id`)
- List users with pagination (GET `/api/v1/users`)
- Update user profile (PUT `/api/v1/users/:id`)
- Delete user - soft delete (DELETE `/api/v1/users/:id`)
- JWT authentication on mutations

#### Database
- Fetches user data from PostgreSQL
- Supports soft deletes (deleted_at timestamp)
- Pagination with limit and offset
- Indexed queries for performance

#### Error Handling
- Validation errors for missing user ID
- Not found errors for missing users
- Auth errors for missing JWT tokens
- Proper HTTP status codes

### 5. API Gateway (`services/api-gateway/`)

Single entry point routing requests to microservices.

#### Structure
```
api-gateway/
├── cmd/main.go                        # Gateway entrypoint
├── go.mod                             # Gateway module
├── Dockerfile                         # Container image
└── internal/
    ├── handler/health_handler.go      # Health check handler
    ├── router/routes.go               # Route definitions
    └── proxy/proxy.go                 # HTTP proxying logic
```

#### Features
- Single entry point on port 8000
- Routes auth requests to auth-service (port 8001)
- Routes user requests to user-service (port 8002)
- Health check endpoint (GET `/health`)
- Middleware stack (CORS, logging, request ID)
- Header forwarding to downstream services

#### No Database
- Gateway is stateless
- All data requests proxied to services
- No caching (future enhancement)

### 6. Database Layer

#### Schema
- **File**: `migrations/001_create_users_table.sql`
- **Table**: `users`
- **Columns**: id, email, password_hash, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at, deleted_at
- **Indexes**: email (unique), deleted_at, created_at

#### Soft Deletes
- Users are marked with `deleted_at` timestamp instead of removed
- Queries filter with `WHERE deleted_at IS NULL`
- Preserves audit trail and referential integrity

#### Connection
- Uses pgx for high-performance database access
- Connection pooling with configurable min/max connections
- Parameterized queries to prevent SQL injection

### 7. Documentation

#### README.md
- Project overview
- Architecture diagram
- Service descriptions
- Getting started guide
- API endpoint documentation
- Common development commands

#### docs/ARCHITECTURE.md
- System design and communication patterns
- Layered architecture explanation
- Error handling strategy
- Authentication flow
- Data models
- Security considerations
- Future enhancements

#### docs/DATABASE.md
- Connection setup (local, Supabase, Docker)
- Schema documentation
- Migration procedures
- Data access patterns
- Performance optimization
- Backup strategies
- Troubleshooting

#### docs/SETUP.md
- Step-by-step setup instructions
- Database configuration options
- Service startup (3 terminals, Docker, binaries)
- API testing examples
- Development commands
- Issue resolution

### 8. Examples and Scripts

#### examples/auth_requests.sh
- Demonstrates sign-up
- Shows sign-in
- Token refresh
- Error handling examples
- Extracts tokens for use with other examples

#### examples/user_requests.sh
- Uses tokens from auth_requests.sh
- Shows user profile retrieval
- Demonstrates listing with pagination
- Update profile examples
- Delete user (soft delete)
- Error cases

## Architecture Highlights

### Clean Layered Design
```
HTTP Request
    ↓
  Handler (validates input, calls service)
    ↓
  Service (business logic, calls repository)
    ↓
  Repository (database queries)
    ↓
  PostgreSQL Database
```

### Middleware Stack
Every service includes:
1. Request ID - for tracing
2. Logger - for request logging
3. CORS - for cross-origin requests
4. Recovery - for panic handling
5. Auth (protected routes) - for JWT validation

### Error Handling
- Custom error types with codes
- Automatic mapping to HTTP status codes
- Consistent response format
- Detailed error messages

### Configuration
- Environment variable driven
- Defaults for development
- Per-service configuration
- Easy for containerization

## Getting Started

### 1. Prerequisites
```bash
# Install Go 1.21+
go version

# Install Docker (optional)
docker --version

# Or install PostgreSQL locally
psql --version
```

### 2. Setup
```bash
# Navigate to project
cd /path/to/project

# Copy environment
cp .env.example .env

# Update .env with database URL
# For Supabase: DATABASE_URL=postgresql://...
# For local: createdb microservices && DATABASE_URL=postgresql://postgres@localhost:5432/microservices
```

### 3. Database
```bash
# Apply migrations
psql $DATABASE_URL < migrations/001_create_users_table.sql

# Or with Docker Compose
docker-compose up -d postgres
docker exec microservices-db psql -U postgres -d microservices -f /docker-entrypoint-initdb.d/001_create_users_table.sql
```

### 4. Run Services

#### Option A: Local (3 terminals)
```bash
# Terminal 1
make run-auth

# Terminal 2
make run-user

# Terminal 3
make run-gateway
```

#### Option B: Docker Compose
```bash
docker-compose up
```

#### Option C: Build and Run Binaries
```bash
make build
./services/api-gateway/bin/gateway
./services/auth-service/bin/auth
./services/user-service/bin/user
```

### 5. Test
```bash
# Sign up
curl -X POST http://localhost:8000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Or use examples
chmod +x examples/auth_requests.sh
./examples/auth_requests.sh
```

## Development Workflow

### Add a New Service
1. Copy `services/auth-service` or `services/user-service` to `services/my-service`
2. Update `services/my-service/go.mod` module path
3. Modify handler, service, and repository for your domain
4. Create Dockerfile
5. Add routes in API Gateway
6. Update docker-compose.yml

### Add Database Columns
1. Create `migrations/002_add_columns.sql`
2. Update model structs in services
3. Update repository queries
4. Run migration with: `psql $DATABASE_URL < migrations/002_add_columns.sql`

### Add New Endpoint
1. Create handler method in service
2. Create service method with business logic
3. Create repository method if database access needed
4. Register route in router
5. Update API Gateway if new service

## Common Commands

```bash
make build          # Build all services
make run-gateway    # Run gateway
make run-auth       # Run auth service
make run-user       # Run user service
make test           # Run tests
make lint           # Lint code
make clean          # Clean artifacts
make docker-up      # Start Docker Compose
make docker-down    # Stop Docker Compose
```

## Key Files Reference

### Entry Points
- `services/api-gateway/cmd/main.go` - Gateway startup
- `services/auth-service/cmd/main.go` - Auth service startup
- `services/user-service/cmd/main.go` - User service startup

### Handlers (HTTP Layer)
- `services/auth-service/internal/handler/auth_handler.go` - Auth endpoints
- `services/user-service/internal/handler/user_handler.go` - User endpoints
- `services/api-gateway/internal/handler/health_handler.go` - Gateway health

### Services (Business Logic)
- `services/auth-service/internal/service/auth_service.go` - Auth logic
- `services/user-service/internal/service/user_service.go` - User logic

### Repositories (Data Layer)
- `services/auth-service/internal/repository/user_repository.go` - Auth queries
- `services/user-service/internal/repository/user_repository.go` - User queries

### Shared Packages
- `shared/pkg/logger/logger.go` - Logging
- `shared/pkg/config/config.go` - Configuration
- `shared/pkg/database/database.go` - Database connection
- `shared/pkg/middleware/middleware.go` - Middleware
- `shared/pkg/response/response.go` - Response helpers
- `shared/pkg/errors/errors.go` - Error types
- `shared/pkg/validator/validator.go` - Validation

## What's Next?

1. **Testing**: Add unit tests in `*_test.go` files
2. **Integration Tests**: Test full flow from API to database
3. **CI/CD**: Set up GitHub Actions or GitLab CI
4. **API Documentation**: Add Swagger/OpenAPI specs
5. **Caching**: Add Redis for session/token caching
6. **Monitoring**: Add Prometheus metrics and Grafana dashboards
7. **Tracing**: Integrate Jaeger for distributed tracing
8. **Rate Limiting**: Protect against abuse
9. **Email Verification**: Add signup email confirmation
10. **Permissions**: Add role-based access control

## Security Notes

- Change `JWT_SECRET` in production
- Use HTTPS in production
- Validate and sanitize all inputs
- Keep dependencies updated: `go get -u`
- Use Supabase for managed database security
- Implement rate limiting
- Add request signing for API-to-API communication
- Use environment-specific secrets

## Performance Optimization

- Database indexes on frequently queried columns (done)
- Connection pooling with pgx (done)
- Structured logging without heavy disk I/O (done)
- Stateless services for easy horizontal scaling
- Future: Add caching layer (Redis)
- Future: Add gRPC for internal service communication

## Troubleshooting

### Port in Use
```bash
lsof -i :8000  # Check what's using port
kill -9 <PID>  # Kill process
# Or change port in .env
```

### Database Connection Failed
```bash
# Verify connection string
echo $DATABASE_URL

# Test connection
psql $DATABASE_URL -c "SELECT 1"

# Check migrations applied
psql $DATABASE_URL -c "\dt"
```

### Module Not Found
```bash
cd /project/root
go mod download
go mod tidy
```

## Support and Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [pgx Driver](https://github.com/jackc/pgx)
- [JWT Guide](https://jwt.io/)
- [Supabase Docs](https://supabase.com/docs)

## License

MIT - Feel free to use this boilerplate for your projects!

---

**Ready to start?** Follow the "Getting Started" section above or check docs/SETUP.md for detailed instructions.
