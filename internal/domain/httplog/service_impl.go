package httplog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

// NewService creates a new instance of the HTTP log service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// LogOutgoingRequest logs an outgoing HTTP request
func (s *service) LogOutgoingRequest(ctx context.Context, log *LogOutgoingRequest) error {
	if log == nil {
		return errors.New("log cannot be nil")
	}

	// Set timestamps if not set
	now := time.Now()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = now
	}

	return s.repo.LogOutgoingRequest(ctx, log)
}

// LogIncomingRequest logs an incoming HTTP request and returns the log ID
func (s *service) LogIncomingRequest(ctx context.Context, log *LogIncomingRequest) (string, error) {
	if log == nil {
		return "", errors.New("log cannot be nil")
	}

	// Set timestamps if not set
	now := time.Now()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = now
	}

	return s.repo.LogIncomingRequest(ctx, log)
}

// LogError logs an error that occurred during request processing
func (s *service) LogError(ctx context.Context, log *LogError) error {
	if log == nil {
		return errors.New("log cannot be nil")
	}

	// Set timestamps if not set
	now := time.Now()
	if log.CreatedAt.IsZero() {
		log.CreatedAt = now
	}
	if log.UpdatedAt.IsZero() {
		log.UpdatedAt = now
	}

	return s.repo.LogError(ctx, log)
}

// LogErrorWithRequest logs an error with a reference to the original request
func (s *service) LogErrorWithRequest(ctx context.Context, requestID string, statusCode int, err error, traces interface{}) error {
	if err == nil {
		return errors.New("error cannot be nil")
	}

	// Try to get the trace ID from context or generate a new one
	traceID, _ := ctx.Value("trace_id").(string)
	if traceID == "" {
		traceID = generateTraceID()
	}

	log := &LogError{
		TraceID:    traceID,
		RequestID:  &requestID,
		StatusCode: statusCode,
		Error:      err.Error(),
		Traces:     traces,
	}

	return s.LogError(ctx, log)
}

// GetRequestLogs retrieves all logs for a specific trace ID
func (s *service) GetRequestLogs(ctx context.Context, traceID string) (*RequestLogs, error) {
	if traceID == "" {
		return nil, errors.New("trace ID cannot be empty")
	}

	// Get incoming request
	incomingReqs, err := s.repo.FindIncomingRequestByTraceID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get incoming requests: %w", err)
	}

	// If no incoming requests found, return empty result
	if len(incomingReqs) == 0 {
		return &RequestLogs{}, nil
	}

	// For simplicity, we'll use the most recent incoming request
	incomingReq := incomingReqs[0]

	// Get outgoing requests
	outgoingReqs, err := s.repo.FindOutgoingRequestByTraceID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get outgoing requests: %w", err)
	}

	// Get errors
	errors, err := s.repo.FindErrorsByTraceID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get errors: %w", err)
	}

	return &RequestLogs{
		IncomingRequest:  incomingReq,
		OutgoingRequests: outgoingReqs,
		Errors:           errors,
	}, nil
}

// GetErrorLogs retrieves error logs for a specific trace ID
func (s *service) GetErrorLogs(ctx context.Context, traceID string) ([]*LogError, error) {
	if traceID == "" {
		return nil, errors.New("trace ID cannot be empty")
	}

	return s.repo.FindErrorsByTraceID(ctx, traceID)
}

// CleanupOldLogs removes logs older than the specified number of days
func (s *service) CleanupOldLogs(ctx context.Context, olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		return 0, errors.New("olderThanDays must be greater than 0")
	}

	return s.repo.CleanupOldLogs(ctx, olderThanDays)
}

// generateTraceID generates a new trace ID
func generateTraceID() string {
	return "trace_" + time.Now().Format("20060102150405") + "_" + uuid.New().String()
}
