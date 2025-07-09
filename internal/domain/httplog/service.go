package httplog

import (
	"context"
)

// Service defines the interface for HTTP log business logic
type Service interface {
	// LogOutgoingRequest logs an outgoing HTTP request
	LogOutgoingRequest(ctx context.Context, log *LogOutgoingRequest) error

	// LogIncomingRequest logs an incoming HTTP request and returns the log ID
	LogIncomingRequest(ctx context.Context, log *LogIncomingRequest) (string, error)

	// LogError logs an error that occurred during request processing
	LogError(ctx context.Context, log *LogError) error

	// LogErrorWithRequest logs an error with a reference to the original request
	LogErrorWithRequest(ctx context.Context, requestID string, statusCode int, err error, traces interface{}) error

	// GetRequestLogs retrieves logs for a specific trace ID
	GetRequestLogs(ctx context.Context, traceID string) (*RequestLogs, error)

	// GetErrorLogs retrieves error logs for a specific trace ID
	GetErrorLogs(ctx context.Context, traceID string) ([]*LogError, error)

	// CleanupOldLogs removes logs older than the specified number of days
	CleanupOldLogs(ctx context.Context, olderThanDays int) (int64, error)
}

// RequestLogs contains all logs related to a specific request
// RequestLogs contains all logs related to a specific request
type RequestLogs struct {
	IncomingRequest  *LogIncomingRequest   `json:"incoming_request"`
	OutgoingRequests []*LogOutgoingRequest `json:"outgoing_requests,omitempty"`
	Errors           []*LogError           `json:"errors,omitempty"`
}
