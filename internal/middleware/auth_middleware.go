package middleware

import (
	"base-code-go-gin-clean/internal/pkg/token"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a middleware that checks for a valid access token
func AuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the access token from the cookie
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Access token is required",
			})
			return
		}

		// Validate the access token
		userID, err := tokenService.ValidateAccessToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired access token",
			})
			return
		}

		// Set the user ID in the context for use in subsequent handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// RoleMiddleware is a middleware that checks if the user has the required role
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real application, you would check the user's roles here
		// For this example, we'll just check if the user has the required role
		// You would typically get the user's roles from the database or token claims
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			return
		}

		// In a real application, you would check the user's roles here
		// For now, we'll just log the user ID and required role
		fmt.Printf("User %s is trying to access a route that requires role: %s\n", userID, requiredRole)

		// Continue to the next handler if the user has the required role
		c.Next()
	}
}
