package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSecureHeaders(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add the security middleware
	router.Use(Secure())
	
	// Add a test route
	router.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})
	
	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check the response headers
	headers := w.Result().Header
	
	// List of security headers we expect to be set
	expectedHeaders := []string{
		"X-Content-Type-Options",
		"X-XSS-Protection",
		"X-Frame-Options",
		"Content-Security-Policy",
		"Referrer-Policy",
		"Permissions-Policy",
		"Strict-Transport-Security",
		"Feature-Policy",
	}
	
	// Check each expected header exists
	for _, header := range expectedHeaders {
		assert.NotEmpty(t, headers.Get(header), "Header %s should be set", header)
	}
}
