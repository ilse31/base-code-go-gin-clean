package telemetry

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer initializes and returns a new TracerProvider configured with Jaeger exporter
func InitTracer(serviceName, jaegerURL string) (func(), error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Create a resource with the service name
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.TelemetrySDKLanguageGo,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create a new tracer provider with the exporter and resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set the global TracerProvider
	otel.SetTracerProvider(tp)

	// Set the propagator to use for distributed tracing
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return a function to flush and shutdown the tracer
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Force flush any remaining spans
		if err := tp.ForceFlush(ctx); err != nil {
			log.Printf("failed to flush tracer: %v", err)
		}

		// Shutdown the tracer provider
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("failed to shutdown tracer provider: %v", err)
			otel.Handle(err)
		}
	}, nil
}

// SpanFromContext returns the span from the context if it exists
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// Start starts a new span with the name derived from the caller's function name.
// It automatically extracts the package and function name to create a meaningful span name.
func Start(ctx context.Context) (context.Context, trace.Span) {
	// Get the caller's function name
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		// If we can't get the caller info, use a generic name
		return otel.Tracer("gin-server").Start(ctx, "unknown_operation")
	}

	// Get the function name
	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")

	// Use the last part of the function name as the operation name
	operation := parts[len(parts)-1]

	// Clean up the operation name (remove package names, etc.)
	operation = strings.TrimSuffix(operation, "-fm") // Remove method receiver suffix if present
	operation = strings.TrimPrefix(operation, "(*")
	operation = strings.TrimSuffix(operation, ")")

	return otel.Tracer("gin-server").Start(ctx, operation)
}
