package telemetry

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerKey  = "otel-tracer"
	tracerName = "github.com/your-org/your-app/telemetry"
)

// Middleware returns a gin middleware that starts a span for each request
func Middleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip tracing for health check endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		// Extract tracing context from the request headers
		ctx := otel.GetTextMapPropagator().Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		// Create a new span
		tracer := otel.Tracer(tracerName)
		spanName := c.FullPath()
		if spanName == "" {
			spanName = c.Request.URL.Path
		}

		// Create attributes for the span
		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", spanName),
			attribute.String("http.url", c.Request.URL.String()),
			attribute.String("http.user_agent", c.Request.UserAgent()),
			attribute.String("http.host", c.Request.Host),
			attribute.String("http.client_ip", c.ClientIP()),
		}

		// Create a new context with the span
		_, span := tracer.Start(
			ctx,
			spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Store the span in the context
		c.Set(tracerKey, span)

		// Create a wrapper for the response writer to capture the status code
		ww := &responseWriter{
			ResponseWriter: c.Writer,
			startTime:      time.Now(),
		}
		c.Writer = ww

		// Process the request
		c.Next()

		// Add response status code to the span
		status := c.Writer.Status()
		span.SetStatus(httpStatusToOtelCode(status), http.StatusText(status))
		span.SetAttributes(attribute.Int("http.status_code", status))

		// Record metrics
		duration := time.Since(ww.startTime)
		span.SetAttributes(semconv.HTTPResponseContentLength(int(ww.size)))

		// Record any errors
		if len(c.Errors) > 0 {
			span.RecordError(c.Errors.Last())
		}

		// Record the duration in milliseconds
		span.SetAttributes(attribute.Int64("http.duration_ms", duration.Milliseconds()))
	}
}

// httpStatusToOtelCode converts an HTTP status code to an OpenTelemetry status code
func httpStatusToOtelCode(code int) codes.Code {
	switch {
	case code >= 200 && code < 400:
		return codes.Ok
	case code == http.StatusUnauthorized:
		return codes.Error // Use Error for all error cases to simplify
	case code == http.StatusForbidden:
		return codes.Error
	case code == http.StatusNotFound:
		return codes.Error
	case code == http.StatusConflict:
		return codes.Error
	case code == http.StatusPreconditionFailed:
		return codes.Error
	default:
		return codes.Error
	}
}

// responseWriter is a wrapper around gin.ResponseWriter that captures the status code
// and other information for tracing purposes.
type responseWriter struct {
	gin.ResponseWriter
	startTime time.Time
	size      int
	status    int
}

// WriteHeader captures the status code and writes it to the response
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written
func (w *responseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// Status returns the captured status code
func (w *responseWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}
