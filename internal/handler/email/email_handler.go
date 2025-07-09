package email

import (
	"context"

	domain "base-code-go-gin-clean/internal/domain/email"
	httpPkg "base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
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
	ctx := c.Request.Context()

	err := telemetry.WithSpan(ctx, "SendEmail", func(ctx context.Context) error {
		var email domain.Email
		if err := c.ShouldBindJSON(&email); err != nil {
			httpPkg.BadRequest(c, "Invalid request payload", nil)
			return err
		}

		// Add email metadata to span (without sensitive content)
		telemetry.AddSpanAttributes(ctx,
			attribute.Int("email.recipients", len(email.To)),
			attribute.String("email.subject", email.Subject),
		)

		if len(email.To) == 0 {
			httpPkg.BadRequest(c, "At least one recipient is required", nil)
			telemetry.AddSpanAttributes(ctx, attribute.String("error.type", "no_recipients"))
			return nil
		}

		if err := h.emailService.SendEmail(&email); err != nil {
			httpPkg.InternalServerError(c, "Failed to send email")
			telemetry.AddSpanAttributes(ctx,
				attribute.String("error.type", "send_failed"),
				attribute.String("error.details", err.Error()),
			)
			return err
		}

		httpPkg.Success(c, map[string]string{
			"message": "Email sent successfully",
		})
		return nil
	})

	// Error is already recorded in the span by WithSpan
	if err != nil {
		return
	}
}
