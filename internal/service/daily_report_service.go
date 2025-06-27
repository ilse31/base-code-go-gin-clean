package service

import (
	"fmt"

	"base-code-go-gin-clean/internal/domain"
	"base-code-go-gin-clean/internal/email"
)

type DailyReportService struct {
	emailService EmailService
}

func NewDailyReportService(emailService EmailService) *DailyReportService {
	return &DailyReportService{
		emailService: emailService,
	}
}

func (s *DailyReportService) GenerateAndSendDailyReport() {
	// Sample data - replace with actual data from your application
	reportData := map[string]interface{}{
		"newUsers":       "5",
		"activeSessions": "27",
		"totalUsers":     "42",
	}

	subject, body, err := email.DailyReportEmail("Admin", reportData)
	if err != nil {
		fmt.Printf("Failed to generate email template: %v\n", err)
		return
	}

	email := &domain.Email{
		To:      []string{"admin@example.com"}, // Replace with actual admin email
		Subject: subject,
		Body:    body,
	}

	
	if err := s.emailService.SendEmail(email); err != nil {
		// Log the error, but don't fail the entire application
		// You might want to use a proper logger here
		fmt.Printf("Failed to send daily report: %v\n", err)
	}
}
