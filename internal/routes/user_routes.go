package routes

import (
	"base-code-go-gin-clean/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all the user routes
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.GET("/:id", userHandler.GetUserByID)
	}
}
