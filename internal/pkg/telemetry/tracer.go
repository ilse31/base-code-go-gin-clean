package telemetry

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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

// NewSpan starts a new span with the given name and context
func NewSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("gin-server").Start(ctx, name)
}

// SpanFromContext returns the span from the context if it exists
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// StartSpan starts a new span with the given name and returns a context with the span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("gin-server").Start(ctx, name)
}

// EndSpan ends the span with the given error (if any)
func EndSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End()
}
