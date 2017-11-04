package config

import (
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var (
	// RateLimit defines if the response on router should be throttled
	RateLimit, _ = strconv.ParseBool(getEnv("RATE_LIMIT", "false"))

	// SourceAddr defines the address used on verifier
	SourceAddr = getEnv("SOURCE_ADDR", "admin@gmail.com")

	// ServeWeb defines if the web static site should be served
	ServeWeb, _ = strconv.ParseBool(getEnv("SERVE_WEB", "false"))

	// Port defines the port used by the api server
	Port = getEnv("PORT", "8000")

	// Env defines the environment where the service is being ran
	Env = getEnv("ENVIRONMENT", "development")
)
