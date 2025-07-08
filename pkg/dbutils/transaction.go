package dbutils

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

// DB is an interface that both *bun.DB and *bun.Tx implement
type DB interface {
	BeginTx(context.Context, *sql.TxOptions) (bun.Tx, error)
	RunInTx(context.Context, *sql.TxOptions, func(context.Context, bun.Tx) error) error
}

// TransactionMiddleware creates a middleware that starts a database transaction
// and stores it in the context. The transaction will be committed if the request
// completes successfully, or rolled back if there's an error.
func TransactionMiddleware(db *bun.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip transaction for read-only methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
			c.Next()
			return
		}

		// Start a new transaction
		tx, err := db.BeginTx(c.Request.Context(), &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		// Store the transaction in the context
		c.Set("db", tx)

		// Process the request
		c.Next()

		// Check if there was an error processing the request
		if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
			// Rollback the transaction if there was an error
			tx.Rollback()
		} else {
			// Commit the transaction if there were no errors
			if err := tx.Commit(); err != nil {
				c.AbortWithError(500, err)
			}
		}

		// Clear the transaction from the context
		c.Set("db", nil)
	}
}

// GetDB retrieves the database connection or transaction from the context
func GetDB(c *gin.Context) DB {
	if tx, exists := c.Get("db"); exists {
		if db, ok := tx.(bun.Tx); ok {
			return db
		}
	}
	// Return the original DB if no transaction is found
	db, _ := c.MustGet("db_conn").(*bun.DB)
	return db
}

// WithTransaction runs the provided function within a transaction
func WithTransaction(ctx context.Context, db DB, fn func(context.Context, bun.Tx) error) error {
	return db.RunInTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}, fn)
}
