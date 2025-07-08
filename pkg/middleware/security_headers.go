package middleware

import "github.com/gin-gonic/gin"

// Secure adds various security headers to the response
func Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Enable XSS filtering
		c.Header("X-XSS-Protection", "1; mode=block")

		// Prevent rendering of page within a frame/iframe
		c.Header("X-Frame-Options", "DENY")

		// Configure Content Security Policy
		// Note: This is a basic CSP. You might need to adjust this based on your application's requirements
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; " +
			"style-src 'self' 'unsafe-inline' https:; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' https: data:; " +
			"connect-src 'self' https:; " +
			"frame-ancestors 'none'; " +
			"form-action 'self'; " +
			"base-uri 'self'; "

		c.Header("Content-Security-Policy", csp)

		// Set Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Set Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Set Strict Transport Security (HSTS)
		// 31536000 seconds = 1 year
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Set Feature Policy (legacy, replaced by Permissions-Policy but still supported by some browsers)
		c.Header("Feature-Policy", "geolocation 'none'; microphone 'none'; camera 'none'")

		// Continue to the next handler
		c.Next()
	}
}
