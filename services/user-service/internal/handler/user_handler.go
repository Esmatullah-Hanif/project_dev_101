package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/user-service/internal/model"
	"github.com/example/microservices/user-service/internal/service"
	appErrors "github.com/example/microservices/shared/pkg/errors"
	"github.com/example/microservices/shared/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser retrieves a user by ID
// @Summary Get user profile
// @Description Retrieve a user's profile information by user ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.UserSuccessResponse "User profile retrieved"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

// ListUsers retrieves a paginated list of all users
// @Summary List all users
// @Description Get a paginated list of all users in the system
// @Tags Users
// @Produce json
// @Param page query int false "Page number (default: 1)" default(1)
// @Param page_size query int false "Page size (default: 10)" default(10)
// @Success 200 {object} model.UserListSuccessResponse "Users retrieved successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid pagination parameters"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "Invalid page parameter", err.Error())
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.BadRequest(c, "Invalid page_size parameter", err.Error())
		return
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		h.handleError(c, err)
		return
	}

	totalPages := (total + pageSize - 1) / pageSize
	meta := response.PaginationMeta{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: totalPages,
	}

	response.SuccessWithPagination(c, "Users retrieved successfully", users, meta)
}

// UpdateUser updates a user's profile
// @Summary Update user profile
// @Description Update user profile information (requires authentication)
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body model.UpdateUserRequest true "Update request"
// @Success 200 {object} model.UserSuccessResponse "User updated successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid request payload"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser deletes a user account
// @Summary Delete user
// @Description Delete a user account (soft delete, requires authentication)
// @Tags Users
// @Param id path string true "User ID"
// @Success 204 "User deleted successfully"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 404 {object} model.ErrorResponse "User not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, "User deleted successfully", nil)
}

// HealthCheck returns the health status of the user service
// @Summary Health check
// @Description Check if user service is running and healthy
// @Tags Health
// @Produce json
// @Success 200 {object} model.SuccessResponse "Service is healthy"
// @Router /health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "User service is healthy", gin.H{
		"service": "user-service",
		"status":  "healthy",
	})
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		switch appErr.Code {
		case appErrors.ValidationError:
			response.BadRequest(c, appErr.Message, appErr.Error())
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
