package health

import (
	"context"

	"base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"
	"base-code-go-gin-clean/internal/service"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health-related HTTP requests
type HealthHandler struct {
	healthService service.HealthService
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(healthService service.HealthService) *HealthHandler {
	return &HealthHandler{
		healthService: healthService,
	}
}

// HealthCheck handles health check requests
// @Summary Show the status of server
// @Description Returns the health status of the API along with version and database information
// @Tags health api
// @Accept json
// @Produce json
// @Success 200 {object} handler.SuccessResponse{data=HealthResponse} "Success response with health status"
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	// Create a new span for the health check
	ctx, span := telemetry.Start(c.Request.Context())
	defer span.End()

	dbStatus := h.performDatabaseCheck(ctx)

	http.Success(c, HealthResponse{
		Status:   "ok",
		Version:  "1.0.0",
		Database: dbStatus,
	})
}

// performDatabaseCheck performs actual database health check
func (h *HealthHandler) performDatabaseCheck(ctx context.Context) DatabaseStatus {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := h.healthService.CheckDBHealth(ctx)
	if err != nil {
		return DatabaseStatus{
			Status:  "error",
			Message: err.Error(),
		}
	}

	return DatabaseStatus{
		Status: "ok",
	}
}
