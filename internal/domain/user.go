package domain

import (
	"context"
)

type User struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Email     string `json:"email"`
	Pssword   string `json:"password"`
	Telephone string `json:"telephone"`
}

// func (doc *Document) Validate() error {
// 	if doc.Value == "" {
// 		return errors.NewInvalidArgument(ErrDocumentValueRequired).
// 			WithMessage("document value is required")
// 	}

// 	return nil
// }

// type DocumentFilter struct {
// 	IsBlacklisted bool
// 	DocType       string
// }

type UserService interface {
	// Create(ctx context.Context, user *User) (*User, error)
	// GetOne(ctx context.Context, ID string) (*User, error)
	// Update(ctx context.Context, user *User) (*User, error)
	// Delete(ctx context.Context, ID string) error
}

type UserStorage interface {
	Insert(ctx context.Context, user *User) (*User, error)
	// FindOne(ctx context.Context, ID string) (*User, error)
	// Set(ctx context.Context, user *User) (*User, error)
	// Remove(ctx context.Context, ID string) error
}
