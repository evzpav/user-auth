package domain

import (
	"html/template"
)

type HTMLTemplate struct {
	Template *template.Template
}

type TemplateService interface {
	RetrieveParsedTemplate(name string) (*HTMLTemplate, error)
}