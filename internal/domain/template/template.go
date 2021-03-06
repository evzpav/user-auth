package template

import (
	"html/template"

	"os"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

type service struct {
	googleMapsClient domain.GoogleMapper
	templatesPath    string
	log              log.Logger
}

func NewService(googleMapsClient domain.GoogleMapper, log log.Logger) *service {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err)
	}

	templatesPath := pwd + "/internal/domain/template/pages/"

	return &service{
		googleMapsClient: googleMapsClient,
		templatesPath:    templatesPath,
		log:              log,
	}
}

func (s *service) RetrieveParsedTemplate(name string) (*domain.HTMLTemplate, error) {
	tpl := template.Must(template.ParseGlob(s.templatesPath + "*"))
	pageTpl, err := tpl.ParseFiles(s.templatesPath+"base.html", s.templatesPath+name+".html")
	if err != nil {
		s.log.Debug().Err(err)
		return nil, err
	}

	return &domain.HTMLTemplate{
		Template: pageTpl,
	}, nil

}

func (s *service) GetAddressSuggestion(input string) (*domain.AutocompletePrediction, error) {
	return s.googleMapsClient.GetAddressSuggestion(input)
}
