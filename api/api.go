package api

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/sdwolfe32/trumail/config"
)

// RegisterEndpoints bind endpoints to the router
func RegisterEndpoints(r Router, s Lookuper) {
	r.HandleEndpoint("/{format:(?:xml|json)}/{email}", s.Lookup).Methods(http.MethodGet)
	r.HandleEndpoint("/healthcheck", s.Healthcheck).Methods(http.MethodGet)
}

// Initialize service Builder and Lookuper
func Initialize(logger *logrus.Logger) (Router, Lookuper) {
	host := retrievePTR()
	router := NewRouter(config.RateLimit)
	service := NewLookuper(logger, host, config.SourceAddr)

	return router, service
}

// retrievePTR attempts to retrieve the PTR record for the IP address retrieved
// via an API call on api.ipify.org
func retrievePTR() string {
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

	// Retrieves the PTR record for our IP and returns without the trailing dot
	names, err := net.LookupAddr(string(data))
	if err != nil {
		return string(data)
	}
	return strings.TrimSuffix(names[0], ".")
}
