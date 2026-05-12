package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type PaginatedResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data"`
	Meta    PaginationMeta  `json:"meta"`
	Error   string          `json:"error,omitempty"`
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessWithPagination(c *gin.Context, message string, data interface{}, meta PaginationMeta) {
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *gin.Context, statusCode int, message string, err string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func BadRequest(c *gin.Context, message string, err string) {
	Error(c, http.StatusBadRequest, message, err)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, "Resource not found")
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, "Unauthorized")
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message, "Internal server error")
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, "Forbidden")
}
