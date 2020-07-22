package domain

import "golang.org/x/oauth2"

type GoogleSigner interface {
	GetLoginURL(state string) string
	GetToken(code string) (*oauth2.Token, error)
	GetProfile(token *oauth2.Token) (*User, error)
}
