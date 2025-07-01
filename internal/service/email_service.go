package service

import (
	"fmt"
	"net/smtp"

	"base-code-go-gin-clean/internal/config"
	domain "base-code-go-gin-clean/internal/domain/email"
)

type emailService struct {
	smtpServer   string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	from         string
}

func NewEmailService(cfg *config.Config) domain.EmailService {
	return &emailService{
		smtpServer:   cfg.Email.SMTPServer,
		smtpPort:     cfg.Email.SMTPPort,
		smtpUsername: cfg.Email.SMTPUsername,
		smtpPassword: cfg.Email.SMTPPassword,
		from:         cfg.Email.From,
	}
}

func (s *emailService) SendEmail(email *domain.Email) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpServer)

	headers := fmt.Sprintf("From: %s\r\n", s.from)
	headers += fmt.Sprintf("To: %s\r\n", email.To[0])
	headers += fmt.Sprintf("Subject: %s\r\n", email.Subject)
	headers += "MIME-version: 1.0;\r\n"
	headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	headers += "\r\n"

	msg := []byte(headers + email.Body)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", s.smtpServer, s.smtpPort),
		auth,
		s.from,
		email.To,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
