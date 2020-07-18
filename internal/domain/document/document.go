package document

import (
	"context"

	"gitlab.com/evzpav/documents/internal/domain"
	"gitlab.com/evzpav/documents/pkg/errors"
)

type service struct {
	storage domain.DocumentStorage
}

func NewService(storage domain.DocumentStorage) *service {
	return &service{
		storage: storage,
	}
}

func (ds *service) Create(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	if err := document.Validate(); err != nil {
		return nil, err
	}

	return ds.storage.Insert(ctx, document)
}

func (ds *service) GetAll(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
	return ds.storage.FindAll(ctx, filter, sort...)
}

func (ds *service) GetOne(ctx context.Context, ID string) (*domain.Document, error) {
	if ID == "" {
		return nil, errors.NewInvalidArgument(domain.ErrDocumentIDRequired).
			WithMessage("document ID is required")
	}

	return ds.storage.FindOne(ctx, ID)
}

func (ds *service) Update(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	err := document.Validate()
	if err != nil {
		return nil, err
	}

	return ds.storage.Set(ctx, document)
}

func (ds *service) Delete(ctx context.Context, ID string) error {
	if ID == "" {
		return errors.NewInvalidArgument(domain.ErrDocumentIDRequired).
			WithMessage("document ID is required")
	}

	return ds.storage.Remove(ctx, ID)
}
