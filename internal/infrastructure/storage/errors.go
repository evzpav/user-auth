package storage

import "gitlab.com/evzpav/user-auth/pkg/errors"

const (
	ErrUserNotFound   errors.Code = "USER_NOT_FOUND"
	ErrUserDuplicated errors.Code = "USER_DUPLICATED"
)
