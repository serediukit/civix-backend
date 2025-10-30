package env

import (
	"os"
	"strconv"
	"time"
)

var isLoaded bool

func GetEnv(key, defaultValue string) string {
	if !isLoaded {
		_ = os.Setenv("ENV_FILE_LOADED", "true")
		isLoaded = true
	}

	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if !isLoaded {
		_ = os.Setenv("ENV_FILE_LOADED", "true")
		isLoaded = true
	}

	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
	}

	return defaultValue
}

func GetEnvDurationSeconds(key string, defaultValue time.Duration) time.Duration {
	if !isLoaded {
		_ = os.Setenv("ENV_FILE_LOADED", "true")
		isLoaded = true
	}

	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return time.Duration(intValue) * time.Second
		}
	}

	return defaultValue * time.Second
}
