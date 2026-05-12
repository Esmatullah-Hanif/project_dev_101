# Swagger API Documentation

This project includes comprehensive Swagger/OpenAPI documentation for both microservices. Each service hosts its own interactive API documentation.

## Accessing the Documentation

### Auth Service
- **URL**: http://localhost:8001/swagger/index.html
- **API Base**: http://localhost:8001
- **Endpoints**: `/api/v1/auth/*`, `/health`

### User Service
- **URL**: http://localhost:8002/swagger/index.html
- **API Base**: http://localhost:8002
- **Endpoints**: `/api/v1/users/*`, `/health`

## Generating Swagger Docs

### Prerequisites

First, install the swag CLI tool:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify installation:

```bash
swag --version
```

### Generate All Documentation

```bash
make swagger
```

This generates Swagger documentation for both services:
- Auth Service docs: `services/auth-service/docs/`
- User Service docs: `services/user-service/docs/`

### Generate Individual Service Docs

**Auth Service only:**
```bash
make swagger-auth
# or
cd services/auth-service && swag init -g cmd/main.go
```

**User Service only:**
```bash
make swagger-user
# or
cd services/user-service && swag init -g cmd/main.go
```

## API Documentation Structure

### Auth Service Endpoints

#### Sign Up
```
POST /api/v1/auth/signup

Request:
{
  "email": "user@example.com",
  "password": "password123"
}

Response (201):
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
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
```
POST /api/v1/auth/signin

Request:
{
  "email": "user@example.com",
  "password": "password123"
}

Response (200):
{
  "success": true,
  "message": "Login successful",
  "data": { ... same as signup ... }
}
```

#### Refresh Token
```
POST /api/v1/auth/refresh

Request:
{
  "refresh_token": "eyJhbGc..."
}

Response (200):
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": { ... auth response with new tokens ... }
}
```

### User Service Endpoints

#### Get User Profile
```
GET /api/v1/users/{id}

Headers:
Authorization: Bearer <access_token>

Response (200):
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
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
}
```

#### List Users
```
GET /api/v1/users?page=1&page_size=10

Response (200):
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [ ... array of users ... ],
  "meta": {
    "page": 1,
    "page_size": 10,
    "total": 50,
    "total_page": 5
  }
}
```

#### Update User Profile
```
PUT /api/v1/users/{id}

Headers:
Authorization: Bearer <access_token>

Request:
{
  "first_name": "Jane",
  "last_name": "Smith",
  "bio": "Senior developer",
  "avatar_url": "https://example.com/new-avatar.jpg"
}

Response (200):
{
  "success": true,
  "message": "User updated successfully",
  "data": { ... updated user object ... }
}
```

#### Delete User
```
DELETE /api/v1/users/{id}

Headers:
Authorization: Bearer <access_token>

Response (204):
No content
```

## Testing with Swagger UI

1. **Open the Swagger UI** for the desired service
2. **Authorize** (if needed for protected endpoints):
   - Click the "Authorize" button
   - Enter the Bearer token obtained from signup/signin
   - Click "Authorize"
3. **Try out an endpoint**:
   - Click "Try it out" button on any endpoint
   - Fill in required parameters
   - Click "Execute"
   - View the response

## Troubleshooting

### Swagger docs not generating

**Problem**: `swag command not found`

**Solution**:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Swagger UI shows 404

**Problem**: Docs package not imported in router

**Ensure these are in your router file:**
```go
import (
    _ "github.com/example/microservices/auth-service/docs"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

// Register route
engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

### Port conflicts

If using custom ports, update the host in main.go:
```go
// @host localhost:CUSTOM_PORT
```

Then regenerate docs:
```bash
make swagger
```

## Documentation Format

The project uses Swagger 2.0 (OpenAPI 2.0) annotations in Go code comments. Each endpoint includes:

- **@Summary**: Brief endpoint description
- **@Description**: Detailed endpoint explanation
- **@Tags**: Grouping category (Auth, Users, Health)
- **@Accept**: Accepted content types (json, form, etc.)
- **@Produce**: Response content types (usually json)
- **@Param**: Request parameters (query, path, body)
- **@Success**: Successful response status and model
- **@Failure**: Error response statuses and models
- **@Security**: Security requirements (BearerAuth)
- **@Router**: HTTP method and path

Example:
```go
// GetUser retrieves a user by ID
// @Summary Get user profile
// @Description Retrieve a user's profile information by user ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.UserSuccessResponse "User profile retrieved"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) { ... }
```

## Customizing Documentation

### Update API Information

Edit the top-level comments in `cmd/main.go`:

```go
// @title Service Name
// @version 1.0
// @description Service description
// @host localhost:PORT
// @BasePath /
// @schemes http https
```

### Add New Endpoints

1. Add handler function with Swagger comments
2. Run `make swagger` to regenerate
3. Swagger UI automatically reflects changes

### Add Custom Models

Create Go structs with `json` tags:

```go
type CustomModel struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
}
```

Reference in Swagger comments:
```go
// @Success 200 {object} model.CustomModel
```

## CI/CD Integration

Add Swagger generation to your CI/CD pipeline:

```bash
# In your build script
go install github.com/swaggo/swag/cmd/swag@latest
make swagger

# Check that docs were generated
if [ ! -d "services/auth-service/docs" ] || [ ! -d "services/user-service/docs" ]; then
    echo "Swagger generation failed"
    exit 1
fi
```

Commit generated docs to version control:

```bash
git add services/*/docs/
git commit -m "Update API documentation"
```

## Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI 2.0 Specification](https://swagger.io/specification/v2/)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

## Next Steps

1. Generate docs: `make swagger`
2. Start services: `make run-auth` and `make run-user`
3. Open Swagger UI in browser
4. Test endpoints interactively
5. Share documentation URL with API consumers
