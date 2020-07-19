package mysql

import (
	"context"
	"strings"

	"github.com/jinzhu/gorm"

	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/internal/infrastructure/storage"
	"gitlab.com/evzpav/user-auth/pkg/errors"
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

func (us *userStorage) Insert(ctx context.Context, inputUser *domain.User) error {
	var user domain.User
	err := us.db.Where("lower(users.email) = (?)", strings.ToLower(inputUser.Email)).Find(&user).Error
	if err == nil || user.ID != 0 {
		return errors.NewDuplicatedRecord(storage.ErrUserDuplicated)
	}

	return us.db.Create(&inputUser).Error
}

func (us *userStorage) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	if err := us.db.Where(`users.email=(?)`, email).Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (us *userStorage) Update(ctx context.Context, inputUser *domain.User) error {
	return us.db.Save(&inputUser).Error
}
