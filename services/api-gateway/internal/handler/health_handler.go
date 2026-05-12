package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/microservices/shared/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response.Success(c, http.StatusOK, "API Gateway is healthy", gin.H{
		"service": "api-gateway",
		"status":  "healthy",
	})
}
