package storage

import "gitlab.com/evzpav/documents/pkg/errors"

const (
	ErrDocumentNotFound   errors.Code = "DOCUMENT_NOT_FOUND"
	ErrDocumentDuplicated errors.Code = "DOCUMENT_DUPLICATED"
)
