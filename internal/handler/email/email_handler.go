package email

import (
	"errors"

	domain "base-code-go-gin-clean/internal/domain/email"
	httpPkg "base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService domain.EmailService
}

func NewEmailHandler(emailService domain.EmailService) *EmailHandler {
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
// @Success 200 {object} map[string]string "Email sent successfully"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /email/send [post]
func (h *EmailHandler) SendEmail(c *gin.Context) {
	// Start a new span for the request
	_, span := telemetry.Start(c.Request.Context())
	defer span.End()

	var email domain.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		httpPkg.BadRequest(c, "Invalid request payload", nil)
		span.RecordError(err)
		return
	}

	if len(email.To) == 0 {
		err := errors.New("no recipients provided")
		span.RecordError(err)
		httpPkg.BadRequest(c, "At least one recipient is required", nil)
		return
	}

	if err := h.emailService.SendEmail(&email); err != nil {
		span.RecordError(err)
		httpPkg.InternalServerError(c, "Failed to send email")
		return
	}

	httpPkg.Success(c, map[string]string{
		"message": "Email sent successfully",
	})
}
