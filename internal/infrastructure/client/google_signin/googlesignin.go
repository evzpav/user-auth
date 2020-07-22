package googlesignin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleClient struct {
	config *oauth2.Config
}

func New(key, secret, redirectURL string) *GoogleClient {
	conf := &oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleClient{
		config: conf,
	}
}

func (c *GoogleClient) GetLoginURL(state string) string {
	return c.config.AuthCodeURL(state)
}

func (c *GoogleClient) GetProfile(code string) (*domain.GoogleUser, error) {
	token, err := c.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %v", err)
	}

	client := c.config.Client(oauth2.NoContext, token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	defer userInfo.Body.Close()
	data, err := ioutil.ReadAll(userInfo.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body in google user info: %v", err)
	}

	var user domain.GoogleUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil

}
