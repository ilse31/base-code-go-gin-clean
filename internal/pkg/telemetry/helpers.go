package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// AddSpanAttributes adds attributes to the current span from context
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attrs...)
	}
}

// RecordError records an error in the current span with optional attributes
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span != nil && err != nil {
		sc := span.SpanContext()
		if sc.IsValid() {
			// Only record errors for valid spans
			span.RecordError(err, trace.WithAttributes(attrs...))
			span.SetStatus(codes.Error, err.Error())
		}
	}
}

// WithSpan creates a new span as a child of the current span (if any)
func WithSpan(ctx context.Context, name string, fn func(context.Context) error, attrs ...attribute.KeyValue) error {
	ctx, span := otel.Tracer("gin-server").Start(ctx, name)
	defer span.End()

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	err := fn(ctx)
	if err != nil {
		sc := span.SpanContext()
		if sc.IsValid() {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}
