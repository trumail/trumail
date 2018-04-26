package main

import (
	"log"
	"net"
	"strings"

	"github.com/labstack/echo"
	"github.com/sdwolfe32/httpclient"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
	"github.com/sdwolfe32/trumail/verifier"
)

func main() {
	// Define all required dependencies
	e := echo.New()
	s := api.NewService(config.HTTPClientTimeout,
		verifier.NewVerifier(retrievePTR(), config.SourceAddr))

	// Bind endpoints to router
	e.GET("/v1/:format/:email", s.Lookup)
	e.GET("/v1/health", api.Healthcheck)

	// Listen and Serve
	log.Fatal(e.Start(":" + config.Port))
}

// retrievePTR attempts to retrieve the PTR record for the IP
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
