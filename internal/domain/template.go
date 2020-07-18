package domain

import (
	"html/template"
)

type HTMLTemplate struct {
	Template *template.Template
}

type TemplateService interface {
	// Login() (*HTMLTemplate, error)
	RetrieveParsedTemplate(name string) (*HTMLTemplate, error)
}