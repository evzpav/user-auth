package mysql

import (
	"context"

	"github.com/jinzhu/gorm"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

type userStorage struct {
	db  *gorm.DB
	log log.Logger
}

func NewUserStorage(db *gorm.DB, log log.Logger) (*userStorage, error) {
	return &userStorage{
		db:  db,
		log: log,
	}, nil
}

func (us *userStorage) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	return nil, nil
}
