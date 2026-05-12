# Microservices Architecture Guide

## Overview

This project implements a modern microservices architecture using Go and the Gin web framework. The design emphasizes scalability, maintainability, and clear separation of concerns.

## System Design

### Services

#### API Gateway (Port 8000)
**Responsibility**: Request routing and aggregation
- Single entry point for all client requests
- Routes requests to appropriate microservices
- Handles CORS headers and request IDs
- Forwards Authorization headers to downstream services

#### Auth Service (Port 8001)
**Responsibility**: User authentication and token management
- User registration with password hashing (bcrypt)
- User login and credential validation
- JWT token generation and validation
- Token refresh mechanism
- Uses persistent user storage in PostgreSQL

#### User Service (Port 8002)
**Responsibility**: User profile management
- Retrieve user information
- List users with pagination
- Update user profiles
- Soft delete users (preserve data integrity)
- Requires JWT authentication for mutations

### Database Layer

**PostgreSQL** (via Supabase or local instance)
- Stores user credentials and profiles
- Supports soft deletes for audit trails
- Efficient indexing on frequently queried fields
- Single source of truth for all user data

## Communication Patterns

### Service-to-Service Communication

```
Client ──→ API Gateway ──→ Auth Service ──→ PostgreSQL
                       ├──→ User Service ──→ PostgreSQL
```

- **Synchronous HTTP/REST** for client-facing APIs
- **HTTP forwarding** through API Gateway for internal communication
- Services communicate directly with database (no service-to-service DB queries)

## Layered Architecture (Per Service)

Each microservice follows a clean layered architecture:

```
┌─────────────────────────┐
│  HTTP Handler Layer     │  (Gin routes, request/response handling)
├─────────────────────────┤
│  Business Logic Layer   │  (Service layer with core logic)
├─────────────────────────┤
│  Data Access Layer      │  (Repository pattern, database queries)
├─────────────────────────┤
│  Database Layer         │  (PostgreSQL via pgx)
└─────────────────────────┘
```

### Handler Layer
- Entry point for HTTP requests
- Validates request format and structure
- Delegates business logic to service layer
- Returns standardized JSON responses
- Handles error mapping to HTTP status codes

Example: `auth-service/internal/handler/auth_handler.go`

### Service Layer
- Contains core business logic
- No HTTP knowledge (pure Go functions)
- Manages domain operations
- Orchestrates calls to repositories
- Performs validation and transformation
- Generates tokens, hashes passwords, etc.

Example: `auth-service/internal/service/auth_service.go`

### Repository Layer
- Encapsulates all database access
- Implements query builders
- Handles connection pooling via pgx
- Manages transactions if needed
- Returns domain models (not raw DB rows)

Example: `auth-service/internal/repository/user_repository.go`

### Model Layer
- Domain objects (User, AuthResponse, etc.)
- DTOs for request/response binding
- Type-safe representations of data

Example: `auth-service/internal/model/user.go`

## Middleware Stack

### Request ID Middleware
- Generates unique ID for each request
- Propagates through logs
- Useful for tracing requests across services

### Logging Middleware
- Logs method, path, status code
- Includes request ID for correlation
- Uses structured logging (zerolog)

### CORS Middleware
- Configurable origins
- Handles preflight requests
- Sets appropriate security headers

### Recovery Middleware
- Catches panics to prevent crashes
- Returns graceful error responses
- Logs panic details

### Auth Middleware (Protected Routes)
- Validates JWT tokens
- Extracts user ID from claims
- Sets authenticated user context
- Returns 401 for missing/invalid tokens

## Error Handling

### Custom Error Types
```go
type AppError struct {
    Code    ErrorCode
    Message string
    Details error
}

const (
    ValidationError ErrorCode = "VALIDATION_ERROR"
    AuthError       ErrorCode = "AUTH_ERROR"
    NotFoundError   ErrorCode = "NOT_FOUND"
    ConflictError   ErrorCode = "CONFLICT"
    InternalError   ErrorCode = "INTERNAL_ERROR"
    ForbiddenError  ErrorCode = "FORBIDDEN"
)
```

### Error Mapping
- Validation errors → 400 Bad Request
- Auth errors → 401 Unauthorized
- Not found → 404 Not Found
- Conflicts → 409 Conflict
- Forbidden → 403 Forbidden
- Internal errors → 500 Internal Server Error

### Response Format
```json
{
  "success": false,
  "message": "Human-readable error message",
  "error": "Detailed error description"
}
```

## Authentication Flow

### Sign Up
```
1. Client sends POST /auth/signup with email/password
2. Auth service validates input
3. Check if email already exists
4. Hash password with bcrypt
5. Insert user into database
6. Generate JWT tokens
7. Return tokens and user data
```

### Sign In
```
1. Client sends POST /auth/signin with email/password
2. Auth service looks up user by email
3. Compare provided password with hash
4. If match, generate JWT tokens
5. Return tokens and user data
```

### Protected Route Access
```
1. Client sends request with Authorization: Bearer <token>
2. Auth middleware validates token signature
3. Extract user ID from token claims
4. Allow request to proceed
5. Handler uses user ID from context
```

## Data Models

### User Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  first_name TEXT DEFAULT '',
  last_name TEXT DEFAULT '',
  bio TEXT DEFAULT '',
  avatar_url TEXT DEFAULT '',
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE  -- Soft delete
);
```

### Indexes
- `idx_users_email` - Fast email lookup for login
- `idx_users_deleted_at` - Efficient soft delete filtering
- `idx_users_created_at` - Timeline queries and sorting

## Configuration Management

### Environment-Driven Configuration
- All config loaded from environment variables
- Fallback defaults for development
- Separated concerns (each service reads own config)

### Config Values
```go
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
```

## Scalability Considerations

### Horizontal Scaling
- Stateless services can be replicated
- Load balancer distributes requests
- Shared database for consistency
- No in-memory session state

### Database Connection Pooling
- pgx pool with configurable min/max connections
- Connection health checks
- Automatic connection recycling
- Prevents connection exhaustion

### Rate Limiting (Future Enhancement)
- Implement per-IP or per-user limits
- Protect against abuse
- Return 429 Too Many Requests when exceeded

### Caching (Future Enhancement)
- Redis for session and token caching
- Reduces database load
- Faster response times for frequently accessed data

## Testing Strategy

### Unit Tests
- Test service layer logic in isolation
- Mock repository interfaces
- Test error scenarios
- Focus on business logic

### Integration Tests
- Test handler → service → repository flow
- Use test database fixtures
- Verify database queries work correctly
- Test error handling end-to-end

### API Tests
- Test HTTP endpoints through full stack
- Verify request/response format
- Test authentication flows
- Check status codes

## Deployment Architecture

### Docker Containers
- Each service runs in separate container
- Multi-stage builds for minimal image size
- Alpine Linux for small footprint
- Health checks for orchestration

### Docker Compose (Local/Development)
```yaml
Services:
- postgres (database)
- api-gateway
- auth-service
- user-service
```

### Production Deployment
- Kubernetes for orchestration
- Helm charts for configuration
- Persistent volumes for database
- Ingress for external routing
- Service mesh for communication (optional)

## Security Architecture

### Authentication
- JWT tokens (HS256 algorithm)
- Short expiration (1 hour)
- Refresh tokens for longer sessions
- Token validation on protected routes

### Password Security
- Bcrypt hashing with salt
- Cost factor of 10 (configurable)
- Never stored in plain text
- Validated before storage

### Data Protection
- HTTPS for all communication
- CORS for origin validation
- Soft deletes preserve audit trails
- No sensitive data in logs

### Future Enhancements
- Rate limiting
- Request validation and sanitization
- Input length limits
- SQL injection prevention (parameterized queries)
- HTTPS enforcement
- API key authentication for service-to-service

## Monitoring and Observability

### Logging
- Structured logs (JSON format possible)
- Request correlation via request ID
- Different log levels (debug, info, warn, error)
- Service and timestamp included

### Metrics (Future Enhancement)
- Request latency
- Error rates
- Database query performance
- Service health status

### Tracing (Future Enhancement)
- Distributed tracing across services
- Request ID propagation
- Timeline visualization
- Performance bottleneck identification

## Decision Rationale

### Why Microservices?
- **Scalability**: Scale individual services based on load
- **Maintainability**: Clear boundaries and responsibilities
- **Technology Flexibility**: Each service can use different tech stack
- **Deployment Independence**: Deploy services without coordinating

### Why Gin Framework?
- **Performance**: Fast HTTP router and middleware support
- **Simplicity**: Minimal boilerplate, clean API
- **Community**: Large ecosystem of middleware
- **Modularity**: Easy to compose middleware

### Why PostgreSQL?
- **Reliability**: ACID transactions, data consistency
- **Features**: JSON support, full-text search, advanced indexing
- **Scalability**: Handles large datasets efficiently
- **Maturity**: Battle-tested in production systems

### Why JWT Tokens?
- **Stateless**: No session store needed
- **Scalable**: Works across multiple service instances
- **Standard**: Well-understood, widely supported
- **Flexible**: Can encode arbitrary claims

## Future Enhancements

1. **Service Mesh**: Istio for advanced routing and security
2. **Event-Driven Architecture**: Message queues (RabbitMQ, Kafka)
3. **gRPC**: Internal service communication (lower latency)
4. **Caching**: Redis for frequently accessed data
5. **Search**: Elasticsearch for advanced user search
6. **Monitoring**: Prometheus + Grafana for observability
7. **Tracing**: Jaeger for distributed tracing
8. **CI/CD**: GitHub Actions for automated testing and deployment
9. **API Gateway**: Kong or Traefik for more advanced features
10. **GraphQL**: Alternative to REST for flexible queries
