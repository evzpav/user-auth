package env

import (
	"os"
	"strconv"

	"gitlab.com/evzpav/documents/pkg/log"
)

// GetString ...
func GetString(envVar string, defaultValue ...string) string {
	value := os.Getenv(envVar)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}

	return value
}

// GetInt ...
func GetInt(envVar string, defaultValue int) int {
	if valueStr := os.Getenv(envVar); valueStr != "" {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}

// CheckRequired ...
func CheckRequired(log log.Logger, envVarArgs ...string) {
	for _, envVar := range envVarArgs {
		if os.Getenv(envVar) == "" {
			log.Fatal().Sendf("Environment variable '%s' is required.", envVar)
			continue
		}

		log.Info().Sendf("Environment variable '%s' is ok.", envVar)
	}
}

// CheckRequiredAnyExists ...
func CheckRequiredAnyExists(log log.Logger, envVarArgs ...string) {
	var exists bool
	for _, envVar := range envVarArgs {
		if os.Getenv(envVar) != "" {
			exists = true
			log.Info().Sendf("Environment variable '%s' is ok.", envVar)
		}
	}

	if !exists {
		log.Fatal().Sendf("Environment variable '%s' is required.", envVarArgs)
	}

}
