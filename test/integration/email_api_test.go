package integration

import (
	"net/http"
	"testing"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler/email"
	emailService "base-code-go-gin-clean/internal/service/email"
	"base-code-go-gin-clean/test"
)

func TestEmailAPI(t *testing.T) {
	// Initialize services
	cfg := config.Config{}
	emailSvc := emailService.NewEmailService(&cfg)

	// Initialize handlers
	emailHandler := email.NewEmailHandler(emailSvc)

	// Setup router
	router := test.SetupTestRouter()
	router.POST("/api/email", emailHandler.SendEmail)

	// Test valid request
	t.Run("successful email send", func(t *testing.T) {
		emailData := map[string]interface{}{
			"to":      []string{"test@example.com"},
			"subject": "Test",
			"body":    "Test body",
		}
		body := test.MakeJSONBody(t, emailData)
		test.MakeTestRequestWithBody(router, "POST", "/api/email", body)
	})

	// Test invalid request
	t.Run("missing required parameters", func(t *testing.T) {
		resp := test.MakeTestRequest(router, "POST", "/api/email")
		test.AssertJSONResponse(t, resp, http.StatusBadRequest)
	})
}
