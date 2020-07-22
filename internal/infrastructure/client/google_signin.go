package googlesignin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleClient struct {
	config *oauth2.Config
}

type User struct {
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

func (c *GoogleClient) GetToken(code string) (*oauth2.Token, error) {
	return c.config.Exchange(oauth2.NoContext, code)
}

func (c *GoogleClient) GetProfile(token *oauth2.Token) (*User, error) {
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

	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %v", err)
	}

	return &user, nil

}
