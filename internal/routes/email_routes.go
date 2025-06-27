package routes

import (
	"base-code-go-gin-clean/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupEmailRoutes configures all email-related routes
func SetupEmailRoutes(router *gin.RouterGroup, emailHandler *handler.EmailHandler) {
	emailGroup := router.Group("/email")
	{
		emailGroup.POST("/send", emailHandler.SendEmail)
	}
}
