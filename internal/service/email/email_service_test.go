// internal/service/email/email_service_test.go
package email

import (
	"testing"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/domain/email"

	"net/smtp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementation
type MockSender struct {
	mock.Mock
}

func (m *MockSender) Send(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}
func TestEmailService_SendEmail(t *testing.T) {
	called := false

	cfg := &config.Config{
		Email: config.EmailConfig{
			SMTPServer:   "localhost",
			SMTPPort:     "1025",
			SMTPUsername: "test",
			SMTPPassword: "test",
			From:         "test@example.com",
		},
	}

	svc := &emailService{
		smtpServer:   cfg.Email.SMTPServer,
		smtpPort:     cfg.Email.SMTPPort,
		smtpUsername: cfg.Email.SMTPUsername,
		smtpPassword: cfg.Email.SMTPPassword,
		from:         cfg.Email.From,
		sendMail: func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			called = true
			assert.Equal(t, "localhost:1025", addr)
			assert.Equal(t, "test@example.com", from)
			assert.Equal(t, []string{"test@example.com"}, to)
			assert.Contains(t, string(msg), "Test Subject")
			return nil
		},
	}

	testEmail := &email.Email{
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test body",
	}

	err := svc.SendEmail(testEmail)

	assert.NoError(t, err)
	assert.True(t, called, "Expected sendMail to be called")
}
