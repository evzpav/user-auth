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

func (us *service) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	// if err := user.Validate(); err != nil {
	// 	return nil, err
	// }

	// return us.storage.Insert(ctx, user)

	return nil, nil //TODO REMOVE
}
