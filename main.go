package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/entrik/httpclient"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/verifier"
)

var (
	// port defines the port used by the api server
	port = getEnv("PORT", "8080")
	// sourceAddr defines the address used on verifier
	sourceAddr = getEnv("SOURCE_ADDR", "admin@gmail.com")
)

func main() {
	// Declare the router
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the API Services
	v := verifier.NewVerifier(retrievePTR(), sourceAddr)

	// Bind the API endpoints to router
	e.GET("/v1/:format/:email", api.LookupHandler(v), authMiddleware)
	e.GET("/v1/health", api.HealthHandler(), authMiddleware)

	// Listen and Serve
	e.Logger.Fatal(e.Start(":" + port))
}

// RetrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	// Request the IP from ipify
	ip, err := httpclient.GetString("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve public IP")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(ip)
	if err != nil {
		return ip
	}
	return strings.TrimSuffix(names[0], ".")
}

// authMiddleware verifies the auth token on the request matches the
// one defined in the environment
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// authToken is the token that must be used on all requests
	authToken := getEnv("AUTH_TOKEN", "")

	// Return the Handlerfunc that asserts the auth token
	return func(c echo.Context) error {
		if authToken != "" {
			if c.Request().Header.Get("X-Auth-Token") == authToken {
				return next(c)
			}
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// getEnv retrieves variables from the environment and falls back
// to a passed fallback variable if it isn't set
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
