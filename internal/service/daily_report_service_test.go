package service_test

import (
	"base-code-go-gin-clean/internal/domain/email"
	svc "base-code-go-gin-clean/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockEmailService struct {
	mock.Mock
}

func (m *mockEmailService) SendEmail(e *email.Email) error {
	args := m.Called(e)
	return args.Error(0)
}

func (m *mockEmailService) On(methodName string, arguments ...interface{}) *mock.Call {
	return m.Mock.On(methodName, arguments...)
}

func (m *mockEmailService) AssertExpectations(t mock.TestingT) bool {
	return m.Mock.AssertExpectations(t)
}

func TestDailyReportService_GenerateAndSendDailyReport(t *testing.T) {
	mockEmailSvc := new(mockEmailService)
	service := svc.NewDailyReportService(mockEmailSvc)

	t.Run("successful report generation", func(t *testing.T) {
		mockEmailSvc.On("SendEmail", mock.AnythingOfType("*email.Email")).Return(nil)

		service.GenerateAndSendDailyReport()

		mockEmailSvc.AssertExpectations(t)
	})

	t.Run("email send failure", func(t *testing.T) {
		mockEmailSvc.On("SendEmail", mock.AnythingOfType("*email.Email")).Return(assert.AnError)

		service.GenerateAndSendDailyReport()

		mockEmailSvc.AssertExpectations(t)
	})
}
