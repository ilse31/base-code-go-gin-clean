package health

import (
	"context"
	"time"

	"base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"

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
	// Create a new span for the health check
	ctx, span := telemetry.Start(c.Request.Context())
	defer span.End()

	// Simulate some work with a child span
	h.performHealthChecks(ctx)

	http.Success(c, HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}

// performHealthChecks simulates performing health checks with its own span
func (h *HealthHandler) performHealthChecks(ctx context.Context) {
	// Create a child span for the health checks
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	// Simulate checking database connection
	{
		_, dbSpan := telemetry.Start(ctx)
		time.Sleep(50 * time.Millisecond) // Simulate database check
		dbSpan.End()
	}

	// Simulate checking cache
	{
		_, cacheSpan := telemetry.Start(ctx)
		time.Sleep(30 * time.Millisecond) // Simulate cache check
		cacheSpan.End()
	}

	// Simulate checking external service
	{
		_, externalSpan := telemetry.Start(ctx)
		time.Sleep(70 * time.Millisecond) // Simulate external service check
		externalSpan.End()
	}
}
