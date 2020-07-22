package domain

type GoogleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

type GoogleSigner interface {
	GetLoginURL(state string) string
	GetProfile(code string) (*GoogleUser, error)
}

type GoogleMapper interface {
	GetAddressSuggestion(input string) (*AutocompletePrediction, error)
}
