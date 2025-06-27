package routes

import (
	"base-code-go-gin-clean/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all the user routes
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	// User routes under /api/v1/users
	router.GET("/users/:id", userHandler.GetUserByID)
}
