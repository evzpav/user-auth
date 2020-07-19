package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/evzpav/user-auth/internal/domain/auth"
	"gitlab.com/evzpav/user-auth/internal/domain/template"
	"gitlab.com/evzpav/user-auth/internal/domain/user"
	"gitlab.com/evzpav/user-auth/pkg/env"
	"gitlab.com/evzpav/user-auth/pkg/log"

	"gitlab.com/evzpav/user-auth/internal/infrastructure/server/http"
	mysql "gitlab.com/evzpav/user-auth/internal/infrastructure/storage/mysql"
)

const (
	envVarHost        = "HOST"
	envVarPort        = "PORT"
	envVarLoggerLevel = "LOGGER_LEVEL"
	envVarMySQLURL    = "MYSQL_URL"

	defaultProjectHost = ""
	defaultProjectPort = "5001"
	defaultLoggerLevel = "info"
)

var (
	version, build, date string
)

func main() {
	log := log.NewZeroLog("user-auth", version, log.Level(getLoggerLevel()))

	log.Info().Sendf("use-auth - build:%s; date:%s", build, date)

	env.CheckRequired(log, envVarMySQLURL)

	db, err := mysql.New(getMySQLURL())
	if err != nil {
		log.Fatal().Err(err).Sendf("failed to connect to mysql: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Sendf("error closing database: %v", err)
		}
	}()
	// if err := mysql.NewMigration(getMySQLURL()).Up(); err != nil {
	// 	log.Fatal().Sendf("Could not run migrations: %v", err)
	// }

	// Check environments

	// ctx := context.Background()

	// storages
	userStorage, err := mysql.NewUserStorage(db, log)
	if err != nil {
		log.Fatal().Err(err).Sendf("error creating storage: %v", err)
	}

	// services
	userService := user.NewService(userStorage)
	authService := auth.NewService(userService)
	templateService := template.NewService()

	// HTTP Server

	handler := http.NewHandler(userService, authService, templateService, log)
	server := http.New(handler, getProjectHost(), getProjectPort(), log)
	server.ListenAndServe()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getProjectHost() string {
	return env.GetString(envVarHost, defaultProjectHost)
}

func getProjectPort() string {
	return env.GetString(envVarPort, defaultProjectPort)
}

func getLoggerLevel() string {
	return env.GetString(envVarLoggerLevel, defaultLoggerLevel)
}

func getMySQLURL() string {
	return env.GetString(envVarMySQLURL)
}
