package email

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	// Templates holds all parsed email templates
	Templates *template.Template
)

func init() {
	// Parse all template files
	tmpl, err := template.New("emails").Funcs(template.FuncMap{
		"safeHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
	}).ParseFS(templateFS, "templates/*.html")

	if err != nil {
		panic("failed to parse email templates: " + err.Error())
	}

	Templates = tmpl
}
