package domain

import (
	"context"
	"fmt"
)

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Address       string `json:"address"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	Token         string `json:"token"`
	RecoveryToken string `json:"recovery_token"`
	GoogleID      string `json:"google_id"`
}

func (u *User) Validate() error {
	if !validateEmail(u.Email) {
		return fmt.Errorf("invalid email")
	}

	if len(u.Password) < 5 {
		return fmt.Errorf("invalid password")
	}

	return nil
}

type UserService interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByToken(ctx context.Context, token string) (*User, error)
	FindByRecoveryToken(ctx context.Context, token string) (*User, error)
	FindByGoogleID(ctx context.Context, token string) (*User, error)
	FindByID(ctx context.Context, ID int) (*User, error)
	Update(ctx context.Context, user *User) error
}

type UserStorage interface {
	Insert(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByToken(ctx context.Context, token string) (*User, error)
	FindByRecoveryToken(ctx context.Context, token string) (*User, error)
	FindByGoogleID(ctx context.Context, token string) (*User, error)
	FindByID(ctx context.Context, ID int) (*User, error)
	Update(ctx context.Context, user *User) error
}
