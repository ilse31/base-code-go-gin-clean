package server

import (
	"base-code-go-gin-clean/pkg/dbutils"
	"base-code-go-gin-clean/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

// setupMiddleware configures global middleware for the server
func (s *Server) setupMiddleware() {
	// Request ID middleware for tracing
	s.router.Use(middleware.RequestID())

	// Logger middleware
	s.router.Use(gin.Logger())

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// CORS middleware
	s.router.Use(middleware.CORS())

	// Security headers middleware
	s.router.Use(middleware.Secure())

	// Rate limiting
	s.router.Use(middleware.RateLimitMiddleware())
}

// setupDatabaseMiddleware adds database-related middleware to the router
func (s *Server) setupDatabaseMiddleware(db *bun.DB) {
	// Add database connection to the Gin context
	s.router.Use(func(c *gin.Context) {
		c.Set("db_conn", db)
		c.Next()
	})

	// Add transaction middleware
	s.router.Use(dbutils.TransactionMiddleware(db))
}
