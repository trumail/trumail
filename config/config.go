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
	// TinystatAppID is the App ID used with Tinystat
	TinystatAppID = getEnv("TINYSTAT_APP_ID", "")
	// TinystatToken is the Token used to authenticate Tinystat requests
	TinystatToken = getEnv("TINYSTAT_TOKEN", "")
	// HerokuAppID is the App ID used with Heroku
	HerokuAppID = getEnv("HEROKU_APP_ID", "")
	// HerokuToken is the Token used to authenticate Heroku requests
	HerokuToken = getEnv("HEROKU_TOKEN", "")
	// HTTPClientTimeout defines the HTTP client timeout used in requests
	HTTPClientTimeout, _ = strconv.Atoi(getEnv("HTTP_CLIENT_TIMEOUT", "20"))
)

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't already set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
