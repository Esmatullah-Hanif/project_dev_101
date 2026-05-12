# Setup and Running Guide

## Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **PostgreSQL 15+** OR **Supabase Account**: [Create Supabase Account](https://supabase.com)
- **Docker and Docker Compose** (Optional, for containerized setup): [Install Docker](https://docs.docker.com/get-docker/)
- **Make**: Usually pre-installed on macOS/Linux, install on Windows via [GnuWin32](http://gnuwin32.sourceforge.net/)

## Quick Start (Local Development)

### 1. Clone and Setup Environment

```bash
cd /path/to/project
cp .env.example .env
```

### 2. Configure Database

#### Option A: Local PostgreSQL
```bash
# Create database
createdb microservices

# Create user if needed
createuser postgres

# Apply migrations
psql microservices < migrations/001_create_users_table.sql

# Update .env
DATABASE_URL=postgresql://postgres@localhost:5432/microservices
```

#### Option B: Supabase (Recommended)
```bash
# 1. Go to https://supabase.com/dashboard
# 2. Create new project
# 3. Copy connection string from Project Settings > Database
# 4. Update .env
DATABASE_URL=postgresql://postgres:[password]@[project].supabase.co:5432/postgres

# 5. Run migrations
psql $DATABASE_URL < migrations/001_create_users_table.sql
```

#### Option C: Docker Compose
```bash
# Start PostgreSQL container
docker-compose up -d postgres

# Wait for health check (5-10 seconds)
docker-compose logs postgres

# Apply migrations (inside container)
docker exec microservices-db psql -U postgres -d microservices -f /docker-entrypoint-initdb.d/001_create_users_table.sql
```

### 3. Update Configuration

Edit `.env` file:

```env
# Database (required)
DATABASE_URL=postgresql://user:password@host:port/database

# Ports (default values shown)
GATEWAY_PORT=8000
AUTH_SERVICE_PORT=8001
USER_SERVICE_PORT=8002

# Security (change in production!)
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=3600
REFRESH_TOKEN_EXPIRATION=604800

# Application
ENVIRONMENT=development
LOG_LEVEL=debug
CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 4. Run Services (Local Development)

Open **3 terminals** and run each in a separate one:

**Terminal 1 - Auth Service:**
```bash
make run-auth
# Output: Auth service starting on port 8001
```

**Terminal 2 - User Service:**
```bash
make run-user
# Output: User service starting on port 8002
```

**Terminal 3 - API Gateway:**
```bash
make run-gateway
# Output: API Gateway starting on port 8000
```

### 5. Generate Swagger Documentation (Optional)

For interactive API documentation, install the swag tool and generate docs:

```bash
# Install swag CLI tool
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation for all services
make swagger

# Or generate for specific services
make swagger-auth   # Auth service only
make swagger-user   # User service only
```

Then access the Swagger UI:
- Auth Service: http://localhost:8001/swagger/index.html
- User Service: http://localhost:8002/swagger/index.html

See [docs/SWAGGER.md](SWAGGER.md) for complete Swagger documentation.

### 6. Verify Setup

Test health endpoint:
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

## Docker Compose Setup

### Start All Services

```bash
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Auth Service (port 8001)
- User Service (port 8002)
- API Gateway (port 8000)

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service
docker-compose logs -f user-service
docker-compose logs -f api-gateway
docker-compose logs -f postgres
```

### Stop Services

```bash
docker-compose down
```

### Remove Data

```bash
# Stop containers and remove volumes
docker-compose down -v
```

## Testing the APIs

### 1. Sign Up

```bash
curl -X POST http://localhost:8000/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
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

Save the `access_token` for subsequent requests.

### 2. Sign In

```bash
curl -X POST http://localhost:8000/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 3. Get User Profile

```bash
export TOKEN="<access_token_from_signup>"

curl -X GET http://localhost:8000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $TOKEN"
```

### 4. List Users

```bash
curl -X GET "http://localhost:8000/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer $TOKEN"
```

### 5. Update User Profile

```bash
curl -X PUT http://localhost:8000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Software developer",
    "avatar_url": "https://example.com/avatar.jpg"
  }'
```

### 6. Refresh Token

```bash
curl -X POST http://localhost:8000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "<refresh_token_from_signup>"
  }'
```

### 7. Delete User

```bash
curl -X DELETE http://localhost:8000/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Authorization: Bearer $TOKEN"
```

## Development Tasks

### Build All Services

```bash
make build
```

Outputs binaries to:
- `services/api-gateway/bin/gateway`
- `services/auth-service/bin/auth`
- `services/user-service/bin/user`

### Run Tests

```bash
make test
```

Tests located in `*_test.go` files within each service.

### Lint Code

```bash
make lint
```

Uses golangci-lint configured in `.golangci.yml`.

### Clean Build Artifacts

```bash
make clean
```

## Database Migrations

### View Current Schema

```bash
psql $DATABASE_URL -c "\dt"  # Tables
psql $DATABASE_URL -c "\di"  # Indexes
```

### Create New Migration

1. Create file: `migrations/002_your_change.sql`
2. Write migration:
```sql
/*
  # Description

  1. Changes
    - List changes here
*/

ALTER TABLE users ADD COLUMN new_column TEXT;
```

3. Apply migration:
```bash
psql $DATABASE_URL -f migrations/002_your_change.sql
```

## Environment Variables

### Required
- `DATABASE_URL`: PostgreSQL connection string

### Ports
- `GATEWAY_PORT`: API Gateway port (default: 8000)
- `AUTH_SERVICE_PORT`: Auth Service port (default: 8001)
- `USER_SERVICE_PORT`: User Service port (default: 8002)

### JWT
- `JWT_SECRET`: Secret for signing tokens (change in production!)
- `JWT_EXPIRATION`: Token lifetime in seconds (default: 3600)
- `REFRESH_TOKEN_EXPIRATION`: Refresh token lifetime (default: 604800)

### Application
- `ENVIRONMENT`: development/staging/production
- `LOG_LEVEL`: debug/info/warn/error
- `CORS_ORIGINS`: Comma-separated list of allowed origins

## Common Issues

### Port Already in Use
```bash
# Find process using port
lsof -i :8000

# Kill process (macOS/Linux)
kill -9 <PID>

# Or change port in .env
GATEWAY_PORT=8080
```

### Database Connection Failed
```bash
# Check connection string
echo $DATABASE_URL

# Verify PostgreSQL is running
pg_isready -h localhost -p 5432

# Test connection
psql $DATABASE_URL -c "SELECT 1"
```

### Module Not Found Error
```bash
# Ensure you're in project root
cd /path/to/project

# Download dependencies
go mod download

# Verify go.mod
cat go.mod
```

### Service Not Starting
```bash
# Check Go installation
go version

# Rebuild
make clean && make build

# Run with verbose output
make run-auth
```

## Next Steps

1. **Explore the Code**: Review service implementations
2. **Add Features**: Extend with new services following the pattern
3. **Write Tests**: Add unit and integration tests
4. **Deploy**: Configure for your deployment platform
5. **Monitor**: Set up logging and monitoring
6. **Secure**: Update JWT_SECRET and review security settings

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Supabase Documentation](https://supabase.com/docs)
- [JWT Guide](https://jwt.io/)
- [Docker Documentation](https://docs.docker.com/)
