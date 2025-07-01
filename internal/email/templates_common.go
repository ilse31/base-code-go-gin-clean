package email

import (
	"bytes"
	"fmt"
	"time"
)

// WelcomeEmail creates a welcome email template
func WelcomeEmail(recipientName, loginURL string) (subject, body string, err error) {
	templateData := TemplateData{
		Subject:  "Welcome to Our Platform!",
		Greeting: "Hello " + recipientName,
		Content: "<p>Thank you for joining our platform! We're excited to have you on board.</p>" +
			"<p>Start exploring all the features we have to offer and make the most of your experience.</p>",
		ButtonURL:   loginURL,
		ButtonText:  "Go to Dashboard",
		Footer:      "If you did not create an account, please contact our support team immediately.",
		CurrentYear: time.Now().Year(),
	}

	body, err = generateEmailFromTemplate("welcome.html", templateData)
	if err != nil {
		return "", "", err
	}

	return templateData.Subject, body, nil
}

// PasswordResetEmail creates a password reset email template
func PasswordResetEmail(recipientName, resetURL string) (subject, body string, err error) {
	templateData := TemplateData{
		Subject:  "Password Reset Request",
		Greeting: "Hello " + recipientName,
		Content: "<p>We received a request to reset your password. Click the button below to set a new password.</p>" +
			"<p>If you didn't request this, you can safely ignore this email.</p>",
		ButtonURL:   resetURL,
		ButtonText:  "Reset Password",
		Footer:      "This password reset link will expire in 24 hours.",
		CurrentYear: time.Now().Year(),
	}

	body, err = generateEmailFromTemplate("password_reset.html", templateData)
	if err != nil {
		return "", "", err
	}

	return templateData.Subject, body, nil
}

// DailyReportEmail creates a daily report email template
func DailyReportEmail(recipientName string, reportData map[string]interface{}) (subject, body string, err error) {
	// Example report data usage
	newUsers := "0"
	if val, ok := reportData["newUsers"]; ok {
		newUsers = val.(string)
	}

	templateData := TemplateData{
		Subject:  "Your Daily Report",
		Greeting: "Hello " + recipientName,
		Content: "<p>Here's your daily activity summary:</p>" +
			"<ul>" +
			"<li>New Users: " + newUsers + "</li>" +
			"</ul>",
		Footer:      "This is an automated message, please do not reply to this email.",
		CurrentYear: time.Now().Year(),
	}

	body, err = generateEmailFromTemplate("daily_report.html", templateData)
	if err != nil {
		return "", "", err
	}
	return templateData.Subject, body, nil
}

// VerifyEmail creates an email verification template
func VerifyEmail(recipientName, verificationURL string) (subject, body string, err error) {
	templateData := TemplateData{
		Subject:  "Verify Your Email Address",
		Greeting: "Hello " + recipientName,
		Content: "<p>Thank you for signing up! Please verify your email address by clicking the button below.</p>" +
			"<p>This link will expire in 24 hours.</p>",
		ButtonURL:   verificationURL,
		ButtonText:  "Verify Email",
		Footer:      "If you didn't create an account, you can safely ignore this email.",
		CurrentYear: time.Now().Year(),
	}

	body, err = generateEmailFromTemplate("verification_email.html", templateData)
	if err != nil {
		return "", "", err
	}

	return templateData.Subject, body, nil
}

// generateEmailFromTemplate is a helper function to render email templates
func generateEmailFromTemplate(templateName string, data TemplateData) (string, error) {
	var buf bytes.Buffer
	err := Templates.ExecuteTemplate(&buf, templateName, data)
	if err != nil {
		return "", err
	}

	fmt.Println("==== BUFFER ====")
	fmt.Println(buf.String())

	return buf.String(), nil
}
