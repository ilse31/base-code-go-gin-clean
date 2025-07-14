package email

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"base-code-go-gin-clean/internal/domain/email"
	"base-code-go-gin-clean/test"
)

// MockEmailService is a mock implementation of EmailService
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(e *email.Email) error {
	args := m.Called(e)
	return args.Error(0)
}

func TestEmailHandler_SendEmail(t *testing.T) {
	// Setup
	mockService := new(MockEmailService)
	handler := NewEmailHandler(mockService)
	router := test.SetupTestRouter()
	router.POST("/email", handler.SendEmail)

	// Prepare test payload
	emailPayload := email.Email{
		To:      []string{"test@example.com"},
		Subject: "Test",
		Body:    "Test body",
	}
	jsonBody, _ := json.Marshal(emailPayload)

	// Expectation
	mockService.On("SendEmail", mock.MatchedBy(func(e *email.Email) bool {
		return e.Subject == "Test" && e.Body == "Test body" && len(e.To) == 1 && e.To[0] == "test@example.com"
	})).Return(nil)

	// Execute
	testReq := test.MakeTestRequestWithBody(router, "POST", "/email", bytes.NewReader(jsonBody))
	testReq.Request.Header.Set("Content-Type", "application/json")

	// Assert
	assert.Equal(t, http.StatusOK, testReq.Response.Code)
	mockService.AssertExpectations(t)
}
