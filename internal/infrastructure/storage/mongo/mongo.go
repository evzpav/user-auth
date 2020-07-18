package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Connect returns a mongo client already connected
func Connect(ctx context.Context, mongoURL string) (*mongo.Client, error) {
	if strings.TrimSpace(mongoURL) == "" {
		return nil, errors.New("mongoURL is required")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, fmt.Errorf("error on MongoDB connection: %q", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error on MongoDB connection: %q", err)
	}

	return client, nil
}

// NewDatabase returns a new mongo database
func NewDatabase(client *mongo.Client, databaseName string) *mongo.Database {
	return client.Database(databaseName)
}

// IsDuplicatedError check if error is deplicated mongo type
func IsDuplicatedError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

// IsNotFoundError check if error is not found mongo type
func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "no documents in result")
}

// CreateIndex create index
func CreateIndex(collection *mongo.Collection, fields []string, unique bool) error {
	keys := bsonx.Doc{}
	for _, field := range fields {
		keys = keys.Append(field, bsonx.Int32(1))
	}

	opt := options.Index()
	opt.SetBackground(true)

	if unique {
		opt.SetUnique(true)
	}

	index := mongo.IndexModel{
		Keys:    keys,
		Options: opt,
	}

	_, err := collection.Indexes().CreateOne(context.Background(), index)

	return err
}
