package config

import (
	"os"
	"strconv"
)

var (
	// SourceAddr defines the address used on verifier
	SourceAddr = getEnv("SOURCE_ADDR", "admin@gmail.com")
	// ServeWeb defines if the web static site should be served
	ServeWeb, _ = strconv.ParseBool(getEnv("SERVE_WEB", "false"))
	// Port defines the port used by the api server
	Port = getEnv("PORT", "8080")
	// Env defines the environment where the service is being ran
	Env = getEnv("ENVIRONMENT", "development")
	// HTTPClientTimeout defines the HTTP client timeout used in requests
	HTTPClientTimeout, _ = strconv.Atoi(getEnv("HTTP_CLIENT_TIMEOUT", "25"))
	// RateLimitCIDRCustom defines an array of cidr you want to exclude from rate limit "IP|max|hours" example: 192.168.0.0/16|0|0,172.16.0.0/12|0|0,10.0.0.0/8|0|0
	RateLimitCIDRCustom = getEnv("RATE_LIMIT_CIDR", "")
	// RateLimitMax is the maximum number of requests allowed in the
	// specified interval
	RateLimitMax, _ = strconv.ParseInt(getEnv("RATE_LIMIT_MAX", ""), 10, 64)
	// RateLimitHours is the interval in which requests will be rate limited
	RateLimitHours, _ = strconv.ParseInt(getEnv("RATE_LIMIT_HOURS", ""), 10, 64)
)

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
