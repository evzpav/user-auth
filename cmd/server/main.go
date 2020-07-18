package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	document "gitlab.com/evzpav/documents/internal/domain/document"
	"gitlab.com/evzpav/documents/pkg/env"
	"gitlab.com/evzpav/documents/pkg/log"

	"gitlab.com/evzpav/documents/internal/infrastructure/server/http"
	"gitlab.com/evzpav/documents/internal/infrastructure/storage/mongo"
)

const (
	envVarDocumentsHost = "DOCUMENTS_HOST"
	envVarDocumentsPort = "DOCUMENTS_PORT"
	envVarLoggerLevel   = "LOGGER_LEVEL"
	envVarMongoURL      = "MONGO_URL"
	envVarDatabaseName  = "DATABASE_NAME"

	defaultDocumentsHost = ""
	defaultDocumentsPort = "5001"
	defaultLoggerLevel   = "info"
	defaultDatabaseName  = "documents"
)

var (
	version, build, date string
)

func main() {
	log := log.NewZeroLog("Documents", version, log.Level(getLoggerLevel()))

	log.Info().Sendf("Documents - build:%s; date:%s", build, date)

	// Check environments
	env.CheckRequired(log, envVarMongoURL)

	// mongo client
	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, getMongoURL())
	if err != nil {
		log.Fatal().Err(err).Sendf("error connecting database: %v", err)
		return
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Error().Err(err).Sendf("error disconnecting database: %v", err)
		}
	}()

	database := mongo.NewDatabase(mongoClient, getDatabaseName())

	// storages
	documentStorage, err := mongo.NewDocumentStorage(database, log)
	if err != nil {
		log.Fatal().Err(err).Sendf("error creating storage: %v", err)
	}

	// services
	documentService := document.NewService(documentStorage)

	// HTTP Server

	handler := http.NewHandler(documentService, log)
	server := http.New(handler, getDocumentsHost(), getDocumentsPort(), log)
	server.ListenAndServe()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getDocumentsHost() string {
	return env.GetString(envVarDocumentsHost, defaultDocumentsHost)
}

func getDocumentsPort() string {
	return env.GetString(envVarDocumentsPort, defaultDocumentsPort)
}

func getLoggerLevel() string {
	return env.GetString(envVarLoggerLevel, defaultLoggerLevel)
}

func getMongoURL() string {
	return env.GetString(envVarMongoURL)
}

func getDatabaseName() string {
	return env.GetString(envVarDatabaseName, defaultDatabaseName)
}
