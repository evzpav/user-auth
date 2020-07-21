package user

import (
	"context"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

type service struct {
	storage domain.UserStorage
	log     log.Logger
}

func NewService(storage domain.UserStorage, log log.Logger) *service {
	return &service{
		storage: storage,
		log:     log,
	}
}

func (us *service) Create(ctx context.Context, user *domain.User) error {
	if err := user.Validate(); err != nil {
		return err
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

func (us *service) FindByToken(ctx context.Context, token string) (*domain.User, error) {
	return us.storage.FindByToken(ctx, token)
}

func (us *service) FindByRecoveryToken(ctx context.Context, token string) (*domain.User, error) {
	return us.storage.FindByRecoveryToken(ctx, token)
}

func (us *service) FindByID(ctx context.Context, id int) (*domain.User, error) {
	return us.storage.FindByID(ctx, id)
}
