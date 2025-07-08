package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow all origins
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Allow specific headers
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Set cache control for preflight requests (24 hours)
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		
		// Continue processing the request
		c.Next()
	}
}
