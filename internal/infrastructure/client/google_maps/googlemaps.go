package googlemaps

import (
	"context"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"googlemaps.github.io/maps"
)

type mapsClient struct {
	client *maps.Client
}

type Suggestion struct {
	AutoCompleteResp maps.AutocompleteResponse
}

func New(apiKey string) (*mapsClient, error) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &mapsClient{
		client: client,
	}, nil
}

func (mp *mapsClient) GetAddressSuggestion(input string) (*domain.AutocompletePrediction, error) {
	r := &maps.QueryAutocompleteRequest{
		Input:    input,
		Language: "en-US",
	}
	resp, err := mp.client.QueryAutocomplete(context.Background(), r)
	if err != nil {
		return nil, err
	}

	if len(resp.Predictions) < 1 {
		return &domain.AutocompletePrediction{
			Suggestion: "",
		}, nil
	}

	return &domain.AutocompletePrediction{
		Suggestion: resp.Predictions[0].Description,
	}, nil
}
