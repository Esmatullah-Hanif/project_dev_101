package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/auth-service/internal/model"
	"github.com/example/microservices/auth-service/internal/service"
	appErrors "github.com/example/microservices/shared/pkg/errors"
	"github.com/example/microservices/shared/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SignUp creates a new user account
// @Summary Create new user account
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.SignUpRequest true "Sign up request"
// @Success 201 {object} model.AuthSuccessResponse "User created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload"
// @Failure 409 {object} model.ErrorResponse "Email already registered"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /api/v1/auth/signup [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req model.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	authResp, err := h.authService.SignUp(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User created successfully", authResp)
}

// SignIn authenticates a user and returns tokens
// @Summary User login
// @Description Authenticate with email and password to receive access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.SignInRequest true "Sign in request"
// @Success 200 {object} model.AuthSuccessResponse "Login successful"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload"
// @Failure 401 {object} model.ErrorResponse "Invalid email or password"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /api/v1/auth/signin [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var req model.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	authResp, err := h.authService.SignIn(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login successful", authResp)
}

// RefreshToken generates new access token from refresh token
// @Summary Refresh access token
// @Description Use refresh token to get a new access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body model.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} model.AuthSuccessResponse "Token refreshed successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload"
// @Failure 401 {object} model.ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	authResp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", authResp)
}

// HealthCheck returns the health status of the auth service
// @Summary Health check
// @Description Check if auth service is running and healthy
// @Tags Health
// @Produce json
// @Success 200 {object} model.SuccessResponse "Service is healthy"
// @Router /health [get]
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "Auth service is healthy", gin.H{
		"service": "auth-service",
		"status":  "healthy",
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		switch appErr.Code {
		case appErrors.ValidationError:
			response.BadRequest(c, appErr.Message, appErr.Error())
		case appErrors.AuthError:
			response.Unauthorized(c, appErr.Message)
		case appErrors.ConflictError:
			response.Error(c, http.StatusConflict, appErr.Message, appErr.Error())
		case appErrors.NotFoundError:
			response.NotFound(c, appErr.Message)
		case appErrors.ForbiddenError:
			response.Forbidden(c, appErr.Message)
		default:
			response.InternalError(c, appErr.Message)
		}
	} else {
		response.InternalError(c, "An unexpected error occurred")
	}
}
