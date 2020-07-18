package mongo_test

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	internalMongo "gitlab.com/evzpav/documents/internal/infrastructure/storage/mongo"
	"gitlab.com/evzpav/documents/pkg/env"
	"gitlab.com/evzpav/documents/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoClient *mongo.Client
	testLog     log.Logger
)

func init() {
	testLog = log.NewZeroLog("document-test", "", log.Error)

	mongoURL := env.GetString("MONGO_URL")
	// mongoURL := "mongodb://172.26.0.2:27017"

	client, err := internalMongo.Connect(context.Background(), mongoURL)
	if err != nil {
		testLog.Fatal().Err(err).Sendf("error on MongoDB connection: %q", err)
	}
	mongoClient = client
}

func NewDatabase() *mongo.Database {
	return mongoClient.Database(primitive.NewObjectID().Hex())
}
