package email

// TemplateData holds the data that will be passed to email templates
type TemplateData struct {
	Subject     string
	Greeting    string
	Content     string
	ButtonURL   string
	ButtonText  string
	Footer      string
	CurrentYear int
}
