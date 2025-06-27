package handler

import (
	"net/http"

	"base-code-go-gin-clean/internal/domain"
	"base-code-go-gin-clean/internal/service"
	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService service.EmailService
}

func NewEmailHandler(emailService service.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// SendEmail godoc
// @Summary Send an email
// @Description Send an email using the configured SMTP server
// @Tags email
// @Accept  json
// @Produce  json
// @Param   email  body      domain.Email  true  "Email details"
// @Success 200 {object} map[string]string "message": "Email sent successfully"
// @Failure 400 {object} map[string]string "error": "Bad request"
// @Failure 500 {object} map[string]string "error": "Internal server error"
// @Router /api/v1/email/send [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	var email domain.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if len(email.To) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one recipient is required"})
		return
	}

	if err := h.emailService.SendEmail(&email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
