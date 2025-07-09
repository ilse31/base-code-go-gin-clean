package httplog

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

// NewRepository creates a new instance of the HTTP log repository
func NewRepository(db *bun.DB) Repository {
	return &repository{db: db}
}

// LogOutgoingRequest logs an outgoing HTTP request
func (r *repository) LogOutgoingRequest(ctx context.Context, log *LogOutgoingRequest) error {
	_, err := r.db.NewInsert().
		Model(log).
		Exec(ctx)
	return err
}

// LogIncomingRequest logs an incoming HTTP request and returns the log ID
func (r *repository) LogIncomingRequest(ctx context.Context, log *LogIncomingRequest) (string, error) {
	_, err := r.db.NewInsert().
		Model(log).
		Returning("id").
		Exec(ctx, &log.ID)

	if err != nil {
		return "", err
	}

	return log.ID, nil
}

// LogError logs an error that occurred during request processing
func (r *repository) LogError(ctx context.Context, log *LogError) error {
	_, err := r.db.NewInsert().
		Model(log).
		Exec(ctx)
	return err
}

// FindOutgoingRequestByTraceID finds outgoing requests by trace ID
func (r *repository) FindOutgoingRequestByTraceID(ctx context.Context, traceID string) ([]*LogOutgoingRequest, error) {
	var logs []*LogOutgoingRequest
	err := r.db.NewSelect().
		Model(&logs).
		Where("trace_id = ?", traceID).
		Order("created_at DESC").
		Scan(ctx)

	return logs, err
}

// FindIncomingRequestByTraceID finds incoming requests by trace ID
func (r *repository) FindIncomingRequestByTraceID(ctx context.Context, traceID string) ([]*LogIncomingRequest, error) {
	var logs []*LogIncomingRequest
	err := r.db.NewSelect().
		Model(&logs).
		Where("trace_id = ?", traceID).
		Order("created_at DESC").
		Scan(ctx)

	return logs, err
}

// FindErrorsByTraceID finds error logs by trace ID
func (r *repository) FindErrorsByTraceID(ctx context.Context, traceID string) ([]*LogError, error) {
	var logs []*LogError
	err := r.db.NewSelect().
		Model(&logs).
		Where("trace_id = ?", traceID).
		Order("created_at DESC").
		Scan(ctx)

	return logs, err
}

// FindRequestWithErrors finds a request and its associated errors
func (r *repository) FindRequestWithErrors(ctx context.Context, requestID string) (*LogIncomingRequest, []*LogError, error) {
	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// Get the request
	var req LogIncomingRequest
	err = tx.NewSelect().
		Model(&req).
		Where("id = ?", requestID).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	// Get associated errors
	var errors []*LogError
	err = tx.NewSelect().
		Model(&errors).
		Where("request_id = ?", requestID).
		Order("created_at DESC").
		Scan(ctx)

	if err != nil {
		return nil, nil, err
	}

	return &req, errors, tx.Commit()
}

// CleanupOldLogs removes logs older than the specified number of days
func (r *repository) CleanupOldLogs(ctx context.Context, olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		return 0, errors.New("olderThanDays must be greater than 0")
	}

	cutoffTime := time.Now().AddDate(0, 0, -olderThanDays)

	// Delete old outgoing requests
	outgoingResult, err := r.db.NewDelete().
		Model((*LogOutgoingRequest)(nil)).
		Where("created_at < ?", cutoffTime).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	// Delete old incoming requests
	incomingResult, err := r.db.NewDelete().
		Model((*LogIncomingRequest)(nil)).
		Where("created_at < ?", cutoffTime).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	// Delete old errors
	errorResult, err := r.db.NewDelete().
		Model((*LogError)(nil)).
		Where("created_at < ?", cutoffTime).
		Exec(ctx)
	if err != nil {
		return 0, err
	}

	outgoingCount, _ := outgoingResult.RowsAffected()
	incomingCount, _ := incomingResult.RowsAffected()
	errorCount, _ := errorResult.RowsAffected()

	return outgoingCount + incomingCount + errorCount, nil
}
