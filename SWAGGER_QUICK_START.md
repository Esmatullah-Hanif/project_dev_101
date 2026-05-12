# Swagger Quick Start Guide

Get Swagger documentation up and running in 5 minutes.

## 1. Install swag CLI (2 minutes)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Verify:
```bash
swag --version
```

## 2. Generate Documentation (1 minute)

```bash
make swagger
```

This generates:
- `services/auth-service/docs/`
- `services/user-service/docs/`

## 3. Start Services (1 minute)

In separate terminals:

```bash
# Terminal 1
make run-auth

# Terminal 2  
make run-user

# Terminal 3 (optional)
make run-gateway
```

## 4. Open in Browser (30 seconds)

- **Auth Service**: http://localhost:8001/swagger/index.html
- **User Service**: http://localhost:8002/swagger/index.html

## 5. Test an Endpoint (30 seconds)

1. Open either Swagger UI link
2. Click "Try it out" on any endpoint
3. Fill in parameters
4. Click "Execute"
5. View response

## Testing Authentication

To test protected endpoints:

1. **Sign up** via `/api/v1/auth/signup` endpoint
   - Enter email and password
   - Copy the `access_token` from response

2. **Authorize** in Swagger UI
   - Click "Authorize" button (top right)
   - Paste: `Bearer <your_access_token>`
   - Click "Authorize"

3. **Test protected endpoints**
   - Now you can test PUT/DELETE on user endpoints
   - Token is automatically included in requests

## Common Tasks

### Regenerate Documentation
```bash
make swagger
```

### View Auth Service Docs Only
```bash
make swagger-auth
# Then open: http://localhost:8001/swagger/index.html
```

### View User Service Docs Only
```bash
make swagger-user
# Then open: http://localhost:8002/swagger/index.html
```

### Update Custom Port

If you changed ports in `.env`:

1. Edit the service's `cmd/main.go`:
   ```go
   // @host localhost:YOUR_PORT
   ```

2. Regenerate docs:
   ```bash
   make swagger
   ```

## Key Endpoints to Test

### Auth Service
- `POST /api/v1/auth/signup` - Create account
- `POST /api/v1/auth/signin` - Login
- `POST /api/v1/auth/refresh` - Get new token

### User Service
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get user profile
- `PUT /api/v1/users/{id}` - Update profile (requires auth)
- `DELETE /api/v1/users/{id}` - Delete user (requires auth)

## Troubleshooting

### swag command not found
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Port already in use
Change port in `.env` and regenerate docs

### Swagger UI shows 404
- Ensure services are running
- Check correct port (8001 for auth, 8002 for user)
- Try: `http://localhost:8001/swagger/index.html`

## Tips

- **Try it out**: Click "Try it out" button to test live
- **Schema**: Click schema name to see structure
- **Response examples**: Check "Response 200" section for examples
- **Authentication**: Use "Authorize" button for bearer tokens
- **History**: Test endpoint data persists during session

## Next Steps

1. Explore endpoints in Swagger UI
2. Review generated `swagger.json` files
3. Share Swagger UI URLs with team
4. Commit docs to version control
5. See `docs/SWAGGER.md` for advanced usage

## Get Help

- Swagger UI built-in help: Click `?` icon
- Detailed docs: `docs/SWAGGER.md`
- Implementation: `SWAGGER_IMPLEMENTATION.md`
- API examples: `examples/auth_requests.sh` and `examples/user_requests.sh`
