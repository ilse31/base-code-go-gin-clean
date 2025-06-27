package health

import (
	"base-code-go-gin-clean/internal/pkg/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health-related HTTP requests
type HealthHandler struct{}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles health check requests
// @Summary Show the status of server
// @Description Returns the health status of the API along with version information
// @Tags health api
// @Accept json
// @Produce json
// @Success 200 {object} handler.SuccessResponse{data=HealthResponse} "Success response with health status"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	http.Success(c, HealthResponse{
		Status:  "ok",
		Version: "1.0.0", // Consider getting this from build flags or config
	})
}
