package domain

import "fmt"

type Profile struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

func (p *Profile) Validate() error {
	if p.ID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if p.Email == "" {
		return fmt.Errorf("invalid email")
	}
	return nil
}

