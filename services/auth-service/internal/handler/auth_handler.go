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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
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
