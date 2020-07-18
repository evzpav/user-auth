package domain

import "gitlab.com/evzpav/user-auth/pkg/errors"

const (
	ErrDocumentValueRequired errors.Code = "DOCUMENT_VALUE_REQUIRED"
	ErrDocumentIDRequired    errors.Code = "DOCUMENT_ID_REQUIRED"
)
