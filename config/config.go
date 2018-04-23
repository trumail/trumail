package config

import (
	"os"
	"strconv"
)

var (
	// Port defines the port used by the api server
	Port = getEnv("PORT", "8080")
	// SourceAddr defines the address used on verifier
	SourceAddr = getEnv("SOURCE_ADDR", "admin@gmail.com")
	// HTTPClientTimeout defines the HTTP client timeout used in requests
	HTTPClientTimeout, _ = strconv.Atoi(getEnv("HTTP_CLIENT_TIMEOUT", "25"))
)

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
