package server

import (
	"base-code-go-gin-clean/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// setupMiddleware configures global middleware for the server
func (s *Server) setupMiddleware() {
	// Request ID middleware for tracing
	s.router.Use(middleware.RequestID())

	// Logger middleware
	s.router.Use(gin.Logger())

	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Add CORS middleware if needed
	// s.router.Use(middleware.CORS())

	// Add security headers
	// s.router.Use(middleware.Secure())

	// Rate limiting
	// s.router.Use(middleware.RateLimit())
}
