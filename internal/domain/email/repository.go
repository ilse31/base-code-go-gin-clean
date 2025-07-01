package email

type EmailService interface {
	SendEmail(email *Email) error
}
