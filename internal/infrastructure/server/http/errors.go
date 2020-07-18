package http

import "gitlab.com/evzpav/user-auth/pkg/errors"

var (
	ErrInvalidBodyRequestCode errors.Code = "INVALID_BODY"

	ErrInvalidBodyRequest = errors.NewInvalidArgument(ErrInvalidBodyRequestCode).
				WithMessage("you have applied a request with an invalid body. Please review the body and check the structure")

	ErrInvalidClientIDRequestCode errors.Code = "INVALID_CLIENT_ID"

	ErrInvalidClientIDRequest = errors.NewInvalidArgument(ErrInvalidClientIDRequestCode).
					WithMessage("clientId parameter must contain only integer values")

	ErrNotAuthorizedRequestCode errors.Code = "NOT_AUTHORIZED"

	ErrNotAuthorizedRequest = errors.NewNotAuthorized(ErrNotAuthorizedRequestCode).
				WithMessage("token not authorized")
)
