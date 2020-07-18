package document_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	document "gitlab.com/evzpav/user-auth/internal/domain/document"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

func TestDocumentService_Insert(t *testing.T) {
	t.Run("should create a document successfully", func(t *testing.T) {
		docForInsert := &domain.Document{
			ID:    "123123123",
			Value: "543543543",
		}

		docStorageMock := &document.StorageMock{
			InsertFn: func(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
				return doc, nil
			},
		}

		docService := document.NewService(docStorageMock)
		doc, err := docService.Create(context.Background(), docForInsert)
		assert.NoError(t, err)
		assert.Equal(t, 1, docStorageMock.InsertInvokedCount)
		assert.Equal(t, "123123123", doc.ID)
		assert.Equal(t, "543543543", doc.Value)
	})

	t.Run("should return error with empty document ID", func(t *testing.T) {
		docForInsert := &domain.Document{}

		docService := document.NewService(nil)
		doc, err := docService.Create(context.Background(), docForInsert)
		assert.Nil(t, doc)
		assert.Equal(t, "<DOCUMENT_VALUE_REQUIRED> document value is required", err.Error())
	})
}

func TestDocumentService_FindOne(t *testing.T) {
	t.Run("should find one document successfully", func(t *testing.T) {
		docStorageMock := &document.StorageMock{
			FindOneFn: func(ctx context.Context, ID string) (*domain.Document, error) {
				return &domain.Document{ID: ID}, nil
			},
		}

		docService := document.NewService(docStorageMock)
		doc, err := docService.GetOne(context.Background(), "docId")

		assert.Equal(t, 1, docStorageMock.FindOneInvokedCount)
		assert.NoError(t, err)
		assert.Equal(t, "docId", doc.ID)
	})

	t.Run("should return error with empty document ID", func(t *testing.T) {
		docService := document.NewService(nil)
		doc, err := docService.GetOne(context.Background(), "")

		assert.Error(t, err)
		assert.Equal(t, "<DOCUMENT_ID_REQUIRED> document ID is required", err.Error())
		assert.Nil(t, doc)
	})

	t.Run("should return error", func(t *testing.T) {
		docStorageMock := &document.StorageMock{
			FindOneFn: func(ctx context.Context, ID string) (*domain.Document, error) {
				return nil, errors.New("internal error")
			},
		}

		docService := document.NewService(docStorageMock)
		doc, err := docService.GetOne(context.Background(), "docId")
		assert.Error(t, err)
		assert.Equal(t, "internal error", err.Error())
		assert.Equal(t, 1, docStorageMock.FindOneInvokedCount)
		assert.Nil(t, doc)
	})
}

func TestDocumentService_FindAll(t *testing.T) {
	t.Run("should find all documents successfully", func(t *testing.T) {
		objs := []*domain.Document{
			{ID: "doc1"},
			{ID: "doc2"},
		}

		docStorageMock := &document.StorageMock{
			FindAllFn: func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
				return objs, nil
			},
		}

		docService := document.NewService(docStorageMock)
		docs, err := docService.GetAll(context.Background(), nil, nil)

		assert.Equal(t, 1, docStorageMock.FindAllInvokedCount)
		assert.NoError(t, err)
		assert.Equal(t, len(docs), 2)
		assert.Equal(t, docs[0].ID, "doc1")
		assert.Equal(t, docs[1].ID, "doc2")
	})

	t.Run("should return internal error", func(t *testing.T) {
		docStorageMock := &document.StorageMock{
			FindAllFn: func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
				return nil, errors.New("internal error")
			},
		}

		docService := document.NewService(docStorageMock)
		docs, err := docService.GetAll(context.Background(), nil, nil)

		assert.Equal(t, 1, docStorageMock.FindAllInvokedCount)
		assert.Equal(t, "internal error", err.Error())
		assert.Nil(t, docs)
	})
}

func TestDocumentService_Update(t *testing.T) {
	t.Run("should update a document successfully", func(t *testing.T) {
		docForUpdate := &domain.Document{
			ID:    "123123123",
			Value: "543543543",
		}

		docStorageMock := &document.StorageMock{
			SetFn: func(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
				return doc, nil
			},
		}

		docService := document.NewService(docStorageMock)
		doc, err := docService.Update(context.Background(), docForUpdate)
		assert.NoError(t, err)
		assert.Equal(t, 1, docStorageMock.SetInvokedCount)
		assert.Equal(t, "123123123", doc.ID)
		assert.Equal(t, "543543543", doc.Value)
	})

	t.Run("should return error with empty document ID", func(t *testing.T) {
		docForUpdate := &domain.Document{}

		docService := document.NewService(nil)
		doc, err := docService.Update(context.Background(), docForUpdate)
		assert.Nil(t, doc)
		assert.Equal(t, "<DOCUMENT_VALUE_REQUIRED> document value is required", err.Error())
	})
}

func TestDocumentService_Remove(t *testing.T) {
	t.Run("should delete a document successfully", func(t *testing.T) {

		docStorageMock := &document.StorageMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				return nil
			},
		}

		docService := document.NewService(docStorageMock)
		err := docService.Delete(context.Background(), "docId")

		assert.Equal(t, 1, docStorageMock.RemoveInvokedCount)
		assert.NoError(t, err)
	})

	t.Run("should return error with empty document ID", func(t *testing.T) {
		docStorageMock := &document.StorageMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				return nil
			},
		}

		docService := document.NewService(docStorageMock)
		err := docService.Delete(context.Background(), "")

		assert.Equal(t, 0, docStorageMock.RemoveInvokedCount)
		assert.Equal(t, "<DOCUMENT_ID_REQUIRED> document ID is required", err.Error())
	})

	t.Run("should return error from document service", func(t *testing.T) {
		docStorageMock := &document.StorageMock{
			RemoveFn: func(ctx context.Context, ID string) error {
				return errors.New("internal error")
			},
		}

		docService := document.NewService(docStorageMock)
		err := docService.Delete(context.Background(), "docId")
		assert.Error(t, err)
		assert.Equal(t, "internal error", err.Error())
	})

}
