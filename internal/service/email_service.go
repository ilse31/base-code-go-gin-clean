package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"base-code-go-gin-clean/internal/domain"
)

type EmailService interface {
	SendEmail(email *domain.Email) error
}

type emailService struct {
	smtpServer   string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	from         string
}

func NewEmailService(smtpServer, smtpPort, smtpUsername, smtpPassword, from string) EmailService {
	return &emailService{
		smtpServer:   smtpServer,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		from:         from,
	}
}

func (s *emailService) SendEmail(email *domain.Email) error {
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpServer)

	// TLS config
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.smtpServer,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", s.smtpServer, s.smtpPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.smtpServer)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Auth
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set the sender and recipient
	if err = client.Mail(s.from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, to := range email.To {
		if err = client.Rcpt(to); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", to, err)
		}
	}

	// Send the email body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	headers := fmt.Sprintf("From: %s\r\n", s.from)
	headers += fmt.Sprintf("To: %s\r\n", email.To[0])
	headers += fmt.Sprintf("Subject: %s\r\n", email.Subject)
	headers += "MIME-version: 1.0;\r\n"
	headers += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	headers += "\r\n"

	if _, err = fmt.Fprint(w, headers+email.Body); err != nil {
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return client.Quit()
}
