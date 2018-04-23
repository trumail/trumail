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
	"github.com/sirupsen/logrus"
)

func main() {
	// Generate a new logrus logger
	logger := logrus.New()
	logger.Level = logrus.DebugLevel
	l := logger.WithField("port", config.Port)

	// Define all required dependencies
	l.Info("Defining all service dependencies")
	e := echo.New()
	v := verifier.NewVerifier(retrievePTR(), config.SourceAddr)
	s := api.NewService(logger, config.HTTPClientTimeout, v)

	// Bind endpoints to router
	l.Info("Binding API endpoints to the router")
	e.GET("/v1/:format/:email", s.Lookup)
	e.GET("/v1/health", s.Health)

	// Listen and Serve
	l.WithField("port", config.Port).Info("Listening and Serving")
	l.Fatal(e.Start(":" + config.Port))
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
