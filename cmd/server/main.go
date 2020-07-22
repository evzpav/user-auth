package main

import (
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/evzpav/user-auth/internal/domain/auth"
	"gitlab.com/evzpav/user-auth/internal/domain/template"
	"gitlab.com/evzpav/user-auth/internal/domain/user"
	googlesignin "gitlab.com/evzpav/user-auth/internal/infrastructure/client/google_signin"
	"gitlab.com/evzpav/user-auth/internal/infrastructure/server/http"
	mysql "gitlab.com/evzpav/user-auth/internal/infrastructure/storage/mysql"
	"gitlab.com/evzpav/user-auth/pkg/env"
	"gitlab.com/evzpav/user-auth/pkg/log"
)

const (
	envVarHost          = "HOST"
	envVarPort          = "PORT"
	envVarPlatformURL   = "PLATFORM_URL"
	envVarLoggerLevel   = "LOGGER_LEVEL"
	envVarMySQLURL      = "DATABASE_URL"
	envVarEmailPassword = "EMAIL_PASSWORD"
	envVarEmailFrom     = "EMAIL_FROM"
	envVarGoogleKey     = "GOOGLE_KEY"
	envVarGoogleSecret  = "GOOGLE_SECRET"
	envVarSessionKey    = "SESSION_KEY"

	defaultProjectPort = "5001"
	defaultLoggerLevel = "info"
)

var (
	version, build, date string
)

func main() {
	log := log.NewZeroLog("user-auth", version, log.Level(getLoggerLevel()))

	log.Info().Sendf("user-auth - build:%s; date:%s", build, date)

	env.CheckRequired(log, envVarMySQLURL, envVarEmailFrom, envVarEmailPassword, envVarGoogleKey, envVarGoogleSecret, envVarPlatformURL)

	db, err := mysql.New(getMySQLURL())
	if err != nil {
		log.Fatal().Err(err).Sendf("failed to connect to mysql: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error().Err(err).Sendf("error closing database: %v", err)
		}
	}()

	// storages
	userStorage, err := mysql.NewUserStorage(db, log)
	if err != nil {
		log.Fatal().Err(err).Sendf("error creating storage: %v", err)
	}

	googleSigninClient := googlesignin.New(getGoogleKey(), getGoogleSecret(), getPlatformURL()+"/login/google/auth")

	// services
	userService := user.NewService(userStorage, log)
	authService := auth.NewService(userService, getEmailFrom(), getEmailPassword(), googleSigninClient, getPlatformURL(), log)
	templateService := template.NewService(log)

	// HTTP Server
	handler := http.NewHandler(userService, authService, templateService, getSessionKey(), log)
	server := http.New(handler, getProjectHost(), getProjectPort(), log)
	server.ListenAndServe()

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getProjectHost() string {
	return env.GetString(envVarHost)
}

func getProjectPort() string {
	return env.GetString(envVarPort, defaultProjectPort)
}

func getPlatformURL() string {
	return env.GetString(envVarPlatformURL)
}

func getLoggerLevel() string {
	return env.GetString(envVarLoggerLevel, defaultLoggerLevel)
}

func getMySQLURL() string {
	return env.GetString(envVarMySQLURL)
}

func getEmailFrom() string {
	return env.GetString(envVarEmailFrom)
}

func getEmailPassword() string {
	return env.GetString(envVarEmailPassword)
}

func getGoogleKey() string {
	return env.GetString(envVarGoogleKey)
}

func getGoogleSecret() string {
	return env.GetString(envVarGoogleSecret)
}

func getSessionKey() string {
	return env.GetString(envVarSessionKey)
}
