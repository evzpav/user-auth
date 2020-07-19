package domain

import (
	"context"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("invalid email")
	}

	if u.Password == "" {
		return fmt.Errorf("invalid password")
	}

	return nil
}

type UserService interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByToken(ctx context.Context, token string) (*User, error)
	FindByID(ctx context.Context, ID int) (*User, error)
	Update(ctx context.Context, user *User) error
}

type UserStorage interface {
	Insert(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByToken(ctx context.Context, token string) (*User, error)
	FindByID(ctx context.Context, ID int) (*User, error)
	Update(ctx context.Context, user *User) error
}
