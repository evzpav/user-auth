package domain

import (
	"context"

	"gitlab.com/evzpav/documents/pkg/errors"
)

type Document struct {
	ID            string `json:"id,omitempty" bson:"_id,omitempty"`
	DocType       string `json:"doc_type" bson:"doc_type"`
	IsBlacklisted bool   `json:"is_blacklisted" bson:"is_blacklisted"`
	Value         string `json:"value" bson:"value"`
}

func (doc *Document) Validate() error {
	if doc.Value == "" {
		return errors.NewInvalidArgument(ErrDocumentValueRequired).
			WithMessage("document value is required")
	}

	return nil
}

type DocumentFilter struct {
	IsBlacklisted bool
	DocType       string
}

type DocumentService interface {
	Create(ctx context.Context, doc *Document) (*Document, error)
	GetAll(ctx context.Context, filter *DocumentFilter, sort ...*StorageSort) ([]*Document, error)
	GetOne(ctx context.Context, ID string) (*Document, error)
	// Count(ctx context.Context, filter *DocumentFilter) (int64, error)
	Update(ctx context.Context, document *Document) (*Document, error)
	Delete(ctx context.Context, ID string) error
}

type DocumentStorage interface {
	Insert(ctx context.Context, doc *Document) (*Document, error)
	FindAll(ctx context.Context, filter *DocumentFilter, sort ...*StorageSort) ([]*Document, error)
	FindOne(ctx context.Context, ID string) (*Document, error)
	// Count(ctx context.Context, filter *DocumentFilter) (int64, error)
	Set(ctx context.Context, document *Document) (*Document, error)
	Remove(ctx context.Context, ID string) error
}
