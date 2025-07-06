package routes

import (
	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler/auth"
	"base-code-go-gin-clean/internal/middleware"
	"base-code-go-gin-clean/internal/pkg/token"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures all the authentication routes
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *auth.AuthHandler, tokenConfig *config.TokenConfig) {
	// Initialize token service with configuration
	tokenService := token.NewTokenService(tokenConfig)

	// Initialize auth middleware with token service
	authMiddleware := middleware.AuthMiddleware(tokenService)

	// Public routes (no authentication required)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		// Protected routes (require valid access token)
		protected := authGroup.Group("")
		protected.Use(authMiddleware)
		{
			// Refresh token endpoint
			protected.POST("/refresh", authHandler.RefreshToken)
			
			// Logout endpoint
			protected.POST("/logout", authHandler.Logout)

			// Example of a protected route with role-based access
			// adminGroup := protected.Group("/admin")
			// adminGroup.Use(middleware.RoleMiddleware("admin"))
			// {
			// 	adminGroup.GET("/dashboard", adminHandler.Dashboard)
			// }
		}
	}
}
