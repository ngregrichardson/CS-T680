package utils

import "os"

func GetEnvironmentVariable(key string, fallback string) string {
	value, exists := os.LookupEnv(key)

	if exists {
		return value
	}

	return fallback
}
