package template

import (
	"go/build"
	"html/template"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

type service struct {
	templatesPath string
}

func NewService() *service {
	return &service{
		templatesPath: build.Default.GOPATH + "/src/gitlab.com/evzpav/user-auth/internal/domain/template/pages/",
	}
}

func (s *service) RetrieveParsedTemplate(name string) (*domain.HTMLTemplate, error) {
	tpl := template.Must(template.ParseGlob(s.templatesPath + "*"))
	pageTpl, err := tpl.ParseFiles(s.templatesPath+"base.html", s.templatesPath+name+".html")
	if err != nil {
		return nil, err
	}

	return &domain.HTMLTemplate{
		Template: pageTpl,
	}, nil

}