package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sdwolfe32/slimhttp"
	"github.com/sdwolfe32/trumail/api"
	"github.com/sdwolfe32/trumail/config"
	"github.com/sirupsen/logrus"
)

func main() {
	// Set default http timeout
	http.DefaultClient = &http.Client{
		Timeout: time.Duration(config.HTTPClientTimeout) * time.Second,
	}

	// Generate a new logrus logger
	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	// Configure the logger based on the environment
	if strings.Contains(config.Env, "prod") {
		logger.Formatter = new(logrus.JSONFormatter)
		logger.Level = logrus.InfoLevel
	}
	l := logger.WithField("port", config.Port)

	// Define all required dependencies
	l.Info("Defining all service dependencies")
	hostname := retrievePTR()
	r := slimhttp.NewRouter()
	s := api.NewTrumailService(logger, hostname, config.SourceAddr)
	h := slimhttp.NewHealthcheckService(logger, hostname)

	// Bind endpoints to router
	l.Info("Binding all endpoints to the router")
	r.HandleJSONEndpoint("/json/{email}", s.Lookup).Methods(http.MethodGet)
	r.HandleXMLEndpoint("/xml/{email}", s.Lookup).Methods(http.MethodGet)
	r.HandleJSONEndpoint("/healthcheck", h.Healthcheck).Methods(http.MethodGet)

	if config.ServeWeb {
		// Set all remaining paths to point to static files (must come after)
		l.Info("Serving web UI on index")
		r.HandleStatic("/", "./web")
	}

	// Listen and Serve
	l.Info("Listening and Serving")
	r.ListenAndServe(config.Port)
}

// retrievePTR attempts to retrieve the PTR record for the IP
// address retrieved via an API call on api.ipify.org
func retrievePTR() string {
	// Request the IP from ipify
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		log.Fatal("Failed to retrieve IP from api.ipify.org")
	}
	defer resp.Body.Close()

	// Decodes the IP response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read IP response body")
	}

	// Retrieve the PTR record for our IP and return without a trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
