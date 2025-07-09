package health

import (
	"context"
	"time"

	"base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
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
	ctx, span := telemetry.StartSpan(c.Request.Context(), "HealthCheck")
	defer telemetry.EndSpan(span, nil)

	// Add custom attributes to the span
	span.SetAttributes(
		attribute.String("health.check.type", "liveness"),
		attribute.String("service.version", "1.0.0"),
	)

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
	ctx, span := telemetry.StartSpan(ctx, "performHealthChecks")
	defer telemetry.EndSpan(span, nil)

	// Simulate checking database connection
	{
		_, dbSpan := telemetry.StartSpan(ctx, "checkDatabase")
		time.Sleep(50 * time.Millisecond) // Simulate database check
		telemetry.EndSpan(dbSpan, nil)
	}

	// Simulate checking cache
	{
		_, cacheSpan := telemetry.StartSpan(ctx, "checkCache")
		time.Sleep(30 * time.Millisecond) // Simulate cache check
		telemetry.EndSpan(cacheSpan, nil)
	}

	// Simulate checking external service
	{
		_, externalSpan := telemetry.StartSpan(ctx, "checkExternalService")
		time.Sleep(70 * time.Millisecond) // Simulate external service check
		telemetry.EndSpan(externalSpan, nil)
	}
}
