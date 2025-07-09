package httplog

import "context"

// Repository defines the interface for HTTP log storage operations
type Repository interface {
	// LogOutgoingRequest logs an outgoing HTTP request
	LogOutgoingRequest(ctx context.Context, log *LogOutgoingRequest) error

	// LogIncomingRequest logs an incoming HTTP request and returns the log ID
	LogIncomingRequest(ctx context.Context, log *LogIncomingRequest) (string, error)

	// LogError logs an error that occurred during request processing
	LogError(ctx context.Context, log *LogError) error

	// FindOutgoingRequestByTraceID finds outgoing requests by trace ID
	FindOutgoingRequestByTraceID(ctx context.Context, traceID string) ([]*LogOutgoingRequest, error)

	// FindIncomingRequestByTraceID finds incoming requests by trace ID
	FindIncomingRequestByTraceID(ctx context.Context, traceID string) ([]*LogIncomingRequest, error)

	// FindErrorsByTraceID finds error logs by trace ID
	FindErrorsByTraceID(ctx context.Context, traceID string) ([]*LogError, error)

	// FindRequestWithErrors finds a request and its associated errors
	FindRequestWithErrors(ctx context.Context, requestID string) (*LogIncomingRequest, []*LogError, error)

	// CleanupOldLogs removes logs older than the specified duration
	CleanupOldLogs(ctx context.Context, olderThanDays int) (int64, error)
}
