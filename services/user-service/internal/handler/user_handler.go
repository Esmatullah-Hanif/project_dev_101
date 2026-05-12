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

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	user, err := h.userService.GetUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", user)
}

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

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.userService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Success(c, http.StatusNoContent, "User deleted successfully", nil)
}

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
