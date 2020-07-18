package mongo

import (
	"context"
	"fmt"

	"gitlab.com/evzpav/documents/pkg/log"

	"gitlab.com/evzpav/documents/internal/domain"
	"gitlab.com/evzpav/documents/internal/infrastructure/storage"
	"gitlab.com/evzpav/documents/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type documentStorage struct {
	collection *mongo.Collection
	log        log.Logger
}

// NewDocumentStorage creates documents collection and returns an instance of documentService
func NewDocumentStorage(db *mongo.Database, log log.Logger) (*documentStorage, error) {
	collection := db.Collection("documents")

	if err := CreateIndex(collection, []string{"value"}, true); err != nil {
		return nil, err
	}

	return &documentStorage{
		collection: collection,
		log:        log,
	}, nil
}

func (ds *documentStorage) Insert(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	if document == nil {
		return nil, fmt.Errorf("failed to insert nil document")
	}
	
	document.ID = primitive.NewObjectID().Hex()

	if _, err := ds.collection.InsertOne(ctx, document); err != nil {
		if IsDuplicatedError(err) {
			return nil, errors.NewDuplicatedRecord(storage.ErrDocumentDuplicated).
				WithMessage("document already exists")
		}

		return nil, err
	}

	return document, nil
}

func (ds *documentStorage) FindAll(ctx context.Context, filter *domain.DocumentFilter, storageSorts ...*domain.StorageSort) ([]*domain.Document, error) {

	fieldsSort := resolveSort(storageSorts)
	queryFilter := resolveDocumentFilters(filter)

	opt := options.Find().SetSort(fieldsSort)

	cur, err := ds.collection.Find(ctx, queryFilter, opt)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cur.Close(ctx); err != nil {
			ds.log.Error().Err(err).Sendf("error closing mongo cursor: %q", err)
		}
	}()

	objs := make([]*domain.Document, 0)

	for cur.Next(ctx) {
		var document domain.Document

		if err = cur.Decode(&document); err != nil {
			return nil, err
		}

		objs = append(objs, &document)
	}

	return objs, nil
}

func (ds *documentStorage) FindOne(ctx context.Context, ID string) (*domain.Document, error) {
	result := &domain.Document{}
	if err := ds.collection.FindOne(ctx, bson.M{"_id": ID}).Decode(result); err != nil {
		if IsNotFoundError(err) {
			return nil, errors.NewNotFound(storage.ErrDocumentNotFound).
				WithMessagef("document not found: %s", ID).
				WithArg("id", ID)

		}

		return nil, err
	}

	return result, nil
}

func (ds *documentStorage) Set(ctx context.Context, document *domain.Document) (*domain.Document, error) {
	key := bson.M{
		"_id": document.ID,
	}

	result, err := ds.collection.UpdateOne(ctx, key, bson.M{"$set": document})
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.NewNotFound(storage.ErrDocumentNotFound).
			WithMessagef("document not found: %s", document.ID)
	}

	return document, nil
}

func (ds *documentStorage) Remove(ctx context.Context, ID string) error {
	result, err := ds.collection.DeleteOne(ctx, bson.M{"_id": ID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.NewNotFound(storage.ErrDocumentNotFound).
			WithMessagef("document not found: %s", ID).
			WithArg("id", ID)
	}

	return nil
}

func resolveSort(storageSorts []*domain.StorageSort) bson.D {
	fieldsSort := make(bson.D, 0)
	for _, storageSort := range storageSorts {
		order := 1
		if storageSort.Type == domain.SortDesc {
			order = -1
		}

		fieldsSort = append(fieldsSort, bson.E{
			Key:   storageSort.Attribute,
			Value: order,
		})
	}
	return fieldsSort
}

func resolveDocumentFilters(filter *domain.DocumentFilter) bson.M {
	filters := bson.M{}

	if filter == nil {
		return filters
	}

	if filter.IsBlacklisted {
		filters["is_blacklisted"] = bson.M{"$in": filter.IsBlacklisted}
	}

	if filter.DocType != "" {
		filters["doc_type"] = filter.DocType
	}

	return filters
}
