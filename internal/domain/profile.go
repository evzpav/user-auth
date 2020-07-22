package domain

import "errors"

type Profile struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func (p *Profile) Validate() error {
	if p.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if !validateEmail(p.Email) {
		return errors.New("invalid email")
	}
	return nil
}

type AutocompletePrediction struct {
	Suggestion string `json:"suggestion"`
}
