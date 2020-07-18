package mongo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/internal/domain"
	"gitlab.com/evzpav/documents/internal/infrastructure/storage/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDocumentStorage(t *testing.T) {
	database := NewDatabase()
	docStorage, err := mongo.NewDocumentStorage(database, testLog)
	assert.Nil(t, err)
	assert.NotNil(t, docStorage)
}

func TestDocumentStorage_Insert(t *testing.T) {
	ctx := context.Background()
	t.Run("should create a new document document on mongoDB", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		docForInsert := &domain.Document{
			ID: "docId",
		}

		doc, err := docStorage.Insert(ctx, docForInsert)
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, docForInsert.ID)

		result := &domain.Document{}
		err = database.Collection("documents").FindOne(ctx, bson.M{"_id": docForInsert.ID}).Decode(result)
		assert.NoError(t, err)
		assert.Equal(t, docForInsert, result)
	})

	t.Run("should return duplicated document error", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		docForInsert := &domain.Document{
			ID: "docId",
		}

		doc, err := docStorage.Insert(ctx, docForInsert)
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, docForInsert.ID)

		doc, err = docStorage.Insert(ctx, docForInsert)
		assert.Error(t, err)
		assert.Equal(t, "<DOCUMENT_DUPLICATED> document already exists", err.Error())
	})

	t.Run("should return error", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		_, err = docStorage.Insert(ctx, nil)
		assert.Error(t, err)
	})
}

func TestDocumentStorage_FindOne(t *testing.T) {
	ctx := context.Background()
	t.Run("should find one document successfully", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		database.Collection("documents").InsertMany(ctx, []interface{}{
			bson.M{"_id": "docId1"},
			bson.M{"_id": "docId2"},
			bson.M{"_id": "docId3"},
		})

		doc1, err := docStorage.FindOne(ctx, "docId1")
		assert.NoError(t, err)
		assert.Equal(t, "docId1", doc1.ID)

		doc2, err := docStorage.FindOne(ctx, "docId3")
		assert.NoError(t, err)
		assert.Equal(t, "docId3", doc2.ID)
	})

	t.Run("should return not found error", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		database.Collection("documents").InsertMany(ctx, []interface{}{
			bson.M{"_id": "docId1"},
			bson.M{"_id": "docId2"},
		})

		doc, err := docStorage.FindOne(ctx, "id_not_exists")
		assert.Error(t, err)
		assert.Errorf(t, err, "<DOCUMENT_NOT_FOUND> document not found: id_not_exists (id: id_not_exists)")
		assert.Nil(t, doc)
	})
}

func TestDocumentStorage_FindAll(t *testing.T) {
	ctx := context.Background()
	t.Run("should retrieve all documents", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		docForInsert := &domain.Document{
			ID: "docId",
		}

		doc, err := docStorage.Insert(ctx, docForInsert)
		assert.NoError(t, err)
		assert.Equal(t, doc.ID, docForInsert.ID)

		docs, err := docStorage.FindAll(ctx, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(docs))
		assert.Equal(t, "docId", docs[0].ID)
	})
}

func TestDocumentStorage_Remove(t *testing.T) {
	ctx := context.Background()
	t.Run("should remove a document successfully", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		_, err = database.Collection("documents").InsertOne(ctx, bson.M{"_id": "docId"})
		assert.NoError(t, err)

		docs, err := docStorage.FindAll(ctx, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(docs))

		err = docStorage.Remove(ctx, "docId")
		assert.NoError(t, err)

		docs, err = docStorage.FindAll(ctx, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(docs))
	})

	t.Run("should return not found error", func(t *testing.T) {
		database := NewDatabase()
		docStorage, err := mongo.NewDocumentStorage(database, testLog)
		assert.Nil(t, err)

		err = docStorage.Remove(ctx, "id_not_exists")
		assert.Error(t, err)
		assert.Errorf(t, err, "<DOCUMENT_NOT_FOUND> document not found: id_not_exists (id: id_not_exists)")
	})
}
