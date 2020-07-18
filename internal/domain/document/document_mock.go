package document

import (
	"context"

	"gitlab.com/evzpav/user-auth/internal/domain"
)

// ServiceMock is the mock implementation of DocumentService
type ServiceMock struct {
	CreateInvokedCount int
	GetAllInvokedCount int
	GetOneInvokedCount int
	UpdateInvokedCount int
	DeleteInvokedCount int
	CreateFn           func(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	GetAllFn           func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error)
	GetOneFn           func(ctx context.Context, ID string) (*domain.Document, error)
	UpdateFn           func(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	DeleteFn           func(ctx context.Context, ID string) error
}

func (sm *ServiceMock) Create(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	sm.CreateInvokedCount++
	return sm.CreateFn(ctx, doc)
}

func (sm *ServiceMock) GetAll(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
	sm.GetAllInvokedCount++
	return sm.GetAllFn(ctx, filter, sort...)
}

func (sm *ServiceMock) GetOne(ctx context.Context, ID string) (*domain.Document, error) {
	sm.GetOneInvokedCount++
	return sm.GetOneFn(ctx, ID)
}

func (sm *ServiceMock) Update(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	sm.UpdateInvokedCount++
	return sm.UpdateFn(ctx, doc)
}

func (sm *ServiceMock) Delete(ctx context.Context, ID string) error {
	sm.DeleteInvokedCount++
	return sm.DeleteFn(ctx, ID)
}

type StorageMock struct {
	InsertInvokedCount  int
	FindAllInvokedCount int
	FindOneInvokedCount int
	SetInvokedCount     int
	RemoveInvokedCount  int
	InsertFn            func(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	FindAllFn           func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error)
	FindOneFn           func(ctx context.Context, ID string) (*domain.Document, error)
	SetFn               func(ctx context.Context, doc *domain.Document) (*domain.Document, error)
	RemoveFn            func(ctx context.Context, ID string) error
}

func (sm *StorageMock) Insert(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	sm.InsertInvokedCount++
	return sm.InsertFn(ctx, doc)
}

func (sm *StorageMock) FindAll(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
	sm.FindAllInvokedCount++
	return sm.FindAllFn(ctx, filter, sort...)
}

func (sm *StorageMock) FindOne(ctx context.Context, ID string) (*domain.Document, error) {
	sm.FindOneInvokedCount++
	return sm.FindOneFn(ctx, ID)
}

func (sm *StorageMock) Set(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	sm.SetInvokedCount++
	return sm.SetFn(ctx, doc)
}

func (sm *StorageMock) Remove(ctx context.Context, ID string) error {
	sm.RemoveInvokedCount++
	return sm.RemoveFn(ctx, ID)
}
