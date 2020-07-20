package domain

import (
	"context"
	"strings"
)

type AuthUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Errors   map[string]string
}

func NewAuthUser(email, password string) *AuthUser {
	return &AuthUser{
		Email:    email,
		Password: password,
		Errors:   make(map[string]string),
	}
}

func (au *AuthUser) ValidateEmail() bool {
	if !validateEmail(au.Email) {
		au.Errors["Email"] = "Please enter a valid email address"
	}

	return len(au.Errors) == 0
}

func (au *AuthUser) Validate() bool {
	au.Errors = make(map[string]string)

	if !validateEmail(au.Email) {
		au.Errors["Email"] = "Please enter a valid email address"
	}

	if len(strings.TrimSpace(au.Password)) < 5 {
		au.Errors["Password"] = "Please enter minimum 5 characters password"
	}

	return len(au.Errors) == 0
}

type AuthService interface {
	Signup(ctx context.Context, authUser *AuthUser) error
	Authenticate(ctx context.Context, authUser *AuthUser) error
	Authenticate2(ctx context.Context, authUser *AuthUser) (*User, error)
	Me(ctx context.Context) (*User, error)
	// AuthenticateToken(ctx context.Context, token string) error
	AuthenticateToken(ctx context.Context, token string) (*User, error)
	SendEmail(ctx context.Context, message, to string) error
}
