package user

import (
	"context"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

type service struct {
	storage domain.UserStorage
}

func NewService(storage domain.UserStorage) *service {
	return &service{
		storage: storage,
	}
}

func (us *service) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return  err
	}

	return us.storage.Insert(ctx, user)
}

func (us *service) Update(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return us.storage.Update(ctx, user)
}

func (us *service) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return us.storage.FindByEmail(ctx, email)
}