package mongo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/internal/infrastructure/storage/mongo"
	"gitlab.com/evzpav/documents/pkg/env"
)

func TestNewMongoSession_Success(t *testing.T) {
	mongoURL := env.GetString("MONGO_URL")
	mongoClient, err := mongo.Connect(context.Background(), mongoURL)
	assert.NoError(t, err)
	assert.NotNil(t, mongoClient)

	database := mongo.NewDatabase(mongoClient, "documents-test")
	assert.NotNil(t, database)
}

func TestNewMongoSession_EmptyURL(t *testing.T) {
	mongoClient, err := mongo.Connect(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, mongoClient)
}

func TestNewMongoSession_InvalidURL(t *testing.T) {
	mongoClient, err := mongo.Connect(context.Background(), "192.168.0.1")
	assert.Error(t, err)
	assert.Nil(t, mongoClient)
}
